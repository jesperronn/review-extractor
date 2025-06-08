package github

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/jesper/review-extractor/pkg/models"
	"github.com/stretchr/testify/assert"
)

// MockClient implements the GitHub client interface for testing
type MockClient struct {
	prs      []*github.PullRequest
	comments []*github.PullRequestComment
	reviews  []*github.PullRequestReview
	diff     string
}

func (m *MockClient) GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
	return m.prs, nil
}

func (m *MockClient) GetPullRequestComments(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestComment, error) {
	return m.comments, nil
}

func (m *MockClient) GetPullRequestReviews(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	return m.reviews, nil
}

func (m *MockClient) GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error) {
	return m.diff, nil
}

func TestExtractReviews(t *testing.T) {
	now := time.Now()
	mockPR := &github.PullRequest{
		Number: github.Int(1),
		Title:  github.String("Test PR"),
		User: &github.User{
			Login: github.String("testuser"),
		},
	}

	mockComment := &github.PullRequestComment{
		ID:        github.Int64(1),
		Body:      github.String("Test comment"),
		Path:      github.String("test.go"),
		Line:      github.Int(10),
		User:      &github.User{Login: github.String("reviewer")},
		CreatedAt: &now,
	}

	mockReview := &github.PullRequestReview{
		ID:          github.Int64(1),
		Body:        github.String("Test review"),
		User:        &github.User{Login: github.String("reviewer")},
		SubmittedAt: &now,
	}

	mockClient := &MockClient{
		prs:      []*github.PullRequest{mockPR},
		comments: []*github.PullRequestComment{mockComment},
		reviews:  []*github.PullRequestReview{mockReview},
		diff:     "test diff",
	}

	extractor := &Extractor{
		client: mockClient,
	}

	reviews, err := extractor.ExtractReviews(context.Background(), "https://github.com/test/repo")
	assert.NoError(t, err)
	assert.Len(t, reviews, 2) // One comment and one review

	// Verify comment
	assert.Equal(t, models.Review{
		PRID:           1,
		PRTitle:        "Test PR",
		PRAuthor:       "testuser",
		Repository:     "repo",
		Provider:       models.ProviderGitHub,
		CommentID:      "1",
		CommentAuthor:  "reviewer",
		CommentText:    "Test comment",
		CommentCreated: now,
		FilePath:       "test.go",
		LineNumber:     10,
		DiffContext:    "test diff",
	}, reviews[0])

	// Verify review
	assert.Equal(t, models.Review{
		PRID:           1,
		PRTitle:        "Test PR",
		PRAuthor:       "testuser",
		Repository:     "repo",
		Provider:       models.ProviderGitHub,
		CommentID:      "1",
		CommentAuthor:  "reviewer",
		CommentText:    "Test review",
		CommentCreated: now,
		FilePath:       "",
		LineNumber:     0,
		DiffContext:    "",
	}, reviews[1])
}

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "valid URL",
			url:       "https://github.com/test/repo",
			wantOwner: "test",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:    "invalid URL",
			url:     "https://github.com/test",
			wantErr: true,
		},
		{
			name:    "non-github URL",
			url:     "https://gitlab.com/test/repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseGitHubURL(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantOwner, owner)
			assert.Equal(t, tt.wantRepo, repo)
		})
	}
}

func TestExtractDiffContext(t *testing.T) {
	tests := []struct {
		name       string
		diff       string
		filePath   string
		lineNumber int
		want       string
	}{
		{
			name: "exact match",
			diff: `diff --git a/test.go b/test.go
@@ -10,7 +10,7 @@
 line1
 line2
 line3
 line4
 line5
 line6
 line7`,
			filePath:   "test.go",
			lineNumber: 10,
			want:       "line1\nline2\nline3\nline4\nline5\nline6\nline7",
		},
		{
			name: "no match",
			diff: `diff --git a/other.go b/other.go
@@ -1,1 +1,1 @@
line1`,
			filePath:   "test.go",
			lineNumber: 10,
			want:       "",
		},
		{
			name:       "empty diff",
			diff:       "",
			filePath:   "test.go",
			lineNumber: 10,
			want:       "",
		},
		{
			name: "invalid line number",
			diff: `diff --git a/test.go b/test.go
@@ -1,1 +1,1 @@
line1`,
			filePath:   "test.go",
			lineNumber: 0,
			want:       "",
		},
		{
			name: "line number out of range",
			diff: `diff --git a/test.go b/test.go
@@ -1,1 +1,1 @@
line1`,
			filePath:   "test.go",
			lineNumber: 100,
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDiffContext(tt.diff, tt.filePath, tt.lineNumber)
			assert.Equal(t, tt.want, got)
		})
	}
}
