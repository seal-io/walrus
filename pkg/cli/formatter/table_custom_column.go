package formatter

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func getCustomColumnRender(group, operation string) []fieldRender {
	key := fmt.Sprintf("%s/%s", group, operation)
	return customColumnRender[key]
}

var customColumnRender = map[string][]fieldRender{
	"Resources/get": {
		resourceFieldType,
	},
	"Resources/list": {
		resourceFieldType,
	},
	"Resources/get-outputs": {
		resourceOutputFieldType,
	},
	"ResourceDefinitions/get-resources": {
		resourceFieldType,
	},
	"Variables/list": {
		variableFieldValue,
	},
}

var (
	resourceFieldType = fieldRender{
		name: "type",
		renderFunc: func(_ string, r gjson.Result) string {
			rtyp := r.Get("type")
			if rtyp.Exists() {
				return rtyp.String()
			}

			rt := r.Get("template")
			if rt.Exists() {
				rn := rt.Get("name")
				rv := rt.Get("version")
				proj := rt.Get("project")

				if proj.Exists() {
					return fmt.Sprintf("%s@%s", rn.String(), rv.String())
				} else {
					return fmt.Sprintf("%s@%s(Global)", rn.String(), rv.String())
				}
			}

			return ""
		},
	}

	resourceOutputFieldType = fieldRender{
		name: "type",
		renderFunc: func(_ string, r gjson.Result) string {
			rtyp := r.Get("type")

			if rtyp.IsArray() {
				if r := rtyp.Array(); len(r) >= 1 {
					return r[0].String()
				}
			}

			return rtyp.String()
		},
	}

	variableFieldValue = fieldRender{
		name: "value",
		renderFunc: func(_ string, r gjson.Result) string {
			rs := r.Get("sensitive")
			if rs.IsBool() && rs.Bool() {
				return "<sensitive>"
			}

			return r.Get("value").String()
		},
	}
)
