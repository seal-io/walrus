package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/templates/openapi"
	"github.com/seal-io/walrus/pkg/templates/translator"
	"github.com/seal-io/walrus/utils/files"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

// TerraformLoader for load terraform template schema and data.
type TerraformLoader struct {
	translator translator.Translator
}

// NewTerraformLoader creates a new terraform loader.
func NewTerraformLoader() SchemaLoader {
	return &TerraformLoader{
		translator: translator.NewTerraformTranslator(),
	}
}

// Load loads the internal template version schema and data.
func (l *TerraformLoader) Load(
	rootDir, templateName string, mode Mode,
) (*types.TemplateVersionSchema, error) {
	if !l.match(rootDir) {
		return nil, nil
	}

	mod, err := l.loadMod(rootDir)
	if err != nil {
		return nil, err
	}

	s, err := l.loadSchema(rootDir, mod, templateName, mode)
	if err != nil {
		return nil, err
	}

	d, err := l.loadData(rootDir, mod)
	if err != nil {
		return nil, err
	}

	return &types.TemplateVersionSchema{
		Schema: types.Schema{
			OpenAPISchema: s,
		},
		TemplateVersionSchemaData: d,
	}, nil
}

// Match checks if the template is a terraform template.
func (l *TerraformLoader) match(rootDir string) bool {
	_, err := os.Stat(filepath.Join(rootDir, "main.tf"))
	return err == nil
}

// loadMod load the terraform module.
func (l *TerraformLoader) loadMod(rootDir string) (*tfconfig.Module, error) {
	mod, diags := tfconfig.LoadModule(rootDir)
	if diags.HasErrors() {
		return nil, diags.Err()
	}

	return mod, nil
}

// loadSchema loads the internal template version schema.
func (l *TerraformLoader) loadSchema(
	rootDir string,
	mod *tfconfig.Module,
	template string,
	mode Mode,
) (*openapi3.T, error) {
	// Variables.
	vs, err := l.getVariableSchema(rootDir, mod, mode)
	if err != nil {
		return nil, err
	}

	// Outputs.
	ots, err := l.getOutputSchemaFromTerraform(mod)
	if err != nil {
		return nil, err
	}

	// Empty schema.
	if vs == nil && ots == nil {
		return nil, nil
	}

	// OpenAPI OpenAPISchema.
	t := &openapi3.T{
		OpenAPI: openapi.OpenAPIVersion,
		Info: &openapi3.Info{
			Title: fmt.Sprintf("OpenAPI schema for template %s", template),
		},
		Components: &openapi3.Components{
			Schemas: openapi3.Schemas{},
		},
	}

	if vs != nil {
		t.Components.Schemas["variables"] = vs.NewRef()
	}

	if ots != nil {
		t.Components.Schemas["outputs"] = ots.NewRef()
	}

	return t, nil
}

func (l *TerraformLoader) getVariableSchema(
	rootDir string,
	mod *tfconfig.Module,
	mode Mode,
) (*openapi3.Schema, error) {
	fromOriginal, err := l.getVariableSchemaFromTerraform(mod)
	if err != nil {
		return nil, err
	}

	switch mode {
	case ModeOriginal:
		return fromOriginal, nil
	case ModeSchemaFile:
		fromFile, err := l.getVariableSchemaFromFile(rootDir)
		if err != nil {
			log.Warnf("error loading schema from file: %v", err)
		}

		return l.injectVariableSchemaExt(fromOriginal, fromFile), nil
	default:
		// Merge, apply customized schema to generated schema.
		fromFile, err := l.getVariableSchemaFromFile(rootDir)
		if err != nil {
			log.Warnf("error loading schema from file: %v", err)
		}

		if fromFile == nil {
			return fromOriginal, nil
		}

		// Generate merged variables in sequence.
		merged, err := openapi.UnionSchema(fromOriginal, fromFile)
		if err != nil {
			return nil, err
		}

		return merged, nil
	}
}

// injectVariableSchemaExt add extension to the customized ui schema in schema.yaml.
func (l *TerraformLoader) injectVariableSchemaExt(generated, customized *openapi3.Schema) *openapi3.Schema {
	if customized == nil {
		return generated
	}

	s := *customized
	if len(s.Extensions) == 0 && len(generated.Extensions) != 0 {
		s.Extensions = generated.Extensions
	}

	for n, v := range s.Properties {
		in := generated.Properties[n]
		if in == nil || in.Value == nil {
			continue
		}

		// Title.
		if v.Value.Title == "" {
			s.Properties[n].Value.Title = generated.Properties[n].Value.Title
		}

		// Extensions.
		if len(in.Value.Extensions) != 0 {
			if len(v.Value.Extensions) == 0 {
				s.Properties[n].Value.Extensions = make(map[string]any)
			}
			s.Properties[n].Value.Extensions[openapi.ExtOriginal] = in.Value.Extensions[openapi.ExtOriginal]
		}
	}

	return &s
}

// getVariableSchemaFromTerraform generate variable schemas from terraform files.
func (l *TerraformLoader) getVariableSchemaFromTerraform(mod *tfconfig.Module) (*openapi3.Schema, error) {
	if len(mod.Variables) == 0 {
		return nil, nil
	}

	var (
		varSchemas = openapi3.NewObjectSchema()
		required   []string
		keys       = make([]string, len(mod.Variables))
	)

	// Variables.
	for i, v := range sortVariables(mod.Variables) {
		// Parse tf expression from type.
		var (
			tfType = cty.DynamicPseudoType
			def    = v.Default
		)

		if v.Type != "" {
			// Type exists.
			expr, diags := hclsyntax.ParseExpression(strs.ToBytes(&v.Type), "", hcl.Pos{Line: 1, Column: 1})
			if diags.HasErrors() {
				return nil, fmt.Errorf("error parsing expression: %w", diags)
			}

			var conDef *typeexpr.Defaults

			tfType, conDef, diags = typeexpr.TypeConstraintWithDefaults(expr)
			if diags.HasErrors() {
				return nil, fmt.Errorf("error getting type: %w", diags)
			}

			if conDef != nil && conDef.DefaultValues != nil && conDef.Children != nil {
				def = conDef
			}
		} else if v.Default != nil {
			// Empty type, use default value to get type.
			b, err := json.Marshal(v.Default)
			if err != nil {
				return nil, fmt.Errorf("error while marshal terraform variable default value: %w", err)
			}

			var sjv ctyjson.SimpleJSONValue

			err = sjv.UnmarshalJSON(b)
			if err != nil {
				return nil, fmt.Errorf("error while unmarshal terraform variable default value: %w", err)
			}
			tfType = sjv.Type()
		}

		// Generate json schema from tf type.
		varSchemas.WithProperty(
			v.Name,
			l.translator.SchemaOfOriginalType(
				tfType,
				v.Name,
				def,
				v.Description,
				v.Sensitive))

		if v.Required {
			required = append(required, v.Name)
		}
		keys[i] = v.Name
	}

	sort.Strings(required)
	varSchemas.Required = required
	varSchemas.Extensions = openapi.NewExt(varSchemas.Extensions).
		SetOriginalVariablesSequence(keys).
		Export()

	return varSchemas, nil
}

// getVariableSchemaFromFile generate variable schemas from schema.yaml.
func (l *TerraformLoader) getVariableSchemaFromFile(rootDir string) (*openapi3.Schema, error) {
	schemaFile := filepath.Join(rootDir, "schema.yaml")
	if !files.Exists(schemaFile) {
		if schemaFile = filepath.Join(rootDir, "schema.yml"); !files.Exists(schemaFile) {
			return nil, nil
		}
	}

	it, err := openapi3.NewLoader().LoadFromFile(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("error loading schema file %s: %w", schemaFile, err)
	}

	if it.Components == nil ||
		it.Components.Schemas == nil ||
		it.Components.Schemas["variables"] == nil ||
		it.Components.Schemas["variables"].Value == nil {
		return nil, nil
	}

	return it.Components.Schemas["variables"].Value, nil
}

// getOutputSchemaFromTerraform generate output schemas from terraform files.
func (l *TerraformLoader) getOutputSchemaFromTerraform(mod *tfconfig.Module) (*openapi3.Schema, error) {
	if len(mod.Outputs) == 0 {
		return nil, nil
	}

	var (
		filenames     = sets.Set[string]{}
		outputSchemas = openapi3.NewObjectSchema()
	)

	for _, v := range sortOutput(mod.Outputs) {
		// Use dynamic type for output.
		outputSchemas.WithProperty(
			v.Name,
			l.translator.SchemaOfOriginalType(
				cty.DynamicPseudoType,
				v.Name,
				nil,
				v.Description,
				v.Sensitive))

		filenames.Insert(v.Pos.Filename)
	}

	values, err := getOutputValues(filenames)
	if err != nil {
		return nil, err
	}

	for n, v := range values {
		outputSchemas.Properties[n].Value.Extensions = openapi.NewExt(outputSchemas.Properties[n].Value.Extensions).
			SetOriginalValueExpression(v).
			SetOriginalType(cty.DynamicPseudoType).
			Export()
	}

	outputSchemas.Extensions = openapi.NewExt(outputSchemas.Extensions).
		Export()

	return outputSchemas, nil
}

// getOutputValues gets the output values from output configuration files.
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

// getOutputValueFromFile gets the output value from output configuration file.
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

// loadData loads the internal template version data.
func (l *TerraformLoader) loadData(rootDir string, mod *tfconfig.Module) (
	data types.TemplateVersionSchemaData, err error,
) {
	// Readme.
	data.Readme, err = l.getReadme(rootDir)
	if err != nil {
		return
	}

	// Providers.
	requiredProviders := l.getRequiredProviders(mod.RequiredProviders)
	data.RequiredProviders = requiredProviders

	return
}

// getReadme gets the readme content.
func (l *TerraformLoader) getReadme(rootDir string) (string, error) {
	path := filepath.Join(rootDir, "README.md")
	if files.Exists(path) {
		content, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}

		return string(content), nil
	}

	return "", nil
}

// getRequiredProviders gets the required providers.
func (l *TerraformLoader) getRequiredProviders(
	m map[string]*tfconfig.ProviderRequirement,
) (s []types.ProviderRequirement) {
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
