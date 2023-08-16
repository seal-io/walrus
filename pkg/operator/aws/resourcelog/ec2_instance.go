package resourcelog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	fetchLogPeriod = 2 * time.Second
)

type ec2Instance struct {
	lastUpdatedTimestamp *time.Time
	cli                  *ec2.Client
	lastContent          string
}

func getEc2Instance(ctx context.Context) (optypes.LoggableResource, error) {
	cli, err := ec2Client(ctx)
	if err != nil {
		return nil, err
	}

	return &ec2Instance{
		cli: cli,
	}, nil
}

func (r *ec2Instance) Log(ctx context.Context, name string, opts optypes.LogOptions) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		resp, err := r.cli.GetConsoleOutput(ctx, &ec2.GetConsoleOutputInput{
			InstanceId: &name,
		})
		if err != nil {
			return fmt.Errorf("error get instance output: %w", err)
		}

		if resp.Timestamp != nil {
			if r.lastUpdatedTimestamp != nil && r.lastUpdatedTimestamp.Equal(*resp.Timestamp) {
				continue
			}
			r.lastUpdatedTimestamp = resp.Timestamp
		}

		if resp.Output != nil {
			var content string

			// The console output is base64 encoded.
			content, err := strs.DecodeBase64(*resp.Output)
			if err != nil {
				return err
			}

			if r.lastContent != "" {
				index := strings.LastIndex(content, r.lastContent)
				if index > 0 {
					content = content[index+1:]
				}
			}

			_, err = opts.Out.Write([]byte(content))
			if err != nil {
				return fmt.Errorf("error write to output: %w", err)
			}

			r.lastContent = strs.LastContent(content, 20)
		}

		if opts.WithoutFollow {
			return nil
		}

		time.Sleep(fetchLogPeriod)
	}
}
