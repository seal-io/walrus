//go:build ginx

package runtime

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/seal-io/seal/utils/json"
)

type H1 struct{}

func (H1) Kind() string {
	return "H1"
}

func (H1) SubResourceHandlers() []IResourceHandler {
	return []IResourceHandler{
		H2{},
	}
}

type (
	H1CreateRequest  struct{}
	H1CreateResponse struct{}
)

func (H1CreateRequest) SetGinContext(*gin.Context) {}

func (H1CreateRequest) SetContext(context.Context) {}

func (H1) Create(H1CreateRequest) (H1CreateResponse, error) {
	return H1CreateResponse{}, nil
}

type (
	H1GetRequest struct {
		ID    string `path:"h1"`
		Watch bool   `query:"watch"`
	}
	H1GetResponse struct{}
)

func (H1) Get(H1GetRequest) (H1GetResponse, error) {
	return H1GetResponse{}, nil
}

type H1UpdateRequest struct {
	json.RawMessage `json:",inline"`
}

func (H1) Update(H1UpdateRequest) error {
	return nil
}

type H1DeleteRequest struct{}

func (H1) Delete(H1DeleteRequest) error {
	return nil
}

type (
	H1CollectionCreateRequest  struct{}
	H1CollectionCreateResponse struct{}
)

func (H1) CollectionCreate(H1CollectionCreateRequest) (H1CollectionCreateResponse, error) {
	return H1CollectionCreateResponse{}, nil
}

type (
	H1CollectionGetRequest  struct{}
	H1CollectionGetResponse struct{}
)

func (*H1CollectionGetRequest) SetStream(RequestUnidiStream) {}

func (H1) CollectionGet(H1CollectionGetRequest) (H1CollectionGetResponse, int, error) {
	return H1CollectionGetResponse{}, 0, nil
}

type H1CollectionUpdateRequest struct{}

func (H1) CollectionUpdate(H1CollectionUpdateRequest) error {
	return nil
}

type H1CollectionDeleteRequest struct{}

func (H1) CollectionDelete(H1CollectionDeleteRequest) error {
	return nil
}

type H1RouteExecRequest struct {
	_ struct{} `route:"GET=/exec"`
}

func (*H1RouteExecRequest) SetStream(RequestBidiStream) {}

func (H1) RouteExec(H1RouteExecRequest) error {
	return nil
}

type H1RouteUpgradeRequest struct {
	_ struct{} `route:"PUT=/upgrade"`
}

func (H1) RouteUpgrade(H1RouteUpgradeRequest) error {
	return nil
}

type H1CollectionRouteFieldsRequest struct {
	_ struct{} `route:"GET=/fields"`
}

func (H1) CollectionRouteFields(H1CollectionRouteFieldsRequest) error {
	return nil
}

type H2 struct{}

func (H2) Kind() string {
	return "H2"
}

type (
	H2GetRequest  struct{}
	H2GetResponse struct{}
)

func (H2GetRequest) SetStream(RequestUnidiStream) {}

func (H2) Get(H2GetRequest) (H2GetResponse, error) {
	return H2GetResponse{}, nil
}

type H3 struct{}

func (H3) Kind() string {
	return "H3"
}

type H3RouteIllegal1Request struct {
	_ struct{} `route:"GET=/"`
}

func (H3) RouteIllegal1(H3RouteIllegal1Request) error {
	return nil
}

type H3RouteIllegal2Request struct {
	_ struct{} `route:"GET=/_/batch"`
}

func (H3) RouteIllegal2(H3RouteIllegal2Request) error {
	return nil
}

type H3RouteAbnormalRequest struct {
	_ struct{} `route:"GET=//./../abnormal"`
}

func (H3) RouteAbnormalX(H3RouteAbnormalRequest) error {
	return nil
}

func (H3) RouteAbnormalX2() error {
	return nil
}

type H3CollectionRouteIllegal1Request struct {
	_ struct{} `route:"GET=/"`
}

func (H3) CollectionRouteIllegal1(H3CollectionRouteIllegal1Request) error {
	return nil
}

type H3CollectionRouteIllegal2Request struct {
	_ struct{} `route:"GET=batch"`
}

func (H3) CollectionRouteIllegal2(H3CollectionRouteIllegal2Request) error {
	return nil
}

type H4 struct{}

func (H4) Kind() string {
	return "H4"
}

type (
	H4CreateRequest  struct{}
	H4CreateResponse struct{}
)

func (H4CreateRequest) SetStream(RequestUnidiStream) {}

func (H4) Create(H4CreateRequest) (H4CreateResponse, error) {
	return H4CreateResponse{}, nil
}

type (
	H4GetRequest  struct{}
	H4GetResponse struct{}
)

func (H4GetRequest) SetStream(RequestBidiStream) {}

func (H4) Get(H4GetRequest) (H4GetResponse, error) {
	return H4GetResponse{}, nil
}

type (
	H4CollectionGetRequest  struct{}
	H4CollectionGetResponse struct{}
)

func (H4CollectionGetRequest) SetStream(RequestUnidiStream) {}

func (H4) CollectionGet(H4CollectionGetRequest) (H4CollectionGetResponse, error) {
	return H4CollectionGetResponse{}, nil
}

type H5 struct{}

func (H5) Kind() string {
	return "H5"
}

func (H5) SubResourceHandlers() []IResourceHandler {
	return []IResourceHandler{
		H5{},
	}
}

func Test_routeResourceHandler(t *testing.T) {
	testCases := []struct {
		name     string
		given    IResourceHandler
		expected []ResourceRoute
	}{
		{
			name:  "full",
			given: H1{},
			expected: []ResourceRoute{
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodPost,
						Path:       "/h1s",
						Collection: false,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage:         "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:            "H1",
					GoFunc:            "Create",
					RequestAttributes: RequestWithGinContext,
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodGet,
						Path:       "/h1s/:h1",
						Collection: false,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage:         "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:            "H1",
					GoFunc:            "Get",
					RequestAttributes: RequestWithBindingPath | RequestWithBindingQuery,
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodPut,
						Path:       "/h1s/:h1",
						Collection: false,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage:         "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:            "H1",
					GoFunc:            "Update",
					RequestAttributes: RequestWithBindingJSON,
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodDelete,
						Path:       "/h1s/:h1",
						Collection: false,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage: "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:    "H1",
					GoFunc:    "Delete",
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodPost,
						Path:       "/h1s/_/batch",
						Collection: true,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage: "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:    "H1",
					GoFunc:    "CollectionCreate",
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodGet,
						Path:       "/h1s",
						Collection: true,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage:          "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:             "H1",
					GoFunc:             "CollectionGet",
					RequestAttributes:  RequestWithUnidiStream,
					ResponseAttributes: ResponseWithPage,
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodPut,
						Path:       "/h1s",
						Collection: true,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage: "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:    "H1",
					GoFunc:    "CollectionUpdate",
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodDelete,
						Path:       "/h1s",
						Collection: true,
						Sub:        false,
						Custom:     false,
						CustomName: "",
					},
					GoPackage: "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:    "H1",
					GoFunc:    "CollectionDelete",
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodGet,
						Path:       "/h1s/:h1/exec",
						Collection: false,
						Sub:        false,
						Custom:     true,
						CustomName: "Exec",
					},
					GoPackage:         "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:            "H1",
					GoFunc:            "RouteExec",
					RequestAttributes: RequestWithBidiStream,
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodPut,
						Path:       "/h1s/:h1/upgrade",
						Collection: false,
						Sub:        false,
						Custom:     true,
						CustomName: "Upgrade",
					},
					GoPackage: "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:    "H1",
					GoFunc:    "RouteUpgrade",
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1"},
							Resources:          []string{"h1s"},
							ResourcePaths:      []string{"h1s"},
							ResourcePathRefers: []string{"h1"},
						},
						Method:     http.MethodGet,
						Path:       "/h1s/_/fields",
						Collection: true,
						Sub:        false,
						Custom:     true,
						CustomName: "Fields",
					},
					GoPackage: "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:    "H1",
					GoFunc:    "CollectionRouteFields",
				},
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H1", "H2"},
							Resources:          []string{"h1s", "h2s"},
							ResourcePaths:      []string{"h1s", "h2s"},
							ResourcePathRefers: []string{"h1", "h2"},
						},
						Method:     http.MethodGet,
						Path:       "/h1s/:h1/h2s/:h2",
						Collection: false,
						Sub:        true,
						Custom:     false,
						CustomName: "",
					},
					GoPackage:         "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:            "H2",
					GoFunc:            "Get",
					RequestAttributes: RequestWithUnidiStream,
				},
			},
		},
		{
			name:  "abnormal custom route",
			given: H3{},
			expected: []ResourceRoute{
				{
					ResourceRouteProfile: ResourceRouteProfile{
						ResourceProfile: ResourceProfile{
							Kinds:              []string{"H3"},
							Resources:          []string{"h3s"},
							ResourcePaths:      []string{"h3s"},
							ResourcePathRefers: []string{"h3"},
						},
						Method:     http.MethodGet,
						Path:       "/h3s/:h3/abnormal",
						Collection: false,
						Sub:        false,
						Custom:     true,
						CustomName: "AbnormalX",
					},
					GoPackage: "github.com/seal-io/seal/pkg/apis/runtime",
					GoType:    "H3",
					GoFunc:    "RouteAbnormalX",
				},
			},
		},
		{
			name:     "illegal stream route",
			given:    H4{},
			expected: nil,
		},
		{
			name:     "circular dependency route",
			given:    H5{},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := routeResourceHandler("", tc.given, ResourceProfile{}, nil)
			// Clear out the fields for comparison.
			for i := range actual {
				actual[i].ResourceRouteProfile.Summary = ""
				actual[i].ResourceRouteProfile.Description = ""
				actual[i].GoCaller = reflect.Value{}
				actual[i].RequestType = nil
				actual[i].ResponseType = nil
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}
