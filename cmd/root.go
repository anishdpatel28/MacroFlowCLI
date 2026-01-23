package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anishdpatel28/macroflow/internal/storage"
	"github.com/spf13/cobra"
)

var store *storage.Storage

var rootCmd = &cobra.Command{
	Use:   "macro",
	Short: "MacroFlow - Directory-scoped command aliases",
	Long: `MacroFlow allows you to create project-specific command aliases (macros) 
that work within specific directories and their subdirectories.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	var err error
	store, err = storage.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		cmdName := os.Args[1]

		knownCommands := map[string]bool{
			"init": true, "add": true, "list": true, "delete": true,
			"delete-name": true, "export": true, "import": true,
			"help": true, "--help": true, "-h": true,
		}

		if !knownCommands[cmdName] {
			macroErr := executeMacro(cmdName, os.Args[2:])
			if macroErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", macroErr)
				os.Exit(1)
			}
			return
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// executeMacro executes a macro by name
func executeMacro(name string, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	project := store.GetProjectByPath(cwd)
	if project == nil {
		return fmt.Errorf("no macro project found for current directory\nUse 'macro init' to create one")
	}

	macro := store.GetMacroByName(project.ID, name)
	if macro == nil {
		return fmt.Errorf("macro '%s' not found in project '%s'\nUse 'macro list' to see available macros", name, project.Name)
	}

	command := macro.Command
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		command = strings.ReplaceAll(command, placeholder, arg)
	}
	command = strings.ReplaceAll(command, "$@", strings.Join(args, " "))

	// Shell wrapper will execute this command in the current shell
	fmt.Printf("%% %s\n", command)

	return nil
}
