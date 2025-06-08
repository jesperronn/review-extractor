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

// RepositoryConfig represents a repository configuration
type RepositoryConfig struct {
        URL      string   `yaml:"url"`
        Provider Provider `yaml:"provider"`
}

// GitHubConfig represents GitHub-specific configuration
type GitHubConfig struct {
        Token string `yaml:"token"`
}

// Config represents the application configuration
type Config struct {
        Repositories []RepositoryConfig `yaml:"repositories"`
        GitHub       GitHubConfig       `yaml:"github"`
        OutputFile   string             `yaml:"output_file"`
        APIToken     string             `yaml:"api_token"`
}

// Statistics represents aggregated review statistics
type Statistics struct {
        TotalReviews    int      `json:"total_reviews"`
        TotalPRs        int      `json:"total_prs"`
        TopReviewers    []string `json:"top_reviewers"`
        TopRepositories []string `json:"top_repositories"`
        AveragePRSize   float64  `json:"average_pr_size"`
        ReviewFrequency float64  `json:"review_frequency"`
}

// ExtractionResult represents the result of a review extraction
type ExtractionResult struct {
        Reviews               []Review   `json:"reviews"`
        Statistics            Statistics `json:"statistics"`
        ExtractedAt           time.Time  `json:"extracted_at"`
        TotalComments         int        `json:"total_comments"`
        RepositoriesProcessed int        `json:"repositories_processed"`
}
