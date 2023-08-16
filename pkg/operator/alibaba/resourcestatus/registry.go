package resourcestatus

import (
	"fmt"

	"github.com/seal-io/walrus/pkg/dao/types/status"
	"github.com/seal-io/walrus/pkg/operator/types"
)

// resourceTypes indicate supported resource type and function to get status.
var resourceTypes map[string]getStatusFunc

// getStatusFunc is function use resource id to get resource status.
type getStatusFunc func(cred types.Credential, typeName, name string) (*status.Status, error)

func init() {
	resourceTypes = map[string]getStatusFunc{
		"alicloud_cdn_domain":            getCdnDomain,
		"alicloud_cs_kubernetes":         getCsKubernetes,
		"alicloud_disk":                  getEcsDisk,
		"alicloud_image":                 getEcsImage,
		"alicloud_instance":              getEcsInstance,
		"alicloud_network_interface":     getEcsNetworkInterface,
		"alicloud_ecs_network_interface": getEcsNetworkInterface,
		"alicloud_snapshot":              getEcsSnapshot,
		"alicloud_ecs_snapshot":          getEcsSnapshot,
		"alicloud_polardb_cluster":       getPolarDBCluster,
		"alicloud_db_instance":           getRdsDBInstance,
		"alicloud_slb":                   getSlbLoadBalancer,
		"alicloud_vpc":                   getVpc,
		"alicloud_eip":                   getVpcEip,
		"alicloud_vswitch":               getVpcVSwitch,
	}
}

// IsSupported indicate whether the resource type is supported.
func IsSupported(typeName string) bool {
	_, ok := resourceTypes[typeName]
	return ok
}

// Get resource status by resource type and name.
func Get(cred types.Credential, typeName, name string) (*status.Status, error) {
	getFunc, exist := resourceTypes[typeName]
	if !exist {
		return nil, fmt.Errorf("unsupported resource type: %s", typeName)
	}

	st, err := getFunc(cred, typeName, name)
	if err != nil {
		return &status.Status{}, err
	}

	return st, nil
}
