# Small-Go

A CLI utility for generating Go project scaffolds with different architectural patterns.

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

#### Interactive Selection
```bash
small-go new <project-name>
```

#### Direct Template Selection
```bash
small-go new <project-name> --template <template-name>
```

## Available Templates

### 1. Hexagonal Architecture
**Template Name**: `hexagonal`

A classic hexagonal architecture with:
- Domain-driven design
- Ports and adapters pattern
- Chi router for HTTP handling
- Uber FX for dependency injection
- In-memory storage for quick development

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
└── README.md
```

### 2. Clean Architecture (Basic)
**Template Name**: `clean`

Clean architecture with domain-driven design principles:
- Pure business logic in core layer
- Infrastructure adapters
- HTTP delivery layer
- Shared utilities
- Public packages for external consumption

**Structure**:
```
.
├── cmd/api/                              # Application entry point
│   ├── main.go                          # Main application
│   └── app/                             # FX application wiring
│       ├── container.go                 # Dependency injection
│       ├── config.go                    # Configuration
│       ├── server.go                    # Server setup
│       └── modules/                     # Domain modules
├── internal/                            # Internal application code
│   ├── core/                            # Business logic and domain
│   │   ├── entities/                    # Domain entities
│   │   ├── services/                    # Business services
│   │   ├── interfaces/                  # Repository interfaces
│   │   └── errors/                      # Domain errors
│   ├── adapters/                        # Infrastructure implementations
│   │   ├── database/                    # Database implementations
│   │   └── external/                    # External service integrations
│   ├── delivery/                        # API layer
│   │   └── http/                        # HTTP API
│   └── shared/                          # Shared utilities
├── pkg/                                 # Public packages
└── go.mod
```

### 3. Clean Architecture with MongoDB
**Template Name**: `clean-mongo`

Clean architecture with MongoDB database implementation:
- All features of basic clean architecture
- MongoDB integration with proper models and mappers
- Production-ready database layer
- Environment-based configuration

**Structure**:
```
.
├── cmd/api/                              # Application entry point
│   ├── main.go                          # Main application
│   └── app/                             # FX application wiring
│       ├── container.go                 # Dependency injection
│       ├── config.go                    # Configuration
│       ├── server.go                    # Server setup
│       └── modules/                     # Domain modules
├── internal/                            # Internal application code
│   ├── core/                            # Business logic and domain
│   │   ├── entities/                    # Domain entities
│   │   ├── services/                    # Business services
│   │   ├── interfaces/                  # Repository interfaces
│   │   └── errors/                      # Domain errors
│   ├── adapters/                        # Infrastructure implementations
│   │   ├── database/mongo/              # MongoDB implementations
│   │   └── external/                    # External service integrations
│   ├── delivery/                        # API layer
│   │   └── http/                        # HTTP API
│   └── shared/                          # Shared utilities
├── pkg/                                 # Public packages
└── go.mod
```

## Features

### Common Features Across All Templates
- **Chi Router**: Modern HTTP routing with middleware support
- **Uber FX**: Dependency injection and lifecycle management
- **Zap Logger**: Structured logging with production-ready configuration
- **Ready to Run**: Compiles and runs immediately with automatic dependency management
- **Clean Architecture**: Strict separation of concerns
- **Modular Design**: Easy to extend and maintain

### Template-Specific Features

#### Hexagonal Architecture
- **In-memory persistence**: Simple in-memory storage for quick development
- **Ports and adapters**: Clear interface definitions
- **Domain isolation**: Pure business logic separation

#### Clean Architecture (Basic)
- **Domain-driven design**: Pure business logic in core layer
- **Infrastructure adapters**: Pluggable database implementations
- **HTTP delivery layer**: RESTful API with DTOs
- **Public packages**: Reusable components for external services

#### Clean Architecture with MongoDB
- **MongoDB integration**: Production-ready database implementation
- **Data mapping**: Clean separation between domain and data models
- **Environment configuration**: Flexible configuration management
- **Connection management**: Proper database connection handling

## Quick Start Examples

### Create a Hexagonal Architecture Project
```bash
small-go new my-hexagonal-project --template hexagonal
cd my-hexagonal-project
go run cmd/server/main.go
```

### Create a Clean Architecture Project
```bash
small-go new my-clean-project --template clean
cd my-clean-project
go run cmd/api/main.go
```

### Create a Clean Architecture Project with MongoDB
```bash
small-go new my-mongo-project --template clean-mongo
cd my-mongo-project
# Start MongoDB (using Docker)
docker run -d -p 27017:27017 --name mongodb mongo:latest
# Run the application
go run cmd/api/main.go
```

## API Endpoints

All templates provide the following endpoints:
- `GET /health` - Health check
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user

## Architecture Benefits

- **Testability**: Easy to unit test domain logic in isolation
- **Flexibility**: Swap implementations without changing core logic
- **Maintainability**: Clear separation of concerns
- **Scalability**: Modular design supports team growth
- **Extensibility**: Easy to add new features and domains

## Contributing

1. Follow the architecture patterns
2. Add tests for new features
3. Update documentation as needed
4. Ensure all tests pass before submitting

## License

This project is licensed under the MIT License.