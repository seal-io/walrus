package openapi

// PreorderTraversal traverses the tree in preorder.
func PreorderTraversal(root Node) []Node {
	if root == nil {
		return nil
	}

	var result []Node
	stack := &Stack{}
	stack.Push(root)

	for stack.Len() > 0 {
		node := stack.Pop()

		result = append(result, node)

		// Push the children to the stack.
		for _, child := range node.Children() {
			stack.Push(child)
		}
	}

	return result
}
