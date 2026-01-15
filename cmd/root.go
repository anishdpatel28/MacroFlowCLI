package cmd

import (
	"fmt"
	"os"
	"os/exec"
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

	// Check if first arg is a known command
	if len(os.Args) > 1 {
		cmdName := os.Args[1]

		// List of known commands
		knownCommands := map[string]bool{
			"init": true, "add": true, "list": true, "delete": true,
			"delete-name": true, "export": true, "import": true,
			"help": true, "--help": true, "-h": true,
		}

		// If not a known command, try as macro
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
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find the project for this directory
	project := store.GetProjectByPath(cwd)
	if project == nil {
		return fmt.Errorf("no macro project found for current directory\nUse 'macro init' to create one")
	}

	// Find the macro
	macro := store.GetMacroByName(project.ID, name)
	if macro == nil {
		return fmt.Errorf("macro '%s' not found in project '%s'\nUse 'macro list' to see available macros", name, project.Name)
	}

	// Replace parameters in command
	command := macro.Command
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		command = strings.ReplaceAll(command, placeholder, arg)
	}

	// Also support $@ for all arguments
	command = strings.ReplaceAll(command, "$@", strings.Join(args, " "))

	// Execute the command
	fmt.Printf("Executing: %s\n", command)

	// Use shell to execute the command
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	cmd := exec.Command(shell, "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cwd

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}
