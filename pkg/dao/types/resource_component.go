package types

const (
	// ResourceComponentModeManaged indicates the resource created to target platform,
	// it is writable(update or delete).
	ResourceComponentModeManaged = "managed"

	// ResourceComponentModeData indicates the resource read from target platform,
	// it is read-only.
	ResourceComponentModeData = "data"

	// ResourceComponentModeDiscovered indicates the resource discovered from target platform,
	// it inherits its composition's characteristic to be writable or not.
	ResourceComponentModeDiscovered = "discovered"
)

const (
	// ResourceComponentShapeClass indicates the resource is a class.
	ResourceComponentShapeClass = "class"
	// ResourceComponentShapeInstance indicates the resource is an instance that implements a class.
	ResourceComponentShapeInstance = "instance"
)

// ResourceComponentRelationshipTypeDependency indicates the relationship between resource component
// and its dependencies.
const ResourceComponentRelationshipTypeDependency = "Dependency"

type ResourceComponentOperationKeys struct {
	// Labels stores label of layer,
	// its length means each key contains levels with the same value as level.
	Labels []string `json:"labels,omitempty"`
	// Keys stores key in tree.
	Keys []ResourceComponentOperationKey `json:"keys,omitempty"`
}

// ResourceComponentOperationKey holds hierarchy query keys.
type ResourceComponentOperationKey struct {
	// Keys indicates the subordinate keys,
	// usually, it should not be valued in leaves.
	Keys []ResourceComponentOperationKey `json:"keys,omitempty"`
	// Name indicates the name of the key.
	Name string `json:"name"`
	// Value indicates the value of the key,
	// usually, it should be valued in leaves.
	Value string `json:"value,omitempty"`
	// Loggable indicates whether to be able to get log.
	Loggable *bool `json:"loggable,omitempty"`
	// Executable indicates whether to be able to execute remote command.
	Executable *bool `json:"executable,omitempty"`
}
