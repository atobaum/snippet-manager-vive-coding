package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sni",
	Short: "Snippet management tool",
	Long: `sni is a CLI tool for managing code snippets, commands, and configuration files.
It provides both command-line interface and web UI for managing your snippets efficiently.`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(serverCmd)
}
