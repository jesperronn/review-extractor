package models

import "time"

// Provider represents the Git platform type
type Provider string

const (
	ProviderBitbucket Provider = "bitbucket"
	ProviderGitHub    Provider = "github"
	ProviderGitLab    Provider = "gitlab"
)

// Review represents a single code review comment with its context
type Review struct {
	PRID           int       `json:"pr_id"`
	PRTitle        string    `json:"pr_title"`
	PRAuthor       string    `json:"pr_author"`
	Repository     string    `json:"repository"`
	Provider       Provider  `json:"provider"`
	CommentID      string    `json:"comment_id"`
	CommentAuthor  string    `json:"comment_author"`
	CommentText    string    `json:"comment_text"`
	CommentCreated time.Time `json:"comment_created"`
	FilePath       string    `json:"file_path"`
	LineNumber     int       `json:"line_number"`
	DiffContext    string    `json:"diff_context"`
}

// ExtractionResult represents the complete output of the review extraction
type ExtractionResult struct {
	ExtractionDate        time.Time `json:"extraction_date"`
	TotalComments         int       `json:"total_comments"`
	RepositoriesProcessed int       `json:"repositories_processed"`
	Reviews              []Review  `json:"reviews"`
	Statistics           Statistics `json:"statistics"`
}

// Statistics represents aggregated data about the reviews
type Statistics struct {
	MostActiveReviewers    []string `json:"most_active_reviewers"`
	CommonCommentTypes     []string `json:"common_comment_types"`
	FilesWithMostComments []string `json:"files_with_most_comments"`
}

// RepositoryConfig represents a single repository configuration
type RepositoryConfig struct {
	Provider Provider `yaml:"provider"`
	URL      string   `yaml:"url"`
}

// Config represents the complete configuration for a customer
type Config struct {
	APIToken    string            `yaml:"api_token"`
	OutputFile  string            `yaml:"output_file"`
	Repositories []RepositoryConfig `yaml:"repositories"`
} 