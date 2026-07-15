# dockerlens

```
go run ./cmd/docklens image analyze nginx  
```


```
docklens/
│
├── cmd/                                           # Contains executable binaries for this repository.
│   └── docklens/
│       └── main.go                                # Starts the CLI application and wires all commands together.
│
├── internal/                                      # Private application packages used only inside DockLens.
│   │
│   ├── cli/                                       # Defines all CLI commands and flags.
│   │   ├── root.go                                # Registers the root `docklens` command.
│   │   ├── image.go                               # Registers all `docklens image` subcommands.
│   │   ├── container.go                           # Registers all `docklens container` subcommands.
│   │   ├── network.go                             # Registers all `docklens network` subcommands.
│   │   ├── volume.go                              # Registers all `docklens volume` subcommands.
│   │   └── system.go                              # Registers all `docklens system` subcommands.
│   │
│   ├── docker/                                    # Communicates with the Docker Engine API.
│   │   ├── client.go                              # Creates and manages the Docker SDK client.
│   │   ├── image.go                               # Retrieves image information from Docker.
│   │   ├── container.go                           # Retrieves container information from Docker.
│   │   ├── network.go                             # Retrieves network information from Docker.
│   │   ├── volume.go                              # Retrieves volume information from Docker.
│   │   └── system.go                              # Retrieves Docker daemon and system information.
│   │
│   ├── image/                                     # Implements image-related business logic.
│   │   ├── analyze.go                             # Performs overall image inspection and analysis.
│   │   ├── layers.go                              # Analyzes image layers and their sizes.
│   │   ├── diff.go                                # Compares two Docker images.
│   │   └── score.go                               # Evaluates image quality using best-practice rules.
│   │
│   ├── container/                                 # Implements container-related business logic.
│   │   ├── analyze.go                             # Performs overall container inspection.
│   │   ├── network.go                             # Analyzes container networking configuration.
│   │   ├── mounts.go                              # Inspects mounted volumes and bind mounts.
│   │   ├── resources.go                           # Analyzes CPU, memory, and resource usage.
│   │   └── processes.go                           # Lists and analyzes running processes.
│   │
│   ├── network/                                   # Implements Docker network analysis.
│   │   ├── analyze.go                             # Inspects Docker network configuration.
│   │   └── topology.go                            # Maps container-to-network relationships.
│   │
│   ├── volume/                                    # Implements Docker volume analysis.
│   │   ├── analyze.go                             # Inspects Docker volume configuration.
│   │   └── usage.go                               # Reports disk usage and utilization.
│   │
│   ├── output/                                    # Formats analysis results for users.
│   │   ├── table.go                               # Renders human-readable terminal tables.
│   │   ├── json.go                                # Renders machine-readable JSON output.
│   │   └── yaml.go                                # Renders YAML output.
│   │
│   └── common/                                    # Shared internal utilities used across packages.
│       ├── errors.go                              # Defines reusable application error types.
│       └── format.go                              # Provides common formatting helper functions.
│
├── docs/                                          # Project documentation, architecture, and design notes.
├── hack/                                          # Build, test, lint, and development automation scripts.
├── test/                                          # Integration and end-to-end test suites.
├── .github/workflows/                             # GitHub Actions CI/CD workflows.
├── Makefile                                       # Standard build, test, and development commands.
├── Dockerfile                                     # Builds the DockLens container image.
├── README.md                                      # Project overview, installation, and usage guide.
├── LICENSE                                        # Open-source license for the project.
├── go.mod                                         # Defines the Go module and project dependencies.
└── go.sum                                         # Records dependency checksums for reproducible builds.

```