package github

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGitHubClient is a mock implementation of the GitHub client interface
type MockGitHubClient struct {
	mock.Mock
}

func (m *MockGitHubClient) GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	args := m.Called(ctx, owner, repo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*github.PullRequest), args.Error(1)
}

func (m *MockGitHubClient) GetPullRequestComments(ctx context.Context, owner, repo string, number int) ([]*github.PullRequestComment, error) {
	args := m.Called(ctx, owner, repo, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*github.PullRequestComment), args.Error(1)
}

func (m *MockGitHubClient) GetPullRequestReviews(ctx context.Context, owner, repo string, number int) ([]*github.PullRequestReview, error) {
	args := m.Called(ctx, owner, repo, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*github.PullRequestReview), args.Error(1)
}

func (m *MockGitHubClient) GetPullRequestDiff(ctx context.Context, owner, repo string, number int) (string, error) {
	args := m.Called(ctx, owner, repo, number)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func TestGetPullRequests(t *testing.T) {
	mockClient := new(MockGitHubClient)
	client := &Client{client: mockClient}

	ctx := context.Background()
	owner := "testowner"
	repo := "testrepo"

	// Test successful case
	expectedPRs := []*github.PullRequest{
		{
			Number: github.Int(1),
			Title:  github.String("Test PR"),
			User: &github.User{
				Login: github.String("testuser"),
			},
		},
	}

	mockClient.On("GetPullRequests", ctx, owner, repo).Return(expectedPRs, nil)

	prs, err := client.GetPullRequests(ctx, owner, repo)
	assert.NoError(t, err)
	assert.Equal(t, expectedPRs, prs)

	// Test error case
	mockClient.On("GetPullRequests", ctx, "error", "repo").Return(nil, assert.AnError)

	_, err = client.GetPullRequests(ctx, "error", "repo")
	assert.Error(t, err)
}

func TestGetPullRequestComments(t *testing.T) {
	mockClient := new(MockGitHubClient)
	client := &Client{client: mockClient}

	ctx := context.Background()
	owner := "testowner"
	repo := "testrepo"
	number := 1

	// Test successful case
	expectedComments := []*github.PullRequestComment{
		{
			ID:        github.Int64(1),
			User:      &github.User{Login: github.String("reviewer")},
			Body:      github.String("Test comment"),
			CreatedAt: &time.Time{},
			Path:      github.String("test.go"),
			Line:      github.Int(10),
		},
	}

	mockClient.On("GetPullRequestComments", ctx, owner, repo, number).Return(expectedComments, nil)

	comments, err := client.GetPullRequestComments(ctx, owner, repo, number)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)

	// Test error case
	mockClient.On("GetPullRequestComments", ctx, "error", "repo", number).Return(nil, assert.AnError)

	_, err = client.GetPullRequestComments(ctx, "error", "repo", number)
	assert.Error(t, err)
}

func TestGetPullRequestReviews(t *testing.T) {
	mockClient := new(MockGitHubClient)
	client := &Client{client: mockClient}

	ctx := context.Background()
	owner := "testowner"
	repo := "testrepo"
	number := 1

	// Test successful case
	expectedReviews := []*github.PullRequestReview{
		{
			ID:    github.Int64(1),
			User:  &github.User{Login: github.String("reviewer")},
			Body:  github.String("Test review"),
			State: github.String("APPROVED"),
		},
	}

	mockClient.On("GetPullRequestReviews", ctx, owner, repo, number).Return(expectedReviews, nil)

	reviews, err := client.GetPullRequestReviews(ctx, owner, repo, number)
	assert.NoError(t, err)
	assert.Equal(t, expectedReviews, reviews)

	// Test error case
	mockClient.On("GetPullRequestReviews", ctx, "error", "repo", number).Return(nil, assert.AnError)

	_, err = client.GetPullRequestReviews(ctx, "error", "repo", number)
	assert.Error(t, err)
}

func TestGetPullRequestDiff(t *testing.T) {
	mockClient := new(MockGitHubClient)
	client := &Client{client: mockClient}

	ctx := context.Background()
	owner := "testowner"
	repo := "testrepo"
	number := 1

	// Test successful case
	expectedDiff := "diff --git a/test.go b/test.go\n@@ -1,1 +1,1 @@\n-test\n+new test"

	mockClient.On("GetPullRequestDiff", ctx, owner, repo, number).Return(expectedDiff, nil)

	diff, err := client.GetPullRequestDiff(ctx, owner, repo, number)
	assert.NoError(t, err)
	assert.Equal(t, expectedDiff, diff)

	// Test error case
	mockClient.On("GetPullRequestDiff", ctx, "error", "repo", number).Return(nil, assert.AnError)

	_, err = client.GetPullRequestDiff(ctx, "error", "repo", number)
	assert.Error(t, err)
}
