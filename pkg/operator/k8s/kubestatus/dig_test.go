package kubestatus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

func Test_digPodErrorReason(t *testing.T) {
	testCases := []struct {
		name     string
		given    core.PodStatus
		expected string
	}{
		{
			name: "init container exit without 0",
			given: core.PodStatus{
				InitContainerStatuses: []core.ContainerStatus{
					{
						Name: "init-volume",
						State: core.ContainerState{
							Terminated: &core.ContainerStateTerminated{
								ExitCode: 0,
							},
						},
					},
					{
						Name: "init-data",
						State: core.ContainerState{
							Terminated: &core.ContainerStateTerminated{
								ExitCode: 1,
							},
						},
					},
				},
			},
			expected: `Init Container "init-data": exit code: 1`,
		},
		{
			name: "init container terminate by signal",
			given: core.PodStatus{
				InitContainerStatuses: []core.ContainerStatus{
					{
						Name: "init-data",
						State: core.ContainerState{
							Terminated: &core.ContainerStateTerminated{
								ExitCode: 1,
								Signal:   9,
							},
						},
					},
				},
			},
			expected: `Init Container "init-data": signal: 9`,
		},
		{
			name: "init container terminate by signal with reason",
			given: core.PodStatus{
				InitContainerStatuses: []core.ContainerStatus{
					{
						Name: "init-data",
						State: core.ContainerState{
							Terminated: &core.ContainerStateTerminated{
								ExitCode: 1,
								Signal:   9,
								Message:  "exhaust all inodes",
								Reason:   "Exhausted",
							},
						},
					},
				},
			},
			expected: `Init Container "init-data": Exhausted, exhaust all inodes`,
		},
		{
			name: "init container waiting",
			given: core.PodStatus{
				InitContainerStatuses: []core.ContainerStatus{
					{
						Name: "init-volume",
						State: core.ContainerState{
							Waiting: &core.ContainerStateWaiting{
								Message: `Back-off pulling image "volume-path-parser:apline"`,
								Reason:  "ImagePullBackOff",
							},
						},
					},
					{
						Name: "init-data",
						State: core.ContainerState{
							Terminated: &core.ContainerStateTerminated{
								ExitCode: 1,
							},
						},
					},
				},
			},
			expected: `Init Container "init-volume": ImagePullBackOff, Back-off pulling image "volume-path-parser:apline"`,
		},
		{
			name: "container terminate by signal",
			given: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name: "wordpress",
						State: core.ContainerState{
							Terminated: &core.ContainerStateTerminated{
								ExitCode: 1,
							},
						},
					},
				},
			},
			expected: `Container "wordpress": exit code: 1`,
		},
		{
			name: "container waiting",
			given: core.PodStatus{
				ContainerStatuses: []core.ContainerStatus{
					{
						Name: "nginx",
						State: core.ContainerState{
							Waiting: &core.ContainerStateWaiting{
								Message: `Back-off pulling image "nginx:apline"`,
								Reason:  "ImagePullBackOff",
							},
						},
					},
				},
			},
			expected: `Container "nginx": ImagePullBackOff, Back-off pulling image "nginx:apline"`,
		},
	}
	for _, tc := range testCases {
		// Convert object to map.
		var given map[string]any
		bs, _ := json.Marshal(tc.given)
		_ = json.Unmarshal(bs, &given)

		t.Run(tc.name, func(t *testing.T) {
			actual := digPodErrorReason(given)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
