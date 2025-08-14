package templates

// HexagonalTemplate represents the hexagonal architecture template
type HexagonalTemplate struct{}

func (h *HexagonalTemplate) Name() string {
	return "hexagonal"
}

func (h *HexagonalTemplate) Description() string {
	return "Hexagonal Architecture (Ports & Adapters) with Uber FX and Chi Router"
}

func (h *HexagonalTemplate) GenerateFiles(projectName string) map[string]string {
	return map[string]string{
		"cmd/server/main.go":                               generateMainGo(projectName),
		"internal/domain/user.go":                          generateDomainUser(),
		"internal/application/user_service.go":             generateApplicationUserService(projectName),
		"internal/ports/inbound/user_service.go":           generateInboundUserService(projectName),
		"internal/ports/outbound/user_repository.go":       generateOutboundUserRepository(projectName),
		"adapters/inbound/http/user_handler.go":            generateHTTPUserHandler(projectName),
		"adapters/inbound/http/router.go":                  generateHTTPRouter(projectName),
		"adapters/outbound/persistence/user_repository.go": generateUserRepository(projectName),
		"initiators/app.go":                                generateAppInitiator(),
		"initiators/http.go":                               generateHTTPInitiator(projectName),
		"initiators/persistence.go":                        generatePersistenceInitiator(projectName),
		"README.md":                                        generateREADME(projectName, "hexagonal"),
	}
}

func (h *HexagonalTemplate) GetDependencies() []string {
	return []string{
		"github.com/go-chi/chi/v5",
		"go.uber.org/fx",
		"go.uber.org/zap",
	}
}
