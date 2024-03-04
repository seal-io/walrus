package reflect_type

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Dummy is the schema for the projects API.
//
// +k8s:crd-gen:resource:categories=["all","walrus"],scope="Namespaced",shortName=["proj"],plural="projects",subResources=["status"]
type Dummy struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec   DummySpec   `json:"spec,omitempty"`
	Status DummyStatus `json:"status,omitempty"`
}

// DummySpec defines the desired state of Dummy.
type DummySpec struct {
	// Primitive.

	// +optional
	// +k8s:validation:minimum=1
	Integer int `json:"integer,omitempty"`
	// +optional
	Float float64 `json:"float,omitempty"`
	// +optional
	// +k8s:validation:enum=["x","y","z"]
	String string `json:"string,omitempty"`
	// +optional
	Bool bool `json:"bool,omitempty"`
	// +k8s:validation:enum=[1,2,3,5,7]
	// +nullable=false
	IntegerPointer *int64 `json:"integerPointer,omitempty"`
	// +k8s:validation:maximum=10
	FloatPointer *float32 `json:"floatPointer,omitempty"`
	// +k8s:validation:default="x"
	StringPointer *string `json:"stringPointer,omitempty"`
	// +default=true
	BoolPointer *bool `json:"boolPointer,omitempty"`

	// WellKnown Extension.

	MicroTime           meta.MicroTime        `json:"microTime,omitempty"`
	Duration            meta.Duration         `json:"duration,omitempty"`
	Time                meta.Time             `json:"time,omitempty"`
	IntOrString         intstr.IntOrString    `json:"intOrString,omitempty"`
	RawExtension        runtime.RawExtension  `json:"rawExtension,omitempty"`
	MicroTimePointer    *meta.MicroTime       `json:"microTimePointer,omitempty"`
	DurationPointer     *meta.Duration        `json:"durationPointer,omitempty"`
	TimePointer         *meta.Time            `json:"timePointer,omitempty"`
	IntOrStringPointer  *intstr.IntOrString   `json:"intOrStringPointer,omitempty"`
	RawExtensionPointer *runtime.RawExtension `json:"rawExtensionPointer,omitempty"`

	// Map.

	// +mapType="atomic"
	// +k8s:validation:maxProperties=3
	MapString      map[string]string  `json:"mapString,omitempty"`
	MapPointString map[string]*string `json:"mapPointString,omitempty"`
	MapObject      map[string]struct {
		A string              `json:"a,omitempty"`
		B int                 `json:"b,omitempty"`
		C bool                `json:"c,omitempty"`
		D float64             `json:"d,omitempty"`
		E map[string]struct{} `json:"e,omitempty"`
		F []struct{}          `json:"f,omitempty"`
	} `json:"mapObject,omitempty"`
	// +mapType="atomic"
	MapInOrString    map[string]intstr.IntOrString `json:"mapInOrString,omitempty"`
	MapStringPointer *map[string]string            `json:"mapStringPointer,omitempty"`
	MapInterface     map[string]any                `json:"mapInterface,omitempty"`

	// Slice.

	// +k8s:validation:maxItems=3
	// +k8s:validation:minItems=1
	SliceString      []string  `json:"sliceString,omitempty"`
	SlicePointString []*string `json:"slicePointString,omitempty"`
	// +k8s:validation:cel[0]:rule="self.b > 0"
	// +k8s:validation:cel[0]:message="b must be greater than 0"
	SliceObject []struct {
		A string              `json:"a,omitempty"`
		B int                 `json:"b,omitempty"`
		C bool                `json:"c,omitempty"`
		D float64             `json:"d,omitempty"`
		E map[string]struct{} `json:"e,omitempty"`
		F []struct{}          `json:"f,omitempty"`
	} `json:"sliceObject,omitempty"`
	SliceBytes        []byte  `json:"sliceBytes,omitempty"`
	SliceBytesPointer *[]byte `json:"sliceBytesPointer,omitempty"`
	SliceInterface    []any   `json:"sliceInterface,omitempty"`

	// Array.

	ArrayString      [3]string  `json:"arrayString,omitempty"`
	ArrayPointString [3]*string `json:"arrayPointString,omitempty"`
	// +k8s:validation:cel[0]:rule> self.b > 0
	// +k8s:validation:cel[0]:message> b must be greater than 0,
	// +k8s:validation:cel[1]:rule> self.e.length() % 2 == 0
	// +k8s:validation:cel[1]:rule>   ? self.a == self.a + ' is even'
	// +k8s:validation:cel[1]:rule>   : self.a == self.a + ' is odd'
	// +k8s:validation:cel[1]:message> mutate a
	// +k8s:validation:cel[1]:reason="FieldValueRequired"
	ArrayObject [3]struct {
		A string              `json:"a,omitempty"`
		B int                 `json:"b,omitempty"`
		C bool                `json:"c,omitempty"`
		D float64             `json:"d,omitempty"`
		E map[string]struct{} `json:"e,omitempty"`
	}

	// Reference.

	SubDummySpecPointer *DummySpec  `json:"subDummySpecPointer,omitempty"`
	SliceSubDummySpec   []DummySpec `json:"sliceSubDummySpec,omitempty"`
}

// DummyStatus defines the observed state of Dummy.
type DummyStatus struct {
}
