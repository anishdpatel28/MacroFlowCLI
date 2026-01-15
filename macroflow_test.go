package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anish/macroflow/internal/models"
	"github.com/anish/macroflow/internal/storage"
)

// TestStorage holds the test storage instance
var testStorage *storage.Storage
var testStoragePath string

// setupTest creates a temporary storage for testing
func setupTest(t *testing.T) func() {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "macroflow-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Set up test storage path
	testStoragePath = filepath.Join(tmpDir, ".macroflow")
	err = os.MkdirAll(testStoragePath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test storage dir: %v", err)
	}

	// Override storage path for testing
	os.Setenv("HOME", tmpDir)

	// Create storage instance
	testStorage, err = storage.New()
	if err != nil {
		t.Fatalf("Failed to create test storage: %v", err)
	}

	// Return cleanup function
	return func() {
		os.RemoveAll(tmpDir)
		os.Unsetenv("HOME")
	}
}

// TestProjectCreation tests creating a new project
func TestProjectCreation(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	projectName := "test-project"
	projectPath := "/tmp/test-project"

	project, err := testStorage.CreateProject(projectName, projectPath)
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	if project.Name != projectName {
		t.Errorf("Expected project name %s, got %s", projectName, project.Name)
	}

	if project.Path != projectPath {
		t.Errorf("Expected project path %s, got %s", projectPath, project.Path)
	}

	if project.ID == "" {
		t.Error("Project ID should not be empty")
	}

	if project.CreatedAt.IsZero() {
		t.Error("Project creation time should not be zero")
	}
}

// TestDuplicateProject tests that creating a duplicate project fails
func TestDuplicateProject(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	projectPath := "/tmp/test-project"

	_, err := testStorage.CreateProject("project1", projectPath)
	if err != nil {
		t.Fatalf("Failed to create first project: %v", err)
	}

	_, err = testStorage.CreateProject("project2", projectPath)
	if err == nil {
		t.Error("Expected error when creating duplicate project at same path")
	}
}

// TestGetProjects tests retrieving all projects
func TestGetProjects(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create multiple projects
	projects := []struct {
		name string
		path string
	}{
		{"project1", "/tmp/project1"},
		{"project2", "/tmp/project2"},
		{"project3", "/tmp/project3"},
	}

	for _, p := range projects {
		_, err := testStorage.CreateProject(p.name, p.path)
		if err != nil {
			t.Fatalf("Failed to create project %s: %v", p.name, err)
		}
	}

	allProjects := testStorage.GetProjects()
	if len(allProjects) != len(projects) {
		t.Errorf("Expected %d projects, got %d", len(projects), len(allProjects))
	}
}

// TestProjectByPath tests finding project by path
func TestProjectByPath(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create nested projects
	project1, _ := testStorage.CreateProject("parent", "/tmp/parent")
	project2, _ := testStorage.CreateProject("child", "/tmp/parent/child")

	tests := []struct {
		path            string
		expectedProject *models.Project
	}{
		{"/tmp/parent", project1},
		{"/tmp/parent/file.txt", project1},
		{"/tmp/parent/child", project2},
		{"/tmp/parent/child/deep", project2},
		{"/tmp/other", nil},
	}

	for _, tt := range tests {
		result := testStorage.GetProjectByPath(tt.path)
		if tt.expectedProject == nil {
			if result != nil {
				t.Errorf("Expected nil for path %s, got project %s", tt.path, result.Name)
			}
		} else {
			if result == nil {
				t.Errorf("Expected project for path %s, got nil", tt.path)
			} else if result.ID != tt.expectedProject.ID {
				t.Errorf("Expected project %s for path %s, got %s", tt.expectedProject.Name, tt.path, result.Name)
			}
		}
	}
}

// TestDeleteProject tests deleting a project
func TestDeleteProject(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")
	
	// Add a macro to the project
	_, err := testStorage.CreateMacro(project.ID, "test-macro", "echo test")
	if err != nil {
		t.Fatalf("Failed to create macro: %v", err)
	}

	// Delete the project
	err = testStorage.DeleteProject(project.ID)
	if err != nil {
		t.Fatalf("Failed to delete project: %v", err)
	}

	// Verify project is gone
	projects := testStorage.GetProjects()
	if len(projects) != 0 {
		t.Error("Project should be deleted")
	}

	// Verify macros are gone
	macros := testStorage.GetAllMacros()
	if len(macros) != 0 {
		t.Error("Macros should be deleted with project")
	}
}

// TestMacroCreation tests creating a macro
func TestMacroCreation(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")

	macroName := "test-macro"
	macroCommand := "echo 'Hello, World!'"

	macro, err := testStorage.CreateMacro(project.ID, macroName, macroCommand)
	if err != nil {
		t.Fatalf("Failed to create macro: %v", err)
	}

	if macro.Name != macroName {
		t.Errorf("Expected macro name %s, got %s", macroName, macro.Name)
	}

	if macro.Command != macroCommand {
		t.Errorf("Expected macro command %s, got %s", macroCommand, macro.Command)
	}

	if macro.ProjectID != project.ID {
		t.Errorf("Expected project ID %s, got %s", project.ID, macro.ProjectID)
	}

	if macro.ID == "" {
		t.Error("Macro ID should not be empty")
	}
}

// TestDuplicateMacro tests that creating duplicate macro fails
func TestDuplicateMacro(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")

	_, err := testStorage.CreateMacro(project.ID, "test", "echo 1")
	if err != nil {
		t.Fatalf("Failed to create first macro: %v", err)
	}

	_, err = testStorage.CreateMacro(project.ID, "test", "echo 2")
	if err == nil {
		t.Error("Expected error when creating duplicate macro")
	}
}

// TestGetMacrosByProject tests retrieving macros for a project
func TestGetMacrosByProject(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project1, _ := testStorage.CreateProject("project1", "/tmp/project1")
	project2, _ := testStorage.CreateProject("project2", "/tmp/project2")

	// Create macros for project1
	testStorage.CreateMacro(project1.ID, "macro1", "echo 1")
	testStorage.CreateMacro(project1.ID, "macro2", "echo 2")

	// Create macro for project2
	testStorage.CreateMacro(project2.ID, "macro3", "echo 3")

	macros1 := testStorage.GetMacrosByProject(project1.ID)
	if len(macros1) != 2 {
		t.Errorf("Expected 2 macros for project1, got %d", len(macros1))
	}

	macros2 := testStorage.GetMacrosByProject(project2.ID)
	if len(macros2) != 1 {
		t.Errorf("Expected 1 macro for project2, got %d", len(macros2))
	}
}

// TestGetMacroByName tests finding a macro by name
func TestGetMacroByName(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")
	createdMacro, _ := testStorage.CreateMacro(project.ID, "test-macro", "echo test")

	foundMacro := testStorage.GetMacroByName(project.ID, "test-macro")
	if foundMacro == nil {
		t.Fatal("Expected to find macro, got nil")
	}

	if foundMacro.ID != createdMacro.ID {
		t.Error("Found macro ID doesn't match created macro ID")
	}

	notFound := testStorage.GetMacroByName(project.ID, "nonexistent")
	if notFound != nil {
		t.Error("Expected nil for nonexistent macro")
	}
}

// TestGetAllMacros tests retrieving all macros
func TestGetAllMacros(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project1, _ := testStorage.CreateProject("project1", "/tmp/project1")
	project2, _ := testStorage.CreateProject("project2", "/tmp/project2")

	testStorage.CreateMacro(project1.ID, "macro1", "echo 1")
	testStorage.CreateMacro(project1.ID, "macro2", "echo 2")
	testStorage.CreateMacro(project2.ID, "macro3", "echo 3")

	allMacros := testStorage.GetAllMacros()
	if len(allMacros) != 3 {
		t.Errorf("Expected 3 total macros, got %d", len(allMacros))
	}
}

// TestDeleteMacro tests deleting a single macro
func TestDeleteMacro(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")
	macro, _ := testStorage.CreateMacro(project.ID, "test", "echo test")

	err := testStorage.DeleteMacro(macro.ID)
	if err != nil {
		t.Fatalf("Failed to delete macro: %v", err)
	}

	macros := testStorage.GetMacrosByProject(project.ID)
	if len(macros) != 0 {
		t.Error("Macro should be deleted")
	}
}

// TestDeleteMacros tests deleting multiple macros
func TestDeleteMacros(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")
	macro1, _ := testStorage.CreateMacro(project.ID, "macro1", "echo 1")
	macro2, _ := testStorage.CreateMacro(project.ID, "macro2", "echo 2")
	macro3, _ := testStorage.CreateMacro(project.ID, "macro3", "echo 3")

	err := testStorage.DeleteMacros([]string{macro1.ID, macro2.ID})
	if err != nil {
		t.Fatalf("Failed to delete macros: %v", err)
	}

	macros := testStorage.GetMacrosByProject(project.ID)
	if len(macros) != 1 {
		t.Errorf("Expected 1 remaining macro, got %d", len(macros))
	}

	if macros[0].ID != macro3.ID {
		t.Error("Wrong macro remained")
	}
}

// TestExportData tests exporting all data
func TestExportData(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")
	testStorage.CreateMacro(project.ID, "macro1", "echo 1")
	testStorage.CreateMacro(project.ID, "macro2", "echo 2")

	data, err := testStorage.ExportData()
	if err != nil {
		t.Fatalf("Failed to export data: %v", err)
	}

	if len(data.Projects) != 1 {
		t.Errorf("Expected 1 project in export, got %d", len(data.Projects))
	}

	if len(data.Macros) != 2 {
		t.Errorf("Expected 2 macros in export, got %d", len(data.Macros))
	}
}

// TestImportData tests importing data
func TestImportData(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create import data
	importData := &models.Database{
		Projects: []models.Project{
			{ID: "proj1", Name: "Imported Project", Path: "/tmp/imported"},
		},
		Macros: []models.Macro{
			{ID: "macro1", ProjectID: "proj1", Name: "imported", Command: "echo imported"},
		},
	}

	// Import (replace)
	err := testStorage.ImportData(importData, false)
	if err != nil {
		t.Fatalf("Failed to import data: %v", err)
	}

	projects := testStorage.GetProjects()
	if len(projects) != 1 {
		t.Errorf("Expected 1 project after import, got %d", len(projects))
	}

	macros := testStorage.GetAllMacros()
	if len(macros) != 1 {
		t.Errorf("Expected 1 macro after import, got %d", len(macros))
	}
}

// TestImportDataMerge tests importing data with merge
func TestImportDataMerge(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create existing data
	existingProject, _ := testStorage.CreateProject("existing", "/tmp/existing")
	testStorage.CreateMacro(existingProject.ID, "existing-macro", "echo existing")

	// Create import data
	importData := &models.Database{
		Projects: []models.Project{
			{ID: "proj1", Name: "Imported Project", Path: "/tmp/imported"},
		},
		Macros: []models.Macro{
			{ID: "macro1", ProjectID: "proj1", Name: "imported", Command: "echo imported"},
		},
	}

	// Import with merge
	err := testStorage.ImportData(importData, true)
	if err != nil {
		t.Fatalf("Failed to merge import data: %v", err)
	}

	projects := testStorage.GetProjects()
	if len(projects) != 2 {
		t.Errorf("Expected 2 projects after merge, got %d", len(projects))
	}

	macros := testStorage.GetAllMacros()
	if len(macros) != 2 {
		t.Errorf("Expected 2 macros after merge, got %d", len(macros))
	}
}

// TestStoragePersistence tests that data persists across storage instances
func TestStoragePersistence(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Create data with first storage instance
	project, _ := testStorage.CreateProject("test", "/tmp/test")
	testStorage.CreateMacro(project.ID, "macro1", "echo 1")

	// Create new storage instance
	newStorage, err := storage.New()
	if err != nil {
		t.Fatalf("Failed to create new storage instance: %v", err)
	}

	// Verify data persisted
	projects := newStorage.GetProjects()
	if len(projects) != 1 {
		t.Error("Data should persist across storage instances")
	}

	macros := newStorage.GetAllMacros()
	if len(macros) != 1 {
		t.Error("Macros should persist across storage instances")
	}
}

// TestInvalidProjectOperations tests operations with invalid project IDs
func TestInvalidProjectOperations(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Try to create macro for nonexistent project
	_, err := testStorage.CreateMacro("invalid-id", "test", "echo test")
	if err == nil {
		t.Error("Expected error when creating macro for nonexistent project")
	}

	// Try to delete nonexistent project
	err = testStorage.DeleteProject("invalid-id")
	if err == nil {
		t.Error("Expected error when deleting nonexistent project")
	}
}

// TestInvalidMacroOperations tests operations with invalid macro IDs
func TestInvalidMacroOperations(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	// Try to delete nonexistent macro
	err := testStorage.DeleteMacro("invalid-id")
	if err == nil {
		t.Error("Expected error when deleting nonexistent macro")
	}
}

// TestJSONMarshaling tests that data can be marshaled to JSON
func TestJSONMarshaling(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	project, _ := testStorage.CreateProject("test", "/tmp/test")
	testStorage.CreateMacro(project.ID, "macro1", "echo 1")

	data, _ := testStorage.ExportData()

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal data to JSON: %v", err)
	}

	// Unmarshal back
	var unmarshaled models.Database
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(unmarshaled.Projects) != 1 {
		t.Error("Project not preserved in JSON marshaling")
	}

	if len(unmarshaled.Macros) != 1 {
		t.Error("Macro not preserved in JSON marshaling")
	}
}

// TestEmptyStorage tests operations on empty storage
func TestEmptyStorage(t *testing.T) {
	cleanup := setupTest(t)
	defer cleanup()

	projects := testStorage.GetProjects()
	if len(projects) != 0 {
		t.Error("New storage should have no projects")
	}

	macros := testStorage.GetAllMacros()
	if len(macros) != 0 {
		t.Error("New storage should have no macros")
	}

	result := testStorage.GetProjectByPath("/any/path")
	if result != nil {
		t.Error("Empty storage should return nil for any path")
	}
}

// BenchmarkProjectCreation benchmarks project creation
func BenchmarkProjectCreation(b *testing.B) {
	// Note: This is a basic benchmark setup
	tmpDir, _ := os.MkdirTemp("", "macroflow-bench-*")
	defer os.RemoveAll(tmpDir)
	
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")
	
	store, _ := storage.New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.CreateProject(fmt.Sprintf("project%d", i), fmt.Sprintf("/tmp/project%d", i))
	}
}

// BenchmarkMacroCreation benchmarks macro creation
func BenchmarkMacroCreation(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "macroflow-bench-*")
	defer os.RemoveAll(tmpDir)
	
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")
	
	store, _ := storage.New()
	project, _ := store.CreateProject("test", "/tmp/test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.CreateMacro(project.ID, fmt.Sprintf("macro%d", i), "echo test")
	}
}

// BenchmarkProjectLookup benchmarks finding project by path
func BenchmarkProjectLookup(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "macroflow-bench-*")
	defer os.RemoveAll(tmpDir)
	
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")
	
	store, _ := storage.New()
	
	// Create several projects
	for i := 0; i < 10; i++ {
		store.CreateProject(fmt.Sprintf("project%d", i), fmt.Sprintf("/tmp/project%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetProjectByPath("/tmp/project5/subdir/file.txt")
	}
}
