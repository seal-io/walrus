package projectsubject

import (
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subjectrolerelationship"
)

type DeleteRequest = model.SubjectRoleRelationshipDeleteInput

type (
	CollectionCreateRequest = model.SubjectRoleRelationshipCreateInputs

	CollectionCreateResponse = []*model.SubjectRoleRelationshipOutput
)

type (
	CollectionGetRequest struct {
		model.SubjectRoleRelationshipQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.SubjectRoleRelationship, subjectrolerelationship.OrderOption,
		] `query:",inline"`
	}

	CollectionGetResponse = []*model.SubjectRoleRelationshipOutput
)

type CollectionDeleteRequest = model.SubjectRoleRelationshipDeleteInputs
