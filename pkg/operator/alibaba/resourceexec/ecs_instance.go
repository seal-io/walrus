package resourceexec

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	assistlog "github.com/aliyun/aliyun_assist_client/agent/log"
	assistclient "github.com/aliyun/aliyun_assist_client/agent/session/plugin"
	"github.com/aliyun/aliyun_assist_client/agent/session/plugin/message"
	"github.com/sirupsen/logrus"

	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/utils/bytespool"
	"github.com/seal-io/walrus/utils/gopool"
	"github.com/seal-io/walrus/utils/log"
)

const (
	defaultTerminalHeight uint16 = 100
	defaultTerminalWidth  uint16 = 100
)

const (
	// These configs are borrowed from aliyun repo
	// nolint: lll
	// https://github.com/aliyun/aliyun_assist_client/blob/d2b430f3fa8aea1abb376d0e4a08d9888729f1c4/agent/session/plugin/client.go#L34
	// maxPackageSend in Byte.
	maxPackageSendSize = 2048
	// DefaultSendSpeed send speed in kbps, this is from aliyun repo.
	defaultPackageSendSpeed = 200
	// DefaultSendInterval in ms, this is calculated from defaultSendSpeed and maxPackageSend.
	defaultPackageSendInterval = 1000 / (defaultPackageSendSpeed * 1024 / 8 / maxPackageSendSize)

	defaultPackageSendIntervalDuration = defaultPackageSendInterval * time.Millisecond
)

func init() {
	logger := logrus.New()
	logger.SetOutput(log.GetLogger())
	assistlog.Log = logger
}

type ecsInstance struct {
	cred          *optypes.Credential
	ecsCli        *ecs.Client
	assistCli     *assistclient.Client
	realConnected bool
}

func getEcsInstance(ctx context.Context) (optypes.ExecutableResource, error) {
	var (
		res = &ecsInstance{}
		err error
	)

	res.cred, err = optypes.CredentialFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	res.ecsCli, err = ecs.NewClientWithAccessKey(res.cred.Region, res.cred.AccessKey, res.cred.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("error create alibaba ecs client %s: %w", res.cred.AccessKey, err)
	}

	return res, err
}

// Supported check whether ecs instance support session manager.
func (r *ecsInstance) Supported(_ context.Context, name string) (bool, error) {
	req := ecs.CreateDescribeCloudAssistantStatusRequest()
	req.Scheme = schemeHttps
	req.InstanceId = &[]string{name}

	resp, err := r.ecsCli.DescribeCloudAssistantStatus(req)
	if err != nil {
		return false, err
	}

	icas := resp.InstanceCloudAssistantStatusSet.InstanceCloudAssistantStatus
	if len(icas) == 0 || !icas[0].SupportSessionManager {
		return false, nil
	}

	return true, nil
}

// Exec support data channel to input and output command.
func (r *ecsInstance) Exec(ctx context.Context, name string, opts optypes.ExecOptions) error {
	wsURL, err := r.getConnectAddress(name)
	if err != nil {
		return err
	}

	readCloser := io.NopCloser(opts.In)

	r.assistCli, err = assistclient.NewClient(
		wsURL, readCloser, opts.Out, false, "", true, true,
	)
	if err != nil {
		return err
	}

	defer func() {
		if !r.realConnected {
			return
		}

		if err := r.assistCli.SendCloseMessage(); err != nil {
			log.Warnf("error send close message: %v", err)
		}
	}()

	ctxWithCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to websocket URL.
	if !r.assistCli.Connected {
		if err = r.assistCli.Connect(); err != nil {
			return fmt.Errorf("error connect: %w", err)
		}

		r.assistCli.Conn.SetCloseHandler(func(int, string) (err error) {
			cancel()
			return
		})
	}

	eg := gopool.GroupWithContextIn(ctxWithCancel)
	eg.Go(func(ctx context.Context) error {
		return r.setTerminalSize(ctx, opts)
	})

	eg.Go(func(ctx context.Context) error {
		return r.writeToConn(ctx)
	})

	eg.Go(func(ctx context.Context) error {
		return r.readFromConn(ctx)
	})

	return eg.Wait()
}

func (r *ecsInstance) setTerminalSize(ctx context.Context, opts optypes.ExecOptions) error {
	set := func(width, height uint16) error {
		// Send resize data to remote connection.
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, int16(height))
		_ = binary.Write(buf, binary.LittleEndian, int16(width))

		err := r.assistCli.SendResizeDataMessage(buf.Bytes())
		if err != nil {
			return fmt.Errorf("error send resize data: %w", err)
		}

		return nil
	}

	// Without resizer.
	if opts.Resizer == nil {
		return set(defaultTerminalWidth, defaultTerminalHeight)
	}

	// With resizer.
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if !r.realConnected {
			continue
		}

		width, height, ok := opts.Resizer.Next()
		if !ok {
			return errors.New("invalid terminal resizer")
		}

		err := set(width, height)
		if err != nil {
			return err
		}
	}
}

func (r *ecsInstance) getConnectAddress(name string) (string, error) {
	req := ecs.CreateStartTerminalSessionRequest()
	req.InstanceId = &[]string{name}

	resp, err := r.ecsCli.StartTerminalSession(req)
	if err != nil {
		return "", fmt.Errorf("error start session for %s: %w", name, err)
	}

	return resp.WebSocketUrl, nil
}

func (r *ecsInstance) writeToConn(ctx context.Context) error {
	buff := bytespool.GetBytes(maxPackageSendSize)
	defer func() { bytespool.Put(buff) }()

	for {
		// Watch done event.
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// Control send speed.
		time.Sleep(defaultPackageSendIntervalDuration)

		// Read from opts in.
		size, err := r.assistCli.Input.Read(buff)
		if err != nil {
			return fmt.Errorf("error read from user input: %w", err)
		}

		// Write to connection.
		if r.realConnected {
			err = r.assistCli.SendStreamDataMessage(buff[:size])
			if err != nil {
				return fmt.Errorf("error send data: %w", err)
			}
		}
	}
}

func (r *ecsInstance) readFromConn(ctx context.Context) error {
	for {
		// Watch done event.
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_, data, err := r.assistCli.Conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("error read message: %w", err)
		}

		msg := message.Message{}

		err = msg.Deserialize(data)
		if err != nil {
			return fmt.Errorf("error deserialize message: %w", err)
		}

		err = msg.Validate()
		if err != nil {
			return fmt.Errorf("error validate message: %w", err)
		}

		switch msg.MessageType {
		case message.OutputStreamDataMessage:
			r.realConnected = true

			_, err = r.assistCli.Output.Write(msg.Payload)
			if err != nil {
				return fmt.Errorf("error write message to output")
			}
		case message.StatusDataChannel:
			err = r.assistCli.ProcessStatusDataChannel(msg.Payload)
			if err != nil {
				return fmt.Errorf("error process status message: %w", err)
			}
		}
	}
}
