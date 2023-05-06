package step

import (
	"context"
	"fmt"
	"strings"

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

func (c GitHubClient) UpsertComment(comment string) (*github.IssueComment, error) {
	ctx := context.Background()
	commentTag := "<!-- codereview-gpt -->"

	taggedComment, err := c.GetFirstCommentWithTag(commentTag)
	if err != nil {
		return nil, err
	}
	if taggedComment != nil {
		issueComment, _, err := c.client.Issues.EditComment(ctx, c.owner, c.repo, *taggedComment.ID, &github.IssueComment{
			Body: &comment,
		})
		return issueComment, err
	}

	taggedBody := fmt.Sprintf("%s\n%s", commentTag, comment)
	issueComment, _, err := c.client.Issues.CreateComment(ctx, c.owner, c.repo, c.prID, &github.IssueComment{
		Body: &taggedBody,
	})

	return issueComment, err
}

func (c GitHubClient) GetFirstCommentWithTag(tag string) (*github.IssueComment, error) {
	ctx := context.Background()
	issueComments, _, err := c.client.Issues.ListComments(ctx, c.owner, c.repo, c.prID, nil)
	if err != nil {
		return nil, err
	}

	for _, issueComment := range issueComments {
		if strings.Contains(*issueComment.Body, tag) {
			return issueComment, nil
		}
	}

	return nil, nil
}
