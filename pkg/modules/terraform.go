package modules

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"

	"github.com/seal-io/seal/pkg/bus/module"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/gopool"
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

	defaultModuleVersion = "0.0.0"

	fileNameMainTf = "main.tf"
	fileNameReadme = "README.md"
)

// GetTerraformModuleFiles parse a full tf configuration to a map using filename as the key.
func GetTerraformModuleFiles(name, content string) (map[string]string, error) {
	modDir := files.TempDir("seal-module*")
	if err := os.WriteFile(filepath.Join(modDir, fileNameMainTf), []byte(content), 0600); err != nil {
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

	var files = make(map[string]string)

	// TODO split the configuration to variables.tf, outputs.tf, etc.
	files[fileNameMainTf] = content
	files[fileNameReadme] = buf.String()

	return files, nil
}

func GetModuleNameByPath(path string) string {
	lastPart := filepath.Base(path)
	if _, err := version.NewVersion(lastPart); err != nil {
		return lastPart
	}

	// lastPart is a version, get the
	nonVersionPath := strings.TrimSuffix(path, lastPart)
	return filepath.Base(nonVersionPath)
}

func SchemaSync(mc model.ClientSet) schemaSyncer {
	return schemaSyncer{mc: mc}
}

type schemaSyncer struct {
	mc model.ClientSet
}

// Do fetches and updates the schema of the given module,
// within 5 mins in the background.
func (s schemaSyncer) Do(_ context.Context, message module.BusMessage) error {
	var logger = log.WithName("module")

	gopool.Go(func() {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		var m = message.Refer

		logger.Debugf("syncing schema for module %s", m.ID)
		var err = syncSchema(ctx, s.mc, m)
		if err == nil {
			return
		}

		logger.Warnf("recording syncing module %s schema failed: %v", m.ID, err)
		m.Status = status.ModuleStatusError
		m.StatusMessage = fmt.Sprintf("sync schema failed: %v", err)
		update, err := dao.ModuleUpdate(s.mc, m)
		if err != nil {
			logger.Errorf("failed to prepare module %s update: %v", m.ID, err)
			return
		}
		err = update.Exec(ctx)
		if err != nil {
			logger.Errorf("failed to update module %s: %v", m.ID, err)
		}
	})

	return nil
}

func syncSchema(ctx context.Context, mc model.ClientSet, m *model.Module) error {
	versions, err := loadTerraformModuleVersions(m)
	if err != nil {
		return err
	}

	return mc.WithTx(ctx, func(tx *model.Tx) error {
		// clean up previous module versions if there's any.
		var _, err = tx.ModuleVersions().Delete().
			Where(moduleversion.ModuleID(m.ID)).
			Exec(ctx)
		if err != nil {
			return err
		}

		// create new module versions.
		creates, err := dao.ModuleVersionCreates(tx, versions...)
		if err != nil {
			return err
		}
		for _, c := range creates {
			_, err = c.Save(ctx)
			if err != nil {
				return err
			}
		}

		// state module.
		m.Status = status.ModuleStatusReady
		update, err := dao.ModuleUpdate(tx, m)
		if err != nil {
			return err
		}

		return update.Exec(ctx)
	})
}

func loadTerraformModuleVersions(m *model.Module) ([]*model.ModuleVersion, error) {
	mRoot := filepath.Join(os.TempDir(), "seal-module-"+strs.String(10))
	defer os.RemoveAll(mRoot)

	if err := getter.Get(mRoot, m.Source); err != nil {
		return nil, err
	}

	versions, err := getVersionsFromRoot(mRoot)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		// try to load a module version from the root
		// This keeps compatible with terraform git source
		schema, err := loadTerraformModuleSchema(mRoot)
		if err != nil {
			return nil, err
		}

		u, err := url.Parse(m.Source)
		if err != nil {
			return nil, err
		}
		v := u.Query().Get("ref")
		if v == "" {
			v = defaultModuleVersion
		}

		return []*model.ModuleVersion{
			{
				ModuleID: m.ID,
				Source:   m.Source,
				Version:  v,
				Schema:   schema,
			},
		}, nil
	}

	// Support reading different module versions in subdirectory
	var mvs []*model.ModuleVersion
	for _, v := range versions {
		schema, err := loadTerraformModuleSchema(filepath.Join(mRoot, v))
		if err != nil {
			return nil, err
		}

		// ModuleVersion.Source is the concatenation of Module.Source and Version
		versionedSource := getVersionedSource(m.Source, v)
		mvs = append(mvs, &model.ModuleVersion{
			ModuleID: m.ID,
			Source:   versionedSource,
			Version:  v,
			Schema:   schema,
		})
	}

	return mvs, nil
}

func getVersionsFromRoot(root string) ([]string, error) {
	var versions []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if _, err = version.NewVersion(entry.Name()); err != nil {
			continue
		}

		versions = append(versions, entry.Name())
	}
	return versions, nil
}

func getVersionedSource(source, version string) string {
	var protocol, base, subdir, query = "", source, "", ""

	if strings.Contains(base, "?") {
		base, query, _ = strings.Cut(base, "?")
	}
	if strings.Contains(base, "://") {
		protocol, base, _ = strings.Cut(base, "://")
	}
	if strings.Contains(base, "//") {
		base, subdir, _ = strings.Cut(base, "//")
	}
	subdir = filepath.Join(subdir, version)

	var result = fmt.Sprintf("%s//%s", base, subdir)
	if protocol != "" {
		result = fmt.Sprintf("%s://%s", protocol, result)
	}
	if query != "" {
		result = fmt.Sprintf("%s?%s", result, query)
	}

	return result
}

func loadTerraformModuleSchema(path string) (*types.ModuleSchema, error) {
	mod, diags := tfconfig.LoadModule(path)
	if diags.HasErrors() {
		return nil, diags.Err()
	}

	readme, err := getReadme(path)
	if err != nil {
		return nil, err
	}

	moduleSchema := &types.ModuleSchema{
		Readme: readme,
	}

	for _, v := range sortVariables(mod.Variables) {
		moduleSchema.Variables = append(moduleSchema.Variables, terraformVariableToModuleVariable(v))
	}

	for _, v := range sortOutput(mod.Outputs) {
		moduleSchema.Outputs = append(moduleSchema.Outputs, types.ModuleOutput{
			Name:        v.Name,
			Description: v.Description,
			Sensitive:   v.Sensitive,
		})
	}

	for name := range mod.RequiredProviders {
		moduleSchema.RequiredConnectorTypes = append(moduleSchema.RequiredConnectorTypes, name)
	}
	sort.Strings(moduleSchema.RequiredConnectorTypes)

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

func sortVariables(m map[string]*tfconfig.Variable) (s []*tfconfig.Variable) {
	s = make([]*tfconfig.Variable, 0, len(m))
	for k := range m {
		s = append(s, m[k])
	}
	sort.SliceStable(s, func(i, j int) bool {
		return judgeSourcePos(&s[i].Pos, &s[j].Pos)
	})
	return
}

func sortOutput(m map[string]*tfconfig.Output) (s []*tfconfig.Output) {
	s = make([]*tfconfig.Output, 0, len(m))
	for k := range m {
		s = append(s, m[k])
	}
	sort.SliceStable(s, func(i, j int) bool {
		return judgeSourcePos(&s[i].Pos, &s[j].Pos)
	})
	return
}

func judgeSourcePos(i, j *tfconfig.SourcePos) bool {
	switch {
	case i.Filename < j.Filename:
		return false
	case i.Filename > j.Filename:
		return true
	}
	return i.Line < j.Line
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
