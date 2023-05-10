package kubestatus

import (
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// The following codes inspired by
// https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go.

func digPodErrorReason(status map[string]interface{}) (r string) {
	initContainerStatuses, _, _ := unstructured.NestedSlice(status, "initContainerStatuses")
	for i := range initContainerStatuses {
		initContainerStatus, ok := initContainerStatuses[i].(map[string]interface{})
		if !ok {
			continue
		}

		name, exist, _ := unstructured.NestedString(initContainerStatus, "name")
		if !exist {
			name = strconv.Itoa(i)
		}

		terminated, exist, _ := unstructured.NestedMap(initContainerStatus, "state", "terminated")
		if exist {
			exitCode, _, _ := unstructured.NestedInt64(terminated, "exitCode")
			if exitCode == 0 {
				continue
			}

			reason, _, _ := unstructured.NestedString(terminated, "reason")
			if reason == "" {
				signal, _, _ := unstructured.NestedInt64(terminated, "signal")
				if signal == 0 {
					return fmt.Sprintf("Init Container %q: exit code: %d", name, exitCode)
				}
				return fmt.Sprintf("Init Container %q: signal: %d", name, signal)
			}
			message, _, _ := unstructured.NestedString(terminated, "message")
			if message == "" {
				return fmt.Sprintf("Init Container %q: %s", name, reason)
			}
			return fmt.Sprintf("Init Container %q: %s, %s", name, reason, message)
		}

		waiting, exist, _ := unstructured.NestedMap(initContainerStatus, "state", "waiting")
		if exist {
			reason, _, _ := unstructured.NestedString(waiting, "reason")
			if reason == "" {
				return fmt.Sprintf("Init Container %q: Failed", name)
			}
			message, _, _ := unstructured.NestedString(waiting, "message")
			if message == "" {
				return fmt.Sprintf("Init Container %q: %s", name, reason)
			}
			return fmt.Sprintf("Init Container %q: %s, %s", name, reason, message)
		}
	}

	containerStatuses, _, _ := unstructured.NestedSlice(status, "containerStatuses")
	for i := len(containerStatuses) - 1; i >= 0; i-- {
		containerStatus, ok := containerStatuses[i].(map[string]interface{})
		if !ok {
			continue
		}

		name, exist, _ := unstructured.NestedString(containerStatus, "name")
		if !exist {
			name = strconv.Itoa(i)
		}

		waiting, exist, _ := unstructured.NestedMap(containerStatus, "state", "waiting")
		if exist {
			reason, _, _ := unstructured.NestedString(waiting, "reason")
			if reason != "" {
				message, _, _ := unstructured.NestedString(waiting, "message")
				if message == "" {
					return fmt.Sprintf("Container %q: %s", name, reason)
				}
				return fmt.Sprintf("Container %q: %s, %s", name, reason, message)
			}
		}

		terminated, exist, _ := unstructured.NestedMap(containerStatus, "state", "terminated")
		if exist {
			reason, _, _ := unstructured.NestedString(terminated, "reason")
			if reason != "" {
				message, _, _ := unstructured.NestedString(terminated, "message")
				if message == "" {
					return fmt.Sprintf("Container %q: %s", name, reason)
				}
				return fmt.Sprintf("Container %q: %s, %s", name, reason, message)
			}

			signal, _, _ := unstructured.NestedInt64(terminated, "signal")
			if signal != 0 {
				return fmt.Sprintf("Container %q: signal: %d", name, signal)
			}

			exitCode, _, _ := unstructured.NestedInt64(terminated, "exitCode")
			return fmt.Sprintf("Container %q: exit code: %d", name, exitCode)
		}
	}

	return
}
