package github

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v28/github"
)

// ArmeriaRepo is a wrapper around the GitHub API client.
type ArmeriaRepo struct {
	client *github.Client
}

// New returns a new GitHub client.
func New() *ArmeriaRepo {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &ArmeriaRepo{
		client: github.NewClient(tc),
	}
}

// CreateIssue creates a new GitHub issue in the heyitsmdr/armeria repo.
func (gh *ArmeriaRepo) CreateIssue(characterName string, issueBody string, rawBug string) (*github.Issue, error) {
	title := fmt.Sprintf("[%s] %s", characterName, rawBug)
	if len(rawBug) > 100 {
		title = fmt.Sprintf("[%s] %s", characterName, rawBug[:100]+"...")
	}

	ir := &github.IssueRequest{
		Title:  &title,
		Body:   &issueBody,
		Labels: &[]string{"triage", "in-game"},
	}

	i, _, err := gh.client.Issues.Create(context.Background(), "heyitsmdr", "armeria", ir)
	if err != nil {
		return nil, err
	}

	return i, nil
}
