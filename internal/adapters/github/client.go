package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// Client implements the ClientInterface for GitHub API operations
type Client struct {
	client *github.Client
}

// NewClient creates a new GitHub API client
func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
	}
}

// GetPullRequests fetches all pull requests for a repository
func (c *Client) GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	var allPRs []*github.PullRequest
	opts := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		prs, resp, err := c.client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}

		allPRs = append(allPRs, prs...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allPRs, nil
}

// GetPullRequestComments fetches all comments for a pull request
func (c *Client) GetPullRequestComments(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestComment, error) {
	var allComments []*github.PullRequestComment
	opts := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		comments, resp, err := c.client.PullRequests.ListComments(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull request comments: %w", err)
		}

		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allComments, nil
}

// GetPullRequestReviews fetches all reviews for a pull request
func (c *Client) GetPullRequestReviews(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	var allReviews []*github.PullRequestReview
	opts := &github.ListOptions{
		PerPage: 100,
	}

	for {
		reviews, resp, err := c.client.PullRequests.ListReviews(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull request reviews: %w", err)
		}

		allReviews = append(allReviews, reviews...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReviews, nil
}

// GetPullRequestDiff fetches the diff for a pull request
func (c *Client) GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	diff, _, err := c.client.PullRequests.GetRaw(ctx, owner, repo, prNumber, github.RawOptions{
		Type: github.Diff,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get pull request diff: %w", err)
	}

	return diff, nil
}
