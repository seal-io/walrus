package pods

import (
	core "k8s.io/api/core/v1"

	"github.com/seal-io/seal/pkg/platformk8s/key"
)

// IsPodReady returns true if Pod is ready.
func IsPodReady(pod *core.Pod) bool {
	if !IsPodRunning(pod) {
		return false
	}
	var c, exist = GetPodCondition(&pod.Status, core.PodReady)
	if exist {
		return c.Status == core.ConditionTrue
	}
	return false
}

// IsPodRunning returns ture if Pod is running.
func IsPodRunning(pod *core.Pod) bool {
	if !IsPodAssigned(pod) {
		return false
	}
	return pod.Status.Phase == core.PodRunning
}

// IsPodAssigned returns true if Pod is assigned.
func IsPodAssigned(pod *core.Pod) bool {
	if pod == nil {
		return false
	}
	return pod.Spec.NodeName != ""
}

// GetPodCondition extracts the provided condition from the given PodStatus and returns that.
func GetPodCondition(status *core.PodStatus, conditionType core.PodConditionType) (c *core.PodCondition, exist bool) {
	if status == nil {
		return
	}
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return &status.Conditions[i], true
		}
	}
	return
}

// ContainerType indicates the type of the Container,
// includes Run, Init, Ephemeral.
type ContainerType = string

const (
	ContainerRun       = "run"
	ContainerInit      = "init"
	ContainerEphemeral = "ephemeral"
)

// Container holds container type and name.
type Container struct {
	Type ContainerType
	Name string
}

// IsContainerRunning returns true if Container is running.
func IsContainerRunning(pod *core.Pod, c Container) bool {
	if pod == nil || c.Name == "" {
		return false
	}

	var css = make([]*[]core.ContainerStatus, 0, 3)
	switch c.Type {
	case ContainerRun:
		css = append(css, &pod.Status.ContainerStatuses)
	case ContainerInit:
		css = append(css, &pod.Status.InitContainerStatuses)
	case ContainerEphemeral:
		css = append(css, &pod.Status.EphemeralContainerStatuses)
	default:
		css = append(css,
			&pod.Status.ContainerStatuses,
			&pod.Status.InitContainerStatuses,
			&pod.Status.EphemeralContainerStatuses,
		)
	}

	for i := 0; i < len(css); i++ {
		var cs = *css[i]
		for j := 0; j < len(cs); j++ {
			if cs[j].Name != c.Name {
				continue
			}
			return cs[j].State.Running != nil
		}
	}
	return false
}

// IsContainerExisted returns true if Container is existed.
func IsContainerExisted(pod *core.Pod, c Container) bool {
	if pod == nil || c.Name == "" {
		return false
	}

	switch c.Type {
	case ContainerRun:
		for i := range pod.Spec.Containers {
			if pod.Spec.Containers[i].Name == c.Name {
				return true
			}
		}
	case ContainerInit:
		for i := range pod.Spec.InitContainers {
			if pod.Spec.InitContainers[i].Name == c.Name {
				return true
			}
		}
	case ContainerEphemeral:
		for i := range pod.Spec.EphemeralContainers {
			if pod.Spec.EphemeralContainers[i].Name == c.Name {
				return true
			}
		}
	default:
		for i := range pod.Spec.Containers {
			if pod.Spec.Containers[i].Name == c.Name {
				return true
			}
		}
		for i := range pod.Spec.InitContainers {
			if pod.Spec.InitContainers[i].Name == c.Name {
				return true
			}
		}
		for i := range pod.Spec.EphemeralContainers {
			if pod.Spec.EphemeralContainers[i].Name == c.Name {
				return true
			}
		}
	}
	return false
}

// ContainerStateType indicates the state type of the Container,
// includes Unknown, Waiting, Running, Terminated.
type ContainerStateType uint8

const (
	ContainerStateUnknown ContainerStateType = iota
	ContainerStateWaiting
	ContainerStateRunning
	ContainerStateTerminated
)

func GetContainerStateType(s core.ContainerState) ContainerStateType {
	switch {
	case s.Waiting != nil:
		return ContainerStateWaiting
	case s.Running != nil:
		return ContainerStateRunning
	case s.Terminated != nil:
		return ContainerStateTerminated
	}
	return ContainerStateUnknown
}

type ContainerState struct {
	Type      ContainerType
	Namespace string
	Pod       string
	ID        string
	Name      string
	State     ContainerStateType
}

func (c ContainerState) String() string {
	return key.Encode(c.Namespace, c.Pod, c.Type, c.Name)
}

// GetContainerStates returns ContainerState list of the given Pod.
func GetContainerStates(pod *core.Pod) (r []ContainerState) {
	if pod == nil {
		return
	}

	var css = []struct {
		Type     ContainerType
		Statuses *[]core.ContainerStatus
	}{
		{
			Type:     ContainerInit,
			Statuses: &pod.Status.InitContainerStatuses,
		},
		{
			Type:     ContainerRun,
			Statuses: &pod.Status.ContainerStatuses,
		},
		{
			Type:     ContainerEphemeral,
			Statuses: &pod.Status.EphemeralContainerStatuses,
		},
	}
	for i := 0; i < len(css); i++ {
		var cs = *css[i].Statuses
		for j := 0; j < len(cs); j++ {
			var s = &cs[j]
			r = append(r, ContainerState{
				Type:      css[i].Type,
				Namespace: pod.Namespace,
				Pod:       pod.Name,
				Name:      s.Name,
				ID:        s.ContainerID,
				State:     GetContainerStateType(s.State),
			})
		}
	}
	return
}
