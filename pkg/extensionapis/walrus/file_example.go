package walrus

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/seal-io/utils/stringx"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/registry/rest"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"

	walrus "github.com/seal-io/walrus/pkg/apis/walrus/v1"
	"github.com/seal-io/walrus/pkg/extensionapi"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/systemsetting"
	"github.com/seal-io/walrus/pkg/vcs"
)

// FileExampleHandler handles v1.FileExample objects.
//
// FileExampleHandler maintains an in-memory v1.FileExampleList object to serve the get/list operation.
//
// During setup, if the serve-walrus-files-url v1.Setting is a remote URL,
// FileExampleHandler will clone the repository to the local file system at first.
// Then, it will load the file examples from the local file system.
//
// The icon of each v1.FileExample object has been encoded as Base64.
type FileExampleHandler struct {
	extensionapi.ObjectInfo
	extensionapi.GetOperation
	extensionapi.ListOperation

	InMemoryListObject walrus.FileExampleList
}

func (h *FileExampleHandler) SetupHandler(
	ctx context.Context,
	opts extensionapi.SetupOptions,
) (gvr schema.GroupVersionResource, srs map[string]rest.Storage, err error) {
	// Declare GVR.
	gvr = walrus.SchemeGroupVersionResource("fileexamples")

	// As storage.
	h.ObjectInfo = &walrus.FileExample{}
	h.GetOperation = extensionapi.WithGet(h)
	h.ListOperation = extensionapi.WithList(nil, h)

	// Get file examples local filesystem.
	var wfFs fs.FS
	{
		var wfUri *url.URL
		wfUri, err = systemsetting.ServeWalrusFilesUrl.ValueURL(ctx)
		if err != nil {
			return
		}
		if wfUri.Scheme != "file" {
			// Git clone from remote.
			// TODO(thxCode): support other remote source?
			dir := system.SubLibDir("files")
			wfFs, err = vcs.GitClone(ctx, dir, vcs.GitCloneOptions{URL: wfUri.String()})
			if err != nil {
				return
			}
		} else {
			// Load from local.
			wfFs = os.DirFS(wfUri.Path)
		}
	}

	// Load file examples in memory.
	h.InMemoryListObject.Items, err = loadFileExamples(wfFs, systemkuberes.SystemNamespaceName)
	if err != nil {
		return
	}

	return
}

var (
	_ rest.Storage = (*FileExampleHandler)(nil)
	_ rest.Lister  = (*FileExampleHandler)(nil)
	_ rest.Getter  = (*FileExampleHandler)(nil)
)

func (h *FileExampleHandler) New() runtime.Object {
	return &walrus.FileExample{}
}

func (h *FileExampleHandler) Destroy() {}

func (h *FileExampleHandler) NewList() runtime.Object {
	return &walrus.FileExampleList{}
}

func (h *FileExampleHandler) OnList(ctx context.Context, opts ctrlcli.ListOptions) (runtime.Object, error) {
	// Support watch with `kubectl get -A`.
	if opts.Namespace == "" {
		opts.Namespace = systemkuberes.SystemNamespaceName
	}

	// Only support list in system namespace.
	if opts.Namespace == systemkuberes.SystemNamespaceName {
		return h.InMemoryListObject.DeepCopy(), nil
	}

	return &walrus.FileExampleList{}, nil
}

func (h *FileExampleHandler) OnGet(ctx context.Context, key types.NamespacedName, opts ctrlcli.GetOptions) (runtime.Object, error) {
	// Only support get in system namespace.
	if key.Namespace == systemkuberes.SystemNamespaceName {
		for _, item := range h.InMemoryListObject.Items {
			if item.Name == key.Name {
				return item.DeepCopy(), nil
			}
		}
	}

	return nil, kerrors.NewNotFound(walrus.SchemeResource("fileexamples"), key.Name)
}

func loadFileExamples(dir fs.FS, namespace string) ([]walrus.FileExample, error) {
	var (
		knownIconFiles = sets.New(
			"icon.png",
			"icon.jpg",
			"icon.jpeg",
			"icon.svg")
		knownReadmeFiles = sets.New(
			"README.md",
			"readme.md")
		knownYamlSuffixes = sets.New(
			".yaml",
			".yml")
	)

	entries, err := fs.ReadDir(dir, ".")
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	// Sort by name.
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	exps := make([]walrus.FileExample, 0, len(entries))
	for i, en := range entries {
		if !en.IsDir() || strings.HasPrefix(en.Name(), ".") {
			continue
		}

		sdir, err := fs.Sub(dir, en.Name())
		if err != nil {
			return nil, fmt.Errorf("sub dir: %w", err)
		}
		sentries, err := fs.ReadDir(sdir, ".")
		if err != nil {
			return nil, fmt.Errorf("read sub dir: %w", err)
		}

		exp := walrus.FileExample{
			ObjectMeta: meta.ObjectMeta{
				Namespace:         namespace,
				Name:              en.Name(),
				UID:               types.UID(uuid.NewMD5(uuid.Nil, []byte(en.Name())).String()), // Create a deterministic UID.
				CreationTimestamp: meta.Now(),
				ResourceVersion:   stringx.FromInt(i + 1),
			},
		}

		for _, sen := range sentries {
			fn := sen.Name()

			switch {
			case knownIconFiles.Has(fn):
				data, err := fs.ReadFile(sdir, fn)
				if err != nil {
					return nil, fmt.Errorf("read icon: %w", err)
				}
				exp.Status.Icon = base64.StdEncoding.EncodeToString(data)
			case knownReadmeFiles.Has(fn):
				data, err := fs.ReadFile(sdir, fn)
				if err != nil {
					return nil, fmt.Errorf("read README: %w", err)
				}
				exp.Status.Readme = string(data)
			case knownYamlSuffixes.Has(filepath.Ext(fn)):
				data, err := fs.ReadFile(sdir, fn)
				if err != nil {
					return nil, fmt.Errorf("read YAML: %w", err)
				}
				exp.Status.Content = string(data)
			}
		}

		exps = append(exps, exp)
	}

	return exps, nil
}
