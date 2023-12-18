package resource

import (
	"fmt"
	"testing"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/property"
)

func TestTopologicalSortResources(t *testing.T) {
	cases := []struct {
		Name     string
		Services model.Resources
		Expected []string
		Error    error
	}{
		{
			Name: "cycle",
			Services: model.Resources{
				&model.Resource{
					Name: "1",
					Attributes: property.Values{
						"attr": []byte("${res.3.attr}"),
					},
				},
				&model.Resource{
					Name: "2",
					Attributes: property.Values{
						"attr": []byte("${res.1.attr}"),
					},
				},
				&model.Resource{
					Name: "3",
					Attributes: property.Values{
						"attr": []byte("${res.2.attr}"),
					},
				},
			},
			Error: fmt.Errorf("cycle detected: 3 -> 1 -> 2 -> 3"),
		},
		{
			Name: "simple dependence",
			Services: model.Resources{
				&model.Resource{
					Name: "6",
					Attributes: property.Values{
						"attr": []byte("${res.5.attr}"),
					},
				},
				&model.Resource{
					Name: "5",
					Attributes: property.Values{
						"attr": []byte("${res.4.attr}"),
					},
				},
				&model.Resource{
					Name: "4",
				},
			},
			Expected: []string{"4", "5", "6"},
		},
		{
			Name: "Complex dependence",
			Services: model.Resources{
				&model.Resource{
					Name: "3",
					Attributes: property.Values{
						"attr":  []byte("${res.2.attr}"),
						"attr2": []byte("${res.4.attr}"),
					},
				},
				&model.Resource{
					Name: "2",
					Attributes: property.Values{
						"attr": []byte("${res.1.attr}"),
					},
				},
				&model.Resource{
					Name: "1",
				},
				&model.Resource{
					Name: "4",
					Attributes: property.Values{
						"attr": []byte("${res.2.attr}"),
					},
				},
				&model.Resource{
					Name: "5",
					Attributes: property.Values{
						"attr": []byte("${res.4.attr}"),
					},
				},
				&model.Resource{
					Name: "6",
				},
				&model.Resource{
					Name: "7",
					Attributes: property.Values{
						"attr": []byte("${res.6.attr}"),
					},
				},
				&model.Resource{
					Name: "8",
					Attributes: property.Values{
						"attr":  []byte("${res.6.attr}"),
						"attr2": []byte("${res.3.attr}"),
					},
				},
				&model.Resource{
					Name: "9",
					Attributes: property.Values{
						"attr":  []byte("${res.7.attr}"),
						"attr2": []byte("${res.8.attr}"),
					},
				},
				&model.Resource{
					Name: "10",
					Attributes: property.Values{
						"attr":  []byte("${res.7.attr}"),
						"attr2": []byte("${res.8.attr}"),
						"attr3": []byte("${res.9.attr}"),
					},
				},
			},
			Expected: []string{
				"1", "6", "2", "7", "4", "3", "5", "8", "9", "10",
			},
		},
	}

	for _, c := range cases {
		sorted, err := TopologicalSortResources(c.Services)
		if err != nil {
			if c.Error != nil {
				continue
			} else {
				t.Fatalf("expected error %v, got %v", c.Error, err)
			}
		}

		if !expectedTopologicalSortResources(sorted, c.Expected) {
			t.Fatalf("case %s failed. Expected %v, got %s", c.Name, c.Expected, getNames(sorted))
		}
	}
}

func expectedTopologicalSortResources(resources model.Resources, expected []string) bool {
	sortNames := make([]string, len(resources))
	for i, resource := range resources {
		sortNames[i] = resource.Name
	}

	return strSliceEqual(sortNames, expected)
}

func strSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if b[i] != v {
			return false
		}
	}

	return true
}

func getNames(s model.Resources) []string {
	sort := []string{}
	for i := 0; i < len(s); i++ {
		sort = append(sort, s[i].Name)
	}

	return sort
}
