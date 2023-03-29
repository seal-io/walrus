package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/seal-io/seal/pkg/i18n/text"
)

//go:generate gotext -srclang=en update -out=catalog.go -lang=en,zh

// GoTextEntry is a temporary workaround for https://github.com/golang/go/issues/58633
// which blocks extracting messages from existing code base.
// Print all i18n text that needs gotext extraction.
func GoTextEntry() {
	p := message.NewPrinter(language.Chinese)

	p.Printf(text.ExampleKubernetesName)
	p.Printf(text.ExampleKubernetesPrompt)
	p.Printf(text.ExampleAlibabaCloudName)
	p.Printf(text.ExampleAlibabaCloudPrompt)
	p.Printf(text.ExampleELKName)
	p.Printf(text.ExampleELKPrompt)
	p.Printf(text.TerraformModuleGenerateSystemMessage)
	p.Printf(text.TerraformModuleExplainSystemMessage)
	p.Printf(text.TerraformModuleCorrectSystemMessageDesc)
	p.Printf(text.TerraformModuleCorrectSystemMessageCorrectedDesc)
	p.Printf(text.TerraformModuleCorrectSystemMessageExplanationDesc)
}
