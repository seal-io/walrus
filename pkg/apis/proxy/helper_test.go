package proxy

import "testing"

func TestIsWhiteListDomain(t *testing.T) {
	type (
		input struct {
			host       string
			whiteLists []string
		}
	)

	testCases := []struct {
		given    input
		expected bool
	}{
		{
			given: input{
				host:       "hub.docker.com",
				whiteLists: []string{"hub.docker.com"},
			},
			expected: true,
		},
		{
			given: input{
				host:       "test.com",
				whiteLists: []string{"hub.docker.com"},
			},
			expected: false,
		},
		{
			given: input{
				host: "ec2.us-west-1.amazonaws.com",
				whiteLists: []string{
					"*.amazonaws.com",
					"*.amazonaws.com.cn",
				},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		if actual := isWhiteListDomain(tc.given.host, tc.given.whiteLists...); actual != tc.expected {
			t.Errorf("Expected %v, got %v", tc.expected, actual)
		}
	}
}
