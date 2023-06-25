package service

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/log"
)

type ServiceNode struct {
	ID       oid.ID         `json:"id"`
	Name     string         `json:"name"`
	Service  *model.Service `json:"service"`
	Children []*ServiceNode `json:"children"`
}

// NewServiceNodes returns a list of service nodes from the given services.
func NewServiceNodes(services model.Services) ([]*ServiceNode, error) {
	var (
		serviceMap = make(map[string]*model.Service)
		nodeMap    = make(map[string]*ServiceNode)
		// The map is a service name to its child dependency services' names.
		dependencyMap = make(map[string][]string)
	)

	for i := range services {
		svc := services[i]

		if _, ok := serviceMap[svc.Name]; ok {
			return nil, fmt.Errorf("duplicate service name: %s", svc.Name)
		}

		if _, ok := dependencyMap[svc.Name]; !ok {
			dependencyMap[svc.Name] = make([]string, 0)
		}

		dependencyNames := dao.GetDependencyNames(svc)

		if len(dependencyNames) > 0 {
			for _, name := range dependencyNames {
				dependencyMap[name] = append(dependencyMap[name], svc.Name)
			}
		}

		serviceMap[svc.Name] = svc
	}

	nodesFromDependency := newServiceNodeFromDependency(dependencyMap, nodeMap)
	nodes := make([]*ServiceNode, 0, len(nodesFromDependency))

	for _, node := range nodesFromDependency {
		if _, ok := serviceMap[node.Name]; !ok {
			continue
		}

		node.Service = serviceMap[node.Name]
		nodes = append(nodes, node)
	}

	if err := checkServiceNodesCycle(nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

// newServiceNodeFromDependency build a list of service nodes from the dependency map.
func newServiceNodeFromDependency(dependencyMap map[string][]string, nodeMap map[string]*ServiceNode) []*ServiceNode {
	logger := log.WithName("service").WithName("topology")

	for name := range dependencyMap {
		nodeMap[name] = &ServiceNode{
			Name:     name,
			Children: make([]*ServiceNode, 0),
		}
	}

	for name, dependencies := range dependencyMap {
		node := nodeMap[name]

		for _, dependencyName := range dependencies {
			if _, ok := nodeMap[dependencyName]; ok {
				dependencyNode := nodeMap[dependencyName]
				node.Children = append(node.Children, dependencyNode)
			} else {
				logger.Warnf("name: %s, %v, dependency %s not found\n", name, dependencies, dependencyName)
			}
		}
	}

	nodes := make([]*ServiceNode, 0, len(nodeMap))

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

// detectServiceNodeCycle using DFS to detect if there is a cycle in the service nodes.
func detectServiceNodeCycle(nodes []*ServiceNode) (bool, []*ServiceNode) {
	visited := make(map[string]bool, len(nodes))
	stack := make(map[string]bool, len(nodes))

	for _, node := range nodes {
		visited[node.Name] = false
		stack[node.Name] = false
	}

	var cyclePath []*ServiceNode

	var dfs func(node *ServiceNode) bool
	dfs = func(node *ServiceNode) bool {
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

// checkServiceNodesCycle checks if there is a cycle in the service nodes.
func checkServiceNodesCycle(nodes []*ServiceNode) error {
	hasCycle, cyclePath := detectServiceNodeCycle(nodes)
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

// TopologicalSortServiceNodes sorts the service nodes by dependencies.
func TopologicalSortServiceNodes(nodes []*ServiceNode) []*ServiceNode {
	inDegree := make(map[*ServiceNode]int, len(nodes))

	for _, node := range nodes {
		inDegree[node] = 0
	}

	for _, node := range nodes {
		for _, dependency := range node.Children {
			inDegree[dependency]++
		}
	}

	queue := []*ServiceNode{}

	for _, node := range nodes {
		if inDegree[node] == 0 {
			queue = append(queue, node)
		}
	}

	result := []*ServiceNode{}

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

// TopologicalSortServices sorts the services by dependencies.
func TopologicalSortServices(services model.Services) (model.Services, error) {
	var sortedServices model.Services

	serviceNodes, err := NewServiceNodes(services)
	if err != nil {
		return nil, err
	}

	sortedServiceNodes := TopologicalSortServiceNodes(serviceNodes)

	for _, node := range sortedServiceNodes {
		sortedServices = append(sortedServices, node.Service)
	}

	return sortedServices, nil
}
