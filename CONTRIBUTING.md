# Contributing to DockLens

Thanks for your interest in contributing to DockLens.

## Ways to contribute

You can contribute by:

- Reporting bugs
- Suggesting new image analysis or verification workflows
- Improving documentation
- Adding integration tests
- Improving CLI UX and output formatting
- Refactoring or extending the analyzer internals

## Local development

Clone the repository:

```bash
git clone https://github.com/<your-user>/docklens.git
cd docklens
```

Install dependencies:

```bash
go mod download
```

Build the CLI:

```bash
go build ./cmd/docklens
```

Run analysis locally:

```bash
go run ./cmd/docklens image analyze nginx
```

## Verification workflows

The verification command compares a local file or directory against a path inside a Docker image using SHA-256 digest comparisons.

File verification:

```bash
go run ./cmd/docklens image verify <image-name> -f <local-path>:<image-path>
```

Directory verification:

```bash
go run ./cmd/docklens image verify <image-name> -d <local-dir>:<image-dir>
```

## Pull request process

1. Create a small, focused branch.
2. Keep the change scoped to a single problem or feature.
3. Please add or update tests when practical.
4. Document any new flags or output behavior in the README.
5. Open a pull request with a short explanation and verification notes.

## Code of conduct

Participation in this repository should remain respectful, constructive, and inclusive.
