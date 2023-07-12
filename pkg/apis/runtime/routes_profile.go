package runtime

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/render"
	"k8s.io/apimachinery/pkg/util/sets"
)

// ProfileType defines the type of the profile.
type ProfileType = string

const (
	ProfileTypeBasic  ProfileType = "basic"
	ProfileTypeArray  ProfileType = "array"
	ProfileTypeObject ProfileType = typeDescriptorObject
)

// ProfileCategory defines the category of the type profile.
type ProfileCategory = string

const (
	ProfileCategoryHeader ProfileCategory = "header"
	ProfileCategoryUri    ProfileCategory = "uri"
	ProfileCategoryForm   ProfileCategory = "form"
	ProfileCategoryQuery  ProfileCategory = "query"
	ProfileCategoryJson   ProfileCategory = "json"

	typeDescriptorInt     = "int"
	typeDescriptorFloat32 = "float32"
	typeDescriptorFloat64 = "float64"
	typeDescriptorString  = "string"
	typeDescriptorRender  = "render.Render"
	typeDescriptorObject  = "object"
)

var categories = []ProfileCategory{
	ProfileCategoryHeader,
	ProfileCategoryUri,
	ProfileCategoryForm,
	ProfileCategoryQuery,
	ProfileCategoryJson,
}

// InputProfile defines the type profile of the input argument.
type InputProfile struct {
	ProfileProperty
}

// OutputProfile defines the type profile of the output result.
type OutputProfile struct {
	ProfileProperty

	// Page specifies whether to envelop the output result as a page.
	Page bool
}

// ProfileRouter defines the route definition that extract from the input argument's `route` tag,
// e.g. type struct { _ struct{} `route:GET=/test` },
// the `GET` should be placed into Method and the `/test` should be placed into SubPath.
type ProfileRouter struct {
	Method  string
	SubPath string
}

// ProfileProperty defines the property of a type/field profile.
type ProfileProperty struct {
	// Required specifies whether the type/field is required.
	Required bool
	// Default specifies the default value of the field.
	Default string
	// Category specifies ProfileCategory of the type/field.
	Category ProfileCategory
	// Name specifies the name of the field.
	Name string
	// Type specifies the ProfileType of the type/field.
	Type ProfileType
	// TypeDescriptor specifies the descriptor of the type/field.
	TypeDescriptor string
	// TypeArrayLength specifies the length of array if the kind of type/field is `reflect.Array`.
	TypeArrayLength int
	// TypeRefer specifies whether the field is referred from other ProfileProperty.
	TypeRefer bool
	// Properties stores the properties of the type.
	Properties []ProfileProperty
	// Extension stores extra data, like cli config.
	Extension map[string]any
}

// State returns the state to describe what categories the ProfileProperty have.
func (p ProfileProperty) State() map[ProfileCategory]bool {
	l := make(map[string]bool)

	for i := 0; i < len(p.Properties); i++ {
		if p.Properties[i].Category == "" {
			continue
		}
		l[p.Properties[i].Category] = true

		for j := 0; j < len(p.Properties[i].Properties); j++ {
			r := p.Properties[i].Properties[j].State()
			for k := range r {
				l[k] = r[k]
			}
		}

		if len(l) == len(categories) {
			break
		}
	}

	return l
}

// Flat flattens the properties of the ProfileProperty with the given categories,
// all leaf nodes have arranged in the result.
func (p ProfileProperty) Flat(categories ...ProfileCategory) []ProfileProperty {
	if len(categories) == 0 {
		return nil
	}

	return p.flat(sets.New(categories...))
}

func (p ProfileProperty) flat(categories sets.Set[string]) []ProfileProperty {
	var l []ProfileProperty

	for i := 0; i < len(p.Properties); i++ {
		if !categories.Has(p.Properties[i].Category) {
			continue
		}

		if len(p.Properties[i].Properties) == 0 {
			if p.Properties[i].Name != "" {
				l = append(l, p.Properties[i])
			}

			continue
		}

		for j := 0; j < len(p.Properties[i].Properties); j++ {
			r := p.Properties[i].Properties[j].flat(categories)
			l = append(l, r...)
		}
	}

	return l
}

// Filter filters the properties of the ProfileProperty with the given category,
// all nodes only keep the same category children in the result.
func (p ProfileProperty) Filter(category ProfileCategory) []ProfileProperty {
	var root []ProfileProperty

	for i := 0; i < len(p.Properties); i++ {
		if p.Properties[i].Category != category {
			continue
		}

		root = append(root, p.Properties[i])

		var props []ProfileProperty

		for j := 0; j < len(p.Properties[i].Properties); j++ {
			if p.Properties[i].Properties[j].Category != category {
				continue
			}

			props = append(props, p.Properties[i].Properties[j])
			props[len(props)-1].Properties = p.Properties[i].Properties[j].Filter(category)
		}
		root[len(root)-1].Properties = props
	}

	return root
}

// GetInputProfile parses the given reflect.Type as an InputProfile.
func GetInputProfile(t reflect.Type) *InputProfile {
	if t == nil {
		return nil
	}

	var p InputProfile

	if t.Kind() == reflect.Func {
		t = t.In(1)
	}
	t = decodeTypePointer(t)
	p.ProfileProperty = getProfileProperty(sets.New[string](), "", "", nil, t)

	for _, category := range categories {
		vs := sets.New[string]()
		p.Properties = append(p.Properties, getProfileProperties(vs, category, t)...)
	}

	return &p
}

// GetOutputProfile parses the given reflect.Type as an OutputProfile.
func GetOutputProfile(t reflect.Type) *OutputProfile {
	if t == nil {
		return nil
	}

	var p OutputProfile

	if t.Kind() == reflect.Func {
		p.Page = t.NumOut() > 2
	}

	if t.Kind() == reflect.Func {
		t = t.Out(0)
	}
	t = decodeTypePointer(t)
	p.ProfileProperty = getProfileProperty(sets.New[string](), "", "", nil, t)

	for _, category := range categories {
		vs := sets.New[string]()
		p.Properties = append(p.Properties, getProfileProperties(vs, category, t)...)
	}

	return &p
}

// getProfileRouter parses the given reflect.Type as an ProfileRouter.
func getProfileRouter(t reflect.Type) *ProfileRouter {
	if t == nil {
		return nil
	}

	t = decodeTypePointer(t)
	if t.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		v := f.Tag.Get("route")
		if !isTagBlank(v) {
			m, sp := getTagAttribute(v)
			if m == "" || sp == "" {
				continue
			}

			m = strings.ToUpper(m)
			switch m {
			default:
				continue
			case http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodGet:
			}

			return &ProfileRouter{
				Method:  m,
				SubPath: path.Join("/", sp),
			}
		}
	}

	return nil
}

func getProfileProperties(vs sets.Set[string], category string, t reflect.Type) []ProfileProperty {
	t = decodeTypePointer(t)

	switch t.Kind() {
	default:
		return []ProfileProperty{getProfileProperty(vs, category, "", nil, t)}
	case reflect.Struct:
	}

	var ps []ProfileProperty

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}

		v := f.Tag.Get(category)
		if isTagBlank(v) {
			continue
		}

		name, attrs := getTagNameAndAttributes(v)
		if isTagInline(attrs) {
			ps = append(ps, getProfileProperties(vs, category, f.Type)...)
			continue
		}

		ps = append(ps, getProfileProperty(vs, category, name, attrs, f.Type))
	}

	return ps
}

func getProfileProperty(
	vs sets.Set[string],
	category string,
	name string,
	attrs []string,
	t reflect.Type,
) ProfileProperty {
	p := ProfileProperty{
		Required:       isTagRequired(attrs),
		Default:        getTagAttributeValue(attrs, "default"),
		Category:       category,
		Name:           name,
		Type:           ProfileTypeBasic,
		TypeDescriptor: t.String(),
		Extension:      getProfilePropertyExtension(attrs),
	}

	// Well known basic type.
	if _, exist := basicSchemas[p.TypeDescriptor]; exist {
		return p
	}

	switch t.Kind() {
	case reflect.Bool:
		p.TypeDescriptor = "bool"
	case reflect.Int:
		p.TypeDescriptor = typeDescriptorInt
	case reflect.Int8:
		p.TypeDescriptor = "int8"
	case reflect.Int16:
		p.TypeDescriptor = "int16"
	case reflect.Int32:
		p.TypeDescriptor = "int32"
	case reflect.Int64:
		p.TypeDescriptor = "int64"
	case reflect.Uint:
		p.TypeDescriptor = "uint"
	case reflect.Uint8:
		p.TypeDescriptor = "uint8"
	case reflect.Uint16:
		p.TypeDescriptor = "uint16"
	case reflect.Uint32:
		p.TypeDescriptor = "uint32"
	case reflect.Uint64:
		p.TypeDescriptor = "uint64"
	case reflect.Float32:
		p.TypeDescriptor = typeDescriptorFloat32
	case reflect.Float64:
		p.TypeDescriptor = typeDescriptorFloat64
	case reflect.Uintptr, reflect.UnsafePointer:
		p.TypeDescriptor = typeDescriptorInt
	case reflect.Complex64:
		p.TypeDescriptor = typeDescriptorFloat32
	case reflect.Complex128:
		p.TypeDescriptor = typeDescriptorFloat64
	case reflect.Chan, reflect.Func:
		p.TypeDescriptor = typeDescriptorString
	case reflect.Map:
		p.Type = ProfileTypeObject

		switch t.Elem().Kind() {
		case reflect.String:
			p.TypeDescriptor = "map[string]string"
		case reflect.Int, reflect.Uint:
			p.TypeDescriptor = "map[string]int"
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
			p.TypeDescriptor = "map[string]int32"
		case reflect.Int64, reflect.Uint64:
			p.TypeDescriptor = "map[string]int64"
		default:
			p.TypeDescriptor = typeDescriptorObject
		}
	case reflect.Interface:
		p.Type = ProfileTypeObject
		p.TypeDescriptor = typeDescriptorObject
	case reflect.Pointer:
		t = decodeTypePointer(t.Elem())
		p = getProfileProperty(vs, category, name, attrs, t)
	case reflect.Struct:
		expected := reflect.TypeOf((*render.Render)(nil)).Elem()
		if t.ConvertibleTo(expected) {
			p.Type = ProfileTypeBasic
			p.TypeDescriptor = typeDescriptorRender

			return p
		}
		p.Type = ProfileTypeObject
		p.TypeRefer = vs.Has(p.TypeDescriptor + "." + p.Category)

		if !p.TypeRefer {
			vs.Insert(p.TypeDescriptor + "." + p.Category)
			p.Properties = getProfileProperties(vs, p.Category, t)
		}
	case reflect.String:
		p.TypeDescriptor = typeDescriptorString
	case reflect.Array, reflect.Slice:
		p.Type = ProfileTypeArray
		if t.Kind() == reflect.Array {
			p.TypeArrayLength = t.Len()
		}
		t = decodeTypePointer(t.Elem())

		p.TypeDescriptor = t.String()
		if p.TypeDescriptor == "byte" {
			p.Type = ProfileTypeBasic
			p.TypeDescriptor = "[]byte"

			return p
		}

		switch t.Kind() {
		case reflect.Bool:
			p.TypeDescriptor = "bool"
		case reflect.Int:
			p.TypeDescriptor = typeDescriptorInt
		case reflect.Int8:
			p.TypeDescriptor = "int8"
		case reflect.Int16:
			p.TypeDescriptor = "int16"
		case reflect.Int32:
			p.TypeDescriptor = "int32"
		case reflect.Int64:
			p.TypeDescriptor = "int64"
		case reflect.Uint:
			p.TypeDescriptor = "uint"
		case reflect.Uint8:
			p.TypeDescriptor = "uint8"
		case reflect.Uint16:
			p.TypeDescriptor = "uint16"
		case reflect.Uint32:
			p.TypeDescriptor = "uint32"
		case reflect.Uint64:
			p.TypeDescriptor = "uint64"
		case reflect.Float32:
			p.TypeDescriptor = typeDescriptorFloat32
		case reflect.Float64:
			p.TypeDescriptor = typeDescriptorFloat64
		case reflect.Uintptr, reflect.UnsafePointer:
			p.TypeDescriptor = typeDescriptorInt
		case reflect.Complex64:
			p.TypeDescriptor = typeDescriptorFloat32
		case reflect.Complex128:
			p.TypeDescriptor = typeDescriptorFloat64
		case reflect.Chan, reflect.Func:
			p.TypeDescriptor = typeDescriptorString
		case reflect.Interface, reflect.Map:
			p.TypeDescriptor = typeDescriptorObject
		case reflect.Struct:
			p.TypeRefer = vs.Has(p.TypeDescriptor + "." + p.Category)
			if !p.TypeRefer {
				vs.Insert(p.TypeDescriptor + "." + p.Category)
				p.Properties = getProfileProperties(vs, p.Category, t)
			}
		case reflect.Array, reflect.Slice:
			p.TypeDescriptor = "array"
			p.Properties = getProfileProperties(vs, p.Category, t)
		case reflect.String:
			p.TypeDescriptor = typeDescriptorString
		}
	}

	return p
}

func isTagBlank(tag string) bool {
	return tag == "" || tag == "-"
}

func getTagNameAndAttributes(tag string) (name string, attrs []string) {
	ss := strings.SplitN(tag, ",", 2)
	name = strings.TrimSpace(ss[0])

	if len(ss) == 2 {
		ss[1] = strings.TrimSpace(ss[1])
		if ss[1] != "" {
			attrs = strings.Split(ss[1], ",")
		}
	}

	return
}

func isTagInline(attrs []string) bool {
	for i := range attrs {
		if attrs[i] == "inline" {
			return true
		}
	}

	return false
}

func isTagRequired(attrs []string) bool {
	for i := range attrs {
		if attrs[i] == "omitempty" || strings.HasPrefix(attrs[i], "default=") {
			return false
		}
	}

	return true
}

func getTagAttributeValue(attrs []string, key string) string {
	for i := range attrs {
		k, v := getTagAttribute(attrs[i])
		if k == key {
			return v
		}
	}

	return ""
}

func getTagAttribute(attr string) (key, value string) {
	ss := strings.SplitN(attr, "=", 2)
	if len(ss) == 1 {
		return strings.TrimSpace(ss[0]), ""
	}

	return strings.TrimSpace(ss[0]), strings.TrimSpace(ss[1])
}

func decodeTypePointer(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	return t
}

func getProfilePropertyExtension(attrs []string) map[string]any {
	ext := make(map[string]any)

	for _, v := range attrs {
		if strings.HasPrefix(v, "cli") {
			key := fmt.Sprintf("x-%s", v)
			ext[key] = true
		}
	}

	return ext
}
