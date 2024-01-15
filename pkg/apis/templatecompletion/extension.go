package templatecompletion

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/drone/go-scm/scm"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/i18n"
	"github.com/seal-io/walrus/pkg/openai"
	"github.com/seal-io/walrus/pkg/templates"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/strs"
)

var builtinExamples = []PromptExample{
	{
		Name:   i18n.ExampleKubernetesName,
		Prompt: i18n.ExampleKubernetesPrompt,
	},
	{
		Name:   i18n.ExampleAlibabaCloudName,
		Prompt: i18n.ExampleAlibabaCloudPrompt,
	},
	{
		Name:   i18n.ExampleELKName,
		Prompt: i18n.ExampleELKPrompt,
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
	prompt := runtime.Translate(req.Context, i18n.TerraformModuleGenerateSystemMessage)

	result, err := openai.CreateCompletion(req.Context, h.modelClient, prompt, req.Text)
	if err != nil {
		return nil, err
	}

	return &CollectionRouteGenerateResponse{
		Text: trimMarkdownCodeBlock(result),
	}, nil
}

func (h Handler) CollectionRouteExplain(req CollectionRouteExplainRequest) (*CollectionRouteExplainResponse, error) {
	prompt := runtime.Translate(req.Context, i18n.TerraformModuleExplainSystemMessage)

	result, err := openai.CreateCompletion(req.Context, h.modelClient, prompt, req.Text)
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
	desc := runtime.Translate(req.Context, i18n.TerraformModuleCorrectSystemMessageDesc)
	correctedDesc := runtime.Translate(req.Context, i18n.TerraformModuleCorrectSystemMessageCorrectedDesc)
	explanationDesc := runtime.Translate(req.Context, i18n.TerraformModuleCorrectSystemMessageExplanationDesc)
	prompt := fmt.Sprintf(`%s\n{\n"corrected": "%s", "explanation": "%s"\n}\n`, desc, correctedDesc, explanationDesc)

	result, err := openai.CreateCompletion(req.Context, h.modelClient, prompt, req.Text)
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
		return nil, errorx.NewHttpError(http.StatusBadRequest, err.Error())
	}

	conn, err := h.modelClient.Connectors().Get(req.Context, req.ConnectorID)
	if err != nil {
		return nil, fmt.Errorf("error getting connector: %w", err)
	}

	if conn.Category != types.ConnectorCategoryVersionControl {
		return nil, errorx.HttpErrorf(http.StatusBadRequest,
			"%q is not a supported version control driver", conn.Type)
	}

	client, err := vcs.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("error creating version control system client: %w", err)
	}

	ref, _, err := client.Git.FindBranch(req.Context, req.Repository, req.Branch)
	if err != nil {
		return nil, errorx.WrapfHttpError(http.StatusBadRequest, err,
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
		return nil, errorx.Wrapf(err, "error creating new commit for repository %s",
			req.Repository)
	}

	stagingBranch := fmt.Sprintf("seal/module-" + strs.String(5))
	refInput := &scm.ReferenceInput{
		Name: stagingBranch,
		Sha:  commit.Sha,
	}
	_, err = client.Git.CreateBranch(req.Context, req.Repository, refInput)

	if err != nil {
		return nil, errorx.Wrapf(err, "error creating new branch %s for repository %s",
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
		return nil, errorx.Wrapf(err, "error creating pull request from branch %s for repository %s",
			prInput.Source, req.Repository)
	}

	return &CollectionRouteCreatePrResponse{
		Link: pr.Link,
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
