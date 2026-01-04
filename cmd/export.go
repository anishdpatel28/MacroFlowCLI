package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var exportFile string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all macros and projects to a JSON file",
	Long: `Export all macros and projects to a JSON file for backup or sharing.

Example:
  macro export --file macros.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if exportFile == "" {
			return fmt.Errorf("--file flag is required")
		}

		data, err := store.ExportData()
		if err != nil {
			return fmt.Errorf("failed to export data: %w", err)
		}

		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}

		if err := os.WriteFile(exportFile, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("✓ Exported %d project(s) and %d macro(s) to %s\n",
			len(data.Projects), len(data.Macros), exportFile)

		return nil
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "Output file path (required)")
	rootCmd.AddCommand(exportCmd)
}
