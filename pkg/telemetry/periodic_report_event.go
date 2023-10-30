package telemetry

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/posthog/posthog-go"

	pkgcatalog "github.com/seal-io/walrus/pkg/catalog"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/catalog"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/perspective"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/version"
)

const (
	Application = "Seal"
)

const (
	propertyApplication = "application"
	propertyVersion     = "version"
	propertyPlatform    = "platform"
	propertyServerID    = "server_id"
)

// EnqueuePeriodicReportEvent enqueues the periodic sync event.
func EnqueuePeriodicReportEvent(ctx context.Context, modelClient model.ClientSet) error {
	if APIKey == "" || !settings.EnableTelemetry.ShouldValueBool(ctx, modelClient) {
		// Disable telemetry.
		return nil
	}

	phClient, err := PhClient()
	if err != nil {
		return err
	}
	defer phClient.Close()

	ds := PeriodicReportEvent{
		modelClient: modelClient,
		phClient:    phClient,
	}

	err = ds.Enqueue(ctx)
	if err != nil {
		return fmt.Errorf("error enqueue periodic report event: %w", err)
	}

	logger.Info("enqueued periodic report event")

	return nil
}

type PeriodicReportEvent struct {
	modelClient model.ClientSet
	phClient    posthog.Client
}

func (i PeriodicReportEvent) EventName() string {
	return "periodic_sync"
}

func (i PeriodicReportEvent) Enqueue(ctx context.Context) error {
	ct, err := i.Capture(ctx)
	if err != nil {
		return err
	}

	return i.phClient.Enqueue(ct)
}

func (i PeriodicReportEvent) Capture(ctx context.Context) (*posthog.Capture, error) {
	type setter = func(context.Context, map[string]any) error

	setters := []setter{
		i.setConnectorStat,
		i.setProjectStat,
		i.setEnvironmentStat,
		i.setServiceStat,
		i.setTemplateStat,
		i.setFinOpsStat,
		i.setUserStat,
	}

	installUUID := settings.InstallationUUID.ShouldValue(ctx, i.modelClient)

	prop := posthog.Properties{
		propertyApplication: Application,
		propertyVersion:     version.Get(),
		propertyPlatform:    runtime.GOARCH,
		propertyServerID:    installUUID,
	}

	for _, s := range setters {
		err := s(ctx, prop)
		if err != nil {
			return nil, err
		}
	}

	return &posthog.Capture{
		DistinctId:       installUUID,
		Event:            i.EventName(),
		Properties:       prop,
		SendFeatureFlags: false,
		Timestamp:        time.Now(),
	}, nil
}

func (i PeriodicReportEvent) setConnectorStat(ctx context.Context, props map[string]any) error {
	globalCount, err := i.modelClient.Connectors().Query().Where(connector.ProjectIDIsNil()).Count(ctx)
	if err != nil {
		return err
	}

	projCount, err := i.modelClient.Connectors().Query().Where(connector.ProjectIDNotNil()).Count(ctx)
	if err != nil {
		return err
	}

	var res []struct {
		Count    int    `json:"count"`
		Category string `json:"category"`
		Type     string `json:"type"`
	}

	err = i.modelClient.Connectors().Query().
		GroupBy(
			connector.FieldCategory,
			connector.FieldType,
		).
		Aggregate(
			model.Count(),
		).
		Scan(ctx, &res)
	if err != nil {
		return err
	}

	var (
		categoryStat = make(map[string]int)
		typesStat    = make(map[string]int)
	)

	for _, r := range res {
		categoryStat[r.Category] += r.Count
		typesStat[r.Type] += r.Count
	}

	props["connector_global_level"] = globalCount
	props["connector_project_level"] = projCount
	props["connector_category_stat"] = categoryStat
	props["connector_type_stat"] = typesStat

	return nil
}

func (i PeriodicReportEvent) setProjectStat(ctx context.Context, props map[string]any) error {
	count, err := i.modelClient.Projects().Query().Count(ctx)
	if err != nil {
		return err
	}

	props["project"] = count

	return nil
}

func (i PeriodicReportEvent) setEnvironmentStat(ctx context.Context, props map[string]any) error {
	count, err := i.modelClient.Environments().Query().Count(ctx)
	if err != nil {
		return err
	}

	props["environment"] = count

	return nil
}

func (i PeriodicReportEvent) setServiceStat(ctx context.Context, props map[string]any) error {
	count, err := i.modelClient.Resources().Query().Count(ctx)
	if err != nil {
		return err
	}

	props["service"] = count

	return nil
}

func (i PeriodicReportEvent) setTemplateStat(ctx context.Context, props map[string]any) error {
	catalogID, err := i.modelClient.Catalogs().Query().
		Where(
			catalog.Name(pkgcatalog.BuiltinCatalog().Name),
		).
		OnlyID(ctx)
	if err != nil && !model.IsNotFound(err) {
		return err
	}

	query := i.modelClient.Templates().Query()
	if catalogID != "" {
		query = query.Where(template.CatalogID(catalogID))
	}

	count, err := query.Count(ctx)
	if err != nil {
		return err
	}

	props["template"] = count

	return nil
}

func (i PeriodicReportEvent) setFinOpsStat(ctx context.Context, props map[string]any) error {
	enabledCount, err := i.modelClient.Connectors().Query().Where(connector.EnableFinOps(true)).Count(ctx)
	if err != nil {
		return err
	}

	pt, err := i.modelClient.Perspectives().Query().
		Where(
			perspective.Builtin(false),
		).
		All(ctx)
	if err != nil {
		return err
	}

	props["finOps_enabled"] = enabledCount
	props["perspective_custom"] = len(pt)

	return nil
}

func (i PeriodicReportEvent) setUserStat(ctx context.Context, props map[string]any) error {
	count, err := i.modelClient.Subjects().Query().Count(ctx)
	if err != nil {
		return err
	}

	props["user"] = count

	return nil
}
