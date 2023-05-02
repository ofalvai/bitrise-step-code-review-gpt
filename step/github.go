package step

import (
	"context"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

type GitHubClient struct {
	client *github.Client
	prID   int
	owner  string
	repo   string
}

func NewGitHubClient(token string, owner, repo string, prID int) GitHubClient {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tokenClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(tokenClient)

	return GitHubClient{
		client: client,
		owner:  owner,
		repo:   repo,
		prID:   prID,
	}
}

func (c GitHubClient) PullRequest(id int) (*github.PullRequest, error) {
	ctx := context.Background()
	pr, _, err := c.client.PullRequests.Get(ctx, c.owner, c.repo, id)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (c GitHubClient) PostComment(comment string) (*github.IssueComment, error) {
	// TODO: use GetFirstCommentWithTag() and UpdateComment() instead of CreateComment()
	// https://github.com/kvvzr/bitrise-step-comment-on-github-pull-request/blob/master/main.go#L29
	ctx := context.Background()
	issueComment, _, err := c.client.Issues.CreateComment(ctx, c.owner, c.repo, c.prID, &github.IssueComment{
		Body: &comment,
	})
	if err != nil {
		return nil, err
	}

	return issueComment, nil
}
