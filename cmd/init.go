package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new macro project in the current directory",
	Long: `Initialize a new macro project in the current directory.
All macros created in this project will be available in this directory and its subdirectories.

The project name is optional - if not provided, the directory name will be used.`,
	Example: `  # Initialize with directory name
  macro init

  # Initialize with custom name
  macro init my-project

  # Initialize in a specific project
  cd ~/projects/website && macro init website`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Use directory name as project name if not provided
		name := filepath.Base(cwd)
		if len(args) > 0 {
			name = args[0]
		}

		project, err := store.CreateProject(name, cwd)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created macro project '%s' at %s\n", project.Name, project.Path)
		fmt.Printf("  Project ID: %s\n", project.ID)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  - Add macros: macro add <name> <command>\n")
		fmt.Printf("  - List macros: macro list\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
