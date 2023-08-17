package openai

import (
	"context"
	"net/http"

	"github.com/sashabaranov/go-openai"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/errorx"
)

const gpt35MaxTokens = 4096

func CreateCompletion(ctx context.Context, mc model.ClientSet, systemMessage, userMessage string) (string, error) {
	apiToken, err := settings.OpenAiApiToken.Value(ctx, mc)
	if err != nil {
		return "", err
	}

	if apiToken == "" {
		return "", errorx.NewHttpError(http.StatusBadRequest,
			"invalid input: OpenAI API token is not configured")
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
		return "", errorx.Errorf("failed to create completion: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}
