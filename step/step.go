package step

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

type Input struct {
	Verbose bool `env:"verbose,required"`

	Model        string          `env:"openai_model,required"`
	OpenAIApiKey stepconf.Secret `env:"openai_api_key,required"`

	GitHubToken stepconf.Secret `env:"github_token,required"`
	PRID        int             `env:"pr_id"`
	RepoOwner   string          `env:"repo_owner"`
	RepoName    string          `env:"repo_name"`

	PRTitle          string `env:"pr_title"`
	PRDescription    string `env:"pr_description"`
	RepoDescription  string `env:"repo_description"`
	PRDetailTemplate string `env:"pr_detail_template,required"`
	SystemPrompt     string `env:"system_prompt,required"`
}

type CodeReviewStep struct {
	logger      log.Logger
	inputParser stepconf.InputParser
	envRepo     env.Repository
}

type PRPromptInventory struct {
	PRTitle          string
	PRDescription    string
	RepoDescription  string
	RepoMainLanguage string
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	envRepo env.Repository,
) CodeReviewStep {
	return CodeReviewStep{
		logger:      logger,
		inputParser: inputParser,
		envRepo:     envRepo,
	}
}

func (step CodeReviewStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	step.logger.Infof("Fetching PR details...")
	ghClient := NewGitHubClient(string(input.GitHubToken), input.RepoOwner, input.RepoName, input.PRID)
	prData, err := ghClient.PullRequest(input.PRID)
	if err != nil {
		return fmt.Errorf("failed to get PR data: %w", err)
	}
	step.logger.Printf("PR title: %s", prData.GetTitle())
	step.logger.Printf("Repo description: %s", prData.GetBase().GetRepo().GetDescription())
	step.logger.Printf("Repo main language: %s", prData.GetBase().GetRepo().GetLanguage())
	step.logger.Donef("Done.")

	templateInventory := PRPromptInventory{
		PRTitle:          prData.GetTitle(),
		PRDescription:    prData.GetBody(),
		RepoDescription:  prData.GetBase().GetRepo().GetDescription(),
		RepoMainLanguage: prData.GetBase().GetRepo().GetLanguage(),
	}
	prompt, err := step.RenderPRPrompt(input.PRDetailTemplate, templateInventory)
	if err != nil {
		return fmt.Errorf("failed to render PR prompt: %w", err)
	}

	client := NewOpenAIClient(string(input.OpenAIApiKey), input.Model, step.logger)
	step.logger.Infof("Generating response...")
	completion, err := client.GetCompletion(input.SystemPrompt, prompt)
	if err != nil {
		return err
	}

	step.logger.Donef("%s", completion)

	return nil
}

func (step CodeReviewStep) RenderPRPrompt(detailTemplate string, inventory PRPromptInventory) (string, error) {
	tmpl, err := template.New("pr_details").Parse(detailTemplate)
	if err != nil {
		return "", fmt.Errorf("error creating template: %w", err)
	}

	resultBuffer := bytes.Buffer{}
	err = tmpl.Execute(&resultBuffer, inventory)
	if err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return resultBuffer.String(), nil
}
