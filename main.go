package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dawit-go/small-go/templates"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "small-go",
		Short: "Small-Go: A CLI utility for generating Go project scaffolds with different architecture patterns",
		Long: `Small-Go is a CLI utility for generating ready-to-use Go project scaffolds 
with different architecture patterns including Hexagonal Architecture and Clean Architecture.

It standardizes Go project layouts and encourages separation of concerns between 
Domain, Application, Ports, and Adapters.`,
	}

	var newCmd = &cobra.Command{
		Use:   "new [project_name]",
		Short: "Create a new Go project with selected architecture pattern",
		Long: `Create a new Go project with selected architecture pattern.
This will generate a complete project scaffold with:
- Domain models and entities
- Application services and use cases
- HTTP handlers and routing
- Infrastructure configurations
- Dependency injection setup
- Documentation`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			templateName, _ := cmd.Flags().GetString("template")

			// If no template specified, show interactive selection
			if templateName == "" {
				templateName = selectTemplate()
			}

			if err := createProject(projectName, templateName); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Successfully created project: %s\n", projectName)
			fmt.Printf("📁 Navigate to the project: cd %s\n", projectName)
			fmt.Printf("🚀 Run the service: go run cmd/server/main.go\n")
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List available architecture templates",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available templates:")
			fmt.Println()
			for i, template := range templates.GetAvailableTemplates() {
				fmt.Printf("  %d. %s: %s\n", i+1, template.Name(), template.Description())
			}
		},
	}

	// Add template flag
	newCmd.Flags().StringP("template", "t", "", "Architecture template to use (hexagonal, clean)")

	rootCmd.AddCommand(newCmd, listCmd)
	rootCmd.Execute()
}

// selectTemplate provides interactive template selection
func selectTemplate() string {
	availableTemplates := templates.GetAvailableTemplates()

	fmt.Println("Select an architecture template:")
	fmt.Println()
	for i, template := range availableTemplates {
		fmt.Printf("  %d. %s: %s\n", i+1, template.Name(), template.Description())
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter your choice (1-2): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		index, err := strconv.Atoi(choice)
		if err != nil || index < 1 || index > len(availableTemplates) {
			fmt.Printf("Please enter a number between 1 and %d\n", len(availableTemplates))
			continue
		}

		selectedTemplate := availableTemplates[index-1]
		fmt.Printf("Selected: %s\n", selectedTemplate.Name())
		return selectedTemplate.Name()
	}
}
