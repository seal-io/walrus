package projectsubject

import "github.com/seal-io/seal/pkg/dao/model"

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (Handler) Kind() string {
	return "ProjectSubject"
}

func (Handler) InternalKind() string {
	// Return the internal kind(DAO) name of the resource.
	return "SubjectRoleRelationship"
}
