package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/seal-io/walrus/pkg/dao/entx"
	"github.com/seal-io/walrus/utils/log"
)

func main() {
	err := generate()
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}
}

// generate produces DAO APIs in a safer and cleaner way.
func generate() error {
	// Prepare.
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting working directory: %w", err)
	}

	header, err := os.ReadFile(filepath.Join(pwd, "/hack/boilerplate/go.txt"))
	if err != nil {
		return err
	}

	// Generate.
	cfg := entx.Config{
		ProjectDir: pwd,
		Project:    "github.com/seal-io/walrus",
		Package:    "github.com/seal-io/walrus/pkg/dao",
		Header:     string(header),
	}

	err = entx.Generate(cfg)
	if err != nil {
		log.Fatalf("error generating: %v", err)
	}

	return nil
}
