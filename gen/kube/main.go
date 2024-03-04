package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/seal-io/utils/stringx"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"

	"github.com/seal-io/walrus/gen/kube/builder"
)

func main() {
	err := generate()
	if err != nil {
		klog.Fatalf("error generating: %v", err)
	}
}

func generate() error {
	// Prepare.
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	header, err := os.ReadFile(filepath.Join(pwd, "/hack/boilerplate/go.txt"))
	if err != nil {
		return err
	}

	// Generate.
	cfg := builder.Config{
		ProjectDir: pwd,
		Project:    "github.com/seal-io/walrus",
		Header:     stringx.FromBytes(ptr.To(bytes.TrimSpace(header))),
		APIs: []string{
			"github.com/seal-io/walrus/pkg/apis/walruscore/v1",
		},
		ExtensionAPIs: []string{
			"github.com/seal-io/walrus/pkg/apis/walrus/v1",
		},
		Webhooks: []string{
			"github.com/seal-io/walrus/pkg/webhooks/walruscore",
		},
		ExternalAPIs: []string{
			"k8s.io/api/admission/v1",
			"k8s.io/api/admissionregistration/v1",
			"k8s.io/api/apps/v1",
			"k8s.io/api/authentication/v1",
			"k8s.io/api/authorization/v1",
			"k8s.io/api/autoscaling/v1",
			"k8s.io/api/autoscaling/v2",
			"k8s.io/api/batch/v1",
			"k8s.io/api/certificates/v1",
			"k8s.io/api/coordination/v1",
			"k8s.io/api/core/v1",
			"k8s.io/api/discovery/v1",
			"k8s.io/api/events/v1",
			"k8s.io/api/rbac/v1",
			"k8s.io/api/scheduling/v1",
			"k8s.io/api/storage/v1",
			"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1",
			"k8s.io/kube-aggregator/pkg/apis/apiregistration/v1",
		},
		PluralExceptions: map[string]string{
			"Endpoints": "Endpoints",
		},
		ProtoImports: []string{
			// NB(thxCode): the go-to-protobuf under code-generator relies on a deprecated project,
			// https://github.com/gogo/protobuf.
			// The upstream already filed an issue about this,
			// https://github.com/kubernetes/kubernetes/issues/96564.
			// In order to support generating protobuf code for extension APIs,
			// we need to tell protoc where to find the gogo/protobuf.
			filepath.Join(pwd, "staging"),
		},
	}

	return builder.Generate(cfg)
}
