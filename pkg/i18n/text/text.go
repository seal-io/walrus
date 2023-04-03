package text

// Note: The internationalization may not be simply translations because the English prompt may work better
// even with user messages in other languages.

const (
	ExampleKubernetesName   = "Create a Kubernetes deployment"
	ExampleKubernetesPrompt = "# Create a Kubernetes deployment. Provide common variables."

	ExampleAlibabaCloudName   = "Create an alibaba cloud virtual machine"
	ExampleAlibabaCloudPrompt = "# Create a resource group, virtual network, subnet and virtual machine on alibaba cloud."

	ExampleELKName   = "Deploy an ELK stack"
	ExampleELKPrompt = "# Deploy an ELK stack using helm chart."

	TerraformModuleGenerateSystemMessage = "You are translating natural language to a Terraform module." +
		" Please do not explain, just write pure terraform HCL code." +
		" Please do not explain, just write pure terraform HCL code." +
		" Please do not explain, just write pure terraform HCL code."

	TerraformModuleExplainSystemMessage = "Please explain the given terraform module."

	TerraformModuleCorrectSystemMessageDesc = "Please Check and fix the given terraform module if there's any mistake.\n" +
		"Strictly respond a valid JSON in the following format:"
	TerraformModuleCorrectSystemMessageCorrectedDesc   = "Terraform code that is fixed. Please do not explain, just write terraform HCL code."
	TerraformModuleCorrectSystemMessageExplanationDesc = "Explanation of the fixes."
)
