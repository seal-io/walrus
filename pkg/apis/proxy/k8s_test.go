package proxy

import "testing"

func TestGetProxyUrl(t *testing.T) {
	type (
		input struct {
			host  string
			path  string
			query string
		}
	)
	testCases := []struct {
		given    input
		expected string
	}{
		{
			given: input{
				host:  "https://127.0.0.1",
				path:  "api/v1",
				query: "foo=bar",
			},
			expected: "https://127.0.0.1/api/v1?foo=bar",
		},
		{
			given: input{
				host: "https://myhost.com",
				path: "api/v1/namespaces",
			},
			expected: "https://myhost.com/api/v1/namespaces",
		},
		{
			given: input{
				host: "https://myhost.com",
				path: "api/v1/secrets",
			},
			expected: "https://myhost.com/api/v1/secrets",
		},
		{
			given: input{
				host: "https://myhost.com",
				path: "api/v1/configmaps",
			},
			expected: "https://myhost.com/api/v1/configmaps",
		},
	}

	for _, tc := range testCases {
		actual, err := getProxyURL(tc.given.host, tc.given.path, tc.given.query)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if actual != tc.expected {
			t.Errorf("Expected %v, got %v", tc.expected, actual)
		}
	}
}
