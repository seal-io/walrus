package modelchange

import (
	"context"
	"time"

	settingbus "github.com/seal-io/walrus/pkg/bus/setting"
	"github.com/seal-io/walrus/pkg/dao/migration"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/migrate"
	"github.com/seal-io/walrus/pkg/dao/model/setting"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/database"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
	"github.com/seal-io/walrus/utils/topic"
)

// TableNames returns the name list of the interested tables for establishing.
func TableNames() []string {
	return tableNameSet.List()
}

type handler struct {
	logger      log.Logger
	modelClient model.ClientSet
	buffer      *eventBuffer
}

// Handle returns an implementation of database.ListenHandler
// for handling the data changes.
func Handle(ctx context.Context, mc model.ClientSet) database.ListenHandler {
	logger := log.WithName("model-change")

	return handler{
		logger:      logger,
		modelClient: mc,
		buffer:      newEventBuffer(ctx, logger),
	}
}

func (handler) Channels() []string {
	return []string{
		migration.ModelChangeChannel,
	}
}

// operation holds the operation type of the database.
type operation string

const (
	operationInsert operation = "insert"
	operationUpdate operation = "update"
	operationDelete operation = "delete"
)

// EventType converts the database operation to the event type.
func (op operation) EventType() EventType {
	switch op {
	case operationInsert:
		return EventTypeCreate
	case operationUpdate:
		return EventTypeUpdate
	case operationDelete:
		return EventTypeDelete
	}

	return _EventTypeUnknown
}

func (h handler) Handle(ctx context.Context, _, payload string) {
	if payload == "" {
		return
	}

	var le struct {
		Timestamp   time.Time `json:"ts"`
		Operation   operation `json:"op"`
		TableSchema string    `json:"tb_s"`
		TableName   string    `json:"tb_n"`
		Rows        []struct {
			ID            object.ID `json:"id"`
			ProjectID     object.ID `json:"project_id"`
			EnvironmentID object.ID `json:"environment_id"`
		} `json:"ids"`
	}

	if err := json.Unmarshal(strs.ToBytes(&payload), &le); err != nil {
		h.logger.Warnf("error unmarshalling payload: %v", err)
		return
	}

	if !tableNameSet.Has(le.TableName) {
		h.logger.Warnf("unknown table name: %s", le.TableName)
		return
	}

	logger := h.logger.WithValues("event", le)

	data := make([]EventData, len(le.Rows))

	for i := range le.Rows {
		data[i].ID = le.Rows[i].ID
		data[i].ProjectID = le.Rows[i].ProjectID
		data[i].EnvironmentID = le.Rows[i].EnvironmentID
	}

	event := Event{
		Type: le.Operation.EventType(),
		Data: data,
	}

	// If the event is a settings update, notify the setting bus.
	if le.TableName == migrate.SettingsTable.Name {
		if le.Operation == operationUpdate {
			settings, err := h.modelClient.Settings().Query().
				Where(setting.IDIn(event.IDs()...)).
				Select(
					setting.FieldID,
					setting.FieldName,
					setting.FieldValue).
				All(ctx)
			if err != nil {
				logger.Errorf("error querying settings: %v", err)
				return
			}

			if err = settingbus.Notify(ctx, h.modelClient, settings); err != nil {
				logger.Errorf("error notifying setting bus: %v", err)
			}
		}

		return
	}

	// Otherwise, write event to the corresponding topic.
	h.buffer.Write(ctx, topic.Topic(le.TableName), event)
}
