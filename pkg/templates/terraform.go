package templates

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
	"sync"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/bus/template"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/property"
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
	attributeHidden  = "hidden"

	defaultTemplateVersion = "0.0.0"

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

func SchemaSync(mc model.ClientSet) schemaSyncer {
	return schemaSyncer{mc: mc}
}

type schemaSyncer struct {
	mc model.ClientSet
}

// Do fetches and updates the schema of the given template,
// within 5 mins in the background.
func (s schemaSyncer) Do(_ context.Context, message template.BusMessage) error {
	logger := log.WithName("template")

	gopool.Go(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		m := message.Refer

		logger.Debugf("syncing schema for template %s", m.ID)

		err := syncSchema(ctx, s.mc, m)
		if err == nil {
			return
		}

		logger.Warnf("recording syncing template %s schema failed: %v", m.ID, err)
		m.Status = status.TemplateStatusError
		m.StatusMessage = fmt.Sprintf("sync schema failed: %v", err)

		update, err := dao.TemplateUpdate(s.mc, m)
		if err != nil {
			logger.Errorf("failed to prepare template %s update: %v", m.ID, err)
			return
		}

		err = update.Exec(ctx)
		if err != nil {
			logger.Errorf("failed to update template %s: %v", m.ID, err)
		}
	})

	return nil
}

func syncSchema(ctx context.Context, mc model.ClientSet, t *model.Template) error {
	versions, err := loadTerraformTemplateVersions(t)
	if err != nil {
		return err
	}

	return mc.WithTx(ctx, func(tx *model.Tx) error {
		// Clean up previous template versions if there's any.
		_, err := tx.TemplateVersions().Delete().
			Where(templateversion.TemplateID(t.ID)).
			Exec(ctx)
		if err != nil {
			return err
		}

		// Create new template versions.
		creates, err := dao.TemplateVersionCreates(tx, versions...)
		if err != nil {
			return err
		}

		for _, c := range creates {
			_, err = c.Save(ctx)
			if err != nil {
				return err
			}
		}

		// State template.
		t.Status = status.TemplateStatusReady

		update, err := dao.TemplateUpdate(tx, t)
		if err != nil {
			return err
		}

		return update.Exec(ctx)
	})
}

func loadTerraformTemplateVersions(t *model.Template) ([]*model.TemplateVersion, error) {
	mRoot := filepath.Join(os.TempDir(), "seal-template-"+strs.String(10))
	defer os.RemoveAll(mRoot)

	if err := getter.Get(mRoot, t.Source); err != nil {
		return nil, err
	}

	versions, err := getVersionsFromRoot(mRoot)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		// Try to load a template version from the root
		// This keeps compatible with terraform git source.
		schema, err := loadTerraformTemplateSchema(mRoot)
		if err != nil {
			return nil, err
		}

		u, err := url.Parse(t.Source)
		if err != nil {
			return nil, err
		}

		v := u.Query().Get("ref")
		if v == "" {
			v = defaultTemplateVersion
		}

		return []*model.TemplateVersion{
			{
				TemplateID: t.ID,
				Source:     t.Source,
				Version:    v,
				Schema:     schema,
			},
		}, nil
	}

	// Support reading different template versions in subdirectory.
	mvs := make([]*model.TemplateVersion, 0, len(versions))

	for _, v := range versions {
		schema, err := loadTerraformTemplateSchema(filepath.Join(mRoot, v))
		if err != nil {
			return nil, err
		}

		// TemplateVersion.Source is the concatenation of Template.Source and Version.
		versionedSource := getVersionedSource(t.Source, v)

		mvs = append(mvs, &model.TemplateVersion{
			TemplateID: t.ID,
			Source:     versionedSource,
			Version:    v,
			Schema:     schema,
		})
	}

	return mvs, nil
}

func getVersionsFromRoot(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0, len(entries))

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
	protocol, base, subdir, query := "", source, "", ""

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

	result := fmt.Sprintf("%s//%s", base, subdir)
	if protocol != "" {
		result = fmt.Sprintf("%s://%s", protocol, result)
	}

	if query != "" {
		result = fmt.Sprintf("%s?%s", result, query)
	}

	return result
}

func loadTerraformTemplateSchema(path string) (*types.TemplateSchema, error) {
	mod, diags := tfconfig.LoadModule(path)
	if diags.HasErrors() {
		return nil, diags.Err()
	}

	readme, err := getReadme(path)
	if err != nil {
		return nil, err
	}

	templateSchema := &types.TemplateSchema{
		Readme: readme,
	}

	for _, v := range sortVariables(mod.Variables) {
		s, err := getVariableSchema(v)
		if err != nil {
			return nil, err
		}

		templateSchema.Variables = append(templateSchema.Variables, s)
	}

	templateSchema.Outputs, err = getOutputsSchema(mod.Outputs)
	if err != nil {
		return nil, err
	}
	templateSchema.RequiredProviders = getRequiredProviders(mod.RequiredProviders)

	return templateSchema, nil
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

func getRequiredProviders(m map[string]*tfconfig.ProviderRequirement) (s []types.ProviderRequirement) {
	if len(m) == 0 {
		return
	}

	for k, v := range m {
		s = append(s, types.ProviderRequirement{
			Name:                k,
			ProviderRequirement: v,
		})
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Name < s[j].Name
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

func getVariableSchema(v *tfconfig.Variable) (property.Schema, error) {
	variable, err := property.GuessSchema(v.Name, v.Type, v.Default)
	if err != nil {
		return property.Schema{}, fmt.Errorf("unresolved variable %s schema: %w", v.Name, err)
	}

	variable.WithDescription(v.Description)

	if v.Required {
		variable = variable.WithRequired()
	}

	if v.Sensitive {
		variable = variable.WithSensitive()
	}

	comments, err := loadComments(v.Pos.Filename, v.Pos.Line)
	if err != nil {
		log.Warnf("failed to load terraform comments for var %s, error: %v", v.Name, err)
		return variable, nil
	}

	extendVariableSchema(&variable, comments)

	return variable, nil
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

func extendVariableSchema(variable *property.Schema, comments []string) {
	const atSign = "@"
	for _, comment := range comments {
		if strings.HasPrefix(comment, atSign) {
			splits := regexp.MustCompile(`\s+`).Split(comment[1:], 2)
			if len(splits) == 1 && splits[0] == attributeHidden {
				variable.Hidden = true
				continue
			}

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

func getOutputsSchema(outputs map[string]*tfconfig.Output) (property.Schemas, error) {
	var (
		filenames     = sets.Set[string]{}
		outputsSchema property.Schemas
	)

	for _, v := range outputs {
		filenames.Insert(v.Pos.Filename)
	}

	values, err := getOutputValues(filenames)
	if err != nil {
		return nil, err
	}

	for _, v := range sortOutput(outputs) {
		s := property.AnySchema(v.Name, nil).
			WithDescription(v.Description)
		if v.Sensitive {
			s = s.WithSensitive()
		}

		if ov, ok := values[v.Name]; ok {
			s = s.WithValue(ov)
		}

		outputsSchema = append(outputsSchema, s)
	}

	return outputsSchema, nil
}

func getOutputValues(filenames sets.Set[string]) (map[string][]byte, error) {
	var (
		mu      sync.Mutex
		logger  = log.WithName("template")
		wg      = gopool.Group()
		outputs = make(map[string][]byte)
	)

	for _, filename := range filenames.UnsortedList() {
		wg.Go(func() error {
			b, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("error read output configuration file %s: %w", filename, err)
			}

			var (
				file   *hcl.File
				diag   hcl.Diagnostics
				parser = hclparse.NewParser()
			)

			if strings.HasSuffix(filename, ".json") {
				file, diag = parser.ParseJSON(b, filename)
			} else {
				file, diag = parser.ParseHCL(b, filename)
			}

			if diag.HasErrors() {
				logger.Warnf("error parse output configuration file %s: %s", filename, diag.Error())
				return nil
			}

			if file == nil {
				return nil
			}

			o := getOutputValueFromFile(file)

			mu.Lock()
			for on, oe := range o {
				outputs[on] = oe
			}
			mu.Unlock()

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	return outputs, nil
}

func getOutputValueFromFile(file *hcl.File) map[string][]byte {
	var (
		rootSchema = &hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type:       "output",
					LabelNames: []string{"name"},
				},
			},
		}
		outputSchema = &hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{
					Name: "value",
				},
			},
		}
	)

	var (
		outputs       = make(map[string][]byte)
		content, _, _ = file.Body.PartialContent(rootSchema)
	)

	for _, block := range content.Blocks {
		if block.Type == "output" {
			ct, _, _ := block.Body.PartialContent(outputSchema)
			name := block.Labels[0]

			if attr, defined := ct.Attributes["value"]; defined {
				outputs[name] = attr.Expr.Range().SliceBytes(file.Bytes)
			}
		}
	}

	return outputs
}
