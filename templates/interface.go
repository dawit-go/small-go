package templates

// Template defines the interface for project templates
type Template interface {
	Name() string
	Description() string
	GenerateFiles(projectName string) map[string]string
	GetDependencies() []string
}

// GetAvailableTemplates returns all available templates
func GetAvailableTemplates() []Template {
	return []Template{
		&HexagonalTemplate{},
		&CleanTemplate{},
	}
}

// GetTemplateByName returns a template by name
func GetTemplateByName(name string) Template {
	for _, template := range GetAvailableTemplates() {
		if template.Name() == name {
			return template
		}
	}
	return nil
}
