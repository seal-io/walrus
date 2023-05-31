package graph

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/seal/pkg/applicationresources"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/utils/strs"
)

// Using regexp to match all module dependencies.
var moduleRegexp = regexp.MustCompile(`\${module\.([^.\s]+)\.[^}]+}`)

// GetModuleLinks returns the module links, e.g. module webservice rely on module database.
func GetModuleLinks(relationships model.ApplicationModuleRelationships) []*Link {
	var (
		links    = make([]*Link, 0)
		linkKeys = sets.NewString()
	)

	relationMap := make(map[string]string, len(relationships))
	for _, r := range relationships {
		relationMap[r.Name] = r.ModuleID
	}

	for _, r := range relationships {
		for _, d := range r.Attributes {
			matches := moduleRegexp.FindAllSubmatch(d, -1)
			for _, m := range matches {
				moduleID, ok := relationMap[string(m[1])]
				if !ok {
					continue
				}

				linkKey := strs.Join("/", moduleID, string(m[1]), r.ModuleID, r.Name)
				if linkKeys.Has(linkKey) {
					continue
				}

				links = append(links, &Link{
					Source: NodeTypeModule.ID(strs.Join("/", moduleID, string(m[1]))),
					Target: NodeTypeModule.ID(strs.Join("/", r.ModuleID, r.Name)),
				})

				linkKeys.Insert(linkKey)
			}
		}
	}

	return links
}

// GetInstanceResourceGraph returns the resource graph of the instance.
func GetInstanceResourceGraph(ctx context.Context, modelClient model.ClientSet, instanceID oid.ID) (*Graph, error) {
	instance, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(instanceID)).
		WithResources(
			func(rq *model.ApplicationResourceQuery) {
				rq.WithConnector(func(cq *model.ConnectorQuery) {
					cq.Select(
						connector.FieldName,
						connector.FieldType,
						connector.FieldCategory,
						connector.FieldConfigVersion,
						connector.FieldConfigData,
					)
				})
			}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	relationships, err := modelClient.ApplicationModuleRelationships().Query().
		Where(applicationmodulerelationship.ApplicationID(instance.ApplicationID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var (
		nodes     = make([]*Node, 0, len(instance.Edges.Resources)+len(relationships)+1)
		moduleMap = make(map[string]*Node)
	)

	links := GetModuleLinks(relationships)
	resources := applicationresources.GetResourcesDetail(ctx, instance.Edges.Resources, false)

	for _, r := range relationships {
		moduleNode := &Node{
			ID:   NodeTypeModule.ID(strs.Join("/", r.ModuleID, r.Name)),
			Name: r.Name,
			Type: NodeTypeModule,
			Data: model.ExposeApplicationModuleRelationship(r),
		}
		nodes = append(nodes, moduleNode)
		moduleMap[r.Name] = moduleNode
	}

	for i, r := range instance.Edges.Resources {
		moduleNode, ok := moduleMap[strings.Split(r.Module, "/")[0]]
		if !ok {
			return nil, fmt.Errorf("module %s not found", r.Module)
		}

		nodes = append(nodes, &Node{
			ID:         NodeTypeResource.ID(r.ID.String()),
			Name:       r.Name,
			Type:       NodeTypeResource,
			Data:       resources[i],
			ParentNode: moduleNode.ID,
		})

		if r.CompositionID != "" {
			links = append(links, &Link{
				Source: NodeTypeResource.ID(r.CompositionID.String()),
				Target: NodeTypeResource.ID(r.ID.String()),
			})
		} else {
			links = append(links, &Link{
				Source: NodeTypeInstance.ID(instance.ID.String()),
				Target: NodeTypeResource.ID(r.ID.String()),
			})
		}
	}

	instance.Edges.Resources = nil
	nodes = append(nodes, &Node{
		ID:   NodeTypeInstance.ID(instance.ID.String()),
		Name: instance.Name,
		Type: NodeTypeInstance,
		Data: model.ExposeApplicationInstance(instance),
	})

	return &Graph{
		Nodes: nodes,
		Links: links,
	}, nil
}
