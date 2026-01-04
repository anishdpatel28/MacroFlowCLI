package models

import "time"

// Project represents a macro project tied to a directory
type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

// Macro represents a command alias
type Macro struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"created_at"`
}

// Database represents the entire storage structure
type Database struct {
	Projects []Project `json:"projects"`
	Macros   []Macro   `json:"macros"`
}
