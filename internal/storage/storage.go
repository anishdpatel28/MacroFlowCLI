package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/anishdpatel28/macroflow/internal/models"
	"github.com/google/uuid"
)

const (
	configDir  = ".macroflow"
	dbFileName = "macroflow.json"
)

// Storage handles all data persistence
type Storage struct {
	dbPath string
	db     *models.Database
}

// New creates a new storage instance
func New() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, configDir)
	dbPath := filepath.Join(configPath, dbFileName)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	s := &Storage{
		dbPath: dbPath,
		db:     &models.Database{},
	}

	// Load existing database or create new one
	if err := s.load(); err != nil {
		// If file doesn't exist, initialize empty database
		if os.IsNotExist(err) {
			s.db.Projects = []models.Project{}
			s.db.Macros = []models.Macro{}
			if err := s.save(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return s, nil
}

// load reads the database from disk
func (s *Storage) load() error {
	data, err := os.ReadFile(s.dbPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s.db)
}

// save writes the database to disk
func (s *Storage) save() error {
	data, err := json.MarshalIndent(s.db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	return os.WriteFile(s.dbPath, data, 0644)
}

// CreateProject creates a new project
func (s *Storage) CreateProject(name, path string) (*models.Project, error) {
	// Check if project already exists at this path
	for _, p := range s.db.Projects {
		if p.Path == path {
			return nil, fmt.Errorf("project already exists at this path: %s", p.Name)
		}
	}

	project := models.Project{
		ID:        uuid.New().String(),
		Name:      name,
		Path:      path,
		CreatedAt: time.Now(),
	}

	s.db.Projects = append(s.db.Projects, project)
	if err := s.save(); err != nil {
		return nil, err
	}

	return &project, nil
}

// GetProjects returns all projects
func (s *Storage) GetProjects() []models.Project {
	return s.db.Projects
}

// GetProjectByPath finds a project that matches the given path
// Returns the most specific (deepest) matching project
func (s *Storage) GetProjectByPath(currentPath string) *models.Project {
	var matchedProject *models.Project
	maxDepth := -1

	// Normalize the current path
	currentPath = filepath.Clean(currentPath)

	for i := range s.db.Projects {
		project := &s.db.Projects[i]
		projectPath := filepath.Clean(project.Path)

		// Check if current path is within or equal to project path
		if currentPath == projectPath || strings.HasPrefix(currentPath, projectPath+string(os.PathSeparator)) {
			// Calculate depth (number of path separators)
			depth := strings.Count(projectPath, string(os.PathSeparator))
			if depth > maxDepth {
				maxDepth = depth
				matchedProject = project
			}
		}
	}

	return matchedProject
}

// DeleteProject deletes a project and all its macros
func (s *Storage) DeleteProject(projectID string) error {
	// Find and remove project
	found := false
	for i, p := range s.db.Projects {
		if p.ID == projectID {
			s.db.Projects = append(s.db.Projects[:i], s.db.Projects[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("project not found")
	}

	// Remove all macros associated with this project
	var filteredMacros []models.Macro
	for _, m := range s.db.Macros {
		if m.ProjectID != projectID {
			filteredMacros = append(filteredMacros, m)
		}
	}
	s.db.Macros = filteredMacros

	return s.save()
}

// CreateMacro creates a new macro
func (s *Storage) CreateMacro(projectID, name, command string) (*models.Macro, error) {
	// Verify project exists
	projectExists := false
	for _, p := range s.db.Projects {
		if p.ID == projectID {
			projectExists = true
			break
		}
	}

	if !projectExists {
		return nil, fmt.Errorf("project not found")
	}

	// Check if macro with same name exists in project
	for _, m := range s.db.Macros {
		if m.ProjectID == projectID && m.Name == name {
			return nil, fmt.Errorf("macro '%s' already exists in this project", name)
		}
	}

	macro := models.Macro{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		Name:      name,
		Command:   command,
		CreatedAt: time.Now(),
	}

	s.db.Macros = append(s.db.Macros, macro)
	if err := s.save(); err != nil {
		return nil, err
	}

	return &macro, nil
}

// GetMacrosByProject returns all macros for a project
func (s *Storage) GetMacrosByProject(projectID string) []models.Macro {
	var macros []models.Macro
	for _, m := range s.db.Macros {
		if m.ProjectID == projectID {
			macros = append(macros, m)
		}
	}
	return macros
}

// GetMacroByName finds a macro by name within a project
func (s *Storage) GetMacroByName(projectID, name string) *models.Macro {
	for i := range s.db.Macros {
		m := &s.db.Macros[i]
		if m.ProjectID == projectID && m.Name == name {
			return m
		}
	}
	return nil
}

// GetAllMacros returns all macros across all projects
func (s *Storage) GetAllMacros() []models.Macro {
	return s.db.Macros
}

// DeleteMacro deletes a macro by ID
func (s *Storage) DeleteMacro(macroID string) error {
	for i, m := range s.db.Macros {
		if m.ID == macroID {
			s.db.Macros = append(s.db.Macros[:i], s.db.Macros[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("macro not found")
}

// DeleteMacros deletes multiple macros by IDs
func (s *Storage) DeleteMacros(macroIDs []string) error {
	idMap := make(map[string]bool)
	for _, id := range macroIDs {
		idMap[id] = true
	}

	var filteredMacros []models.Macro
	for _, m := range s.db.Macros {
		if !idMap[m.ID] {
			filteredMacros = append(filteredMacros, m)
		}
	}

	s.db.Macros = filteredMacros
	return s.save()
}

// ExportData exports the entire database
func (s *Storage) ExportData() (*models.Database, error) {
	return s.db, nil
}

// ImportData imports database, optionally merging with existing data
func (s *Storage) ImportData(data *models.Database, merge bool) error {
	if !merge {
		s.db = data
	} else {
		// Merge projects (skip duplicates by path)
		existingPaths := make(map[string]bool)
		for _, p := range s.db.Projects {
			existingPaths[p.Path] = true
		}

		for _, p := range data.Projects {
			if !existingPaths[p.Path] {
				s.db.Projects = append(s.db.Projects, p)
			}
		}

		// Merge macros (regenerate IDs to avoid conflicts)
		for _, m := range data.Macros {
			m.ID = uuid.New().String()
			s.db.Macros = append(s.db.Macros, m)
		}
	}

	return s.save()
}
