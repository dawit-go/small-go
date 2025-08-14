# Small-Go

A CLI utility for generating Go project scaffolds with different architecture patterns.

## Overview

Small-Go standardizes Go project layouts and encourages separation of concerns between **Domain**, **Application**, **Ports**, and **Adapters**. It makes it quick to bootstrap a new Go service with best practices built-in.

## Installation

```bash
go install github.com/dawit-go/small-go@latest
```

## Usage

### List Available Templates

```bash
small-go list
```

### Create a New Project

#### Interactive Template Selection
```bash
small-go new <project_name>
```

#### Direct Template Specification
```bash
small-go new <project_name> --template <template_name>
```

This will:
1. Create a new folder named `<project_name>`
2. Initialize a Go module inside (`go mod init <project_name>`)
3. Generate a complete project scaffold with the selected architecture
4. Automatically run `go mod tidy` to download dependencies

## Available Templates

### 1. Hexagonal Architecture (`hexagonal`)
**Description**: Hexagonal Architecture (Ports & Adapters) with Uber FX and Chi Router

**Structure**:
```
.
├── cmd/server/main.go                    # Application entry point
├── internal/                             # Inner Hexagon (Domain, Application, Ports)
│   ├── domain/                           # Pure Domain Models
│   ├── application/                      # Application Services
│   └── ports/                            # Ports (Inbound & Outbound)
├── adapters/                             # Outer Hexagon (Adapters)
│   ├── inbound/http/                     # HTTP handlers with Chi router
│   └── outbound/persistence/             # Repository implementation
├── initiators/                           # Dependency Injection & Lifecycle
├── go.mod
├── go.sum
└── README.md
```

**Features**:
- **Hexagonal Architecture**: Strict separation between domain, application, and infrastructure
- **Chi Router**: Modern HTTP routing with middleware support
- **Uber FX**: Dependency injection and lifecycle management
- **Zap Logger**: Structured logging with production-ready configuration
- **In-memory persistence**: Simple in-memory storage for quick development
- **Clean architecture**: Strict separation of concerns
- **Ready to run**: Compiles and runs immediately with automatic dependency management

### 2. Clean Architecture (`clean`)
**Description**: Clean Architecture with Domain-Driven Design (DDD) principles

**Structure**:
```
.
├── cmd/server/main.go                    # Application entry point
├── internal/                             # Internal application layers
│   ├── domain/                           # Domain layer (entities & services)
│   │   ├── entity/                       # Domain entities
│   │   └── service/                      # Domain services
│   ├── storage/                          # Data access layer
│   │   ├── interfaces/                   # Repository interfaces
│   │   └── mongo/                        # MongoDB implementations
│   ├── handler/                          # HTTP handlers
│   │   ├── rest/                         # REST API handlers
│   │   │   ├── dto/                      # Data Transfer Objects
│   │   │   ├── http/                     # HTTP handlers
│   │   │   └── mapper/                   # Entity-DTO mappers
│   │   └── middleware/                   # HTTP middleware
│   └── glue/                             # Application glue
│       └── routing/                      # Route definitions
├── initiator/                            # Dependency injection
├── platform/                             # Platform utilities
│   ├── utils/                            # Utility functions
│   └── mongo/                            # MongoDB utilities
├── go.mod
├── go.sum
└── README.md
```

**Features**:
- **Clean Architecture**: Domain-Driven Design with clear layer separation
- **MongoDB Integration**: Production-ready MongoDB repository implementation
- **DTO Pattern**: Clean data transfer objects with validation
- **Mapper Pattern**: Entity-DTO mapping for clean API responses
- **Middleware Support**: Extensible middleware architecture
- **Structured Logging**: Production-ready logging with Zap
- **Dependency Injection**: Uber FX for clean dependency management

## Architecture Benefits

- **Testability**: Easy to unit test domain logic in isolation
- **Flexibility**: Swap implementations without changing core logic
- **Maintainability**: Clear separation of concerns
- **Scalability**: Modular design supports team growth

## Development

### Building from source

```bash
git clone <repository>
cd small-go
go build -o small-go .
```

### Running locally

```bash
go run main.go new my-project
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.