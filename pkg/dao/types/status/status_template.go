package status

const (
	TemplateStatusInitialized ConditionType = "Initialized"
	TemplateStatusReady       ConditionType = "Ready"
)

// templateStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Initialized      | Unknown                 | Initializing          | Transitioning         |
//	| Initialized      | False                   | InitializeFailed      | Error                 |
//	| Initialized      | True                    | Initialized           |                       |
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | NotReady              | Error                 |
//	| Ready            | True                    | Ready                 |                       |
var templateStatusPaths = NewWalker(
	[][]ConditionType{
		{
			TemplateStatusInitialized,
			TemplateStatusReady,
		},
	},
)

func WalkTemplate(st *Status) *Summary {
	return templateStatusPaths.Walk(st)
}
