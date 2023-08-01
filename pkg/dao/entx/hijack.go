package entx

import (
	"entgo.io/ent/entc/gen"

	"github.com/seal-io/seal/utils/strs"
)

func fixDefaultTemplates() {
	// Allow loading update additional template.
	for i := range gen.Templates {
		if gen.Templates[i].Name != "update" {
			continue
		}

		gen.Templates[i].ExtendPatterns = append(
			gen.Templates[i].ExtendPatterns,
			"dialect/sql/update/additional/*",
			"dialect/sql/update/fields/additional/*")

		break
	}
}

func fixDefaultTemplateFuncs() {
	// Overwrite.
	gen.Funcs["camel"] = strs.CamelizeDownFirst
	gen.Funcs["snake"] = strs.Underscore
	gen.Funcs["pascal"] = strs.Camelize
	gen.Funcs["singular"] = strs.Singularize
	gen.Funcs["plural"] = strs.Pluralize
}

func fixDefaultTemplateRulesetAcronyms() {
	for _, acronym := range strs.Acronyms() {
		gen.AddAcronym(acronym)
	}
}
