package perspective

import (
	"github.com/seal-io/walrus/pkg/dao/types"
)

var (
	builtinFilterFields = []Field{
		getFilterField("Connector", types.FilterFieldConnectorID),
		getFilterField("Project", types.FilterFieldProject),
		getFilterField("Environment", types.FilterFieldEnvironmentPath),
		getFilterField("Service", types.FilterFieldServicePath),
		getFilterField("Cluster Name", types.FilterFieldClusterName),
		getFilterField("Namespace", types.FilterFieldNamespace),
		getFilterField("Node", types.FilterFieldNode),
		getFilterField("Controller", types.FilterFieldController),
		getFilterField("Controller Kind", types.FilterFieldControllerKind),
		getFilterField("Pod", types.FilterFieldPod),
		getFilterField("Container", types.FilterFieldContainer),
		getFilterField("Name", types.FilterFieldName),
	}

	builtinGroupFields = []Field{
		getGroupField("Connector", types.GroupByFieldConnectorID),
		getGroupField("Project", types.GroupByFieldProject),
		getGroupField("Environment", types.GroupByFieldEnvironmentPath),
		getGroupField("Service", types.GroupByFieldServicePath),
		getGroupField("Cluster Name", types.GroupByFieldClusterName),
		getGroupField("Namespace", types.GroupByFieldNamespace),
		getGroupField("Node", types.GroupByFieldNode),
		getGroupField("Controller", types.GroupByFieldController),
		getGroupField("Controller Kind", types.GroupByFieldControllerKind),
		getGroupField("Pod", types.GroupByFieldPod),
		getGroupField("Container", types.GroupByFieldContainer),
		getGroupField("Workload", types.GroupByFieldWorkload),
		getGroupField("Day", types.GroupByFieldDay),
		getGroupField("Week", types.GroupByFieldWeek),
		getGroupField("Month", types.GroupByFieldMonth),
	}

	builtinStepFields = []Field{
		getStepField("Cumulative", ""),
		getStepField("Daily", types.StepDay),
		getStepField("Weekly", types.StepWeek),
		getStepField("Monthly", types.StepMonth),
	}
)

type (
	Field struct {
		Label     string `json:"label"`
		FieldName string `json:"fieldName"`
	}

	Value struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}
)

func getFilterField(label string, fieldName types.FilterField) Field {
	return Field{
		Label:     label,
		FieldName: string(fieldName),
	}
}

func getGroupField(label string, fieldName types.GroupByField) Field {
	return Field{
		Label:     label,
		FieldName: string(fieldName),
	}
}

func getStepField(label string, fieldName types.Step) Field {
	return Field{
		Label:     label,
		FieldName: string(fieldName),
	}
}
