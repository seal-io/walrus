package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/parser"
	"sigs.k8s.io/yaml"
)

func Test_reflectPackage(t *testing.T) {
	var (
		dir = filepath.Join("testdata", "reflect_package")
		pkg = "github.com/seal-io/code-generator/cmd/apireg-gen/generators/testdata/reflect_package"
	)

	b := parser.New()
	if err := b.AddDir(pkg); err != nil {
		t.Fatalf("failed to parse package %q: %v", pkg, err)
	}

	u, err := b.FindTypes()
	if err != nil {
		t.Fatalf("failed to find types: %v", err)
	}

	td := reflectPackage(u.Package(pkg))

	actualBytes, err := yaml.Marshal(td)
	if err != nil {
		t.Fatalf("failed to marshal type definition: %v", err)
	}

	expectedBytes, err := os.ReadFile(filepath.Join(dir, "expected.yaml"))
	if err != nil {
		t.Fatalf("failed to read expected file: %v", err)
	}

	assert.Equal(t, string(expectedBytes), string(actualBytes))
}
