package migration

import (
	"bytes"
	"fmt"
	"text/template"
)

func executeTemplate(tpl *template.Template, data any) string {
	var buff bytes.Buffer

	if err := tpl.Execute(&buff, data); err != nil {
		panic(fmt.Errorf("error execute template: %w", err))
	}

	return buff.String()
}
