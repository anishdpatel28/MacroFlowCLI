package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listAll      bool
	listProjects bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List macros and projects",
	Long: `List macros in the current project, all macros, or all projects.

Use flags to control what is listed:
  (no flags) - Macros in current project
  --all      - All macros across all projects
  --projects - All macro projects`,
	Example: `  # List macros in current project
  macro list

  # List all macros across all projects (grouped by project)
  macro list --all
  macro list -a

  # List all macro projects with details
  macro list --projects
  macro list -p`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if listProjects {
			return listAllProjects()
		}

		if listAll {
			return listAllMacros()
		}

		return listCurrentProjectMacros()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "List all macros across all projects")
	listCmd.Flags().BoolVarP(&listProjects, "projects", "p", false, "List all projects")
	rootCmd.AddCommand(listCmd)
}

func listCurrentProjectMacros() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	project := store.GetProjectByPath(cwd)
	if project == nil {
		return fmt.Errorf("no macro project found for current directory\nUse 'macro init' to create one")
	}

	macros := store.GetMacrosByProject(project.ID)

	fmt.Printf("Project: %s\n", project.Name)
	fmt.Printf("Path: %s\n", project.Path)
	fmt.Printf("\n")

	if len(macros) == 0 {
		fmt.Println("No macros found.")
		fmt.Println("\nCreate one with: macro add <name> <command>")
		return nil
	}

	// Sort macros by name
	sort.Slice(macros, func(i, j int) bool {
		return macros[i].Name < macros[j].Name
	})

	fmt.Printf("Macros (%d):\n", len(macros))
	fmt.Println(strings.Repeat("-", 80))

	for _, macro := range macros {
		fmt.Printf("  %-15s → %s\n", macro.Name, macro.Command)
		fmt.Printf("  %s ID: %s\n", strings.Repeat(" ", 15), macro.ID)
		fmt.Println()
	}

	return nil
}

func listAllMacros() error {
	projects := store.GetProjects()
	allMacros := store.GetAllMacros()

	if len(allMacros) == 0 {
		fmt.Println("No macros found.")
		return nil
	}

	// Create project lookup map
	projectMap := make(map[string]string)
	for _, p := range projects {
		projectMap[p.ID] = p.Name
	}

	// Group macros by project
	macrosByProject := make(map[string][]string)
	for _, m := range allMacros {
		projectName := projectMap[m.ProjectID]
		if projectName == "" {
			projectName = "Unknown Project"
		}
		info := fmt.Sprintf("  %-15s → %s (ID: %s)", m.Name, m.Command, m.ID)
		macrosByProject[projectName] = append(macrosByProject[projectName], info)
	}

	fmt.Printf("All Macros (%d total):\n", len(allMacros))
	fmt.Println(strings.Repeat("=", 80))

	// Sort project names
	var projectNames []string
	for name := range macrosByProject {
		projectNames = append(projectNames, name)
	}
	sort.Strings(projectNames)

	for _, projectName := range projectNames {
		fmt.Printf("\n%s:\n", projectName)
		for _, info := range macrosByProject[projectName] {
			fmt.Println(info)
		}
	}

	return nil
}

func listAllProjects() error {
	projects := store.GetProjects()

	if len(projects) == 0 {
		fmt.Println("No projects found.")
		fmt.Println("\nCreate one with: macro init")
		return nil
	}

	// Sort projects by name
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	fmt.Printf("Macro Projects (%d):\n", len(projects))
	fmt.Println(strings.Repeat("=", 80))

	for _, project := range projects {
		macros := store.GetMacrosByProject(project.ID)
		fmt.Printf("\n%s\n", project.Name)
		fmt.Printf("  Path: %s\n", project.Path)
		fmt.Printf("  ID: %s\n", project.ID)
		fmt.Printf("  Macros: %d\n", len(macros))
		fmt.Printf("  Created: %s\n", project.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}
