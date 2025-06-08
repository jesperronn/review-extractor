# Review Extractor

A modular Python tool for extracting code review comments and diff context from pull requests across multiple Git hosting platforms. Designed to prepare structured data for AI-powered code review automation.

## 🚀 Features

- **Multi-platform support**: Bitbucket Server, GitHub, and GitLab
- **Comprehensive extraction**: Pull request comments, inline reviews, and diff context
- **Customer-configurable**: Per-customer configuration with multiple repositories
- **AI-ready output**: Structured JSON format optimized for machine learning workflows
- **Modular architecture**: Easy to extend with new platforms
- **Statistics generation**: Analysis of review patterns and comment frequency

## 📋 Requirements

- Python 3.8+
- API access to your Git hosting platforms
- Valid API tokens or credentials

## 🛠️ Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd review-extractor
```

2. Install dependencies:
```bash
pip install -r requirements.txt
```

## ⚙️ Configuration

Create a YAML configuration file for each customer in the `config/` directory:

```yaml
# config/customer-a.config
api_token: "your-api-token-here"
output_file: "customer-a-reviews.json"

repositories:
  - provider: bitbucket
    url: "https://bitbucket.example.com/projects/PROJ/repos/web-service"
  - provider: github  
    url: "https://github.com/customer-a/mobile-app"
  - provider: gitlab
    url: "https://gitlab.com/customer-a/backend-api"
```

### Configuration Options

| Field | Description | Required |
|-------|-------------|----------|
| `api_token` | Authentication token for the Git platform | Yes |
| `output_file` | Path for the generated JSON output | No (defaults to `reviews.json`) |
| `repositories` | List of repositories to extract from | Yes |
| `repositories[].provider` | Platform type: `bitbucket`, `github`, or `gitlab` | Yes |
| `repositories[].url` | Full repository URL | Yes |

## 🚀 Usage

Extract reviews for a specific customer:

```bash
python main.py --config config/customer-a.config
```

The tool will:
1. Connect to each configured repository
2. Fetch all pull requests (open, merged, declined)
3. Extract review comments with diff context
4. Generate a structured JSON file with the results

## 📊 Output Format

The tool generates JSON output with the following structure:

```json
{
  "extraction_date": "2024-06-08T10:30:00Z",
  "total_comments": 1247,
  "repositories_processed": 3,
  "reviews": [
    {
      "pr_id": 123,
      "pr_title": "Fix authentication timeout",
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
  ],
  "statistics": {
    "most_active_reviewers": ["jane.reviewer", "bob.senior"],
    "common_comment_types": ["naming", "performance", "security"],
    "files_with_most_comments": ["auth.py", "utils.js"]
  }
}
```

## 🏗️ Architecture

```
review_extractor/
├── main.py                 # Entry point and CLI interface
├── config/                 # Customer configuration files
├── adapters/               # Platform-specific API implementations
│   ├── bitbucket.py       # Bitbucket Server adapter
│   ├── github.py          # GitHub adapter
│   └── gitlab.py          # GitLab adapter
├── core/                   # Core business logic
│   ├── extractor.py       # Main extraction orchestration
│   ├── formatter.py       # Output formatting and statistics
│   └── utils.py           # Shared utilities and helpers
└── tests/                  # Unit tests
```

## 🔧 API Authentication

### Bitbucket Server
- Use personal access tokens or app passwords
- Ensure token has repository read permissions

### GitHub
- Generate a personal access token with `repo` scope
- For GitHub Enterprise, ensure API access is enabled

### GitLab
- Create a personal access token with `read_repository` scope
- For self-hosted GitLab, verify API endpoint accessibility

## 🧪 Testing

Run the test suite:

```bash
pytest tests/
```

Run tests with coverage:

```bash
pytest --cov=. tests/
```

## 📈 Statistics & Analysis

The tool automatically generates statistics including:

- **Comment frequency analysis**: Most common review patterns
- **Reviewer activity**: Who provides the most feedback
- **Code hotspots**: Files and functions that attract the most comments
- **Review density**: Comments per lines of code changed

## 🤖 AI Integration

The structured output is designed for AI workflows:

- **Training data**: Use historical reviews to train custom models
- **Few-shot prompting**: Provide examples for consistent review styles
- **Pattern recognition**: Identify common issues and suggestions
- **Automated suggestions**: Generate review comments for new PRs

## 🔒 Security Notes

- Store API tokens securely (consider environment variables)
- Review permissions before granting repository access
- Limit token scope to minimum required permissions
- Regularly rotate API tokens

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

[Add your license information here]

## 🆘 Troubleshooting

### Common Issues

**Authentication failures**
- Verify API token validity and permissions
- Check repository access rights
- Ensure correct API endpoints for self-hosted instances

**Rate limiting**
- The tool respects API rate limits automatically
- For large repositories, extraction may take time
- Consider running during off-peak hours

**Missing comments or PRs**
- Verify repository URLs are correct
- Check if PRs are accessible with the provided credentials
- Some platforms may limit historical data access

### Support

For issues and questions:
1. Check the troubleshooting section
2. Review configuration examples
3. Open an issue with detailed error information