package i18n

import (
	"io"

	"golang.org/x/text/message"
)

// Note: The internationalization may not be simply translations because the English prompt may work better
// even with user messages in other languages.

const (
	ExampleKubernetesName   = "Create a Kubernetes deployment"
	ExampleKubernetesPrompt = "# Create a Kubernetes deployment. Provide common variables."

	ExampleAlibabaCloudName   = "Create an alibaba cloud virtual machine"
	ExampleAlibabaCloudPrompt = "# Create a resource group, virtual network," +
		" subnet and virtual machine on alibaba cloud."

	ExampleELKName   = "Deploy an ELK stack"
	ExampleELKPrompt = "# Deploy an ELK stack using helm chart."

	TerraformModuleGenerateSystemMessage = "You are translating natural language to a Terraform module." +
		" Please do not explain, just write pure terraform HCL code." +
		" Please do not explain, just write pure terraform HCL code." +
		" Please do not explain, just write pure terraform HCL code."

	TerraformModuleExplainSystemMessage = "Please explain the given terraform module."

	TerraformModuleCorrectSystemMessageDesc = "Please Check and fix the given terraform module" +
		" if there's any mistake.\n" +
		"Strictly respond a valid JSON in the following format:"
	TerraformModuleCorrectSystemMessageCorrectedDesc = "Terraform code that is fixed." +
		" Please do not explain, just write terraform HCL code."
	TerraformModuleCorrectSystemMessageExplanationDesc = "Explanation of the fixes."
)

// Fprintf is a helper function to print all the messages to the given writer.
//
// Fprintf also displays all i18n text for gotext extraction.
func Fprintf(p *message.Printer, w io.Writer) {
	_, _ = p.Fprintf(w, ExampleKubernetesName)
	_, _ = p.Fprintf(w, ExampleKubernetesPrompt)
	_, _ = p.Fprintf(w, ExampleAlibabaCloudName)
	_, _ = p.Fprintf(w, ExampleAlibabaCloudPrompt)
	_, _ = p.Fprintf(w, ExampleELKName)
	_, _ = p.Fprintf(w, ExampleELKPrompt)
	_, _ = p.Fprintf(w, TerraformModuleGenerateSystemMessage)
	_, _ = p.Fprintf(w, TerraformModuleExplainSystemMessage)
	_, _ = p.Fprintf(w, TerraformModuleCorrectSystemMessageDesc)
	_, _ = p.Fprintf(w, TerraformModuleCorrectSystemMessageCorrectedDesc)
	_, _ = p.Fprintf(w, TerraformModuleCorrectSystemMessageExplanationDesc)
}
