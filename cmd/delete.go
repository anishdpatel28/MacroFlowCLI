package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteProject bool

var deleteCmd = &cobra.Command{
	Use:   "delete <id> [id...]",
	Short: "Delete macro(s) or a project by ID",
	Long: `Delete one or more macros by their IDs, or delete a project with --project flag.

Examples:
  macro delete abc-123                    # Delete a single macro
  macro delete abc-123 def-456 ghi-789   # Delete multiple macros
  macro delete abc-123 --project         # Delete a project and all its macros`,
	Aliases: []string{"rm", "remove"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteProject {
			if len(args) > 1 {
				return fmt.Errorf("can only delete one project at a time")
			}
			return deleteProjectByID(args[0])
		}

		return deleteMacrosByIDs(args)
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteProject, "project", "p", false, "Delete a project instead of macros")
	rootCmd.AddCommand(deleteCmd)
}

func deleteMacrosByIDs(ids []string) error {
	if len(ids) == 1 {
		if err := store.DeleteMacro(ids[0]); err != nil {
			return err
		}
		fmt.Printf("✓ Deleted macro with ID: %s\n", ids[0])
	} else {
		if err := store.DeleteMacros(ids); err != nil {
			return err
		}
		fmt.Printf("✓ Deleted %d macros\n", len(ids))
	}

	return nil
}

func deleteProjectByID(id string) error {
	// Get project info before deleting
	projects := store.GetProjects()
	var projectName string
	var macroCount int

	for _, p := range projects {
		if p.ID == id {
			projectName = p.Name
			macros := store.GetMacrosByProject(p.ID)
			macroCount = len(macros)
			break
		}
	}

	if projectName == "" {
		return fmt.Errorf("project not found with ID: %s", id)
	}

	// Confirm deletion
	fmt.Printf("Are you sure you want to delete project '%s'?\n", projectName)
	fmt.Printf("This will also delete %d macro(s).\n", macroCount)
	fmt.Print("Type 'yes' to confirm: ")

	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "yes" {
		fmt.Println("Deletion cancelled.")
		return nil
	}

	if err := store.DeleteProject(id); err != nil {
		return err
	}

	fmt.Printf("✓ Deleted project '%s' and %d macro(s)\n", projectName, macroCount)
	return nil
}

// Helper command to delete by name in current project
var deleteNameCmd = &cobra.Command{
	Use:   "delete-name <macro-name>",
	Short: "Delete a macro by name in the current project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		project := store.GetProjectByPath(cwd)
		if project == nil {
			return fmt.Errorf("no macro project found for current directory")
		}

		macro := store.GetMacroByName(project.ID, args[0])
		if macro == nil {
			return fmt.Errorf("macro '%s' not found", args[0])
		}

		if err := store.DeleteMacro(macro.ID); err != nil {
			return err
		}

		fmt.Printf("✓ Deleted macro '%s'\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteNameCmd)
}
