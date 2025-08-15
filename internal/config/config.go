package config

import (
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	ConfigDir   string
	SnippetFile string
	ServerPort  int
}

// DefaultConfig returns the default configuration
func DefaultConfig() (*Config, error) {
	// Try to get config directory from environment variable first
	configDir := os.Getenv("SNI_CONFIG_DIR")
	if configDir == "" {
		// Fall back to project directory (current working directory)
		workDir, err := os.Getwd()
		if err != nil {
			// Final fallback to user home directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			configDir = filepath.Join(homeDir, ".config", "sni")
		} else {
			configDir = filepath.Join(workDir, ".sni")
		}
	}

	snippetFile := filepath.Join(configDir, "snippets.yaml")

	return &Config{
		ConfigDir:   configDir,
		SnippetFile: snippetFile,
		ServerPort:  8080,
	}, nil
}
