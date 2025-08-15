package cli

import (
	"strings"

	"github.com/fatih/color"
)

// Color functions for consistent styling
var (
	// Title colors
	TitleColor  = color.New(color.FgCyan, color.Bold)
	HeaderColor = color.New(color.FgYellow, color.Bold)

	// Content colors
	NameColor    = color.New(color.FgGreen, color.Bold)
	DescColor    = color.New(color.FgWhite)
	TagColor     = color.New(color.FgBlue)
	CommandColor = color.New(color.FgHiBlack)

	// Status colors
	SuccessColor = color.New(color.FgGreen)
	ErrorColor   = color.New(color.FgRed)
	WarningColor = color.New(color.FgYellow)
	InfoColor    = color.New(color.FgCyan)

	// Special colors
	NumberColor  = color.New(color.FgMagenta, color.Bold)
	BracketColor = color.New(color.FgHiBlack)
)

// ColorizeSnippetName formats snippet name with color
func ColorizeSnippetName(name string) string {
	return NameColor.Sprintf("🔹 %s", name)
}

// ColorizeDescription formats description with color
func ColorizeDescription(desc string) string {
	if desc == "" {
		return ""
	}
	return "   " + DescColor.Sprintf("Description: %s", desc)
}

// ColorizeLanguage formats language with color
func ColorizeLanguage(language string) string {
	if language == "" {
		return ""
	}
	return "   " + InfoColor.Sprintf("Language: %s", language)
}

// ColorizeTags formats tags with color
func ColorizeTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	var coloredTags []string
	for _, tag := range tags {
		coloredTags = append(coloredTags, TagColor.Sprintf("#%s", tag))
	}

	return "   " + DescColor.Sprint("Tags: ") + strings.Join(coloredTags, " ")
}

// ColorizeCommand formats command with color and truncation
func ColorizeCommand(command string, maxLen int) string {
	// Truncate if too long
	if len(command) > maxLen {
		command = command[:maxLen] + "..."
	}

	// Replace newlines with spaces for display
	command = strings.ReplaceAll(command, "\n", " ")

	return "   " + CommandColor.Sprintf("Command: %s", command)
}

// ColorizeNumber formats number with color
func ColorizeNumber(num int) string {
	return NumberColor.Sprintf("%d.", num)
}

// ColorizeTitle formats main title
func ColorizeTitle(title string) string {
	return TitleColor.Sprintf("📝 %s", title)
}

// ColorizeHeader formats section header
func ColorizeHeader(header string) string {
	return HeaderColor.Sprintf("📋 %s", header)
}

// ColorizeSuccess formats success message
func ColorizeSuccess(msg string) string {
	return SuccessColor.Sprintf("✅ %s", msg)
}

// ColorizeError formats error message
func ColorizeError(msg string) string {
	return ErrorColor.Sprintf("❌ %s", msg)
}

// ColorizeWarning formats warning message
func ColorizeWarning(msg string) string {
	return WarningColor.Sprintf("⚠️  %s", msg)
}

// ColorizeInfo formats info message
func ColorizeInfo(msg string) string {
	return InfoColor.Sprintf("ℹ️  %s", msg)
}

// EnableColors enables color output
func EnableColors(enable bool) {
	color.NoColor = !enable
}
