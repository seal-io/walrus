package openapi

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Stack is a stack of nodes.
type Stack struct {
	data []Node
}

// Push pushes a node to the stack.
func (s *Stack) Push(node Node) {
	s.data = append(s.data, node)
}

// Pop pops a node from the stack.
func (s *Stack) Pop() Node {
	if len(s.data) == 0 {
		return nil
	}

	node := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return node
}

// Len return the current item length in stack.
func (s *Stack) Len() int {
	return len(s.data)
}

// Node is a node in the tree.
type Node interface {
	GetID() string
	GetParentID() string
	Children() []Node
	GetDefault() any
}

// DefaultPatchNode is a implement of Node to generate default value tree.
type DefaultPatchNode struct {
	*openapi3.Schema

	id  string
	def any
}

// GetID returns the id of the node.
func (n *DefaultPatchNode) GetID() string {
	return n.id
}

// GetParentID returns the id of the ancestor nodes.
func (n *DefaultPatchNode) GetParentID() string {
	arr := strings.Split(n.id, ".")
	if len(arr) < 2 {
		return ""
	}

	return strings.Join(arr[:len(arr)-1], ".")
}

// Children returns the children of the node.
func (n *DefaultPatchNode) Children() []Node {
	switch n.Schema.Type {
	case openapi3.TypeObject:
		children := make([]Node, 0, len(n.Schema.Properties))

		for pn, prop := range n.Schema.Properties {
			sdn := &DefaultPatchNode{
				Schema: prop.Value,
				id:     getID(n.GetID(), pn),
				def:    prop.Value.Default,
			}

			children = append(children, sdn)
		}

		return children
	case openapi3.TypeArray:
		if n.Default == nil {
			return nil
		}

		return []Node{
			&DefaultPatchNode{
				Schema: n.Schema.Items.Value,
				id:     n.GetID() + ".0",
				def:    n.Default,
			},
		}
	}

	return nil
}

// GetDefault returns the default value of the node.
func (n *DefaultPatchNode) GetDefault() any {
	return n.def
}

func getID(pid, cid string) string {
	if pid != "" {
		return pid + "." + cid
	}

	return cid
}
