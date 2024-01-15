package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

func (r *Server) initServeURLIfDockerMode(ctx context.Context, opts initOptions) error {
	localEnvironmentMode, err := settings.LocalEnvironmentMode.Value(ctx, opts.ModelClient)
	if err != nil {
		return err
	}

	if localEnvironmentMode != localEnvironmentModeDocker {
		return nil
	}

	// Get the first non-loopback IPv4 address and sets the ServeUrl setting.
	interfaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to get network interfaces: %w", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				log.Warnf("failed to get addresses of interface %s: %v", iface.Name, err)
				continue
			}

			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err == nil && ip.To4() != nil {
					return settings.ServeUrl.Set(ctx, opts.ModelClient, fmt.Sprintf("https://%s", ip))
				}
			}
		}
	}

	return errors.New("failed to get ip address of docker network")
}
