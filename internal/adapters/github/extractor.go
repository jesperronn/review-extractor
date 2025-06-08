package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jesper/review-extractor/pkg/models"
)

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
	if diff == "" || filePath == "" || lineNumber <= 0 {
		return ""
	}

	// Split diff into files
	fileDiffs := strings.Split(diff, "diff --git")

	// Find the relevant file diff
	for _, fileDiff := range fileDiffs {
		if strings.Contains(fileDiff, filePath) {
			// Split into lines
			lines := strings.Split(fileDiff, "\n")

			// Find the hunk that contains our line
			var currentLine int
			var inTargetHunk bool
			var contextLines []string

			for _, line := range lines {
				// Skip empty lines
				if line == "" {
					continue
				}

				// Check for hunk header
				if strings.HasPrefix(line, "@@") {
					// Parse the hunk header to get the starting line
					parts := strings.Split(line, " ")
					if len(parts) >= 2 {
						lineInfo := strings.Split(parts[1], ",")
						if len(lineInfo) >= 1 {
							// Remove the + or - from the line number
							lineNum := strings.TrimLeft(lineInfo[0], "+-")
							if num, err := strconv.Atoi(lineNum); err == nil {
								currentLine = num
								inTargetHunk = false
							}
						}
					}
					continue
				}

				// Skip diff header lines
				if strings.HasPrefix(line, "diff --git") || strings.HasPrefix(line, "index ") ||
					strings.HasPrefix(line, "---") || strings.HasPrefix(line, "+++") {
					continue
				}

				// Count lines in the hunk
				if !strings.HasPrefix(line, " ") {
					currentLine++
				}

				// If we're in the target line's hunk, collect context
				if currentLine >= lineNumber-3 && currentLine <= lineNumber+3 {
					inTargetHunk = true
					// Remove the + or - prefix and trim spaces
					cleanLine := strings.TrimSpace(strings.TrimLeft(line, "+- "))
					contextLines = append(contextLines, cleanLine)
				} else if inTargetHunk {
					// We've passed the context window
					break
				}
			}

			if len(contextLines) > 0 {
				return strings.Join(contextLines, "\n")
			}
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
