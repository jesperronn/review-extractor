package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/jesper/review-extractor/pkg/models"
)

// ClientInterface defines the interface for GitHub API operations
type ClientInterface interface {
	GetPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error)
	GetPullRequestComments(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestComment, error)
	GetPullRequestReviews(ctx context.Context, owner, repo string, prNumber int) ([]*github.PullRequestReview, error)
	GetPullRequestDiff(ctx context.Context, owner, repo string, prNumber int) (string, error)
}

// Extractor implements the core.Extractor interface for GitHub
type Extractor struct {
	client ClientInterface
}

// NewExtractor creates a new GitHub extractor
func NewExtractor(token string) *Extractor {
	return &Extractor{
		client: NewClient(token),
	}
}

// ExtractReviews implements the core.Extractor interface
func (e *Extractor) ExtractReviews(ctx context.Context, repoURL string) ([]models.Review, error) {
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("invalid GitHub URL: %w", err)
	}

	// Get all pull requests
	prs, err := e.client.GetPullRequests(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests: %w", err)
	}

	var allReviews []models.Review

	// Process each pull request
	for _, pr := range prs {
		// Get comments
		comments, err := e.client.GetPullRequestComments(ctx, owner, repo, pr.GetNumber())
		if err != nil {
			return nil, fmt.Errorf("failed to get comments for PR #%d: %w", pr.GetNumber(), err)
		}

		// Get reviews
		reviews, err := e.client.GetPullRequestReviews(ctx, owner, repo, pr.GetNumber())
		if err != nil {
			return nil, fmt.Errorf("failed to get reviews for PR #%d: %w", pr.GetNumber(), err)
		}

		// Get diff for context
		diff, err := e.client.GetPullRequestDiff(ctx, owner, repo, pr.GetNumber())
		if err != nil {
			return nil, fmt.Errorf("failed to get diff for PR #%d: %w", pr.GetNumber(), err)
		}

		// Process comments
		for _, comment := range comments {
			review := models.Review{
				PRID:           pr.GetNumber(),
				PRTitle:        pr.GetTitle(),
				PRAuthor:       pr.GetUser().GetLogin(),
				Repository:     repo,
				Provider:       models.ProviderGitHub,
				CommentID:      fmt.Sprintf("%d", comment.GetID()),
				CommentAuthor:  comment.GetUser().GetLogin(),
				CommentText:    comment.GetBody(),
				CommentCreated: comment.GetCreatedAt(),
				FilePath:       comment.GetPath(),
				LineNumber:     comment.GetLine(),
				DiffContext:    extractDiffContext(diff, comment.GetPath(), comment.GetLine()),
			}
			allReviews = append(allReviews, review)
		}

		// Process reviews
		for _, review := range reviews {
			if review.GetBody() == "" {
				continue
			}

			reviewModel := models.Review{
				PRID:           pr.GetNumber(),
				PRTitle:        pr.GetTitle(),
				PRAuthor:       pr.GetUser().GetLogin(),
				Repository:     repo,
				Provider:       models.ProviderGitHub,
				CommentID:      fmt.Sprintf("%d", review.GetID()),
				CommentAuthor:  review.GetUser().GetLogin(),
				CommentText:    review.GetBody(),
				CommentCreated: review.GetSubmittedAt(),
				// Note: Reviews don't have file/line context by default
				FilePath:    "",
				LineNumber:  0,
				DiffContext: "",
			}
			allReviews = append(allReviews, reviewModel)
		}
	}

	return allReviews, nil
}

// parseGitHubURL extracts owner and repo from a GitHub URL
func parseGitHubURL(url string) (owner, repo string, err error) {
	// Remove protocol and domain
	parts := strings.Split(url, "github.com/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format")
	}

	// Split owner/repo
	pathParts := strings.Split(parts[1], "/")
	if len(pathParts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format")
	}

	return pathParts[0], pathParts[1], nil
}

// extractDiffContext extracts the relevant diff context around a specific line
func extractDiffContext(diff, filePath string, lineNumber int) string {
	// Split diff into files
	fileDiffs := strings.Split(diff, "diff --git")

	// Find the relevant file diff
	for _, fileDiff := range fileDiffs {
		if strings.Contains(fileDiff, filePath) {
			// Split into lines
			lines := strings.Split(fileDiff, "\n")

			// Find the context around the line
			start := max(0, lineNumber-3)
			end := min(len(lines), lineNumber+3)

			return strings.Join(lines[start:end], "\n")
		}
	}

	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
