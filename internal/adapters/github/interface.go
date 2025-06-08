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
