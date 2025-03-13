# AI Resume Builder Backend

This is the backend service for the AI Resume Builder application. It provides APIs for analyzing job postings, generating tailored resumes, and creating PDF documents.

## Features

- Job posting analysis using AI
- Smart resume generation based on job requirements
- PDF resume generation with professional formatting
- RESTful API design
- Structured logging
- Configuration management
- Error handling

## Prerequisites

- Go 1.21 or later
- Make (for using Makefile commands)
- golangci-lint (for linting)

## Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Install dependencies:
   ```bash
   make deps
   ```
4. Run the application:
   ```bash
   make run
   ```

## API Endpoints

### POST /api/v1/analyze
Analyzes a job posting and extracts key requirements.

### POST /api/v1/generate
Generates a tailored resume based on job requirements.

### POST /api/v1/pdf
Generates a PDF version of the resume.

## Development

### Available Make Commands

- `make build`: Build the application
- `make run`: Run the application
- `make test`: Run tests
- `make clean`: Clean build artifacts
- `make deps`: Install dependencies
- `make lint`: Run linter
- `make mocks`: Generate mock files (if needed)
- `make all`: Clean, build, and test

## Project Structure

```
📦 backend
├── 📂 cmd
│   └── 📂 server            # Entry point
├── 📂 config                # Configuration
├── 📂 internal
│   ├── 📂 handlers         # HTTP handlers
│   ├── 📂 services         # Business logic
│   ├── 📂 repositories     # Data access
│   ├── 📂 utils           # Helper functions
│   └── 📂 models          # Data models
├── 📂 pkg                  # Shared packages
└── 📂 tests               # Test files
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 