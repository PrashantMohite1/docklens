# DockLens

DockLens is a Docker image inspection and verification utility written in Go. It currently exposes image-level analysis commands and a verification workflow that compares local files or directories against paths inside a container image.

This repository is intentionally structured as an open-source project. The project already includes a project license and source layout, and the README now provides a route for new contributors and users to build, run, and verify the CLI tooling.

## Overview

DockLens is available under the license in the repository root. Anyone can clone, use, study, fork, and modify this project according to the terms described in the license.

## Prerequisites

Before running the CLI, make sure the following are available:

- Go 1.22 or later
- Docker Engine available on the host
- Access to Docker images that can be inspected or started as verification targets

## Quick Start

Build the CLI:

```bash
go build ./cmd/docklens
```

Analyze an image:

```bash
go run ./cmd/docklens image analyze nginx
```

You can also run the built binary directly:

```bash
./docklens image analyze nginx
```

### File Check Quick Guide

For quick setup, you can build an Alpine-style verification image from the repository docs folder:

```bash
cd ./docs/
docker build -t <image-name> .
```

A file check is useful when you want to confirm that a single file or a small group of files in your workspace matches the same SHA-256 content inside a container image.

Format:

```bash
docklens image verify <image-name> -f <local-filepath>:<image-filepath>,<local-filepath>:<image-filepath>
```

Example:

```bash
go run ./cmd/docklens image verify alpine:latest -f ./docs/test-files/first-file.txt:/first-file.txt
```

The command calculates a SHA-256 digest for the local file and runs `sha256sum` inside the container image to compare the digest values.



### Directory Check Quick Guide

The verify command also accepts a directory mapping through the `-d` flag. This performs a deterministic recursive file hash over the local directory and receives the image directory hash through a container command.

Format:

```bash
docklens image verify <image-name> -d <local-directory>:<image-directory>
```

Example:

```bash
go run ./cmd/docklens image verify alpine:latest -d ./docs/test-files:/app/test-files
```

Multiple directory mappings can be supplied in one call by separating each mapping with a comma:

```bash
go run ./cmd/docklens image verify alpine:latest -d ./docs/test-files:/app/test-files,./docs/test-files/temp:/app/temp
```


The implementation reads the local directory recursively, hashes each file, sorts the hashes, creates a directory manifest, and compares the generated digest against the result from `sha256sum` inside the container.

## Open-Source Contribution Flow

If you want to contribute:

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Run available tests or compile checks.
5. Open a pull request with a short explanation of the problem and fix.

The repository already keeps the source in the Go module and under the `internal` package layout for implementation details. `cmd` contains the CLI entry point.

## Repository Layout

```text
docklens/
├── cmd/                  # CLI command entry points
├── internal/             # Analyzer, CLI, Docker, and verification logic
├── docs/                 # Additional project documentation and examples
├── LICENSE               # Open-source project license
├── README.md             # User and contributor introduction
└── go.mod                # Go module definition
```
