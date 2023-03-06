package modules

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/reader"
	"github.com/seal-io/seal/utils/strs"
)

const (
	attributeLabel   = "label"
	attributeGroup   = "group"
	attributeOptions = "options"
	attributeShowIf  = "show_if"
)

func loadTerraformModuleSchema(source string) (*types.ModuleSchema, error) {
	tmpDir := filepath.Join(os.TempDir(), "seal-module-"+strs.String(10))
	defer os.RemoveAll(tmpDir)

	if err := getter.Get(tmpDir, source); err != nil {
		return nil, err
	}
	mod, diags := tfconfig.LoadModule(tmpDir)
	if diags.HasErrors() {
		return nil, diags.Err()
	}

	readme, err := getReadme(tmpDir)
	if err != nil {
		return nil, err
	}

	moduleSchema := &types.ModuleSchema{
		Readme: readme,
	}

	for _, v := range mod.Variables {
		moduleSchema.Variables = append(moduleSchema.Variables, terraformVariableToModuleVariable(v))
	}

	for _, v := range mod.Outputs {
		moduleSchema.Outputs = append(moduleSchema.Outputs, types.ModuleOutput{
			Name:        v.Name,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		})
	}

	for name := range mod.RequiredProviders {
		moduleSchema.RequiredConnectorTypes = append(moduleSchema.RequiredConnectorTypes, name)
	}

	return moduleSchema, nil
}

func getReadme(dir string) (string, error) {
	path := filepath.Join(dir, "README.md")
	if files.Exists(path) {
		content, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(content), nil
	}

	return "", nil
}

func terraformVariableToModuleVariable(v *tfconfig.Variable) types.ModuleVariable {
	variable := types.ModuleVariable{
		Name:        v.Name,
		Type:        v.Type,
		Description: v.Description,
		Default:     v.Default,
		Required:    v.Required,
		Sensitive:   v.Sensitive,
	}

	comments, err := loadComments(v.Pos.Filename, v.Pos.Line)
	if err != nil {
		log.Warnf("failed to load terraform comments for var %s, error: %v", v.Name, err)
		return variable
	}
	setTerraformVariableExtensions(&variable, comments)
	return variable
}

func loadComments(filename string, lineNum int) ([]string, error) {
	lines := reader.Lines{
		FileName: filename,
		LineNum:  lineNum,
		Condition: func(line string) bool {
			return strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//")
		},
		Parser: func(line string) (string, bool) {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimPrefix(line, "//")
			line = strings.TrimSpace(line)
			return line, true
		},
	}
	return lines.Read()
}

func setTerraformVariableExtensions(variable *types.ModuleVariable, comments []string) {
	const atSign = "@"
	for _, comment := range comments {
		if strings.HasPrefix(comment, atSign) {
			splits := regexp.MustCompile(`\s+`).Split(comment[1:], 2)
			if len(splits) < 2 {
				continue
			}
			key, value := splits[0], strings.TrimSpace(splits[1])
			var attr any
			switch key {
			case attributeLabel:
				attr = &variable.Label
			case attributeGroup:
				attr = &variable.Group
			case attributeShowIf:
				attr = &variable.ShowIf
			case attributeOptions:
				attr = &variable.Options
			default:
				log.Warnf("unrecognized variable attribute %q in comment: %s", key, comment)
				continue
			}
			if err := json.Unmarshal([]byte(value), attr); err != nil {
				log.Warnf("failed to unmarshal schema in hcl comment: %s, %v", comment, err)
				continue
			}
		}
	}
}
