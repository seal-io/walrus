package modulecompletion

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/drone/go-scm/scm"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"

	"github.com/seal-io/seal/pkg/apis/modulecompletion/view"
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/connectors/types"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/i18n/text"
	"github.com/seal-io/seal/pkg/modules"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/pkg/vcs"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

var examples = []view.ModuleCompletionPromptExample{
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

const (
	gpt35MaxTokens = 4096
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

func (h Handler) CollectionRouteExamples(c *gin.Context, _ view.ExampleRequest) (view.ExampleResponse, error) {
	var translated []view.ModuleCompletionPromptExample

	for _, e := range examples {
		translated = append(translated, view.ModuleCompletionPromptExample{
			Name:   runtime.Translate(c, e.Name),
			Prompt: runtime.Translate(c, e.Prompt),
		})
	}

	return translated, nil
}

func (h Handler) CollectionRouteGenerate(ctx *gin.Context, req view.GenerateRequest) (*view.GenerateResponse, error) {
	prompt := runtime.Translate(ctx, text.TerraformModuleGenerateSystemMessage)
	result, err := h.createCompletion(ctx, prompt, req.Text)
	if err != nil {
		return nil, err
	}

	return &view.GenerateResponse{
		Text: trimMarkdownCodeBlock(result),
	}, nil
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

func (h Handler) CollectionRouteExplain(ctx *gin.Context, req view.ExplainRequest) (*view.ExplainResponse, error) {
	prompt := runtime.Translate(ctx, text.TerraformModuleExplainSystemMessage)
	result, err := h.createCompletion(ctx, prompt, req.Text)
	if err != nil {
		return nil, err
	}
	return &view.ExplainResponse{
		Text: result,
	}, nil
}

func (h Handler) CollectionRouteCorrect(ctx *gin.Context, req view.CorrectRequest) (*view.CorrectResponse, error) {
	// gotext cannot handle brackets in messages, see https://github.com/golang/go/issues/27849.
	// we need to split the text as a workaround.
	desc := runtime.Translate(ctx, text.TerraformModuleCorrectSystemMessageDesc)
	correctedDesc := runtime.Translate(ctx, text.TerraformModuleCorrectSystemMessageCorrectedDesc)
	explanationDesc := runtime.Translate(ctx, text.TerraformModuleCorrectSystemMessageExplanationDesc)
	prompt := fmt.Sprintf(`%s\n{"corrected": "%s", "explanation": "%s"}\n`, desc, correctedDesc, explanationDesc)
	result, err := h.createCompletion(ctx, prompt, req.Text)
	if err != nil {
		return nil, err
	}
	correctResp := &view.CorrectResponse{}
	if err := json.Unmarshal([]byte(result), correctResp); err != nil {
		log.Debugf("correction message is not in the format requested by the prompt. output:\n%v", result)
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
		return "", runtime.Error(http.StatusBadRequest, "invalid input: OpenAI API token is not configured")
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

func (h Handler) CollectionRouteCreatePr(ctx *gin.Context, req view.CreatePrRequest) (*view.CreatePrResponse, error) {
	moduleName := modules.GetModuleNameByPath(req.Path)
	moduleFiles, err := modules.GetTerraformModuleFiles(moduleName, req.Content)
	if err != nil {
		return nil, runtime.Error(http.StatusBadRequest, err)
	}
	conn, err := h.modelClient.Connectors().Get(ctx, req.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("error getting connector: %w", err)
	}

	if !types.IsVCS(conn) {
		return nil, runtime.Errorf(http.StatusBadRequest, "%q is not a supported version control driver", conn.Type)
	}

	client, err := vcs.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating version control system client: %w", err)
	}

	ref, _, err := client.Git.FindBranch(ctx, req.Repository, req.Branch)
	if err != nil {
		return nil, runtime.Errorpf(http.StatusBadRequest, err, "error indexing branch %s from repository %s",
			req.Branch, req.Repository)
	}

	var commitInput = &scm.CommitInput{
		Message: "Module generated from Seal",
		Base:    ref.Sha,
	}
	for name, content := range moduleFiles {
		commitInput.Blobs = append(commitInput.Blobs, scm.Blob{
			Path:    filepath.Join(req.Path, name),
			Mode:    "100644",
			Content: content,
		})
	}
	commit, _, err := client.Git.CreateCommit(ctx, req.Repository, commitInput)
	if err != nil {
		return nil, runtime.Errorwf(err, "error creating new commit for repository %s",
			req.Repository)
	}

	stagingBranch := fmt.Sprintf("seal/module-" + strs.String(5))
	var refInput = &scm.ReferenceInput{
		Name: stagingBranch,
		Sha:  commit.Sha,
	}
	_, err = client.Git.CreateBranch(ctx, req.Repository, refInput)
	if err != nil {
		return nil, runtime.Errorwf(err, "error creating new branch %s for repository %s",
			refInput.Name, req.Repository)
	}

	// TODO more informative PR body. e.g., let chatGPT generate it.
	var prInput = &scm.PullRequestInput{
		Title:  fmt.Sprintf("Add module %s", moduleName),
		Body:   "This is a module proposed from Seal.",
		Source: stagingBranch,
		Target: req.Branch,
	}
	pr, _, err := client.PullRequests.Create(ctx, req.Repository, prInput)
	if err != nil {
		return nil, runtime.Errorwf(err, "error creating pull request from branch %s for repository %s",
			prInput.Source, req.Repository)
	}
	return &view.CreatePrResponse{
		Link: pr.Link,
	}, nil
}
