package apis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseBindAddress(t *testing.T) {
	type (
		input struct {
			ip   string
			port int
			dual bool
		}
		output struct {
			network string
			address string
			err     error
		}
	)

	testCases := []struct {
		name     string
		given    input
		expected output
	}{
		{
			name: "loopback",
			given: input{
				ip:   "127.0.0.1",
				port: 80,
				dual: true,
			},
			expected: output{
				network: "tcp",
				address: "127.0.0.1:80",
			},
		},
		{
			name: "unspecified without dual",
			given: input{
				ip:   "::",
				port: 443,
				dual: false,
			},
			expected: output{
				network: "tcp6",
				address: "[::]:443",
			},
		},
		{
			name: "ipv4",
			given: input{
				ip:   "192.168.50.23",
				port: 80,
				dual: true,
			},
			expected: output{
				network: "tcp",
				address: "192.168.50.23:80",
			},
		},
		{
			name: "ipv4 without dual",
			given: input{
				ip:   "192.168.50.23",
				port: 80,
				dual: false,
			},
			expected: output{
				network: "tcp4",
				address: "192.168.50.23:80",
			},
		},
		{
			name: "ipv6",
			given: input{
				ip:   "240e:3b3:30b0:51d0:c45d:65be:e7eb:f611",
				port: 443,
				dual: true,
			},
			expected: output{
				network: "tcp",
				address: "[240e:3b3:30b0:51d0:c45d:65be:e7eb:f611]:443",
			},
		},
		{
			name: "ipv6 without dual",
			given: input{
				ip:   "240e:3b3:30b0:51d0:c45d:65be:e7eb:f611",
				port: 443,
				dual: false,
			},
			expected: output{
				network: "tcp6",
				address: "[240e:3b3:30b0:51d0:c45d:65be:e7eb:f611]:443",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actual output
			actual.network, actual.address, actual.err = parseBindAddress(tc.given.ip, tc.given.port, tc.given.dual)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
