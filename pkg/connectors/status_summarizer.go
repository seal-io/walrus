package connectors

import (
	"github.com/seal-io/seal/pkg/dao/types/status"
	pkgstatus "github.com/seal-io/seal/pkg/status"
)

var (
	StatusSummarizer *pkgstatus.Summarizer
)

func init() {
	statusMapping := []map[status.ConditionType]string{
		{
			status.ConnectorStatusProvisioned: status.ConnectorStatusProvisionedTransitioning,
		},
		{
			status.ConnectorStatusToolsDeployed: status.ConnectorStatusToolsDeployedTransitioning,
		},
		{
			status.ConnectorStatusCostSynced: status.ConnectorStatusCostSyncedTransitioning,
		},
		{
			status.ConnectorStatusReady: string(status.ConnectorStatusReady),
		},
	}

	StatusSummarizer = pkgstatus.NewSummarizer(status.ConnectorStatusReady)
	addErrorFalseTransitioningUnknown(StatusSummarizer, statusMapping)
}

func addErrorFalseTransitioningUnknown(summarizer *pkgstatus.Summarizer, mappings []map[status.ConditionType]string) {
	for _, mapping := range mappings {
		for st, tr := range mapping {
			summarizer.AddErrorFalseTransitioningUnknown(st, tr)
		}
	}
}
