# Review Extractor - Requirements

## Objective
Build a modular Go application that extracts code review comments and diff context from pull requests across multiple Git platforms (Bitbucket Server, GitHub, GitLab) to prepare data for AI-powered code review automation.

## Core Functionality

### Data Extraction
- Extract from **all repositories** across multiple platforms
- For each pull request (all states: open, merged, declined):
  - All review comments (inline + general)
  - Diff context around each inline comment
  - Metadata: PR ID, title, author, repository, file path, line number, timestamp
  - Comment author and text

### Platform Support
- **Bitbucket Server** (on-premises/self-hosted)
- **GitHub** (github.com or GitHub Enterprise)
- **GitLab** (gitlab.com or self-hosted)

## Architecture

### Project Structure
```
review-extractor/
├── main.go                     # Entry point with config loading
├── cmd/                        # CLI commands
│   └── extract.go
├── config/                     # Per-customer configuration
│   ├── customer-a.yaml
│   └── customer-b.yaml
├── internal/
│   ├── adapters/              # Platform-specific API logic
│   │   ├── bitbucket/
│   │   │   ├── client.go
│   │   │   └── extractor.go
│   │   ├── github/
│   │   │   ├── client.go
│   │   │   └── extractor.go
│   │   └── gitlab/
│   │       ├── client.go
│   │       └── extractor.go
│   ├── core/                  # Shared business logic
│   │   ├── extractor.go       # Main extraction orchestration
│   │   ├── formatter.go       # Output formatting
│   │   └── types.go           # Shared data structures
│   └── utils/                 # Shared utilities
│       ├── http.go
│       └── config.go
├── pkg/                       # Public API interfaces
│   └── models/
│       └── review.go
├── test/                      # Integration tests
├── go.mod
├── go.sum
└── Makefile
```

### Configuration Format (YAML)
```yaml
# customer-a.yaml
api_token: "your-api-token-here"
output_file: "customer-a-reviews.json"

repositories:
  - provider: bitbucket
    url: "https://bitbucket.example.com/projects/PROJ/repos/repo1"
  - provider: github  
    url: "https://github.com/customer-a/project1"
  - provider: gitlab
    url: "https://gitlab.com/customer-a/project2"
```

## Output Format

### Primary Output: JSON
Structured data ready for AI consumption:
```json
{
  "reviews": [
    {
      "pr_id": 123,
      "pr_title": "Fix authentication bug",
      "pr_author": "john.doe",
      "repository": "web-service",
      "provider": "github",
      "comment_id": "456",
      "comment_author": "jane.reviewer",
      "comment_text": "Consider using a constant instead of magic number",
      "comment_created": "2024-11-10T14:30:00Z",
      "file_path": "src/auth.py",
      "line_number": 42,
      "diff_context": "- timeout = 30\n+ timeout = 300\n  return authenticate(user)"
    }
  ]
}
```

## Technical Requirements

### Dependencies
- Go 1.21+ - Core language and standard library
- `gopkg.in/yaml.v3` - YAML configuration parsing
- `github.com/stretchr/testify` - Testing framework
- `github.com/spf13/cobra` - CLI framework (optional)
- Standard library `net/http` - HTTP API calls

### API Integration
- Handle pagination for large repositories
- Proper error handling and rate limiting
- Support for different authentication methods (tokens, basic auth)

### Quality Assurance
- Unit tests using Go's built-in testing framework
- Integration tests for each adapter
- Modular design with clear interfaces for easy platform addition
- Git version control with incremental commits

## Usage
```bash
# Build the application
go build -o review-extractor ./cmd

# Run extraction
./review-extractor --config config/customer-a.yaml

# Alternative with go run
go run ./cmd --config config/customer-a.yaml
```

## Future Enhancements (Optional)

### Statistics Generation
- Most frequent review comment types
- Active reviewer analysis
- Comment density metrics (comments per diff line)
- File/language hotspot analysis

### AI Training Preparation
The extracted data will be used to:
- Train AI models on historical review patterns
- Create automated review comment suggestions
- Generate consistency guidelines for code review processes

## Success Criteria
1. Successfully extract comments from all three platforms using Go
2. Include meaningful diff context for each comment
3. Generate clean, structured JSON output suitable for AI processing
4. Configurable per customer with multiple repositories
5. Maintainable, tested Go codebase with clear separation of concerns
6. Efficient memory usage and concurrent processing where appropriate