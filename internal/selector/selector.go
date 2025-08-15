package selector

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/atobaum/snippet-manager/internal/snippet"
)

// Selector interface for different selection methods
type Selector interface {
	Select(snippets []snippet.Snippet, prompt string) (*snippet.Snippet, error)
}

// FzfSelector uses fzf for selection
type FzfSelector struct {
	colorEnabled bool
}

// NumberSelector uses number-based selection as fallback
type NumberSelector struct {
	colorEnabled bool
}

// NewSelector creates the best available selector
func NewSelector(colorEnabled bool) Selector {
	if IsFzfAvailable() {
		return &FzfSelector{colorEnabled: colorEnabled}
	}
	return &NumberSelector{colorEnabled: colorEnabled}
}

// IsFzfAvailable checks if fzf is installed
func IsFzfAvailable() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

// Select using fzf
func (f *FzfSelector) Select(snippets []snippet.Snippet, prompt string) (*snippet.Snippet, error) {
	if len(snippets) == 0 {
		return nil, fmt.Errorf("no snippets available")
	}

	// Prepare fzf input - simpler format
	var lines []string
	for i, s := range snippets {
		// Show name and description for easier selection
		line := fmt.Sprintf("%d: %s", i, s.Name)
		if s.Description != "" {
			line += fmt.Sprintf(" - %s", s.Description)
		}
		lines = append(lines, line)
	}

	// Create temporary file for input
	inputStr := strings.Join(lines, "\n")

	// Run fzf with simple setup
	cmd := exec.Command("fzf",
		"--height=40%",
		"--layout=reverse",
		"--border",
		"--prompt=Select snippet: ",
	)

	cmd.Stdin = strings.NewReader(inputStr)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		// User cancelled or error
		return nil, nil
	}

	// Parse selection
	selectedLine := strings.TrimSpace(string(output))
	if selectedLine == "" {
		return nil, nil
	}

	// Extract index from "index: name - description" format
	parts := strings.SplitN(selectedLine, ":", 2)
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid fzf output")
	}

	index, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || index < 0 || index >= len(snippets) {
		return nil, fmt.Errorf("invalid selection index")
	}

	return &snippets[index], nil
}

// Select using number input (fallback)
func (n *NumberSelector) Select(snippets []snippet.Snippet, prompt string) (*snippet.Snippet, error) {
	if len(snippets) == 0 {
		return nil, fmt.Errorf("no snippets available")
	}

	// Display snippets with numbers
	fmt.Println(prompt)
	fmt.Println()
	for i, s := range snippets {
		fmt.Printf("%d. %s\n", i+1, s.Name)
		if s.Description != "" {
			fmt.Printf("   Description: %s\n", s.Description)
		}
		if len(s.Tags) > 0 {
			fmt.Printf("   Tags: %s\n", strings.Join(s.Tags, ", "))
		}
		// Show command preview
		commandPreview := strings.ReplaceAll(s.Command, "\n", " ")
		if len(commandPreview) > 100 {
			commandPreview = commandPreview[:100] + "..."
		}
		fmt.Printf("   Command: %s\n", commandPreview)
		fmt.Println()
	}

	// Get user selection
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter number (or 'q' to quit): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "quit" {
		return nil, nil
	}

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(snippets) {
		return nil, fmt.Errorf("invalid selection: %s", input)
	}

	return &snippets[num-1], nil
}
