# CodeReviewGPT

[![Step changelog](https://shields.io/github/v/release/ofalvai/bitrise-step-code-review-gpt?include_prereleases&label=changelog&color=blueviolet)](https://github.com/ofalvai/bitrise-step-code-review-gpt/releases)

Let a language model give early feedback on your PRs automatically.

<details>
<summary>Description</summary>

Let a language model give early feedback on your PRs automatically.

This step uses the GitHub and OpenAI APIs to review a PR and post a comment with reminders and thought-provoking questions the PR author might have missed.

The comments are generated based on the PR title, description, and a system prompt. The step **does not** access code, only PR metadata (makes you write better PR descriptions 😜).

Limitations:
- bring your own OpenAI API key
- GitHub only

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `github_token` | GitHub API token with **write** access to the repository | required, sensitive |  |
| `pr_id` | The ID of the Pull Request to review |  | `$BITRISE_PULL_REQUEST` |
| `repo_owner` | The owner of the repository to review. Used for fetching PR metadata and posting the comment. |  | `$BITRISEIO_GIT_REPOSITORY_OWNER` |
| `repo_name` | The name of the repository to review. Used for fetching PR metadata and posting the comment. |  | `$BITRISEIO_GIT_REPOSITORY_SLUG` |
| `pr_title` | Title of Pull Request |  |  |
| `pr_description` | Description of Pull Request |  |  |
| `repo_description` | Description of the repository related to the PR. Use this field to give more context about the PR |  |  |
| `system_prompt` | Initial prompt for the language model | required | `You are an assistant helping pull request authors improve the changeset by asking questions from the author. Your questions may include potential edge-cases, side-effects to consider, or testing strategies, but you are encouraged to ask about other concerns if they are relevant to the PR. Your questions are creative and open-ended. Be polite, but use simple language and short sentences. Your response is a numbered list of comments. Limit your response to at most 5 comments. You may use Markdown formatting, but you must not include links in your output.` |
| `pr_detail_template` | Information about the Pull Request | required | `About this repo: {{ .RepoDescription }} Main language of this repo: {{ .RepoMainLanguage }}  PR title: {{ .PRTitle }}  PR description: {{ .PRDescription }}` |
| `openai_api_key` | OpenAI API key to use for requests | required, sensitive |  |
| `openai_model` | OpenAI model to use | required | `gpt-4o` |
| `verbose` | Enable logging additional information for troubleshooting | required | `false` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/ofalvai/bitrise-step-code-review-gpt/pulls) and [issues](https://github.com/ofalvai/bitrise-step-code-review-gpt/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
