package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// createProject creates a new Go project with simplified Hexagonal Architecture
func createProject(projectName string) error {
	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Change to project directory
	if err := os.Chdir(projectName); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Initialize Go module
	if err := runGoModInit(projectName); err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}

	// Create simplified directory structure
	dirs := []string{
		"cmd/server",
		"internal/domain",
		"internal/application",
		"internal/ports/inbound",
		"internal/ports/outbound",
		"adapters/inbound/http",
		"adapters/outbound/persistence",
		"initiators",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate core files
	if err := generateCoreFiles(projectName); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Run go mod tidy
	if err := runGoModTidy(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

// runGoModInit initializes a new Go module
func runGoModInit(projectName string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runGoModTidy runs go mod tidy to download dependencies
func runGoModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// generateCoreFiles creates the essential boilerplate files
func generateCoreFiles(projectName string) error {
	files := map[string]string{
		"cmd/server/main.go":                               generateMainGo(projectName),
		"internal/domain/user.go":                          generateDomainUser(),
		"internal/application/user_service.go":             generateApplicationUserService(projectName),
		"internal/ports/inbound/user_service.go":           generateInboundUserService(projectName),
		"internal/ports/outbound/user_repository.go":       generateOutboundUserRepository(projectName),
		"adapters/inbound/http/user_handler.go":            generateHTTPUserHandler(projectName),
		"adapters/inbound/http/router.go":                  generateHTTPRouter(projectName),
		"adapters/outbound/persistence/user_repository.go": generateUserRepository(projectName),
		"initiators/app.go":                                generateAppInitiator(projectName),
		"initiators/http.go":                               generateHTTPInitiator(projectName),
		"initiators/persistence.go":                        generatePersistenceInitiator(projectName),
		"README.md":                                        generateREADME(projectName),
	}

	for filePath, content := range files {
		if err := writeFile(filePath, content); err != nil {
			return fmt.Errorf("failed to write %s: %w", filePath, err)
		}
	}

	return nil
}

// writeFile writes content to a file
func writeFile(filePath, content string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(filePath, []byte(content), 0644)
}

// generateMainGo creates the main.go file
func generateMainGo(projectName string) string {
	return fmt.Sprintf(`package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"%s/initiators"
)

func main() {
	app := fx.New(
		fx.Provide(
			initiators.NewLogger,
			initiators.NewUserRepository,
			initiators.NewUserService,
			initiators.NewHTTPHandler,
		),
		fx.Invoke(initiators.StartServer),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return fxevent.NopLogger
		}),
	)

	app.Run()
}
`, projectName)
}

// generateDomainUser creates the domain user model
func generateDomainUser() string {
	return `package domain

import (
	"time"
)

// User represents a user entity in the domain
type User struct {
	ID        string    ` + "`json:\"id\"`" + `
	Email     string    ` + "`json:\"email\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

// NewUser creates a new user instance
func NewUser(email, name string) *User {
	now := time.Now()
	return &User{
		Email:     email,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}
`
}

// generateApplicationUserService creates the application service
func generateApplicationUserService(projectName string) string {
	return fmt.Sprintf(`package application

import (
	"context"
	"fmt"

	"%s/internal/domain"
	"%s/internal/ports/inbound"
	"%s/internal/ports/outbound"
)

// UserService implements the user application service
type UserService struct {
	userRepo outbound.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo outbound.UserRepository) inbound.UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	user := domain.NewUser(email, name)
	
	// Save to repository
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %%w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %%w", err)
	}

	return user, nil
}
`, projectName, projectName, projectName)
}

// generateInboundUserService creates the inbound port interface
func generateInboundUserService(projectName string) string {
	return fmt.Sprintf(`package inbound

import (
	"context"

	"%s/internal/domain"
)

// UserService defines the inbound port for user operations
type UserService interface {
	CreateUser(ctx context.Context, email, name string) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
}
`, projectName)
}

// generateOutboundUserRepository creates the outbound repository port
func generateOutboundUserRepository(projectName string) string {
	return fmt.Sprintf(`package outbound

import (
	"context"

	"%s/internal/domain"
)

// UserRepository defines the outbound port for user persistence
type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}
`, projectName)
}

// generateHTTPUserHandler creates the HTTP user handler
func generateHTTPUserHandler(projectName string) string {
	return fmt.Sprintf(`package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"%s/internal/ports/inbound"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService inbound.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService inbound.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email string `+"`json:\"email\"`"+`
	Name  string `+"`json:\"name\"`"+`
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
`, projectName)
}

// generateHTTPRouter creates the HTTP router
func generateHTTPRouter(projectName string) string {
	return fmt.Sprintf(`package http

import (
	"net/http"

	"%s/internal/ports/inbound"
)

// Router sets up HTTP routes
func NewRouter(userService inbound.UserService) http.Handler {
	mux := http.NewServeMux()
	
	// Initialize handlers
	userHandler := NewUserHandler(userService)

	// User routes
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users", userHandler.GetUser)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`+"`{\"status\":\"ok\"}`"+`))
	})

	return mux
}
`, projectName)
}

// generateUserRepository creates the user repository adapter
func generateUserRepository(projectName string) string {
	return fmt.Sprintf(`package persistence

import (
	"context"
	"fmt"

	"%s/internal/domain"
	"%s/internal/ports/outbound"
)

// UserRepository implements UserRepository using in-memory storage
type UserRepository struct {
	users map[string]*domain.User
}

// NewUserRepository creates a new user repository
func NewUserRepository() outbound.UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

// Save saves a user to storage
func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	// Simple ID generation (in a real app, use UUID)
	if user.ID == "" {
		user.ID = fmt.Sprintf("user_%%d", len(r.users)+1)
	}
	
	r.users[user.ID] = user
	fmt.Printf("Saved user: %%s\n", user.Email)
	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID] = user
	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(r.users, id)
	return nil
}
`, projectName, projectName)
}

// generateREADME creates the README.md
func generateREADME(projectName string) string {
	return fmt.Sprintf(`# %s

A Go service built with Hexagonal Architecture (Ports & Adapters) pattern.

## Project Structure

This project follows the Hexagonal Architecture pattern with clear separation of concerns:

- **Domain**: Pure business logic and entities
- **Application**: Use case implementations and orchestration
- **Ports**: Interfaces defining contracts (inbound/outbound)
- **Adapters**: External implementations (HTTP, databases, etc.)

## Quick Start

### Prerequisites

- Go 1.21 or later

### Running the Service

1. **Navigate to the project:**
   `+"```bash"+`
   cd %s
   `+"```"+`

2. **Install dependencies:**
   `+"```bash"+`
   go mod tidy
   `+"```"+`

3. **Run the service:**
   `+"```bash"+`
   go run cmd/server/main.go
   `+"```"+`

The service will be available at `+"`http://localhost:8080`"+`

## API Endpoints

- `+"`GET /health`"+` - Health check
- `+"`POST /users`"+` - Create a new user
- `+"`GET /users?id=<user_id>`"+` - Get user by ID

## Architecture Overview

`+"```"+`
.
├── cmd/server/main.go          # Application entry point
├── internal/                   # Inner Hexagon (Domain, Application, Ports)
│   ├── domain/                 # Pure Domain Models
│   ├── application/            # Application Services
│   └── ports/                  # Ports (Inbound & Outbound)
├── adapters/                   # Outer Hexagon (Adapters)
│   ├── inbound/                # Driving Adapters (HTTP)
│   └── outbound/               # Driven Adapters (Databases)
└── README.md                   # Project documentation
`+"```"+`

## Testing

`+"```bash"+`
go test ./...
`+"```"+`

## Building

`+"```bash"+`
go build -o bin/server cmd/server/main.go
`+"```"+`

## Contributing

1. Follow the Hexagonal Architecture pattern
2. Add tests for new features
3. Update documentation as needed
4. Ensure all tests pass before submitting

## License

This project is licensed under the MIT License.
`, projectName, projectName)
}

// generateAppInitiator creates the main app initiator
func generateAppInitiator(projectName string) string {
	return fmt.Sprintf(`package initiators

import (
	"context"
	"net/http"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// StartServer starts the HTTP server
func StartServer(lifecycle fx.Lifecycle, logger *zap.Logger, handler http.Handler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("Starting HTTP server", zap.String("port", port))
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("Server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}
`)
}

// generateHTTPInitiator creates the HTTP initiator
func generateHTTPInitiator(projectName string) string {
	return fmt.Sprintf(`package initiators

import (
	"net/http"

	httphandler "%s/adapters/inbound/http"
	"%s/internal/ports/inbound"
)

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(userService inbound.UserService) http.Handler {
	return httphandler.NewRouter(userService)
}
`, projectName, projectName)
}

// generatePersistenceInitiator creates the persistence initiator
func generatePersistenceInitiator(projectName string) string {
	return fmt.Sprintf(`package initiators

import (
	"go.uber.org/zap"

	"%s/adapters/outbound/persistence"
	"%s/internal/application"
	"%s/internal/ports/inbound"
	"%s/internal/ports/outbound"
)

// NewUserRepository creates a new user repository
func NewUserRepository() outbound.UserRepository {
	return persistence.NewUserRepository()
}

// NewUserService creates a new user service
func NewUserService(userRepo outbound.UserRepository) inbound.UserService {
	return application.NewUserService(userRepo)
}

// NewLogger creates a new logger
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
`, projectName, projectName, projectName, projectName)
}
