package modulecompletion

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"

	"github.com/seal-io/seal/pkg/apis/modulecompletion/view"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/json"
)

var examples = []view.ModuleCompletionPromptExample{
	{
		Name:   "Create a Kubernetes deployment",
		Prompt: "# Create a Kubernetes deployment. Provide common variables.",
	},
	{
		Name:   "Create an alibaba cloud virtual machine",
		Prompt: "# Create a resource group, virtual network, subnet and virtual machine on alibaba cloud.",
	},
	{
		Name:   "Deploy an ELK stack",
		Prompt: "# Deploy an ELK stack using helm chart.",
	},
}

const (
	gpt35MaxTokens = 4096

	terraformModuleGenerateSystemMessage = "You are translating natural language to a Terraform module." +
		" Please do not explain, just write terraform code." +
		" Please do not explain, just write terraform code." +
		" Please do not explain, just write terraform code."

	terraformModuleExplainSystemMessage = "Please explain the given terraform module."

	terraformModuleCorrectSystemMessage = "Please Check and fix the given terraform module if there's any mistake.\n" +
		"Output in the following JSON format:\n" +
		`
		{
			"corrected": "The corrected terraform code.",
			"explanation": "Explanation of the fixes."
		}
		`
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "ModuleCompletion"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs

// Extensional APIs

func (h Handler) CollectionRouteExamples(_ *gin.Context, _ view.ExampleRequest) (view.ExampleResponse, error) {
	return examples, nil
}

func (h Handler) CollectionRouteGenerate(ctx *gin.Context, req view.GenerateRequest) (*view.GenerateResponse, error) {
	result, err := h.createCompletion(ctx, terraformModuleGenerateSystemMessage, req.Text)
	if err != nil {
		return nil, err
	}
	return &view.GenerateResponse{
		Text: result,
	}, nil
}

func (h Handler) CollectionRouteExplain(ctx *gin.Context, req view.ExplainRequest) (*view.ExplainResponse, error) {
	result, err := h.createCompletion(ctx, terraformModuleExplainSystemMessage, req.Text)
	if err != nil {
		return nil, err
	}
	return &view.ExplainResponse{
		Text: result,
	}, nil
}

func (h Handler) CollectionRouteCorrect(ctx *gin.Context, req view.CorrectRequest) (*view.CorrectResponse, error) {
	result, err := h.createCompletion(ctx, terraformModuleCorrectSystemMessage, req.Text)
	if err != nil {
		return nil, err
	}
	correctResp := &view.CorrectResponse{}
	if err := json.Unmarshal([]byte(result), correctResp); err != nil {
		logrus.Debugf("correction output from openAI: %v", result)
		return nil, errors.New("failed to parse correction advice")
	}

	return correctResp, nil
}

func (h Handler) createCompletion(ctx *gin.Context, systemMessage, userMessage string) (string, error) {
	apiToken, err := settings.OpenAiApiToken.Value(ctx, h.modelClient)
	if err != nil {
		return "", err
	}

	if apiToken == "" {
		return "", errors.New("openAI API token is not configured in settings")
	}

	client := openai.NewClient(apiToken)

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemMessage,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			// TODO Roughly reserve 1000 for the input for now. Update when a tokenizer golang library is available.
			// The tokens from the input message and the completion message cannot exceed the gpt35MaxTokens.
			MaxTokens: gpt35MaxTokens - 1000,
			// Here's an Empirical value. Tunable.
			Temperature: 0.2,
		})

	if err != nil {
		return "", fmt.Errorf("failed to create completion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}
