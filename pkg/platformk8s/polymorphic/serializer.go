package polymorphic

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructuredscheme"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializerjson "k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/seal-io/seal/utils/json"
)

type metaFactory struct{}

func (metaFactory) Interpret(data []byte) (*schema.GroupVersionKind, error) {
	var gvk runtime.TypeMeta
	if err := json.Unmarshal(data, &gvk); err != nil {
		return nil, fmt.Errorf("error unmarsalling runtime type meta: %w", err)
	}
	gv, err := schema.ParseGroupVersion(gvk.APIVersion)
	if err != nil {
		return nil, fmt.Errorf("error pasring runtime group version: %w", err)
	}
	return &schema.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: gvk.Kind}, nil
}

var (
	mf      = metaFactory{}
	creator = unstructuredscheme.NewUnstructuredCreator()
	typer   = unstructuredscheme.NewUnstructuredObjectTyper()
)

func JsonSerializer() runtime.Serializer {
	opts := serializerjson.SerializerOptions{Yaml: false, Pretty: false, Strict: false}
	return serializerjson.NewSerializerWithOptions(mf, creator, typer, opts)
}

func YamlSerializer() runtime.Serializer {
	opts := serializerjson.SerializerOptions{Yaml: true, Pretty: false, Strict: false}
	return serializerjson.NewSerializerWithOptions(mf, creator, typer, opts)
}
