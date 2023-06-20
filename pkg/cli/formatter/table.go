package formatter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"

	"github.com/alexeyco/simpletable"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/utils/json"
)

// watchFields represent the common headers for list response.
var watchFields = []string{"id", "name", "createTime"}

// TableFormatter use to convert response to table format.
type TableFormatter struct{}

func (f *TableFormatter) Format(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	// Response status is not 200.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if len(body) == 0 {
			return nil, fmt.Errorf("unexpected status code %d", resp.StatusCode)
		}

		data := errorResponse{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code %d: %w", resp.StatusCode, err)
		}

		return []byte(fmt.Sprintf("unexpected status code %d: %s", resp.StatusCode, data.Message)), nil
	}

	if len(body) == 0 {
		return nil, nil
	}

	var data interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		m, ok := data.(map[string]interface{})
		if !ok {
			return nil, errors.New("can't decode response in table, use json or yaml format instead")
		}

		if len(m) == 0 {
			return []byte{}, nil
		}

		its, ok := m["items"]
		if ok {
			items, ok := its.([]interface{})
			if ok {
				formatted := f.resourceItems(items)
				return []byte(formatted), nil
			}
		}

		return []byte(f.resourceItem(m)), nil
	case reflect.Slice, reflect.Array:
		m, ok := data.([]interface{})
		if !ok {
			return nil, errors.New("can't decode response in table, use json or yaml format instead")
		}

		if len(m) == 0 {
			return []byte{}, nil
		}

		return []byte(f.generalItems(m)), nil
	}

	return nil, nil
}

func (f *TableFormatter) generalItems(data []interface{}) string {
	if len(data) == 0 {
		return ""
	}

	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	var (
		headers []string
		item    = data[0]
	)

	switch reflect.TypeOf(item).Kind() {
	case reflect.Map:
		it, ok := item.(map[string]interface{})
		if !ok {
			return ""
		}

		for h := range it {
			headers = append(headers, h)
		}

		sort.SliceStable(headers, func(i, j int) bool {
			if headers[i] == "id" {
				return true
			}

			if headers[j] == "id" {
				return false
			}

			if headers[i] == "name" {
				return true
			}

			if headers[j] == "name" {
				return false
			}

			return headers[i] < headers[j]
		})

		for _, v := range headers {
			table.Header.Cells = append(table.Header.Cells, &simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  v,
			})
		}

	case reflect.Array:
		table.Header.Cells = []*simpletable.Cell{
			{
				Align: simpletable.AlignRight,
				Text:  "value",
			},
		}
	}

	for _, it := range data {
		switch reflect.TypeOf(it).Kind() {
		case reflect.Map:
			itm, ok := it.(map[string]interface{})
			if !ok {
				continue
			}

			var r []*simpletable.Cell

			for _, v := range headers {
				d, ok := itm[v]
				if !ok {
					continue
				}

				r = append(r, &simpletable.Cell{
					Align: simpletable.AlignRight,
					Text:  fmt.Sprintf("%v", d),
				})
			}

			table.Body.Cells = append(table.Body.Cells, r)
		case reflect.Array:
			items, ok := it.([]interface{})
			if !ok {
				continue
			}

			for _, v := range items {
				table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
					{
						Align: simpletable.AlignRight,
						Text:  fmt.Sprintf("%v", v),
					},
				})
			}
		}
	}

	return table.String()
}

func (f *TableFormatter) resourceItems(data []interface{}) string {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	exitedFields := sets.Set[string]{}

	for _, it := range data {
		item, ok := it.(map[string]interface{})
		if !ok {
			continue
		}

		var r []*simpletable.Cell

		for _, v := range watchFields {
			d, ok := item[v]
			if !ok {
				continue
			}

			r = append(r, &simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  fmt.Sprintf("%v", d),
			})

			exitedFields.Insert(v)
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	for _, v := range watchFields {
		if exitedFields.Has(v) {
			table.Header.Cells = append(
				table.Header.Cells,
				&simpletable.Cell{
					Align: simpletable.AlignRight,
					Text:  v,
				})
		}
	}

	return table.String()
}

func (f *TableFormatter) resourceItem(data map[string]interface{}) string {
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)

	var (
		bodyCells = make([]*simpletable.Cell, 0)
		headCells = make([]*simpletable.Cell, 0)
	)

	for _, v := range watchFields {
		d, ok := data[v]
		if !ok {
			continue
		}

		headCells = append(headCells,
			&simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  v,
			})

		bodyCells = append(bodyCells,
			&simpletable.Cell{
				Align: simpletable.AlignRight,
				Text:  fmt.Sprintf("%v", d),
			})
	}

	table.Header = &simpletable.Header{
		Cells: headCells,
	}
	table.Body = &simpletable.Body{
		Cells: [][]*simpletable.Cell{
			bodyCells,
		},
	}

	return table.String()
}

type errorResponse struct {
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
}
