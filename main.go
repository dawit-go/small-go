package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "small-go",
		Short: "Small-Go: A CLI utility for generating Go project scaffolds with Hexagonal Architecture",
		Long: `Small-Go is a CLI utility for generating ready-to-use Go project scaffolds 
that follow a strict Hexagonal Architecture (Ports & Adapters) structure.

It standardizes Go project layouts and encourages separation of concerns between 
Domain, Application, Ports, and Adapters.`,
	}

	var newCmd = &cobra.Command{
		Use:   "new [project_name]",
		Short: "Create a new Go project with Hexagonal Architecture",
		Long: `Create a new Go project with Hexagonal Architecture structure.
This will generate a complete project scaffold with:
- Domain models and entities
- Application services and use cases
- Inbound and outbound ports
- HTTP and gRPC adapters
- Infrastructure configurations
- Docker setup and documentation`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]
			if err := createProject(projectName); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Successfully created project: %s\n", projectName)
			fmt.Printf("📁 Navigate to the project: cd %s\n", projectName)
			fmt.Printf("🚀 Run the service: go run cmd/server/main.go\n")
		},
	}

	rootCmd.AddCommand(newCmd)
	rootCmd.Execute()
}
