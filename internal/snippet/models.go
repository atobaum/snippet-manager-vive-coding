package snippet

import (
	"time"
)

// Snippet represents a single code snippet
type Snippet struct {
	Name        string    `yaml:"name,omitempty" json:"name"`
	Description string    `yaml:"description" json:"description"`
	Tags        []string  `yaml:"tags" json:"tags"`
	Command     string    `yaml:"command" json:"command"`
	CreatedAt   time.Time `yaml:"created_at,omitempty" json:"created_at"`
	UpdatedAt   time.Time `yaml:"updated_at,omitempty" json:"updated_at"`
}

// SnippetsFile represents the structure of the snippets.yaml file
type SnippetsFile struct {
	Snippets map[string]Snippet `yaml:"snippets" json:"snippets"`
}

// NewSnippet creates a new snippet with the given parameters
func NewSnippet(name, description, command string, tags []string) Snippet {
	now := time.Now()
	return Snippet{
		Name:        name,
		Description: description,
		Tags:        tags,
		Command:     command,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update updates the snippet with new values
func (s *Snippet) Update(description, command string, tags []string) {
	if description != "" {
		s.Description = description
	}
	if command != "" {
		s.Command = command
	}
	if tags != nil {
		s.Tags = tags
	}
	s.UpdatedAt = time.Now()
}
