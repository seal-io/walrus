package view

import (
	"github.com/seal-io/seal/pkg/dao/types"
)

var BuiltInPerspectiveFilterFields = []PerspectiveField{
	FilterPerspectiveField("Connector", types.FilterFieldConnectorID),
	FilterPerspectiveField("Project", types.FilterFieldProject),
	FilterPerspectiveField("Environment", types.FilterFieldEnvironmentPath),
	FilterPerspectiveField("Service", types.FilterFieldServicePath),
	FilterPerspectiveField("Cluster Name", types.FilterFieldClusterName),
	FilterPerspectiveField("Namespace", types.FilterFieldNamespace),
	FilterPerspectiveField("Node", types.FilterFieldNode),
	FilterPerspectiveField("Controller", types.FilterFieldController),
	FilterPerspectiveField("Controller Kind", types.FilterFieldControllerKind),
	FilterPerspectiveField("Pod", types.FilterFieldPod),
	FilterPerspectiveField("Container", types.FilterFieldContainer),
}

var BuiltInPerspectiveGroupFields = []PerspectiveField{
	GroupByPerspectiveField("Connector", types.GroupByFieldConnectorID),
	GroupByPerspectiveField("Project", types.GroupByFieldProject),
	GroupByPerspectiveField("Environment", types.GroupByFieldEnvironmentPath),
	GroupByPerspectiveField("Service", types.GroupByFieldServicePath),
	GroupByPerspectiveField("Cluster Name", types.GroupByFieldClusterName),
	GroupByPerspectiveField("Namespace", types.GroupByFieldNamespace),
	GroupByPerspectiveField("Node", types.GroupByFieldNode),
	GroupByPerspectiveField("Controller", types.GroupByFieldController),
	GroupByPerspectiveField("Controller Kind", types.GroupByFieldControllerKind),
	GroupByPerspectiveField("Pod", types.GroupByFieldPod),
	GroupByPerspectiveField("Container", types.GroupByFieldContainer),
	GroupByPerspectiveField("Workload", types.GroupByFieldWorkload),
	GroupByPerspectiveField("Day", types.GroupByFieldDay),
	GroupByPerspectiveField("Week", types.GroupByFieldWeek),
	GroupByPerspectiveField("Month", types.GroupByFieldMonth),
}

var BuiltInPerspectiveStepFields = []PerspectiveField{
	StepPerspectiveField("Cumulative", ""),
	StepPerspectiveField("Daily", types.StepDay),
	StepPerspectiveField("Weekly", types.StepWeek),
	StepPerspectiveField("Monthly", types.StepMonth),
}

type PerspectiveField struct {
	Label     string `json:"label"`
	FieldName string `json:"fieldName"`
}

type PerspectiveValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func FilterPerspectiveField(label string, fieldName types.FilterField) PerspectiveField {
	return PerspectiveField{
		Label:     label,
		FieldName: string(fieldName),
	}
}

func GroupByPerspectiveField(label string, fieldName types.GroupByField) PerspectiveField {
	return PerspectiveField{
		Label:     label,
		FieldName: string(fieldName),
	}
}

func StepPerspectiveField(label string, fieldName types.Step) PerspectiveField {
	return PerspectiveField{
		Label:     label,
		FieldName: string(fieldName),
	}
}

func LabelKeyToPerspectiveField(labelKey string) PerspectiveField {
	return PerspectiveField{
		Label:     labelKey,
		FieldName: types.LabelPrefix + labelKey,
	}
}
