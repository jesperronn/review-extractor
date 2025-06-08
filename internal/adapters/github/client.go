package github

import (
	"context"

	"github.com/google/go-github/v45/github"
)

// ClientInterface defines the interface for GitHub API operations
type ClientInterface interface {
        GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error)
        GetPullRequestComments(ctx context.Context, owner, repo string, number int) ([]*github.PullRequestComment, error)
        GetPullRequestReviews(ctx context.Context, owner, repo string, number int) ([]*github.PullRequestReview, error)
        GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (string, error)
}

// Client wraps the GitHub API client
type Client struct {
        client ClientInterface
}

// NewClient creates a new GitHub client
func NewClient(token string) *Client {
        client := github.NewClient(nil)
        if token != "" {
                client = github.NewClient(nil).WithAuthToken(token)
	}
        return &Client{client: client}
}

// GetPullRequests fetches pull requests for a repository
func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
        return c.client.GetPullRequests(ctx, owner, repo)
}

// GetPullRequestComments fetches comments for a pull request
func (c *Client) GetPullRequestComments(ctx context.Context, owner, repo string, number int) ([]*github.PullRequestComment, error) {
        return c.client.GetPullRequestComments(ctx, owner, repo, number)
}

// GetPullRequestReviews fetches reviews for a pull request
func (c *Client) GetPullRequestReviews(ctx context.Context, owner, repo string, number int) ([]*github.PullRequestReview, error) {
        return c.client.GetPullRequestReviews(ctx, owner, repo, number)
}

// GetPullRequestDiff fetches the diff for a pull request
func (c *Client) GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (string, error) {
        return c.client.GetPullRequestDiff(ctx, owner, repo, number)
}
