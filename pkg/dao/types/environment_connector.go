package types

// EnvironmentConnector is a snapshot of model.EnvironmentConnectorRelationship to avoid cycle importing.
type EnvironmentConnector struct {
	// ID of connector that configure to the environment.
	ConnectorID ID `json:"connectorID"`
}
