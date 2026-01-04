package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/anish/macroflow/internal/models"
	"github.com/spf13/cobra"
)

var (
	importFile  string
	importMerge bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import macros and projects from a JSON file",
	Long: `Import macros and projects from a JSON file.

By default, this will replace all existing data. Use --merge to add to existing data.

Examples:
  macro import --file macros.json           # Replace all data
  macro import --file macros.json --merge   # Merge with existing data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if importFile == "" {
			return fmt.Errorf("--file flag is required")
		}

		fileData, err := os.ReadFile(importFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var data models.Database
		if err := json.Unmarshal(fileData, &data); err != nil {
			return fmt.Errorf("failed to parse file: %w", err)
		}

		if !importMerge {
			fmt.Print("Warning: This will replace all existing macros and projects.\nType 'yes' to confirm: ")
			var confirmation string
			fmt.Scanln(&confirmation)

			if confirmation != "yes" {
				fmt.Println("Import cancelled.")
				return nil
			}
		}

		if err := store.ImportData(&data, importMerge); err != nil {
			return fmt.Errorf("failed to import data: %w", err)
		}

		if importMerge {
			fmt.Printf("✓ Merged %d project(s) and %d macro(s) from %s\n",
				len(data.Projects), len(data.Macros), importFile)
		} else {
			fmt.Printf("✓ Imported %d project(s) and %d macro(s) from %s\n",
				len(data.Projects), len(data.Macros), importFile)
		}

		return nil
	},
}

func init() {
	importCmd.Flags().StringVarP(&importFile, "file", "f", "", "Input file path (required)")
	importCmd.Flags().BoolVarP(&importMerge, "merge", "m", false, "Merge with existing data instead of replacing")
	rootCmd.AddCommand(importCmd)
}
