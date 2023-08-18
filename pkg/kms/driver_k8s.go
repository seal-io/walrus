package kms

import (
	"context"
	"errors"
	"fmt"
	stdpath "path"
	"time"

	"github.com/dustin/go-humanize"
	core "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/hash"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

type KubernetesOptions struct {
	// Namespace indicates the working namespace.
	Namespace string
	// Config indicates the kubernetes rest config.
	Config *rest.Config
	// Logger indicates the logger used by the driver,
	// uses default logger if not set.
	Logger log.Logger
	// RaiseNotFound indicates the function to raise not found error,
	// uses default function if not set.
	RaiseNotFound func(key string) error
}

const (
	k8sManagedLabel = "walrus.seal.io/kms"

	k8sManagedKeyLabel       = "walrus.seal.io/kms-key"
	k8sManagedKeyHashLabel   = "walrus.seal.io/kms-key-hash"
	k8sManagedValueHashLabel = "walrus.seal.io/kms-value-hash"
	k8sManagedValueKey       = "value"
)

func NewKubernetes(ctx context.Context, opts KubernetesOptions) (*KubernetesDriver, error) {
	if opts.Namespace == "" {
		return nil, errors.New("namespace is required")
	}

	logger := log.WithName("kms").WithName("k8s")
	if opts.Logger != nil {
		logger = opts.Logger.WithName("kms").WithName("k8s")
	}

	raiseNotFound := func(key string) error {
		return fmt.Errorf("not found key %s", key)
	}
	if opts.RaiseNotFound != nil {
		raiseNotFound = opts.RaiseNotFound
	}

	cfg := rest.CopyConfig(opts.Config)
	cfg.ContentType = "application/vnd.kubernetes.protobuf"

	cli, err := coreclient.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	sCli := cli.Secrets(opts.Namespace)
	sLw := func() *cache.ListWatch {
		const labelSelector = k8sManagedLabel + "=true"

		return &cache.ListWatch{
			ListFunc: func(options meta.ListOptions) (runtime.Object, error) {
				options.ResourceVersion = "0"
				options.LabelSelector = labelSelector
				return sCli.List(ctx, options)
			},
			WatchFunc: func(options meta.ListOptions) (watch.Interface, error) {
				options.LabelSelector = labelSelector
				return sCli.Watch(ctx, options)
			},
		}
	}()

	sInf := cache.NewSharedIndexInformer(sLw, &core.Secret{}, 1*time.Hour,
		map[string]cache.IndexFunc{
			"_": func(obj any) ([]string, error) {
				s, ok := obj.(*core.Secret)
				if !ok {
					return nil, errors.New("object is not a secret")
				}

				if s.DeletionTimestamp != nil ||
					s.Type != core.SecretTypeOpaque ||
					s.Labels == nil || s.Data == nil {
					return nil, nil
				}

				if s.Labels[k8sManagedKeyLabel] == "" || s.Labels[k8sManagedKeyHashLabel] == "" ||
					s.Data[k8sManagedValueKey] == nil || s.Labels[k8sManagedValueHashLabel] == "" {
					return nil, nil
				}

				key, err := decodeK8sKey(s.Labels[k8sManagedKeyLabel])
				if err != nil {
					logger.Warnf("failed to decode key %q: %v", s.Labels[k8sManagedKeyLabel], err)
					return nil, nil
				}

				if hashK8sKey(key) != s.Labels[k8sManagedKeyHashLabel] ||
					hashK8sValue(s.Data[k8sManagedValueKey]) != s.Labels[k8sManagedValueHashLabel] {
					logger.Warnf("invalid key %q", key)
					return nil, nil
				}

				ps := []string{stdpath.Join("/", key)}
				for p := stdpath.Dir(ps[len(ps)-1]); p != "/"; p = stdpath.Dir(p) {
					ps = append(ps, p)
				}
				ps = append(ps, "/")

				return ps, nil
			},
		})

	gopool.Go(func() {
		sInf.Run(ctx.Done())
	})

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if !cache.WaitForCacheSync(ctx.Done(), sInf.HasSynced) {
		return nil, fmt.Errorf("failed to sync informer: %w", err)
	}

	return &KubernetesDriver{
		cli:           sCli,
		inf:           sInf,
		logger:        logger,
		raiseNotFound: raiseNotFound,
	}, nil
}

type KubernetesDriver struct {
	cli           coreclient.SecretInterface
	inf           cache.SharedIndexInformer
	logger        log.Logger
	raiseNotFound func(key string) error
}

func (d KubernetesDriver) Get(ctx context.Context, key string) ([]byte, error) {
	if key == "" {
		return nil, d.raiseNotFound(key)
	}

	key = normalize(key)

	// Get existed secret.
	sec := d.get(ctx, key)
	if sec == nil {
		return nil, d.raiseNotFound(key)
	}

	return sec.Data[k8sManagedValueKey], nil
}

func (d KubernetesDriver) get(ctx context.Context, key string) *core.Secret {
	secs, err := d.inf.GetIndexer().ByIndex("_", key)
	if err != nil {
		d.logger.Warnf("error indexing cached secrets: %v", err)
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
	d.logger.WarnS("found multiple cached secrets with the same key, going to clean",
		"names", getK8sNamespacedName(secs...))

	for i := range secs {
		_ = d.delete(ctx, secs[i].(*core.Secret))
	}

	return nil
}

func (d KubernetesDriver) Put(ctx context.Context, key string, v []byte) error {
	backoff := wait.Backoff{
		Duration: 100 * time.Millisecond,
		Factor:   2,
		Steps:    3,
	}

	return wait.ExponentialBackoff(backoff, func() (bool, error) {
		if err := d.put(ctx, key, v); err != nil {
			if errors.Is(err, errRetry) {
				err = nil
			}

			return false, err
		}

		return true, nil
	})
}

var errRetry = errors.New("retry")

func (d KubernetesDriver) put(ctx context.Context, key string, v []byte) error {
	if key == "" || len(v) == 0 {
		return nil
	}

	key = normalize(key)

	store := func(key string, value []byte, sec *core.Secret) *core.Secret {
		// Configure labels.
		if sec.Labels == nil {
			sec.Labels = map[string]string{}
		}
		sec.Labels[k8sManagedLabel] = "true"
		sec.Labels[k8sManagedKeyLabel] = encodeK8sKey(key)
		sec.Labels[k8sManagedKeyHashLabel] = hashK8sKey(key)
		sec.Labels[k8sManagedValueHashLabel] = hashK8sValue(value)

		// Configure data.
		if sec.Data == nil {
			sec.Data = map[string][]byte{}
		}
		sec.Data[k8sManagedValueKey] = value

		return sec
	}

	// Get existed secret.
	sec := d.get(ctx, key)

	var err error

	// Update existed secret.
	if sec != nil && sec.DeletionTimestamp == nil {
		sec, err = d.cli.Update(ctx, store(key, v, sec), meta.UpdateOptions{})
		if err != nil {
			if !kerrors.IsConflict(err) && !kerrors.IsNotAcceptable(err) {
				return fmt.Errorf("error updating secret: %w", err)
			}

			// Retry if conflicted.
			return errRetry
		}

		d.logger.V(5).InfoS("updated secret",
			"namespace", sec.Namespace, "name", sec.Name, "revision", sec.ResourceVersion)

		return nil
	}

	// Create new secret.
	sec = &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			GenerateName: "walrus-kms-",
		},
	}

	sec, err = d.cli.Create(ctx, store(key, v, sec), meta.CreateOptions{})
	if err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return fmt.Errorf("error creating secret: %w", err)
		}

		// Retry if already existed.
		return errRetry
	}

	d.logger.V(5).InfoS("created secret",
		"namespace", sec.Namespace, "name", sec.Name, "revision", sec.ResourceVersion)

	return nil
}

func (d KubernetesDriver) Delete(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}

	key = normalize(key)

	// Get existed secret.
	sec := d.get(ctx, key)

	// Return directly if not existed or deleted.
	if sec == nil || sec.DeletionTimestamp != nil {
		return nil
	}

	return d.delete(ctx, sec)
}

func (d KubernetesDriver) delete(ctx context.Context, sec *core.Secret) error {
	// Delete existed secret.
	opts := meta.DeleteOptions{
		PropagationPolicy: point(meta.DeletePropagationBackground),
	}

	err := d.cli.Delete(ctx, sec.Name, opts)
	if err != nil && !kerrors.IsNotFound(err) {
		return fmt.Errorf("error deleting secret: %w", err)
	}

	d.logger.V(5).InfoS("deleted secret",
		"namespace", sec.Namespace, "name", sec.Name, "revision", sec.ResourceVersion)

	return nil
}

func (d KubernetesDriver) List(_ context.Context, path string) ([]KeyValue, error) {
	path = normalize(path)

	secs, err := d.inf.GetIndexer().ByIndex("_", path)
	if err != nil {
		return nil, fmt.Errorf("error indexing cached secrets: %w", err)
	}

	kvs := make([]KeyValue, 0, len(secs))

	for i := range secs {
		s := secs[i].(*core.Secret)

		key, err := decodeK8sKey(s.Labels[k8sManagedKeyLabel])
		if err != nil {
			d.logger.WarnS("found cached secret with invalid key",
				"name", getK8sNamespacedName(s))
			continue
		}

		kvs = append(kvs, KeyValue{
			Path:      stdpath.Dir(key),
			Key:       stdpath.Base(key),
			ValueHash: s.Labels[k8sManagedValueHashLabel],
			ValueSize: humanize.Bytes(uint64(len(s.Data[k8sManagedValueKey]))),
		})
	}

	return kvs, nil
}

func encodeK8sKey(k string) string {
	return strs.EncodeBase64(k)
}

func decodeK8sKey(k string) (string, error) {
	return strs.DecodeBase64(k)
}

func hashK8sKey(k string) string {
	return "fnv64a_" + hash.SumFnv64a(strs.ToBytes(&k))
}

func hashK8sValue(v []byte) string {
	return "sha224_" + hash.SumSHA224(v)
}

func getK8sNamespacedName(objs ...any) any {
	if len(objs) == 0 {
		return nil
	}

	ns := make([]string, 0, len(objs))

	for _, obj := range objs {
		n, _ := cache.MetaNamespaceKeyFunc(obj)
		ns = append(ns, n)
	}

	if len(ns) == 1 {
		return ns[0]
	}

	return ns
}
