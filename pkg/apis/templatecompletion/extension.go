package templatecompletion

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/sashabaranov/go-openai"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/i18n/text"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

const gpt35MaxTokens = 4096

var builtinExamples = []PromptExample{
	{
		Name:   text.ExampleKubernetesName,
		Prompt: text.ExampleKubernetesPrompt,
	},
	{
		Name:   text.ExampleAlibabaCloudName,
		Prompt: text.ExampleAlibabaCloudPrompt,
	},
	{
		Name:   text.ExampleELKName,
		Prompt: text.ExampleELKPrompt,
	},
}

func (h Handler) CollectionRouteGetExamples(
	req CollectionRouteGetExampleRequest,
) (CollectionRouteGetExampleResponse, error) {
	resp := make([]PromptExample, 0, len(builtinExamples))

	for _, e := range builtinExamples {
		resp = append(resp, PromptExample{
			Name:   runtime.Translate(req.Context, e.Name),
			Prompt: runtime.Translate(req.Context, e.Prompt),
		})
	}

	return resp, nil
}

func (h Handler) CollectionRouteGenerate(req CollectionRouteGenerateRequest) (*CollectionRouteGenerateResponse, error) {
	prompt := runtime.Translate(req.Context, text.TerraformModuleGenerateSystemMessage)

	result, err := createCompletion(req.Context, h.modelClient, prompt, req.Text)
	if err != nil {
		return nil, err
	}

	return &CollectionRouteGenerateResponse{
		Text: trimMarkdownCodeBlock(result),
	}, nil
}

func (h Handler) CollectionRouteExplain(req CollectionRouteExplainRequest) (*CollectionRouteExplainResponse, error) {
	prompt := runtime.Translate(req.Context, text.TerraformModuleExplainSystemMessage)

	result, err := createCompletion(req.Context, h.modelClient, prompt, req.Text)
	if err != nil {
		return nil, err
	}

	return &CollectionRouteExplainResponse{
		Text: result,
	}, nil
}

func (h Handler) CollectionRouteCorrect(req CollectionRouteCorrectRequest) (*CollectionRouteCorrectResponse, error) {
	// gotext cannot handle brackets in messages, see https://github.com/golang/go/issues/27849.
	// we need to split the text as a workaround.
	desc := runtime.Translate(req.Context, text.TerraformModuleCorrectSystemMessageDesc)
	correctedDesc := runtime.Translate(req.Context, text.TerraformModuleCorrectSystemMessageCorrectedDesc)
	explanationDesc := runtime.Translate(req.Context, text.TerraformModuleCorrectSystemMessageExplanationDesc)
	prompt := fmt.Sprintf(`%s\n{\n"corrected": "%s", "explanation": "%s"\n}\n`, desc, correctedDesc, explanationDesc)

	result, err := createCompletion(req.Context, h.modelClient, prompt, req.Text)
	if err != nil {
		return nil, err
	}

	var resp CollectionRouteCorrectResponse

	if err = json.Unmarshal([]byte(result), &resp); err != nil {
		log.Debugf("correction message is not in the format requested by the prompt. output:\n%v", result)
		return nil, errors.New("failed to parse correction advice")
	}

	return &resp, nil
}

func (h Handler) CollectionRouteCreatePullRequest(
	req CollectionRouteCreatePrRequest,
) (*CollectionRouteCreatePrResponse, error) {
	moduleName := templates.GetTemplateNameByPath(req.Path)

	moduleFiles, err := templates.GetTerraformTemplateFiles(moduleName, req.Content)
	if err != nil {
		return nil, runtime.Error(http.StatusBadRequest, err)
	}

	conn, err := h.modelClient.Connectors().Get(req.Context, req.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("error getting connector: %w", err)
	}

	if conn.Category != types.ConnectorCategoryVersionControl {
		return nil, runtime.Errorf(http.StatusBadRequest,
			"%q is not a supported version control driver", conn.Type)
	}

	client, err := vcs.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating version control system client: %w", err)
	}

	ref, _, err := client.Git.FindBranch(req.Context, req.Repository, req.Branch)
	if err != nil {
		return nil, runtime.Errorpf(http.StatusBadRequest, err,
			"error indexing branch %s from repository %s",
			req.Branch, req.Repository)
	}

	commitInput := &scm.CommitInput{
		Message: "Template generated from Seal",
		Base:    ref.Sha,
	}
	for name, content := range moduleFiles {
		commitInput.Blobs = append(commitInput.Blobs, scm.Blob{
			Path:    filepath.Join(req.Path, name),
			Mode:    "100644",
			Content: content,
		})
	}

	commit, _, err := client.Git.CreateCommit(req.Context, req.Repository, commitInput)
	if err != nil {
		return nil, runtime.Errorwf(err, "error creating new commit for repository %s",
			req.Repository)
	}

	stagingBranch := fmt.Sprintf("seal/module-" + strs.String(5))
	refInput := &scm.ReferenceInput{
		Name: stagingBranch,
		Sha:  commit.Sha,
	}
	_, err = client.Git.CreateBranch(req.Context, req.Repository, refInput)

	if err != nil {
		return nil, runtime.Errorwf(err, "error creating new branch %s for repository %s",
			refInput.Name, req.Repository)
	}

	// TODO more informative PR body. E.g., let chatGPT generate it.
	prInput := &scm.PullRequestInput{
		Title:  fmt.Sprintf("Add module %s", moduleName),
		Body:   "This is a module proposed from Seal.",
		Source: stagingBranch,
		Target: req.Branch,
	}

	pr, _, err := client.PullRequests.Create(req.Context, req.Repository, prInput)
	if err != nil {
		return nil, runtime.Errorwf(err, "error creating pull request from branch %s for repository %s",
			prInput.Source, req.Repository)
	}

	return &CollectionRouteCreatePrResponse{
		Link: pr.Link,
	}, nil
}

func createCompletion(ctx context.Context, mc model.ClientSet, systemMessage, userMessage string) (string, error) {
	apiToken, err := settings.OpenAiApiToken.Value(ctx, mc)
	if err != nil {
		return "", err
	}

	if apiToken == "" {
		return "", runtime.Error(http.StatusBadRequest,
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
		return "", fmt.Errorf("failed to create completion: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}

// trimMarkdownCodeBlock trims the beginning/ending markdown code block annotations(```)
// ChatGPT loves to reply codes with those and there's no deterministic way to tell it not to do that.
func trimMarkdownCodeBlock(s string) string {
	const (
		codeAnnotation = "```"
		newline        = "\n"
	)

	if strings.HasPrefix(s, codeAnnotation) && strings.HasSuffix(s, codeAnnotation) {
		s = strings.TrimPrefix(s, codeAnnotation)
		s = strings.TrimSuffix(s, codeAnnotation)
		s = strings.TrimPrefix(s, newline)
		s = strings.TrimSuffix(s, newline)
	}

	return s
}
