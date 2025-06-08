package models

import "time"

// Provider represents a code review platform
type Provider string

const (
	ProviderGitHub Provider = "github"
	ProviderGitLab Provider = "gitlab"
)

// Review represents a code review comment
type Review struct {
	PRID           int       `json:"pr_id" yaml:"pr_id"`
	PRTitle        string    `json:"pr_title" yaml:"pr_title"`
	PRAuthor       string    `json:"pr_author" yaml:"pr_author"`
	Repository     string    `json:"repository" yaml:"repository"`
	Provider       Provider  `json:"provider" yaml:"provider"`
	CommentID      string    `json:"comment_id" yaml:"comment_id"`
	CommentAuthor  string    `json:"comment_author" yaml:"comment_author"`
	CommentText    string    `json:"comment_text" yaml:"comment_text"`
	CommentCreated time.Time `json:"comment_created" yaml:"comment_created"`
	FilePath       string    `json:"file_path" yaml:"file_path"`
	LineNumber     int       `json:"line_number" yaml:"line_number"`
	DiffContext    string    `json:"diff_context" yaml:"diff_context"`
}

// ExtractionResult represents the result of extracting reviews
type ExtractionResult struct {
	Reviews     []Review   `json:"reviews" yaml:"reviews"`
	Statistics  Statistics `json:"statistics" yaml:"statistics"`
	ExtractedAt time.Time  `json:"extracted_at" yaml:"extracted_at"`
}

// Statistics represents statistics about the extracted reviews
type Statistics struct {
	TotalReviews    int      `json:"total_reviews" yaml:"total_reviews"`
	TotalPRs        int      `json:"total_prs" yaml:"total_prs"`
	TopReviewers    []string `json:"top_reviewers" yaml:"top_reviewers"`
	TopRepositories []string `json:"top_repositories" yaml:"top_repositories"`
	AveragePRSize   float64  `json:"average_pr_size" yaml:"average_pr_size"`
	ReviewFrequency float64  `json:"review_frequency" yaml:"review_frequency"`
}

// RepositoryConfig represents configuration for a repository
type RepositoryConfig struct {
	URL      string   `json:"url" yaml:"url"`
	Provider Provider `json:"provider" yaml:"provider"`
}

// GitHubConfig represents GitHub-specific configuration
type GitHubConfig struct {
	Token string `json:"token" yaml:"token"`
}

// Config represents the overall configuration
type Config struct {
	Repositories []RepositoryConfig `json:"repositories" yaml:"repositories"`
	GitHub       GitHubConfig       `json:"github" yaml:"github"`
}
