package status

const (
	CatalogStatusInitialized ConditionType = "Initialized"
	CatalogStatusReady       ConditionType = "Ready"
)

// catalogStatusPaths makes the following decision.
//
//	|  Condition Type  |     Condition Status    | Human Readable Status | Human Sensible Status |
//	| ---------------- | ----------------------- | --------------------- | --------------------- |
//	| Initialized      | Unknown                 | Initializing          | Transitioning         |
//	| Initialized      | False                   | InitializeFailed      | Error                 |
//	| Initialized      | True                    | Initialized           |                       |
//	| Ready            | Unknown                 | Preparing             | Transitioning         |
//	| Ready            | False                   | NotReady              | Error                 |
//	| Ready            | True                    | Ready                 |                       |
var catalogStatusPaths = NewWalker(
	[][]ConditionType{
		{
			CatalogStatusInitialized,
			CatalogStatusReady,
		},
	},
)

func WalkCatalog(st *Status) *Summary {
	return catalogStatusPaths.Walk(st)
}
