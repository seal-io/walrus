package resourcestatus

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/api/option"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	gtypes "github.com/seal-io/walrus/pkg/operator/google/types"
	"github.com/seal-io/walrus/pkg/operator/types"
)

func getSQLDatabaseInstance(ctx context.Context, resourceType, name string) (*status.Status, error) {
	cred, ok := ctx.Value(types.CredentialKey).(*gtypes.Credential)
	if !ok {
		return nil, errors.New("not found credential from context")
	}

	service, err := sqladmin.NewService(ctx, option.WithCredentialsJSON([]byte(cred.Credentials)))
	if err != nil {
		return nil, err
	}

	resp, err := service.Instances.Get(cred.Project, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get google resource %s %s: %w", resourceType, name, err)
	}

	return sqlDatabaseInstanceStatusConverter.Convert(strings.ToLower(resp.State), ""), nil
}
