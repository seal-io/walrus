package resources

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/log"
)

type ResourceNode struct {
	ID       object.ID       `json:"id"`
	Name     string          `json:"name"`
	Resource *model.Resource `json:"resource"`
	Children []*ResourceNode `json:"children"`
}

// NewResourceNodes returns a list of resource nodes from the given resources.
func NewResourceNodes(resources model.Resources) ([]*ResourceNode, error) {
	var (
		resourceMap = make(map[string]*model.Resource)
		nodeMap     = make(map[string]*ResourceNode)
		// The map is a resource name to its child dependency resources' names.
		dependencyMap = make(map[string][]string)
	)

	for i := range resources {
		res := resources[i]

		if _, ok := resourceMap[res.Name]; ok {
			return nil, fmt.Errorf("duplicate resource name: %s", res.Name)
		}

		if _, ok := dependencyMap[res.Name]; !ok {
			dependencyMap[res.Name] = make([]string, 0)
		}

		dependencyNames := dao.ResourceRelationshipGetDependencyNames(res)

		if len(dependencyNames) > 0 {
			for _, name := range dependencyNames {
				dependencyMap[name] = append(dependencyMap[name], res.Name)
			}
		}

		resourceMap[res.Name] = res
	}

	nodesFromDependency := newResourceNodeFromDependency(dependencyMap, nodeMap)
	nodes := make([]*ResourceNode, 0, len(nodesFromDependency))

	for _, node := range nodesFromDependency {
		if _, ok := resourceMap[node.Name]; !ok {
			continue
		}

		node.Resource = resourceMap[node.Name]
		nodes = append(nodes, node)
	}

	if err := checkResourceNodesCycle(nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

// newResourceNodeFromDependency build a list of resource nodes from the dependency map.
func newResourceNodeFromDependency(
	dependencyMap map[string][]string,
	nodeMap map[string]*ResourceNode,
) []*ResourceNode {
	logger := log.WithName("resource").WithName("topology")

	for name := range dependencyMap {
		nodeMap[name] = &ResourceNode{
			Name:     name,
			Children: make([]*ResourceNode, 0),
		}
	}

	for name, dependencies := range dependencyMap {
		node := nodeMap[name]

		for _, dependencyName := range dependencies {
			if _, ok := nodeMap[dependencyName]; ok {
				dependencyNode := nodeMap[dependencyName]
				node.Children = append(node.Children, dependencyNode)
			} else {
				logger.Warnf("name: %s, %v, dependency %s not found", name, dependencies, dependencyName)
			}
		}
	}

	nodes := make([]*ResourceNode, 0, len(nodeMap))

	names := make([]string, 0, len(nodeMap))
	for _, node := range nodeMap {
		names = append(names, node.Name)
	}

	slices.Sort(names)

	for _, name := range names {
		nodes = append(nodes, nodeMap[name])
	}

	return nodes
}

// detectResourceNodeCycle using DFS to detect if there is a cycle in the resource nodes.
func detectResourceNodeCycle(nodes []*ResourceNode) (bool, []*ResourceNode) {
	visited := make(map[string]bool, len(nodes))
	stack := make(map[string]bool, len(nodes))

	for _, node := range nodes {
		visited[node.Name] = false
		stack[node.Name] = false
	}

	var cyclePath []*ResourceNode

	var dfs func(node *ResourceNode) bool
	dfs = func(node *ResourceNode) bool {
		visited[node.Name] = true
		stack[node.Name] = true

		if len(node.Children) != 0 {
			for _, child := range node.Children {
				if !visited[child.Name] {
					if dfs(child) {
						cyclePath = append(cyclePath, child)
						return true
					}
				} else if stack[child.Name] {
					cyclePath = append(cyclePath, child)
					return true
				}
			}
		}

		stack[node.Name] = false

		return false
	}

	for _, node := range nodes {
		if !visited[node.Name] {
			if dfs(node) {
				cyclePath = append(cyclePath, node)
				return true, cyclePath
			}
		}
	}

	return false, nil
}

// checkResourceNodesCycle checks if there is a cycle in the resource nodes.
func checkResourceNodesCycle(nodes []*ResourceNode) error {
	hasCycle, cyclePath := detectResourceNodeCycle(nodes)
	if hasCycle {
		path := ""

		for i := len(cyclePath) - 1; i >= 0; i-- {
			if i == 0 {
				path += cyclePath[i].Name
			} else {
				path += cyclePath[i].Name + " -> "
			}
		}

		return fmt.Errorf("cycle detected: %s", path)
	}

	return nil
}

// TopologicalSortResourceNodes sorts the resource nodes by dependencies.
func TopologicalSortResourceNodes(nodes []*ResourceNode) []*ResourceNode {
	inDegree := make(map[*ResourceNode]int, len(nodes))

	for _, node := range nodes {
		inDegree[node] = 0
	}

	for _, node := range nodes {
		for _, dependency := range node.Children {
			inDegree[dependency]++
		}
	}

	var queue []*ResourceNode

	for _, node := range nodes {
		if inDegree[node] == 0 {
			queue = append(queue, node)
		}
	}

	var result []*ResourceNode

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		result = append(result, current)

		for _, dependency := range current.Children {
			inDegree[dependency]--

			if inDegree[dependency] == 0 {
				queue = append(queue, dependency)
			}
		}
	}

	return result
}

// TopologicalSortResources sorts the resources by dependencies.
func TopologicalSortResources(resources model.Resources) (model.Resources, error) {
	var sortedResources model.Resources

	resourceNodes, err := NewResourceNodes(resources)
	if err != nil {
		return nil, err
	}

	sortResourceNodes := TopologicalSortResourceNodes(resourceNodes)

	for _, node := range sortResourceNodes {
		sortedResources = append(sortedResources, node.Resource)
	}

	return sortedResources, nil
}

// ReverseTopologicalSortResources sorts the resource by dependencies in reverse order.
func ReverseTopologicalSortResources(resources model.Resources) (model.Resources, error) {
	var sortedResources model.Resources

	sortResources, err := TopologicalSortResources(resources)
	if err != nil {
		return nil, err
	}

	for i := len(sortResources) - 1; i >= 0; i-- {
		sortedResources = append(sortedResources, sortResources[i])
	}

	return sortedResources, nil
}
