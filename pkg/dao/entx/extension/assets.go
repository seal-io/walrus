package extension

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"
)

// Assets writes the given content to corresponding paths,
// borrows from entgo.io/ent/entc/gen/graph.go.
type Assets struct {
	dirs  map[string]struct{}
	files map[string][]byte
}

// Add adds content which should write to the given path later.
func (a *Assets) Add(path string, content []byte) {
	if a.dirs == nil {
		a.dirs = make(map[string]struct{})
	}

	if a.files == nil {
		a.files = make(map[string][]byte)
	}
	a.dirs[filepath.Dir(path)] = struct{}{}
	a.files[path] = content
}

// Write outputs added content to corresponding paths.
func (a Assets) Write() error {
	for dir := range a.dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("create dir %q: %w", dir, err)
		}
	}

	for path, content := range a.files {
		if err := os.WriteFile(path, content, 0o600); err != nil {
			return fmt.Errorf("write file %q: %w", path, err)
		}
	}

	return nil
}

// Format runs "goimports" on all outputs.
func (a Assets) Format() error {
	for path, content := range a.files {
		src, err := imports.Process(path, content, nil)
		if err != nil {
			return fmt.Errorf("format file %s: %w", path, err)
		}

		if err = os.WriteFile(path, src, 0o600); err != nil {
			return fmt.Errorf("write file %s: %w", path, err)
		}
	}

	return nil
}
