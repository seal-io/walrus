package polymorphic

import (
	"fmt"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/seal-io/walrus/utils/json"
)

// SelectorsForObject returns the Pod label selector for a given object.
func SelectorsForObject(obj *unstructured.Unstructured) (ns string, s labels.Selector, err error) {
	switch obj.GetKind() {
	default:
		return "", nil, fmt.Errorf("selector for %s not implemented", obj.GetKind())
	case "ReplicaSet",
		"StatefulSet", "DaemonSet", "Deployment",
		"Job":
		ns = obj.GetNamespace()

		lso, exist, _ := unstructured.NestedFieldNoCopy(obj.Object, "spec", "selector")
		if !exist {
			return "", nil, fmt.Errorf("%s defined without a selector",
				obj.GetKind())
		}

		// Any -> bs -> structure.
		lsob, err := json.Marshal(lso)
		if err != nil {
			return "", nil, fmt.Errorf("failed %s marshall, %w",
				obj.GetKind(), err)
		}

		var ls meta.LabelSelector
		if err = json.Unmarshal(lsob, &ls); err != nil {
			return "", nil, fmt.Errorf("failed %s unmarshall, %w",
				obj.GetKind(), err)
		}

		s, err = meta.LabelSelectorAsSelector(&ls)
		if err != nil {
			return "", nil, fmt.Errorf("invalid label selector, %w", err)
		}
	case "ReplicationController", "Service":
		ns = obj.GetNamespace()

		ss, exist, _ := unstructured.NestedStringMap(obj.Object, "spec", "selector")
		if !exist {
			return "", nil, fmt.Errorf("%s defined without a selector",
				obj.GetKind())
		}
		s = labels.SelectorFromSet(ss)
	}

	return ns, s, nil
}
