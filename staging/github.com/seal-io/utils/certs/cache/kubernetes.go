package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	"github.com/seal-io/utils/certs"
	"github.com/seal-io/utils/pools/gopool"
	"github.com/seal-io/utils/stringx"
)

const (
	k8sManagedLabel      = "certs.seal.io/managed"
	k8sManagedGroupLabel = "certs.seal.io/group"

	k8sManagedNameSumAnno  = "certs.seal.io/name-sum"
	k8sManagedNameAnno     = "certs.seal.io/name"
	k8sManagedValueSumAnno = "certs.seal.io/value-sum"
	k8sManagedValueKey     = "value"
)

// k8sCache implements certs.Cache using the Kubernetes Secret to store the certificate data.
type k8sCache struct {
	logger klog.Logger
	cli    certs.SecretInterface
	inf    cache.SharedIndexInformer
	grp    string
}

// NewK8sCache creates a new k8sCache instance with the given client.
func NewK8sCache(ctx context.Context, group string, cli certs.SecretInterface) (certs.Cache, error) {
	lg := klog.Background().WithName("certs").WithName("k8s")

	lw := func() *cache.ListWatch {
		labelSelector := labels.FormatLabels(map[string]string{
			k8sManagedLabel:      "true",
			k8sManagedGroupLabel: group,
		})

		return &cache.ListWatch{
			ListFunc: func(options meta.ListOptions) (runtime.Object, error) {
				options.ResourceVersion = "0"
				options.LabelSelector = labelSelector
				return cli.List(ctx, options)
			},
			WatchFunc: func(options meta.ListOptions) (watch.Interface, error) {
				options.LabelSelector = labelSelector
				return cli.Watch(ctx, options)
			},
		}
	}()

	inf := cache.NewSharedIndexInformer(lw, &core.Secret{}, 1*time.Hour,
		map[string]cache.IndexFunc{
			"_": func(obj any) ([]string, error) {
				s, ok := obj.(*core.Secret)
				if !ok {
					return nil, errors.New("object is not a secret")
				}

				if s.DeletionTimestamp != nil ||
					s.Type != core.SecretTypeOpaque ||
					s.Annotations == nil || s.Data == nil {
					return nil, nil
				}

				annos, data := s.Annotations, s.Data

				if annos[k8sManagedNameAnno] == "" || annos[k8sManagedNameSumAnno] == "" ||
					data[k8sManagedValueKey] == nil || annos[k8sManagedValueSumAnno] == "" {
					return nil, nil
				}

				if sumName(annos[k8sManagedNameAnno]) != annos[k8sManagedNameSumAnno] ||
					sumValue(data[k8sManagedValueKey]) != annos[k8sManagedValueSumAnno] {
					lg.Error(nil, "invalid key %q", annos[k8sManagedNameAnno])
					return nil, nil
				}

				return []string{annos[k8sManagedNameAnno]}, nil
			},
		})

	gopool.Go(func() {
		inf.Run(ctx.Done())
	})

	// Wait for the informer to sync.
	{
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if !cache.WaitForCacheSync(ctx.Done(), inf.HasSynced) {
			return k8sCache{}, errors.New("sync informer")
		}
	}

	return k8sCache{
		logger: lg,
		cli:    cli,
		inf:    inf,
		grp:    group,
	}, nil
}

// Get reads a certificate data from the specified secret name.
func (k k8sCache) Get(ctx context.Context, name string) ([]byte, error) {
	if name == "" {
		return nil, certs.ErrCacheMiss
	}

	// Get existed secret.
	sec := k.get(ctx, name)
	if sec == nil || sec.DeletionTimestamp != nil {
		return nil, certs.ErrCacheMiss
	}

	return sec.Data[k8sManagedValueKey], nil
}

// Put writes the certificate data to specified secret name.
func (k k8sCache) Put(ctx context.Context, name string, data []byte) error {
	if name == "" || len(data) == 0 {
		return nil
	}

	sec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			GenerateName: "seal-cert-",
			Annotations: map[string]string{
				k8sManagedNameAnno:     name,
				k8sManagedNameSumAnno:  sumName(name),
				k8sManagedValueSumAnno: sumValue(data),
			},
			Labels: map[string]string{
				k8sManagedLabel:      "true",
				k8sManagedGroupLabel: k.grp,
			},
		},
		Data: map[string][]byte{
			k8sManagedValueKey: data,
		},
	}

	// Update existed secret if found.
	if asec := k.get(ctx, name); asec != nil && asec.Name != "" && asec.DeletionTimestamp == nil {
		asecCopy := asec.DeepCopy()
		asecCopy.Annotations = sec.Annotations
		asecCopy.Labels = sec.Labels
		asecCopy.Data = sec.Data
		if reflect.DeepEqual(asecCopy, asec) {
			return nil
		}

		var err error
		asec, err = k.cli.Update(ctx, asecCopy, meta.UpdateOptions{})
		if err != nil {
			if !kerrors.IsConflict(err) && !kerrors.IsNotAcceptable(err) {
				return fmt.Errorf("update secret: %w", err)
			}
			// Retry if conflict or not acceptable.
			return k.Put(ctx, name, data)
		}

		k.logger.V(5).Info("updated secret", "object", klog.KObj(asec))

		return nil
	}

	// Otherwise, create new secret.
	sec, err := k.cli.Create(ctx, sec, meta.CreateOptions{})
	if err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return fmt.Errorf("create secret: %w", err)
		}
		// Retry if already existed.
		return k.Put(ctx, name, data)
	}

	k.logger.V(5).Info("created secret", "object", klog.KObj(sec))

	return nil
}

// Delete removes the specified secret name.
func (k k8sCache) Delete(ctx context.Context, name string) error {
	if name == "" {
		return nil
	}

	// Get existed secret.
	sec := k.get(ctx, name)
	if sec == nil || sec.DeletionTimestamp != nil {
		return nil
	}

	return k.delete(ctx, sec)
}

func (k k8sCache) get(ctx context.Context, name string) *core.Secret {
	if name == "" {
		return nil
	}

	secs, err := k.inf.GetIndexer().ByIndex("_", name)
	if err != nil {
		k.logger.Error(err, "get indexed cached secrets")
	}

	switch len(secs) {
	case 0:
		// Not found.
		return nil
	case 1:
		// Found.
		return secs[0].(*core.Secret)
	default:
		// Found multiple.
	}

	// Clean up multiple secrets with the same key.
	k.logger.Error(nil, "found multiple cached secrets with the same key, going to clean",
		"objects", klog.KObjSlice(secs))

	for i := range secs {
		_ = k.delete(ctx, secs[i].(*core.Secret))
	}

	// Not found.
	return nil
}

func (k k8sCache) delete(ctx context.Context, sec *core.Secret) error {
	if sec == nil {
		return nil
	}

	// Delete existed secret.
	opts := meta.DeleteOptions{
		PropagationPolicy: ptr.To(meta.DeletePropagationBackground),
	}

	err := k.cli.Delete(ctx, sec.Name, opts)
	if err != nil && !kerrors.IsNotFound(err) {
		return fmt.Errorf("delete secret: %w", err)
	}

	k.logger.V(5).Info("deleted secret",
		"object", klog.KObj(sec))

	return nil
}

func sumName(k string) string {
	return "fnv64a:" + stringx.SumByFNV64a(k)
}

func sumValue(v []byte) string {
	return "sha224:" + stringx.SumBytesBySHA224(v)
}
