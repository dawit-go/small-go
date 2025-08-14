package templates

// CleanTemplate represents the clean architecture template
type CleanTemplate struct{}

func (c *CleanTemplate) Name() string {
	return "clean"
}

func (c *CleanTemplate) Description() string {
	return "Clean Architecture with Domain-Driven Design (DDD) principles"
}

func (c *CleanTemplate) GenerateFiles(projectName string) map[string]string {
	return map[string]string{
		"cmd/server/main.go":                                    generateCleanMainGo(projectName),
		"internal/domain/entity/user.go":                        generateCleanDomainEntity(),
		"internal/domain/service/user_service.go":               generateCleanDomainService(projectName),
		"internal/storage/interfaces/user_repository.go":        generateCleanStorageInterface(projectName),
		"internal/storage/mongo/user_repository.go":             generateCleanMongoRepository(projectName),
		"internal/handler/rest/dto/user_dto.go":                 generateCleanUserDTO(projectName),
		"internal/handler/rest/http/user_handler.go":            generateCleanUserHandler(projectName),
		"internal/handler/rest/mapper/user_mapper.go":           generateCleanUserMapper(projectName),
		"internal/handler/middleware/auth.go":                   generateCleanAuthMiddleware(),
		"internal/glue/routing/routes.go":                       generateCleanRoutes(projectName),
		"initiator/initiator.go":                                generateCleanInitiator(projectName),
		"initiator/service.go":                                  generateCleanServiceInitiator(projectName),
		"initiator/persistence.go":                              generateCleanPersistenceInitiator(projectName),
		"initiator/handler.go":                                  generateCleanHandlerInitiator(projectName),
		"initiator/config.go":                                   generateCleanConfigInitiator(),
		"initiator/logger.go":                                   generateCleanLoggerInitiator(),
		"platform/utils/response.go":                            generateCleanResponseUtils(),
		"platform/mongo/connection.go":                          generateCleanMongoConnection(),
		"README.md":                                             generateREADME(projectName, "clean"),
	}
}

func (c *CleanTemplate) GetDependencies() []string {
	return []string{
		"github.com/go-chi/chi/v5",
		"go.uber.org/fx",
		"go.uber.org/zap",
		"go.mongodb.org/mongo-driver/mongo",
		"go.mongodb.org/mongo-driver/bson",
	}
} 