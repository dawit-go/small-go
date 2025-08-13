# Small-Go

A CLI utility for generating Go project scaffolds with Hexagonal Architecture (Ports & Adapters) pattern.

## Overview

Small-Go standardizes Go project layouts and encourages separation of concerns between **Domain**, **Application**, **Ports**, and **Adapters**. It makes it quick to bootstrap a new Go service with best practices built-in.

## Installation

```bash
go install github.com/my-org/small-go@latest
```

## Usage

Create a new Go project with Hexagonal Architecture:

```bash
small-go new <project_name>
```

This will:
1. Create a new folder named `<project_name>`
2. Initialize a Go module inside (`go mod init <project_name>`)
3. Generate a complete project scaffold with Hexagonal Architecture

## Generated Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/                       // Inner Hexagon — Domain, Application, Ports
│   ├── domain/                 // Pure Domain Models (Entities, VOs)
│   │   └── user.go
│   ├── application/            // Application Services (Use Case Implementations)
│   │   └── user_service.go
│   └── ports/                  // Ports (Inbound & Outbound split)
│       ├── inbound/            // Inbound Ports (Interfaces offered by the Core)
│       │   └── user_service.go
│       └── outbound/           // Outbound Ports (Core's expectations from outside)
│           └── user_repository.go
├── adapters/                   // Outer Hexagon — Adapters Implementing Ports
│   ├── inbound/                // Driving Adapters (Input into Core)
│   │   └── http/
│   │       ├── user_handler.go
│   │       └── router.go
│   └── outbound/               // Driven Adapters (Infra Implementations)
│       └── persistence/
│           └── user_repository.go
├── initiators/                 // Dependency Injection & Lifecycle Management
│   ├── app.go                  // Main application lifecycle
│   ├── http.go                 // HTTP handler initialization
│   └── persistence.go          // Repository & service initialization
├── go.mod
├── go.sum
└── README.md
```

## Features

- **Hexagonal Architecture**: Strict separation between domain, application, and infrastructure
- **Ready-to-use boilerplate**: Minimal but functional code that compiles immediately
- **Chi Router**: Modern HTTP routing with middleware support
- **Uber FX**: Dependency injection and lifecycle management
- **Zap Logger**: Structured logging with production-ready configuration
- **In-memory persistence**: Simple in-memory storage for quick development
- **Clean architecture**: Strict separation of concerns
- **Ready to run**: Compiles and runs immediately with automatic dependency management

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