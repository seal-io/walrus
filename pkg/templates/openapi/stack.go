package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tidwall/sjson"

	"github.com/seal-io/walrus/utils/log"
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
	Children() []Node
	GetPatch() []byte
}

// DefaultPatchNode is a implement of Node to generate default value tree.
type DefaultPatchNode struct {
	*openapi3.Schema

	id    string
	patch []byte
}

// GetID returns the id of the node.
func (n *DefaultPatchNode) GetID() string {
	return n.id
}

// Children returns the children of the node.
func (n *DefaultPatchNode) Children() []Node {
	var children []Node

	switch n.Schema.Type {
	case openapi3.TypeObject:
		for pn, prop := range n.Schema.Properties {
			var (
				p   = getID(n.GetID(), pn)
				sdn = &DefaultPatchNode{
					Schema: prop.Value,
					id:     p,
				}
			)

			if prop.Value.Default != nil {
				var b []byte

				patch, err := sjson.SetBytes(b, p, prop.Value.Default)
				if err != nil {
					log.Errorf("failed to set default value for %s: %v", p, err)
				} else {
					sdn.patch = patch
				}
			}

			children = append(
				children,
				sdn,
			)
		}
	case openapi3.TypeArray:
		if n.Default != nil {
			var (
				p   = n.GetID() + ".0"
				sdn = &DefaultPatchNode{
					Schema: n.Schema.Items.Value,

					id: p,
				}
			)

			var b []byte

			patch, err := sjson.SetBytes(b, p, n.Default)
			if err != nil {
				log.Errorf("failed to set default value for %s: %v", p, err)
			} else {
				sdn.patch = patch
			}

			children = append(children, sdn)
		}
	}

	return children
}

// GetPatch returns the patch of the node.
func (n *DefaultPatchNode) GetPatch() []byte {
	return n.patch
}

func getID(pid, cid string) string {
	if pid != "" {
		return pid + "." + cid
	}

	return cid
}
