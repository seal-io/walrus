package resourcelog

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/strs"
)

const (
	fetchLogPeriod = 2 * time.Second
)

type ecsInstance struct {
	lastUpdatedTimestamp string
	cli                  *ecs.Client
	lastContent          string
}

func getEcsInstance(ctx context.Context) (optypes.LoggableResource, error) {
	cli, err := ecsClient(ctx)
	if err != nil {
		return nil, err
	}

	return &ecsInstance{
		cli: cli,
	}, nil
}

func (r *ecsInstance) Log(ctx context.Context, name string, opts optypes.LogOptions) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		req := ecs.CreateGetInstanceConsoleOutputRequest()
		req.InstanceId = name
		req.Scheme = schemeHttps

		resp, err := r.cli.GetInstanceConsoleOutput(req)
		if err != nil {
			return fmt.Errorf("error get console output for %s: %w", name, err)
		}

		if resp.LastUpdateTime != "" {
			if r.lastUpdatedTimestamp != "" && r.lastUpdatedTimestamp == resp.LastUpdateTime {
				continue
			}
			r.lastUpdatedTimestamp = resp.LastUpdateTime
		}

		if resp.ConsoleOutput != "" {
			var content string

			// The console output is base64 encoded.
			content, err := strs.DecodeBase64(resp.ConsoleOutput)
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
