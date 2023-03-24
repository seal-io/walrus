package connectors

import (
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
)

// IsVCS checks if the given connector is a version control system.
func IsVCS(conn *model.Connector) bool {
	return conn.Category == types.ConnectorCategoryVersionControl
}

// IsOperator checks if the given connector is a known operator.
func IsOperator(conn *model.Connector) bool {
	return conn.Type == types.ConnectorCategoryKubernetes
}
