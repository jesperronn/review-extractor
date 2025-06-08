# Review Extractor - Requirements

## Objective
Build a modular Python script that extracts code review comments and diff context from pull requests across multiple Git platforms (Bitbucket Server, GitHub, GitLab) to prepare data for AI-powered code review automation.

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
review_extractor/
├── main.py                     # Entry point with config loading
├── config/                     # Per-customer configuration
│   ├── customer-a.config
│   └── customer-b.config
├── adapters/                   # Platform-specific API logic
│   ├── bitbucket.py
│   ├── github.py
│   └── gitlab.py
├── core/                       # Shared business logic
│   ├── extractor.py           # Main extraction orchestration
│   ├── formatter.py           # Output formatting
│   └── utils.py               # Shared utilities
├── tests/                      # Unit tests for all modules
└── requirements.txt
```

### Configuration Format (YAML)
```yaml
# customer-a.config
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
- `requests` - HTTP API calls
- `PyYAML` - Configuration parsing
- `pytest` - Unit testing

### API Integration
- Handle pagination for large repositories
- Proper error handling and rate limiting
- Support for different authentication methods (tokens, basic auth)

### Quality Assurance
- Unit tests for each adapter and core module
- Modular design for easy platform addition
- Git version control with incremental commits

## Usage
```bash
python main.py --config config/customer-a.config
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
1. Successfully extract comments from all three platforms
2. Include meaningful diff context for each comment
3. Generate clean, structured output suitable for AI processing
4. Configurable per customer with multiple repositories
5. Maintainable, tested codebase with clear separation of concerns