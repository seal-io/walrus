package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/walrus/utils/files"
)

const (
	fileNameMainTf = "main.tf"
	fileNameReadme = "README.md"
)

// GetTerraformTemplateFiles parse a full tf configuration to a map using filename as the key.
func GetTerraformTemplateFiles(name, content string) (map[string]string, error) {
	modDir := files.TempDir("seal-template*")
	if err := os.WriteFile(filepath.Join(modDir, fileNameMainTf), []byte(content), 0o600); err != nil {
		return nil, err
	}

	defer os.RemoveAll(modDir)

	mod, diags := tfconfig.LoadModule(modDir)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse given terraform module: %w", diags.Err())
	}

	// Replace the tmp path by name. It will be rendered in the readme title.
	mod.Path = name

	buf := &bytes.Buffer{}
	if err := tfconfig.RenderMarkdown(buf, mod); err != nil {
		return nil, fmt.Errorf("failed to render module readme: %w", err)
	}

	files := make(map[string]string)

	// TODO split the configuration to variables.tf, outputs.tf, etc.
	files[fileNameMainTf] = content
	files[fileNameReadme] = buf.String()

	return files, nil
}

func GetTemplateNameByPath(path string) string {
	lastPart := filepath.Base(path)
	if _, err := version.NewVersion(lastPart); err != nil {
		return lastPart
	}

	// LastPart is a version, get the.
	nonVersionPath := strings.TrimSuffix(path, lastPart)

	return filepath.Base(nonVersionPath)
}
