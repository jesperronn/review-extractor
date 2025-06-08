package github

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the GitHub client
type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	args := m.Called(ctx, owner, repo)
	return args.Get(0).([]*github.PullRequest), args.Error(1)
}

func (m *MockClient) GetPullRequestComments(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestComment, error) {
	args := m.Called(ctx, owner, repo, prNumber)
	return args.Get(0).([]*github.PullRequestComment), args.Error(1)
}

func (m *MockClient) GetPullRequestReviews(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	args := m.Called(ctx, owner, repo, prNumber)
	return args.Get(0).([]*github.PullRequestReview), args.Error(1)
}

func (m *MockClient) GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	args := m.Called(ctx, owner, repo, prNumber)
	return args.String(0), args.Error(1)
}

func TestExtractReviews(t *testing.T) {
	// Create mock client
	mockClient := new(MockClient)

	// Create test data
	pr := &github.PullRequest{
		Number: github.Int(1),
		Title:  github.String("Test PR"),
		User: &github.User{
			Login: github.String("testuser"),
		},
	}

	comment := &github.PullRequestComment{
		ID:        github.Int64(1),
		Body:      github.String("Test comment"),
		Path:      github.String("test.go"),
		Line:      github.Int(10),
		User:      &github.User{Login: github.String("reviewer")},
		CreatedAt: &github.Timestamp{Time: time.Now()},
	}

	review := &github.PullRequestReview{
		ID:          github.Int64(1),
		Body:        github.String("Test review"),
		User:        &github.User{Login: github.String("reviewer")},
		SubmittedAt: &github.Timestamp{Time: time.Now()},
	}

	// Set up expectations
	mockClient.On("GetPullRequests", mock.Anything, "testowner", "testrepo").Return([]*github.PullRequest{pr}, nil)
	mockClient.On("GetPullRequestComments", mock.Anything, "testowner", "testrepo", 1).Return([]*github.PullRequestComment{comment}, nil)
	mockClient.On("GetPullRequestReviews", mock.Anything, "testowner", "testrepo", 1).Return([]*github.PullRequestReview{review}, nil)
	mockClient.On("GetPullRequestDiff", mock.Anything, "testowner", "testrepo", 1).Return("test diff", nil)

	// Create extractor with mock client
	extractor := &Extractor{client: mockClient}

	// Test extraction
	reviews, err := extractor.ExtractReviews(context.Background(), "https://github.com/testowner/testrepo")
	assert.NoError(t, err)
	assert.Len(t, reviews, 2) // One comment and one review

	// Verify comment
	assert.Equal(t, 1, reviews[0].PRID)
	assert.Equal(t, "Test PR", reviews[0].PRTitle)
	assert.Equal(t, "testuser", reviews[0].PRAuthor)
	assert.Equal(t, "testrepo", reviews[0].Repository)
	assert.Equal(t, "1", reviews[0].CommentID)
	assert.Equal(t, "reviewer", reviews[0].CommentAuthor)
	assert.Equal(t, "Test comment", reviews[0].CommentText)
	assert.Equal(t, "test.go", reviews[0].FilePath)
	assert.Equal(t, 10, reviews[0].LineNumber)

	// Verify review
	assert.Equal(t, 1, reviews[1].PRID)
	assert.Equal(t, "Test PR", reviews[1].PRTitle)
	assert.Equal(t, "testuser", reviews[1].PRAuthor)
	assert.Equal(t, "testrepo", reviews[1].Repository)
	assert.Equal(t, "1", reviews[1].CommentID)
	assert.Equal(t, "reviewer", reviews[1].CommentAuthor)
	assert.Equal(t, "Test review", reviews[1].CommentText)
}

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		owner    string
		repo     string
		hasError bool
	}{
		{
			name:     "valid URL",
			url:      "https://github.com/testowner/testrepo",
			owner:    "testowner",
			repo:     "testrepo",
			hasError: false,
		},
		{
			name:     "invalid URL format",
			url:      "https://github.com/testowner",
			owner:    "",
			repo:     "",
			hasError: true,
		},
		{
			name:     "invalid domain",
			url:      "https://gitlab.com/testowner/testrepo",
			owner:    "",
			repo:     "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseGitHubURL(tt.url)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.owner, owner)
				assert.Equal(t, tt.repo, repo)
			}
		})
	}
}

func TestExtractDiffContext(t *testing.T) {
	diff := `diff --git a/test.go b/test.go
index abc123..def456 100644
--- a/test.go
+++ b/test.go
@@ -10,6 +10,7 @@ func main() {
 	fmt.Println("Hello")
+	fmt.Println("World")
 	return
 }`

	context := extractDiffContext(diff, "test.go", 11)
	assert.Contains(t, context, "fmt.Println(\"Hello\")")
	assert.Contains(t, context, "fmt.Println(\"World\")")
} 