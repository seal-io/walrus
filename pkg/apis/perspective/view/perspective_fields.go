package view

import (
	"github.com/seal-io/seal/pkg/dao/types"
)

var BuiltInPerspectiveFields = []PerspectiveField{
	PF("Connector", types.FilterFieldConnectorID),
	PF("Project", types.FilterFieldProject),
	PF("Environment", types.FilterFieldEnvironment),
	PF("Application", types.FilterFieldApplication),
	PF("Cluster Name", types.FilterFieldClusterName),
	PF("Namespace", types.FilterFieldNamespace),
	PF("Node", types.FilterFieldNode),
	PF("Controller", types.FilterFieldController),
	PF("Controller Kind", types.FilterFieldControllerKind),
	PF("Pod", types.FilterFieldPod),
	PF("Container", types.FilterFieldContainer),
}

type PerspectiveField struct {
	Label     string            `json:"label"`
	FieldName types.FilterField `json:"fieldName"`
}

type PerspectiveValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func PF(label string, fieldName types.FilterField) PerspectiveField {
	return PerspectiveField{
		Label:     label,
		FieldName: fieldName,
	}
}

func LabelKeyToPF(labelKey string) PerspectiveField {
	return PerspectiveField{
		Label:     labelKey,
		FieldName: types.FilterField(types.LabelPrefix + labelKey),
	}
}
