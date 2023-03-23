package connectors

import (
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/scm/driver/github"
	"github.com/seal-io/seal/pkg/scm/driver/gitlab"
)

// TODO check by category

// IsVCS checks if the given connector is a version control system.
func IsVCS(conn *model.Connector) bool {
	return conn.Type == github.Driver || conn.Type == gitlab.Driver
}

// IsOperator checks if the given connector is a known operator.
func IsOperator(conn *model.Connector) bool {
	return conn.Type == platformk8s.OperatorType
}
