package service

import (
	"fmt"
	"testing"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

func TestTopologicalSortServices(t *testing.T) {
	cases := []struct {
		Name     string
		Services model.Services
		Expected []string
		Error    error
	}{
		{
			Name: "cycle",
			Services: model.Services{
				&model.Service{
					Name: "1",
					Attributes: property.Values{
						"attr": []byte("${service.3.attr}"),
					},
				},
				&model.Service{
					Name: "2",
					Attributes: property.Values{
						"attr": []byte("${service.1.attr}"),
					},
				},
				&model.Service{
					Name: "3",
					Attributes: property.Values{
						"attr": []byte("${service.2.attr}"),
					},
				},
			},
			Error: fmt.Errorf("cycle detected: 3 -> 1 -> 2 -> 3"),
		},
		{
			Name: "simple dependence",
			Services: model.Services{
				&model.Service{
					Name: "6",
					Attributes: property.Values{
						"attr": []byte("${service.5.attr}"),
					},
				},
				&model.Service{
					Name: "5",
					Attributes: property.Values{
						"attr": []byte("${service.4.attr}"),
					},
				},
				&model.Service{
					Name: "4",
				},
			},
			Expected: []string{"4", "5", "6"},
		},
		{
			Name: "Complex dependence",
			Services: model.Services{
				&model.Service{
					Name: "3",
					Attributes: property.Values{
						"attr":  []byte("${service.2.attr}"),
						"attr2": []byte("${service.4.attr}"),
					},
				},
				&model.Service{
					Name: "2",
					Attributes: property.Values{
						"attr": []byte("${service.1.attr}"),
					},
				},
				&model.Service{
					Name: "1",
				},
				&model.Service{
					Name: "4",
					Attributes: property.Values{
						"attr": []byte("${service.2.attr}"),
					},
				},
				&model.Service{
					Name: "5",
					Attributes: property.Values{
						"attr": []byte("${service.4.attr}"),
					},
				},
				&model.Service{
					Name: "6",
				},
				&model.Service{
					Name: "7",
					Attributes: property.Values{
						"attr": []byte("${service.6.attr}"),
					},
				},
				&model.Service{
					Name: "8",
					Attributes: property.Values{
						"attr":  []byte("${service.6.attr}"),
						"attr2": []byte("${service.3.attr}"),
					},
				},
				&model.Service{
					Name: "9",
					Attributes: property.Values{
						"attr":  []byte("${service.7.attr}"),
						"attr2": []byte("${service.8.attr}"),
					},
				},
				&model.Service{
					Name: "10",
					Attributes: property.Values{
						"attr":  []byte("${service.7.attr}"),
						"attr2": []byte("${service.8.attr}"),
						"attr3": []byte("${service.9.attr}"),
					},
				},
			},
			Expected: []string{
				"1", "6", "2", "7", "4", "3", "5", "8", "9", "10",
			},
		},
	}

	for _, c := range cases {
		sorted, err := TopologicalSortServices(c.Services)
		if err != nil {
			if c.Error != nil {
				continue
			} else {
				t.Fatalf("expected error %v, got %v", c.Error, err)
			}
		}

		if !expectedTopologicalSortServices(sorted, c.Expected) {
			t.Fatalf("expected %v, got %s", c.Expected, getNames(sorted))
		}
	}
}

func expectedTopologicalSortServices(services model.Services, expected []string) bool {
	sortNames := make([]string, len(services))
	for i, service := range services {
		sortNames[i] = service.Name
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

func getNames(s model.Services) []string {
	sort := []string{}
	for i := 0; i < len(s); i++ {
		sort = append(sort, s[i].Name)
	}

	return sort
}
