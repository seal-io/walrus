package types

const (
	ConnectorTypeKubernetes string = "Kubernetes"
	ConnectorTypeAlibaba    string = "Alibaba"
	ConnectorTypeAWS        string = "AWS"
)

const (
	ConnectorCategoryKubernetes     string = "Kubernetes"
	ConnectorCategoryCustom         string = "Custom"
	ConnectorCategoryVersionControl string = "VersionControl"
	ConnectorCategoryCloudProvider  string = "CloudProvider"
)

// FinOpsCustomPricing used to config opencost.
type FinOpsCustomPricing struct {
	// CPU describing cost per core-month of CPU.
	CPU string `json:"cpu"`
	// CPU describing cost per core-month of CPU for spot nodes.
	SpotCPU string `json:"spotCPU"`
	// RAM describing cost per GiB-month of RAM/memory.
	RAM string `json:"ram"`
	// SpotRAM describing cost per GiB-month of RAM/memory for spot nodes.
	SpotRAM string `json:"spotRAM"`
	GPU     string `json:"gpu"`
	SpotGPU string `json:"spotGPU"`
	// Storage describing cost per GB-month of storage (e.g. PV, disk) resources.
	Storage string `json:"storage"`
}

func (c *FinOpsCustomPricing) IsZero() bool {
	return c == nil ||
		c.CPU == "" &&
			c.SpotCPU == "" &&
			c.RAM == "" &&
			c.SpotRAM == "" &&
			c.GPU == "" &&
			c.SpotGPU == "" &&
			c.Storage == ""
}

func DefaultFinOpsCustomPricing() *FinOpsCustomPricing {
	// Opencost will treat custom price from configMap as month cost, and divide / 730 to calculate hourly cost,
	// the default pricing from opencost is hourly cost, so we multiply by 730 and place here.
	//nolint: lll
	// https://github.com/opencost/opencost/blob/d7958021bff300610b4585de0fac4e7289d3b5b9/pkg/cloud/providerconfig.go#L214
	return &FinOpsCustomPricing{
		CPU:     "23.07603",
		SpotCPU: "4.85815",
		RAM:     "3.09301",
		SpotRAM: "0.65116",
		GPU:     "693.5",
		SpotGPU: "224.84",
		Storage: "0.0399999996",
	}
}
