package resourceexec

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/mmmorris1975/ssm-session-client/datachannel"

	opawstypes "github.com/seal-io/seal/pkg/operator/aws/types"
	optypes "github.com/seal-io/seal/pkg/operator/types"
	"github.com/seal-io/seal/utils/gopool"
)

const (
	defaultTerminalHeight uint16 = 100
	defaultTerminalWidth  uint16 = 100
)

type ec2Instance struct {
	awsCfg *aws.Config
}

func getEc2Instance(ctx context.Context) (optypes.ExecutableResource, error) {
	awsCfg, err := opawstypes.ConfigFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return &ec2Instance{
		awsCfg: awsCfg,
	}, nil
}

// Supported check whether the session manager configured for current instance.
func (r *ec2Instance) Supported(ctx context.Context, name string) (bool, error) {
	ssmCli := ssm.NewFromConfig(*r.awsCfg)

	req := &ssm.DescribeInstanceInformationInput{
		InstanceInformationFilterList: []ssmtypes.InstanceInformationFilter{
			{
				Key:      ssmtypes.InstanceInformationFilterKey("InstanceIds"),
				ValueSet: []string{name},
			},
		},
	}

	resp, err := ssmCli.DescribeInstanceInformation(ctx, req)
	if err != nil {
		return false, err
	}

	if len(resp.InstanceInformationList) == 0 ||
		resp.InstanceInformationList[0].PingStatus != ssmtypes.PingStatusOnline {
		return false, nil
	}

	return true, nil
}

// Exec support read and write operations with the connection.
func (r *ec2Instance) Exec(ctx context.Context, name string, opts optypes.ExecOptions) error {
	c := new(datachannel.SsmDataChannel)
	if err := c.Open(*r.awsCfg, &ssm.StartSessionInput{Target: aws.String(name)}); err != nil {
		return err
	}
	defer c.Close()

	eg := gopool.Group()
	eg.Go(func() error {
		return r.setTerminalSize(ctx, opts, c)
	})
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, err := io.Copy(c, opts.In)
			if err != nil {
				return fmt.Errorf("error write to remote: %w", err)
			}

			return nil
		}
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, err := io.Copy(opts.Out, c)
			if err != nil {
				return fmt.Errorf("error read from remote: %w", err)
			}

			return nil
		}
	})

	return eg.Wait()
}

func (r *ec2Instance) setTerminalSize(
	ctx context.Context,
	opts optypes.ExecOptions,
	c *datachannel.SsmDataChannel,
) error {
	// Without resizer.
	if opts.Resizer == nil {
		err := c.SetTerminalSize(uint32(defaultTerminalHeight), uint32(defaultTerminalWidth))
		if err != nil {
			return fmt.Errorf("error set terminal size")
		}

		return nil
	}

	// With resizer.
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		width, height, ok := opts.Resizer.Next()
		if !ok {
			return errors.New("invalid terminal resizer")
		}

		// Send resize data to remote connection.
		err := c.SetTerminalSize(uint32(height), uint32(width))
		if err != nil {
			return err
		}
	}
}
