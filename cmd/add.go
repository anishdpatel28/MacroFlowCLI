package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name> <command>",
	Short: "Add a new macro to the current project",
	Long: `Add a new macro to the current project.

The command can include parameters:
  $1, $2, $3... - Individual positional parameters
  $@            - All parameters as a single string

Quotes are important for commands with spaces or special characters.`,
	Example: `  # Simple alias
  macro add dev "npm run dev"
  macro add test "npm test"

  # With single parameter
  macro add goto "cd $1"
  macro add kill "kill -9 $1"

  # With multiple parameters
  macro add copy "cp $1 $2"
  
  # With all parameters
  macro add commit "git commit -m \"$@\""
  macro add echo "echo $@"

  # Complex commands
  macro add deploy "npm run build && npm run deploy"
  macro add fresh "rm -rf node_modules && npm install"`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// Find the project for this directory
		project := store.GetProjectByPath(cwd)
		if project == nil {
			return fmt.Errorf("no macro project found for current directory\nUse 'macro init' to create one")
		}

		macroName := args[0]
		command := strings.Join(args[1:], " ")

		macro, err := store.CreateMacro(project.ID, macroName, command)
		if err != nil {
			return err
		}

		fmt.Printf("✓ Created macro '%s' in project '%s'\n", macro.Name, project.Name)
		fmt.Printf("  Command: %s\n", macro.Command)
		fmt.Printf("  Macro ID: %s\n", macro.ID)
		fmt.Printf("\nUsage: macro %s\n", macroName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
