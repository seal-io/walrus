package platform

import (
	"context"
	"fmt"

	"github.com/seal-io/seal/pkg/operatorany"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platform/operator"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/platformtf"
)

// GetDeployer returns deployer.Deployer with the given deployer.CreateOptions.
func GetDeployer(ctx context.Context, opts deployer.CreateOptions) (deployer.Deployer, error) {
	f, exist := dpCreators[opts.Type]
	if !exist {
		return nil, fmt.Errorf("unknown deployer: %s", opts.Type)
	}
	dp, err := f(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error connecting %s deployer: %w", opts.Type, err)
	}
	return dp, nil
}

var dpCreators map[deployer.Type]deployer.Creator

func init() {
	// Register deployer creators as below.
	dpCreators = map[deployer.Type]deployer.Creator{
		platformtf.DeployerType: platformtf.NewDeployer,
	}
}

// GetOperator returns operator.Operator with the given operator.CreateOptions.
func GetOperator(ctx context.Context, opts operator.CreateOptions) (op operator.Operator, err error) {
	f, exist := opCreators[opts.Connector.Type]
	if !exist {
		// Try to create an any operator.
		op, err = operatorany.NewOperator(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("unknown operator: %s", opts.Connector.Type)
		}
	} else {
		op, err = f(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("error connecting %s operator: %w", opts.Connector.Type, err)
		}
	}
	return op, nil
}

var opCreators map[operator.Type]operator.Creator

func init() {
	// Register operator creators as below.
	opCreators = map[operator.Type]operator.Creator{
		platformk8s.OperatorType: platformk8s.NewOperator,
	}
}
