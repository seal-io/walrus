package netx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPv4_Overlap(t *testing.T) {
	type input struct {
		l, r IPv4
	}

	testCases := []struct {
		given    input
		expected bool
	}{
		{
			given: input{
				l: MustIPv4FromCIDR("10.24.0.0/16"),
				r: MustIPv4FromCIDR("10.24.0.0/16"),
			},
			expected: true,
		},
		{
			given: input{
				l: MustIPv4FromCIDR("192.168.0.0/16"),
				r: MustIPv4FromCIDR("192.168.0.0/24"),
			},
			expected: true,
		},
		{
			given: input{
				l: MustIPv4FromCIDR("172.16.0.0/17"),
				r: MustIPv4FromCIDR("172.16.1.0/17"),
			},
			expected: true,
		},
		{
			given: input{
				l: MustIPv4FromCIDR("172.16.0.0/17"),
				r: MustIPv4FromCIDR("172.16.128.0/17"),
			},
			expected: false,
		},
		{
			given: input{
				l: MustIPv4FromCIDR("192.168.0.0/16"),
				r: MustIPv4FromCIDR("10.24.0.0/16"),
			},
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.given.l.String()+" "+tc.given.r.String(), func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.given.l.Overlap(tc.given.r))
		})
	}
}

func TestIPv4_Next(t *testing.T) {
	testCases := []struct {
		given    IPv4
		expected IPv4
	}{
		{
			given:    MustIPv4FromCIDR("192.168.0.0/16"),
			expected: MustIPv4FromCIDR("192.169.0.0/16"),
		},
		{
			given:    MustIPv4FromCIDR("192.169.0.0/16"),
			expected: MustIPv4FromCIDR("192.170.0.0/16"),
		},
		{
			given:    MustIPv4FromCIDR("192.255.0.0/16"),
			expected: MustIPv4FromCIDR("193.0.0.0/16"),
		},
		{
			given:    MustIPv4FromCIDR("172.16.0.0/17"),
			expected: MustIPv4FromCIDR("172.16.128.0/17"),
		},
		{
			given:    MustIPv4FromCIDR("172.16.0.0/18"),
			expected: MustIPv4FromCIDR("172.16.64.0/18"),
		},
		{
			given:    MustIPv4FromCIDR("172.16.64.0/18"),
			expected: MustIPv4FromCIDR("172.16.128.0/18"),
		},
		{
			given:    MustIPv4FromCIDR("172.16.128.0/18"),
			expected: MustIPv4FromCIDR("172.16.192.0/18"),
		},
		{
			given:    MustIPv4FromCIDR("172.16.192.0/18"),
			expected: MustIPv4FromCIDR("172.17.0.0/18"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.given.String(), func(t *testing.T) {
			assert.Equal(t, tc.expected.String(), tc.given.Next().String())
		})
	}
}
