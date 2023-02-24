package types

const (
	// LabelPrefix is used for generate label's field names.
	LabelPrefix = "label:"
	// UnallocatedLabel indicate the cost for the resources unallocated.
	UnallocatedLabel = "__unallocated__"
)

// built-in labels.
const (
	LabelSealProject     string = "seal.io/project"
	LabelSealEnvironment string = "seal.io/environment"
	LabelSealApplication string = "seal.io/app"
)
