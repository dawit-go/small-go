package templates

import "fmt"

// Hexagonal Architecture Generators

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

func generateHTTPRouter(projectName string) string {
	return fmt.Sprintf(`package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"%s/internal/ports/inbound"
)

// Router sets up HTTP routes using Chi
func NewRouter(userService inbound.UserService) http.Handler {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Initialize handlers
	userHandler := NewUserHandler(userService)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`+"`{\"status\":\"ok\"}`"+`))
	})

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUser)
	})

	return r
}
`, projectName)
}

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

func generateAppInitiator() string {
	return `package initiators

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
`
}

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

// Clean Architecture Generators

func generateCleanMainGo(projectName string) string {
	return fmt.Sprintf(`package main

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"%s/initiator"
)

func main() {
	app := fx.New(
		fx.Provide(
			initiator.NewLogger,
			initiator.NewConfig,
			initiator.NewMongoConnection,
			initiator.NewUserRepository,
			initiator.NewUserService,
			initiator.NewUserMapper,
			initiator.NewUserHandler,
			initiator.NewRoutes,
		),
		fx.Invoke(initiator.StartServer),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return fxevent.NopLogger
		}),
	)

	app.Run()
}
`, projectName)
}

func generateCleanDomainEntity() string {
	return `package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user entity in the domain
type User struct {
	ID        primitive.ObjectID ` + "`bson:\"_id,omitempty\" json:\"id\"`" + `
	Email     string             ` + "`bson:\"email\" json:\"email\"`" + `
	Name      string             ` + "`bson:\"name\" json:\"name\"`" + `
	CreatedAt time.Time          ` + "`bson:\"created_at\" json:\"created_at\"`" + `
	UpdatedAt time.Time          ` + "`bson:\"updated_at\" json:\"updated_at\"`" + `
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

func generateCleanDomainService(projectName string) string {
	return fmt.Sprintf(`package service

import (
	"context"
	"fmt"

	"%s/internal/domain/entity"
	"%s/internal/storage/interfaces"
)

// UserService implements the user domain service
type UserService struct {
	userRepo interfaces.UserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*entity.User, error) {
	user := entity.NewUser(email, name)
	
	// Save to repository
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %%w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %%w", err)
	}

	return user, nil
}
`, projectName, projectName)
}

func generateCleanStorageInterface(projectName string) string {
	return fmt.Sprintf(`package interfaces

import (
	"context"

	"%s/internal/domain/entity"
)

// UserRepository defines the repository interface for user persistence
type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
}
`, projectName)
}

func generateCleanMongoRepository(projectName string) string {
	return fmt.Sprintf(`package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"%s/internal/domain/entity"
	"%s/internal/storage/interfaces"
)

// UserRepository implements UserRepository using MongoDB
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new MongoDB user repository
func NewUserRepository(collection *mongo.Collection) interfaces.UserRepository {
	return &UserRepository{
		collection: collection,
	}
}

// Save saves a user to MongoDB
func (r *UserRepository) Save(ctx context.Context, user *entity.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// FindByID finds a user by ID in MongoDB
func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format")
	}

	var user entity.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// FindByEmail finds a user by email in MongoDB
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// Update updates a user in MongoDB
func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

// Delete deletes a user from MongoDB
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
`, projectName, projectName)
}

func generateCleanUserDTO(projectName string) string {
	return fmt.Sprintf(`package dto

import (
	"%s/internal/domain/entity"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email string `+"`json:\"email\" validate:\"required,email\"`"+`
	Name  string `+"`json:\"name\" validate:\"required\"`"+`
}

// UserResponse represents the user response
type UserResponse struct {
	ID        string `+"`json:\"id\"`"+`
	Email     string `+"`json:\"email\"`"+`
	Name      string `+"`json:\"name\"`"+`
	CreatedAt string `+"`json:\"created_at\"`"+`
	UpdatedAt string `+"`json:\"updated_at\"`"+`
}

// ToEntity converts CreateUserRequest to entity.User
func (req *CreateUserRequest) ToEntity() *entity.User {
	return entity.NewUser(req.Email, req.Name)
}
`, projectName)
}

func generateCleanUserHandler(projectName string) string {
	return fmt.Sprintf(`package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"%s/internal/domain/service"
	"%s/internal/handler/rest/dto"
	"%s/internal/handler/rest/mapper"
	"%s/platform/utils"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService *service.UserService
	userMapper  *mapper.UserMapper
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, userMapper *mapper.UserMapper) *UserHandler {
	return &UserHandler{
		userService: userService,
		userMapper:  userMapper,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := h.userMapper.ToResponse(user)
	utils.SendSuccessResponse(w, response, http.StatusCreated)
}

// GetUser handles GET /users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		utils.SendErrorResponse(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	response := h.userMapper.ToResponse(user)
	utils.SendSuccessResponse(w, response, http.StatusOK)
}
`, projectName, projectName, projectName, projectName)
}

func generateCleanUserMapper(projectName string) string {
	return fmt.Sprintf(`package mapper

import (
	"time"

	"%s/internal/domain/entity"
	"%s/internal/handler/rest/dto"
)

// UserMapper handles mapping between entities and DTOs
type UserMapper struct{}

// NewUserMapper creates a new user mapper
func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// ToResponse converts entity.User to dto.UserResponse
func (m *UserMapper) ToResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
`, projectName, projectName)
}

func generateCleanAuthMiddleware() string {
	return `package middleware

import (
	"net/http"
)

// AuthMiddleware provides basic authentication middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement authentication logic
		// For now, just pass through
		next.ServeHTTP(w, r)
	})
}
`
}

func generateCleanRoutes(projectName string) string {
	return fmt.Sprintf(`package routing

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	userhandler "%s/internal/handler/rest/http"
	authmiddleware "%s/internal/handler/middleware"
)

// Routes sets up all HTTP routes
func Routes(userHandler *userhandler.UserHandler) http.Handler {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(authmiddleware.AuthMiddleware)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`+"`{\"status\":\"ok\"}`"+`))
	})

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUser)
	})

	return r
}
`, projectName, projectName)
}

func generateCleanInitiator(projectName string) string {
	return `package initiator

import (
	"context"
	"net/http"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// StartServer starts the HTTP server
func StartServer(lifecycle fx.Lifecycle, logger *zap.Logger, routes http.Handler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: routes,
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
`
}

func generateCleanServiceInitiator(projectName string) string {
	return fmt.Sprintf(`package initiator

import (
	"%s/internal/domain/service"
	"%s/internal/storage/interfaces"
)

// NewUserService creates a new user service
func NewUserService(userRepo interfaces.UserRepository) *service.UserService {
	return service.NewUserService(userRepo)
}
`, projectName, projectName)
}

func generateCleanPersistenceInitiator(projectName string) string {
	return fmt.Sprintf(`package initiator

import (
	"%s/internal/storage/interfaces"
	mongorepo "%s/internal/storage/mongo"
	mongoplatform "%s/platform/mongo"
)

// NewUserRepository creates a new user repository
func NewUserRepository(connection *mongoplatform.Connection) interfaces.UserRepository {
	collection := connection.GetCollection("users")
	return mongorepo.NewUserRepository(collection)
}

// NewMongoConnection creates a new MongoDB connection
func NewMongoConnection(config *Config) (*mongoplatform.Connection, error) {
	return mongoplatform.NewConnection(config.MongoURI)
}
`, projectName, projectName, projectName)
}

func generateCleanHandlerInitiator(projectName string) string {
	return fmt.Sprintf(`package initiator

import (
	"net/http"

	"%s/internal/domain/service"
	userhandler "%s/internal/handler/rest/http"
	"%s/internal/handler/rest/mapper"
	"%s/internal/glue/routing"
)

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService, userMapper *mapper.UserMapper) *userhandler.UserHandler {
	return userhandler.NewUserHandler(userService, userMapper)
}

// NewUserMapper creates a new user mapper
func NewUserMapper() *mapper.UserMapper {
	return mapper.NewUserMapper()
}

// NewRoutes creates new routes
func NewRoutes(userHandler *userhandler.UserHandler) http.Handler {
	return routing.Routes(userHandler)
}
`, projectName, projectName, projectName, projectName)
}

func generateCleanConfigInitiator() string {
	return `package initiator

import (
	"os"
)

// Config represents application configuration
type Config struct {
	MongoURI string
	Port     string
}

// NewConfig creates a new configuration
func NewConfig() *Config {
	return &Config{
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Port:     getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`
}

func generateCleanLoggerInitiator() string {
	return `package initiator

import (
	"go.uber.org/zap"
)

// NewLogger creates a new logger
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}
`
}

func generateCleanResponseUtils() string {
	return `package utils

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        ` + "`json:\"success\"`" + `
	Message string      ` + "`json:\"message,omitempty\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
	Error   string      ` + "`json:\"error,omitempty\"`" + `
}

// SendSuccessResponse sends a success response
func SendSuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	response := Response{
		Success: true,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SendErrorResponse sends an error response
func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := Response{
		Success: false,
		Error:   message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
`
}

func generateCleanMongoConnection() string {
	return `package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection represents MongoDB connection
type Connection struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// NewConnection creates a new MongoDB connection
func NewConnection(mongoURI string) (*Connection, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database("myapp")

	return &Connection{
		Client: client,
		DB:     db,
	}, nil
}

// GetCollection returns a collection by name
func (c *Connection) GetCollection(name string) *mongo.Collection {
	return c.DB.Collection(name)
}
`
}

func generateREADME(projectName, templateType string) string {
	var structure string
	var features string

	switch templateType {
	case "clean":
		structure = `
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
└── README.md`
		features = `
- **Clean Architecture**: Domain-Driven Design with clear layer separation
- **MongoDB Integration**: Production-ready MongoDB repository implementation
- **DTO Pattern**: Clean data transfer objects with validation
- **Mapper Pattern**: Entity-DTO mapping for clean API responses
- **Middleware Support**: Extensible middleware architecture
- **Structured Logging**: Production-ready logging with Zap
- **Dependency Injection**: Uber FX for clean dependency management`
	default:
		structure = `
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
└── README.md`
		features = `
- **Hexagonal Architecture**: Strict separation between domain, application, and infrastructure
- **Chi Router**: Modern HTTP routing with middleware support
- **Uber FX**: Dependency injection and lifecycle management
- **Zap Logger**: Structured logging with production-ready configuration
- **In-memory persistence**: Simple in-memory storage for quick development
- **Clean architecture**: Strict separation of concerns
- **Ready to run**: Compiles and runs immediately with automatic dependency management`
	}

	return fmt.Sprintf(`# %s

A Go service built with %s.

## Project Structure

This project follows the %s pattern with clear separation of concerns:

%s

## Quick Start

### Prerequisites

- Go 1.21 or later
- MongoDB (for clean architecture template)

### Running the Service

1. **Navigate to the project:**
   `+"```bash"+`
   cd %s
   `+"```"+`

2. **Run the service (dependencies are automatically managed):**
   `+"```bash"+`
   go run cmd/server/main.go
   `+"```"+`

The service will be available at `+"`http://localhost:8080`"+`

## API Endpoints

- `+"`GET /health`"+` - Health check
- `+"`POST /users`"+` - Create a new user
- `+"`GET /users/{id}`"+` - Get user by ID

## Features

%s

## Architecture Benefits

- **Testability**: Easy to unit test domain logic in isolation
- **Flexibility**: Swap implementations without changing core logic
- **Maintainability**: Clear separation of concerns
- **Scalability**: Modular design supports team growth

## Testing

`+"```bash"+`
go test ./...
`+"```"+`

## Building

`+"```bash"+`
go build -o bin/server cmd/server/main.go
`+"```"+`

## Contributing

1. Follow the architecture pattern
2. Add tests for new features
3. Update documentation as needed
4. Ensure all tests pass before submitting

## License

This project is licensed under the MIT License.
`, projectName, templateType, templateType, structure, projectName, features)
}
