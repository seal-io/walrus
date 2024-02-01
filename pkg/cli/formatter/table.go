package formatter

import (
	"io"
	"net/http"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/pkg/cli/common"
)

const (
	fieldID          = "id"
	fieldName        = "name"
	fieldProject     = "project.name"
	fieldEnvironment = "environment.name"
	fieldStatus      = "status.summaryStatus"
	fieldCreateTime  = "createTime"
)

// builtinDisplayFields represent the common headers for list response.
var builtinDisplayFields = []string{
	fieldID,
	fieldName,
	fieldProject,
	fieldEnvironment,
	fieldStatus,
	fieldCreateTime,
}

var fieldAlias = map[string]string{
	fieldStatus:      "STATUS",
	fieldCreateTime:  "CREATED",
	fieldProject:     "PROJECT",
	fieldEnvironment: "ENVIRONMENT",
}

// TableFormatter use to convert response to table format.
type TableFormatter struct {
	Columns   []string
	Group     string
	Operation string
}

func (f *TableFormatter) Format(resp *http.Response) ([]byte, error) {
	err := common.CheckResponseStatus(resp)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	var (
		r         = gjson.ParseBytes(body)
		formatted string
	)

	switch {
	case r.IsObject():
		result := gjson.GetBytes(body, "items")
		if result.Exists() {
			formatted = f.resourceItems([]byte(result.Raw))
		} else {
			formatted = f.resourceItem(body)
		}

	case r.IsArray():
		formatted = f.resourceItems(body)
	}

	return []byte(formatted), nil
}

func (f *TableFormatter) resourceItems(body []byte) string {
	data := gjson.ParseBytes(body).Array()
	if len(data) == 0 {
		return ""
	}

	var (
		totalColumns = append(builtinDisplayFields, f.Columns...)
		columnSet    = sets.NewString(totalColumns...)
		existColumns = sets.NewString()
	)

	for _, arr := range data {
		for _, k := range columnSet.UnsortedList() {
			if arr.Get(k).Exists() {
				existColumns.Insert(k)
			}
		}

		if existColumns.HasAll(columnSet.UnsortedList()...) {
			break
		}
	}

	columns := make([]fieldRender, 0)

	for _, v := range totalColumns {
		if existColumns.Has(v) {
			columns = append(columns, fieldRender{
				name:       v,
				renderFunc: defaultRenderFunc,
			})
		}
	}

	columns = append(columns, getCustomColumnRender(f.Group, f.Operation)...)

	return f.renderColumns(columns, data)
}

func (f *TableFormatter) resourceItem(body []byte) string {
	data := gjson.ParseBytes(body)
	if !data.Exists() {
		return ""
	}

	columns := make([]fieldRender, 0)

	for _, v := range builtinDisplayFields {
		if data.Get(v).Exists() {
			columns = append(columns, fieldRender{
				name:       v,
				renderFunc: defaultRenderFunc,
			})
		}
	}

	for _, v := range f.Columns {
		if data.Get(v).Exists() {
			columns = append(columns, fieldRender{
				name:       v,
				renderFunc: defaultRenderFunc,
			})
		}
	}

	columns = append(columns, getCustomColumnRender(f.Group, f.Operation)...)

	return f.renderColumns(columns, []gjson.Result{data})
}

func (f *TableFormatter) renderColumns(columns []fieldRender, data []gjson.Result) string {
	var (
		header        = make([]*simpletable.Cell, 0, len(columns))
		rows          = make([][]*simpletable.Cell, 0, len(data))
		sortedColumns = customSort(columns)
	)

	for _, v := range sortedColumns {
		header = append(header,
			&simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  columnDisplayName(v.name),
			})
	}

	for _, arr := range data {
		var row []*simpletable.Cell

		for _, v := range sortedColumns {
			row = append(row, &simpletable.Cell{
				Align: simpletable.AlignLeft,
				Text:  v.renderFunc(v.name, arr),
			})
		}

		rows = append(rows, row)
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)
	table.Header.Cells = header
	table.Body.Cells = rows

	return table.String()
}

type fieldRender struct {
	name       string
	renderFunc func(name string, r gjson.Result) string
}

var defaultRenderFunc = func(name string, r gjson.Result) string {
	return r.Get(name).String()
}

func columnDisplayName(n string) string {
	an := fieldAlias[n]
	if an != "" {
		return an
	}

	return strings.ToUpper(n)
}

func customSort(fields []fieldRender) []fieldRender {
	frontFields := []string{fieldID, fieldName, fieldProject, fieldEnvironment}
	endFields := []string{fieldStatus, fieldCreateTime}

	var (
		existFrontFields = make([]fieldRender, 0)
		existEndFields   = make([]fieldRender, 0)
		unknownFields    = make([]fieldRender, 0)
	)

	for _, field := range fields {
		switch {
		case slices.Contains(frontFields, field.name):
			existFrontFields = append(existFrontFields, field)
		case slices.Contains(endFields, field.name):
			existEndFields = append(existEndFields, field)
		default:
			unknownFields = append(unknownFields, field)
		}
	}

	sortedFields := make([]fieldRender, 0, len(existFrontFields)+len(existEndFields)+len(unknownFields))
	sortedFields = append(sortedFields, existFrontFields...)
	sortedFields = append(sortedFields, unknownFields...)
	sortedFields = append(sortedFields, existEndFields...)

	return sortedFields
}
