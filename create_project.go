package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dawit-go/small-go/templates"
)

// createProject creates a new Go project with the selected template
func createProject(projectName, templateName string) error {
	// Get the selected template
	template := templates.GetTemplateByName(templateName)
	if template == nil {
		return fmt.Errorf("unknown template: %s. Use 'small-go list' to see available templates", templateName)
	}

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

	// Generate files using the selected template
	if err := generateTemplateFiles(projectName, template); err != nil {
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

// generateTemplateFiles generates files using the selected template
func generateTemplateFiles(projectName string, template templates.Template) error {
	files := template.GenerateFiles(projectName)

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
