package graph

const (
	// NodeTypeResource is the resource node type.
	NodeTypeResource NodeType = "resource"
	// NodeTypeModule is the module node type.
	NodeTypeModule NodeType = "module"
	// NodeTypeInstance is the instance node type.
	NodeTypeInstance NodeType = "instance"
)

type LinkType string

// Link define the relationship between two nodes.
type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type NodeType string

func (t NodeType) ID(id string) string {
	return string(t) + "/" + id
}

// Node define the node of the graph.
type Node struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        NodeType `json:"type"`
	Labels      []string `json:"labels,omitempty"`
	Description string   `json:"description,omitempty"`
	// Data is the data of the node.
	// It can be a resource, a module, or an instance.
	Data any `json:"data,omitempty"`
	// ParentNode is the parent node of the node.
	ParentNode string `json:"parentNode,omitempty"`
}

type Graph struct {
	Nodes []*Node `json:"nodes"`
	Links []*Link `json:"links"`
}
