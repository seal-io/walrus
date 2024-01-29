package types

import (
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/templates/openapi"
)

const WalrusContextVariableName = openapi.WalrusContextVariableName

// Context indicates the walrus-related metadata,
// will set to attribute context while user module include this attribute.
type Context struct {
	// Project indicates the project metadata.
	Project struct {
		Name string    `json:"name,omitempty"`
		ID   object.ID `json:"id,omitempty"`
	} `json:"project,omitempty"`

	// Environment indicate the environment metadata.
	Environment struct {
		Name string    `json:"name,omitempty"`
		ID   object.ID `json:"id,omitempty"`
		// Namespace is the managed namespace name of an environment,
		// valid when Kubernetes connector is used in the environment.
		Namespace string `json:"namespace,omitempty"`
	} `json:"environment,omitempty"`

	// Resource indicates the resource metadata.
	Resource struct {
		Name string    `json:"name,omitempty"`
		ID   object.ID `json:"id,omitempty"`
	} `json:"resource,omitempty"`
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) SetProject(id object.ID, name string) *Context {
	c.Project.ID = id
	c.Project.Name = name

	return c
}

func (c *Context) SetEnvironment(id object.ID, name, namespace string) *Context {
	c.Environment.ID = id
	c.Environment.Name = name
	c.Environment.Namespace = namespace

	return c
}

func (c *Context) SetResource(id object.ID, name string) *Context {
	c.Resource.ID = id
	c.Resource.Name = name

	return c
}
