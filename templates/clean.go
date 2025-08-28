package templates

// CleanTemplate represents the clean architecture template
type CleanTemplate struct{}

func (c *CleanTemplate) Name() string {
	return "clean"
}

func (c *CleanTemplate) Description() string {
	return "Clean Architecture with Domain-Driven Design (DDD) principles - Basic structure"
}

func (c *CleanTemplate) GenerateFiles(projectName string) map[string]string {
	return map[string]string{
		// Application entry points
		"cmd/api/main.go":             generateCleanAPIMain(projectName),
		"cmd/api/app/container.go":    generateCleanBasicAPIContainer(projectName),
		"cmd/api/app/config.go":       generateCleanAPIConfig(),
		"cmd/api/app/server.go":       generateCleanAPIServer(projectName),
		"cmd/api/app/modules/user.go": generateCleanUserModule(projectName),

		// Core domain layer
		"internal/core/entities/user.go":                            generateCleanUserEntity(),
		"internal/core/services/user_service.go":                    generateCleanUserService(projectName),
		"internal/core/interfaces/repositories/user_repository.go":  generateCleanUserRepositoryInterface(projectName),
		"internal/core/interfaces/services/notification_service.go": generateCleanNotificationInterface(),
		"internal/core/errors/errors.go":                            generateCleanDomainErrors(),

		// Adapters layer - Database (placeholder for future implementations)
		"internal/adapters/database/inmemory/user_repository.go": generateCleanInMemoryUserRepository(projectName),

		// Adapters layer - External services
		"internal/adapters/external/notification/service.go": generateCleanNotificationService(),

		// Delivery layer - HTTP
		"internal/delivery/http/v1/dto/user.go":              generateCleanUserDTO(projectName),
		"internal/delivery/http/v1/mappers/user.go":          generateCleanUserMapper(projectName),
		"internal/delivery/http/v1/handlers/user_handler.go": generateCleanUserHandler(projectName),
		"internal/delivery/http/v1/router.go":                generateCleanHTTPRouter(projectName),
		"internal/delivery/http/middleware/auth.go":          generateCleanAuthMiddleware(),
		"internal/delivery/http/middleware/logging.go":       generateCleanLoggingMiddleware(),
		"internal/delivery/http/middleware/error_handler.go": generateCleanErrorHandler(),

		// Shared utilities
		"internal/shared/config/config.go": generateCleanSharedConfig(),
		"internal/shared/logger/logger.go": generateCleanSharedLogger(),
		"internal/shared/utils/id.go":      generateCleanIDUtils(),
		"internal/shared/utils/time.go":    generateCleanTimeUtils(),

		// Public packages
		"pkg/user/client.go":       generateCleanUserClient(projectName),
		"pkg/user/types.go":        generateCleanUserTypes(),
		"pkg/common/pagination.go": generateCleanPagination(),
		"pkg/common/response.go":   generateCleanResponse(),

		// Project files
		"go.mod":     generateGoMod(projectName),
		"Makefile":   generateMakefile(),
		"README.md":  generateCleanREADME(projectName),
		".gitignore": generateGitignore(),
	}
}

func (c *CleanTemplate) GetDependencies() []string {
	return []string{
		"github.com/go-chi/chi/v5",
		"go.uber.org/fx",
		"go.uber.org/zap",
	}
}

// CleanMongoTemplate represents clean architecture with MongoDB
type CleanMongoTemplate struct{}

func (cm *CleanMongoTemplate) Name() string {
	return "clean-mongo"
}

func (cm *CleanMongoTemplate) Description() string {
	return "Clean Architecture with MongoDB database implementation"
}

func (cm *CleanMongoTemplate) GenerateFiles(projectName string) map[string]string {
	files := map[string]string{
		// Application entry points
		"cmd/api/main.go":             generateCleanAPIMain(projectName),
		"cmd/api/app/container.go":    generateCleanMongoAPIContainer(projectName),
		"cmd/api/app/config.go":       generateCleanMongoAPIConfig(),
		"cmd/api/app/server.go":       generateCleanAPIServer(projectName),
		"cmd/api/app/modules/user.go": generateCleanMongoUserModule(projectName),

		// Core domain layer
		"internal/core/entities/user.go":                            generateCleanUserEntity(),
		"internal/core/services/user_service.go":                    generateCleanUserService(projectName),
		"internal/core/interfaces/repositories/user_repository.go":  generateCleanUserRepositoryInterface(projectName),
		"internal/core/interfaces/services/notification_service.go": generateCleanNotificationInterface(),
		"internal/core/errors/errors.go":                            generateCleanDomainErrors(),

		// Adapters layer - MongoDB
		"internal/adapters/database/mongo/models/user.go":                  generateCleanMongoUserModel(),
		"internal/adapters/database/mongo/mappers/user.go":                 generateCleanMongoUserMapper(projectName),
		"internal/adapters/database/mongo/repositories/user_repository.go": generateCleanMongoUserRepository(projectName),
		"internal/adapters/database/mongo/connection.go":                   generateCleanMongoConnection(),

		// Adapters layer - External services
		"internal/adapters/external/notification/service.go": generateCleanNotificationService(),

		// Delivery layer - HTTP
		"internal/delivery/http/v1/dto/user.go":              generateCleanUserDTO(projectName),
		"internal/delivery/http/v1/mappers/user.go":          generateCleanUserMapper(projectName),
		"internal/delivery/http/v1/handlers/user_handler.go": generateCleanUserHandler(projectName),
		"internal/delivery/http/v1/router.go":                generateCleanHTTPRouter(projectName),
		"internal/delivery/http/middleware/auth.go":          generateCleanAuthMiddleware(),
		"internal/delivery/http/middleware/logging.go":       generateCleanLoggingMiddleware(),
		"internal/delivery/http/middleware/error_handler.go": generateCleanErrorHandler(),

		// Shared utilities
		"internal/shared/config/config.go": generateCleanSharedConfig(),
		"internal/shared/logger/logger.go": generateCleanSharedLogger(),
		"internal/shared/utils/id.go":      generateCleanIDUtils(),
		"internal/shared/utils/time.go":    generateCleanTimeUtils(),

		// Public packages
		"pkg/user/client.go":       generateCleanUserClient(projectName),
		"pkg/user/types.go":        generateCleanUserTypes(),
		"pkg/common/pagination.go": generateCleanPagination(),
		"pkg/common/response.go":   generateCleanResponse(),

		// Project files
		"go.mod":     generateCleanMongoGoMod(projectName),
		"Makefile":   generateMakefile(),
		"README.md":  generateCleanMongoREADME(projectName),
		".gitignore": generateGitignore(),
	}

	return files
}

func (cm *CleanMongoTemplate) GetDependencies() []string {
	return []string{
		"github.com/go-chi/chi/v5",
		"go.uber.org/fx",
		"go.uber.org/zap",
		"go.mongodb.org/mongo-driver/mongo",
		"go.mongodb.org/mongo-driver/bson",
	}
}
