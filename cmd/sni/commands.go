package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/atobaum/snippet-manager/internal/cli"
	"github.com/atobaum/snippet-manager/internal/selector"
	"github.com/atobaum/snippet-manager/internal/server"
	"github.com/atobaum/snippet-manager/internal/snippet"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		svc, err := snippet.NewService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing service: %v\n", err)
			return
		}

		// Interactive input
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Description: ")
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)

		fmt.Print("Language (e.g., bash, go, python, javascript): ")
		language, _ := reader.ReadString('\n')
		language = strings.TrimSpace(language)

		fmt.Print("Tags (comma separated): ")
		tagsInput, _ := reader.ReadString('\n')
		tagsInput = strings.TrimSpace(tagsInput)
		var tags []string
		if tagsInput != "" {
			tags = strings.Split(tagsInput, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		fmt.Println("Command/Content (end with Ctrl+D on empty line):")
		var commandLines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			commandLines = append(commandLines, line)
		}
		command := strings.Join(commandLines, "")
		command = strings.TrimSpace(command)

		if err := svc.CreateSnippet(name, description, command, language, tags); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating snippet: %v\n", err)
			return
		}

		fmt.Printf("✅ Snippet '%s' created successfully!\n", name)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snippets",
	Run: func(cmd *cobra.Command, args []string) {
		colorEnabled, _ := cmd.Flags().GetBool("color")
		cli.EnableColors(colorEnabled)

		svc, err := snippet.NewService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", cli.ColorizeError(fmt.Sprintf("Error initializing service: %v", err)))
			return
		}

		snippets, err := svc.ListSnippets()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", cli.ColorizeError(fmt.Sprintf("Error listing snippets: %v", err)))
			return
		}

		if len(snippets) == 0 {
			fmt.Println(cli.ColorizeWarning("No snippets found. Create one with 'sni new <name>'"))
			return
		}

		fmt.Printf("%s\n\n", cli.ColorizeTitle(fmt.Sprintf("Found %d snippet(s):", len(snippets))))
		for _, s := range snippets {
			fmt.Println(cli.ColorizeSnippetName(s.Name))
			if desc := cli.ColorizeDescription(s.Description); desc != "" {
				fmt.Println(desc)
			}
			if lang := cli.ColorizeLanguage(s.Language); lang != "" {
				fmt.Println(lang)
			}
			if tags := cli.ColorizeTags(s.Tags); tags != "" {
				fmt.Println(tags)
			}
			fmt.Println()
		}
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Search snippets by keyword",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := args[0]
		colorEnabled, _ := cmd.Flags().GetBool("color")

		cli.EnableColors(colorEnabled)

		svc, err := snippet.NewService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", cli.ColorizeError(fmt.Sprintf("Error initializing service: %v", err)))
			return
		}

		snippets, err := svc.SearchSnippets(keyword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", cli.ColorizeError(fmt.Sprintf("Error searching snippets: %v", err)))
			return
		}

		if len(snippets) == 0 {
			fmt.Println(cli.ColorizeWarning(fmt.Sprintf("No snippets found for keyword: %s", keyword)))
			return
		}

		fmt.Printf("%s\n\n", cli.ColorizeTitle(fmt.Sprintf("🔍 Found %d snippet(s) for '%s':", len(snippets), keyword)))
		for _, s := range snippets {
			fmt.Println(cli.ColorizeSnippetName(s.Name))
			if desc := cli.ColorizeDescription(s.Description); desc != "" {
				fmt.Println(desc)
			}
			if lang := cli.ColorizeLanguage(s.Language); lang != "" {
				fmt.Println(lang)
			}
			if tags := cli.ColorizeTags(s.Tags); tags != "" {
				fmt.Println(tags)
			}
			fmt.Println()
		}
	},
}

var useCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Output snippet content to terminal",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		svc, err := snippet.NewService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing service: %v\n", err)
			return
		}

		snippet, err := svc.GetSnippet(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting snippet: %v\n", err)
			return
		}

		// Output the command content directly
		fmt.Print(snippet.Command)
	},
}

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		svc, err := snippet.NewService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing service: %v\n", err)
			return
		}

		// Get existing snippet
		existing, err := svc.GetSnippet(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting snippet: %v\n", err)
			return
		}

		// Interactive editing
		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("Description [%s]: ", existing.Description)
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)
		if description == "" {
			description = existing.Description
		}

		fmt.Printf("Language [%s]: ", existing.Language)
		language, _ := reader.ReadString('\n')
		language = strings.TrimSpace(language)
		if language == "" {
			language = existing.Language
		}

		fmt.Printf("Tags [%s]: ", strings.Join(existing.Tags, ", "))
		tagsInput, _ := reader.ReadString('\n')
		tagsInput = strings.TrimSpace(tagsInput)
		var tags []string
		if tagsInput == "" {
			tags = existing.Tags
		} else {
			tags = strings.Split(tagsInput, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
		}

		fmt.Println("Command/Content (current content shown, edit and end with Ctrl+D):")
		fmt.Println("--- Current Content ---")
		fmt.Print(existing.Command)
		fmt.Println("\n--- Enter New Content ---")

		var commandLines []string
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			commandLines = append(commandLines, line)
		}
		command := strings.Join(commandLines, "")
		command = strings.TrimSpace(command)
		if command == "" {
			command = existing.Command
		}

		if err := svc.UpdateSnippet(name, description, command, language, tags); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating snippet: %v\n", err)
			return
		}

		fmt.Printf("✅ Snippet '%s' updated successfully!\n", name)
	},
}

var rmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "Remove a snippet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		svc, err := snippet.NewService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing service: %v\n", err)
			return
		}

		// Confirm deletion
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to delete snippet '%s'? (y/N): ", name)
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(strings.ToLower(confirmation))

		if confirmation != "y" && confirmation != "yes" {
			fmt.Println("Deletion cancelled.")
			return
		}

		if err := svc.DeleteSnippet(name); err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting snippet: %v\n", err)
			return
		}

		fmt.Printf("✅ Snippet '%s' deleted successfully!\n", name)
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the web UI server",
	Run: func(cmd *cobra.Command, args []string) {
		devMode, _ := cmd.Flags().GetBool("dev")
		port, _ := cmd.Flags().GetInt("port")

		srv, err := server.NewServer(port, devMode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating server: %v\n", err)
			return
		}

		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
			return
		}
	},
}

func init() {
	serverCmd.Flags().BoolP("dev", "d", false, "Run in development mode (proxy to Svelte dev server)")
	serverCmd.Flags().IntP("port", "p", 8080, "Port to run server on")
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a snippet interactively",
	Run: func(cmd *cobra.Command, args []string) {
		tagFilter, _ := cmd.Flags().GetString("tag")
		colorEnabled, _ := cmd.Flags().GetBool("color")
		cli.EnableColors(colorEnabled)

		svc, err := snippet.NewService()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", cli.ColorizeError(fmt.Sprintf("Error initializing service: %v", err)))
			return
		}

		var snippets []snippet.Snippet
		if tagFilter != "" {
			snippets, err = svc.SearchSnippets(tagFilter)
		} else {
			snippets, err = svc.ListSnippets()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", cli.ColorizeError(fmt.Sprintf("Error loading snippets: %v", err)))
			return
		}

		if len(snippets) == 0 {
			fmt.Println(cli.ColorizeWarning("No snippets found."))
			return
		}

		// Use selector to choose snippet
		sel := selector.NewSelector(colorEnabled)
		prompt := cli.ColorizeTitle("Select a snippet to execute:")
		selectedSnippet, err := sel.Select(snippets, prompt)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", cli.ColorizeError(fmt.Sprintf("Selection error: %v", err)))
			return
		}

		if selectedSnippet == nil {
			// User cancelled
			return
		}

		// Copy to clipboard and show info
		fmt.Printf("\n%s\n", cli.InfoColor.Sprintf("📋 Selected: %s", selectedSnippet.Name))
		if selectedSnippet.Description != "" {
			fmt.Printf("%s\n", cli.ColorizeDescription(selectedSnippet.Description))
		}
		fmt.Printf("%s\n", cli.CommandColor.Sprintf("Command: %s", selectedSnippet.Command))

		// Try to copy to clipboard
		err = copyToClipboard(selectedSnippet.Command)
		if err != nil {
			fmt.Printf("\n%s\n", cli.ColorizeWarning("Could not copy to clipboard. Here's the command:"))
			fmt.Println(selectedSnippet.Command)
		} else {
			fmt.Printf("\n%s\n", cli.ColorizeSuccess("✅ Command copied to clipboard! Paste it in your terminal."))
		}
	},
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure sni settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔧 sni Configuration")
		fmt.Println()
		fmt.Println("Configuration is managed through environment variables:")
		fmt.Println()
		fmt.Println("📁 SNI_CONFIG_DIR    - Custom config directory")
		fmt.Printf("   Current: %s\n", getConfigInfo())
		fmt.Println()
		fmt.Println("Example usage:")
		fmt.Println("  export SNI_CONFIG_DIR=\"/path/to/config\"")
		fmt.Println("  sni list")
		fmt.Println()
		fmt.Println("🌐 Web UI settings:")
		fmt.Println("  sni server --port 9090         # Custom port")
		fmt.Println("  sni server --dev               # Development mode")
		fmt.Println()
	},
}

// Helper functions
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func getConfigInfo() string {
	if configDir := os.Getenv("SNI_CONFIG_DIR"); configDir != "" {
		return configDir + " (from SNI_CONFIG_DIR)"
	}

	workDir, err := os.Getwd()
	if err != nil {
		return "~/.config/sni (fallback)"
	}
	return workDir + "/.sni (default)"
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else {
			return fmt.Errorf("no clipboard utility found (install xclip or xsel)")
		}
	case "windows":
		cmd = exec.Command("clip")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func init() {
	execCmd.Flags().StringP("tag", "t", "", "Filter snippets by tag")
	execCmd.Flags().Bool("color", false, "Enable colorized output")

	listCmd.Flags().Bool("color", false, "Enable colorized output")
	searchCmd.Flags().Bool("color", false, "Enable colorized output")
}
