package kuberest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seal-io/utils/json"
	"github.com/seal-io/utils/waitx"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// WaitUntilAvailable waits until the Kubernetes to be available,
// or returns an error if the context is canceled.
func WaitUntilAvailable(ctx context.Context, cli rest.Interface) error {
	if cli == nil {
		return errors.New("rest client is invalid")
	}

	return waitx.PollUntilContextCancel(ctx, time.Second, true,
		func(ctx context.Context) error {
			return IsAvailable(ctx, cli)
		},
	)
}

// WaitConfigUntilAvailable is similar to WaitUntilAvailable,
// accepts rest.Config as input.
func WaitConfigUntilAvailable(ctx context.Context, cfg *rest.Config) error {
	if cfg == nil {
		return errors.New("rest config is invalid")
	}

	cli, err := rest.UnversionedRESTClientFor(dynamic.ConfigFor(cfg))
	if err != nil {
		return fmt.Errorf("create rest client to check kubernetes cluster: %w", err)
	}

	return WaitUntilAvailable(ctx, cli)
}

// IsAvailable returns true if the Kubernetes cluster is available.
func IsAvailable(ctx context.Context, cli rest.Interface) error {
	if cli == nil {
		return errors.New("rest client is invalid")
	}

	body, err := cli.Get().
		AbsPath("/version").
		Do(ctx).
		Raw()
	if err != nil {
		return err
	}

	var info struct {
		Major    string `json:"major"`
		Minor    string `json:"minor"`
		Compiler string `json:"compiler"`
		Platform string `json:"platform"`
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return fmt.Errorf("unable to parse the server version: %w", err)
	}

	return nil
}
