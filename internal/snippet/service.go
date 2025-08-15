package snippet

import (
	"fmt"
	"os"
	"strings"

	"github.com/atobaum/snippet-manager/internal/config"
	"gopkg.in/yaml.v3"
)

// Service handles snippet operations
type Service struct {
	config *config.Config
}

// NewService creates a new snippet service
func NewService() (*Service, error) {
	cfg, err := config.DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(cfg.ConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &Service{
		config: cfg,
	}, nil
}

// LoadSnippets loads all snippets from the YAML file
func (s *Service) LoadSnippets() (*SnippetsFile, error) {
	// If file doesn't exist, return empty snippets
	if _, err := os.Stat(s.config.SnippetFile); os.IsNotExist(err) {
		return &SnippetsFile{
			Snippets: make(map[string]Snippet),
		}, nil
	}

	data, err := os.ReadFile(s.config.SnippetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read snippets file: %w", err)
	}

	var snippetsFile SnippetsFile
	if err := yaml.Unmarshal(data, &snippetsFile); err != nil {
		return nil, fmt.Errorf("failed to parse snippets file: %w", err)
	}

	// Initialize map if nil
	if snippetsFile.Snippets == nil {
		snippetsFile.Snippets = make(map[string]Snippet)
	}

	return &snippetsFile, nil
}

// SaveSnippets saves all snippets to the YAML file
func (s *Service) SaveSnippets(snippetsFile *SnippetsFile) error {
	data, err := yaml.Marshal(snippetsFile)
	if err != nil {
		return fmt.Errorf("failed to marshal snippets: %w", err)
	}

	if err := os.WriteFile(s.config.SnippetFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write snippets file: %w", err)
	}

	return nil
}

// CreateSnippet creates a new snippet
func (s *Service) CreateSnippet(name, description, command, language string, tags []string) error {
	snippetsFile, err := s.LoadSnippets()
	if err != nil {
		return err
	}

	if _, exists := snippetsFile.Snippets[name]; exists {
		return fmt.Errorf("snippet '%s' already exists", name)
	}

	snippet := NewSnippet(name, description, command, language, tags)
	snippetsFile.Snippets[name] = snippet

	return s.SaveSnippets(snippetsFile)
}

// GetSnippet retrieves a snippet by name
func (s *Service) GetSnippet(name string) (*Snippet, error) {
	snippetsFile, err := s.LoadSnippets()
	if err != nil {
		return nil, err
	}

	snippet, exists := snippetsFile.Snippets[name]
	if !exists {
		return nil, fmt.Errorf("snippet '%s' not found", name)
	}

	return &snippet, nil
}

// UpdateSnippet updates an existing snippet
func (s *Service) UpdateSnippet(name, description, command, language string, tags []string) error {
	snippetsFile, err := s.LoadSnippets()
	if err != nil {
		return err
	}

	snippet, exists := snippetsFile.Snippets[name]
	if !exists {
		return fmt.Errorf("snippet '%s' not found", name)
	}

	snippet.Update(description, command, language, tags)
	snippetsFile.Snippets[name] = snippet

	return s.SaveSnippets(snippetsFile)
}

// DeleteSnippet removes a snippet
func (s *Service) DeleteSnippet(name string) error {
	snippetsFile, err := s.LoadSnippets()
	if err != nil {
		return err
	}

	if _, exists := snippetsFile.Snippets[name]; !exists {
		return fmt.Errorf("snippet '%s' not found", name)
	}

	delete(snippetsFile.Snippets, name)
	return s.SaveSnippets(snippetsFile)
}

// ListSnippets returns all snippets
func (s *Service) ListSnippets() ([]Snippet, error) {
	snippetsFile, err := s.LoadSnippets()
	if err != nil {
		return nil, err
	}

	snippets := make([]Snippet, 0, len(snippetsFile.Snippets))
	for name, snippet := range snippetsFile.Snippets {
		snippet.Name = name // Ensure name is set
		snippets = append(snippets, snippet)
	}

	return snippets, nil
}

// SearchSnippets searches snippets by keyword
func (s *Service) SearchSnippets(keyword string) ([]Snippet, error) {
	snippets, err := s.ListSnippets()
	if err != nil {
		return nil, err
	}

	keyword = strings.ToLower(keyword)
	var results []Snippet

	for _, snippet := range snippets {
		// Search in name, description, tags, and command
		if strings.Contains(strings.ToLower(snippet.Name), keyword) ||
			strings.Contains(strings.ToLower(snippet.Description), keyword) ||
			strings.Contains(strings.ToLower(snippet.Command), keyword) ||
			s.containsTag(snippet.Tags, keyword) {
			results = append(results, snippet)
		}
	}

	return results, nil
}

// containsTag checks if any tag contains the keyword
func (s *Service) containsTag(tags []string, keyword string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), keyword) {
			return true
		}
	}
	return false
}
