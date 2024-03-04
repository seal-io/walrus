package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/parser"
	"sigs.k8s.io/yaml"
)

func Test_reflectType(t *testing.T) {
	var (
		dir = filepath.Join("testdata", "reflect_type")
		pkg = "github.com/seal-io/code-generator/cmd/webhook-gen/generators/testdata/reflect_type"
		typ = "DummyReconciler"
	)

	b := parser.New()
	if err := b.AddDir(pkg); err != nil {
		t.Fatalf("failed to parse package %q: %v", pkg, err)
	}

	u, err := b.FindTypes()
	if err != nil {
		t.Fatalf("failed to find types: %v", err)
	}

	ty := u.Package(pkg).Type(typ)
	if ty == nil {
		t.Fatalf("failed to find type %q in package %q", typ, pkg)
	}

	td := reflectType(ty)

	actualBytes, err := yaml.Marshal(td)
	if err != nil {
		t.Fatalf("failed to marshal typed definition: %v", err)
	}

	expectedBytes, err := os.ReadFile(filepath.Join(dir, "expected.yaml"))
	if err != nil {
		t.Fatalf("failed to read expected file: %v", err)
	}

	assert.Equal(t, string(expectedBytes), string(actualBytes))
}
