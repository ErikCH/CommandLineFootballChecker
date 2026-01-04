package formatter

import (
	"fmt"
	"strings"

	"nfl-scores/models"

	"github.com/charmbracelet/lipgloss"
)

// Icons (Nerd Font)
const (
	iconFootball = "󰈍" // nf-md-football
	iconLive     = ""  // nf-fa-circle (for live indicator)
	iconClock    = ""  // nf-fa-clock_o
	iconCheck    = ""  // nf-fa-check
	iconAt       = "@"
)

// TerminalFormatter handles terminal output formatting
type TerminalFormatter struct {
	width int
	plain bool
}

// NewTerminalFormatter creates a formatter with specified width
func NewTerminalFormatter(width int, plain bool) *TerminalFormatter {
	if width <= 0 {
		width = 80
	}
	return &TerminalFormatter{width: width, plain: plain}
}

// FormatScoreboard renders games as formatted terminal output
func (f *TerminalFormatter) FormatScoreboard(games []models.Game) string {
	if f.plain {
		return f.formatPlain(games)
	}
	return f.formatStyled(games)
}

// formatPlain renders without colors/icons
func (f *TerminalFormatter) formatPlain(games []models.Game) string {
	if len(games) == 0 {
		return "No NFL games are currently scheduled."
	}

	var sb strings.Builder
	line := strings.Repeat("=", 76)

	sb.WriteString("\n" + line + "\n")
	sb.WriteString("                            NFL SCORES\n")
	sb.WriteString(line + "\n\n")

	for _, game := range games {
		awayName := truncate(game.AwayTeam.Name, 18)
		homeName := truncate(game.HomeTeam.Name, 18)
		status := truncate(game.StatusText, 12)
		if status == "" {
			status = game.Status.String()
		}
		sb.WriteString(fmt.Sprintf("  %-18s %3d  @  %-18s %3d  [%-12s]\n",
			awayName, game.AwayTeam.Score, homeName, game.HomeTeam.Score, status))
	}

	sb.WriteString("\n" + line + "\n")
	return sb.String()
}

// formatStyled renders with colors and icons
func (f *TerminalFormatter) formatStyled(games []models.Game) string {
	if len(games) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("No NFL games are currently scheduled.")
	}

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("235")).
		Padding(0, 2).
		Width(76).
		Align(lipgloss.Center)

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	teamStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	scoreStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226"))

	liveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	scheduledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	finalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("40"))

	var sb strings.Builder
	border := borderStyle.Render(strings.Repeat("━", 76))

	sb.WriteString("\n")
	sb.WriteString(border + "\n")
	sb.WriteString(headerStyle.Render(iconFootball+"  NFL SCORES  "+iconFootball) + "\n")
	sb.WriteString(border + "\n\n")

	for _, game := range games {
		awayName := truncate(game.AwayTeam.Name, 18)
		homeName := truncate(game.HomeTeam.Name, 18)

		// Format scores
		awayScore := scoreStyle.Render(fmt.Sprintf("%3d", game.AwayTeam.Score))
		homeScore := scoreStyle.Render(fmt.Sprintf("%3d", game.HomeTeam.Score))

		// Format status with icon
		var statusStr string
		statusText := truncate(game.StatusText, 12)
		if statusText == "" {
			statusText = game.Status.String()
		}

		switch game.Status {
		case models.StatusInProgress:
			statusStr = liveStyle.Render(iconLive + " " + statusText)
		case models.StatusFinal:
			statusStr = finalStyle.Render(iconCheck + " " + statusText)
		default:
			statusStr = scheduledStyle.Render(iconClock + " " + statusText)
		}

		// Build line
		line := fmt.Sprintf("  %s %s  %s  %s %s  %s",
			teamStyle.Render(fmt.Sprintf("%-18s", awayName)),
			awayScore,
			lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(iconAt),
			teamStyle.Render(fmt.Sprintf("%-18s", homeName)),
			homeScore,
			statusStr,
		)
		sb.WriteString(line + "\n")
	}

	sb.WriteString("\n" + border + "\n")
	return sb.String()
}

// truncate shortens a string to max length
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// FormatError renders error messages for terminal display
func (f *TerminalFormatter) FormatError(err error) string {
	msg := err.Error()

	var userMsg string
	if strings.Contains(msg, "connection refused") {
		userMsg = "NFL data service is unavailable. Please try again later."
	} else if strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline exceeded") {
		userMsg = "Unable to connect to NFL data service. Please check your internet connection."
	} else if strings.Contains(msg, "parse") || strings.Contains(msg, "json") {
		userMsg = "Received invalid data from NFL service. Please try again."
	} else {
		userMsg = "An unexpected error occurred. Please try again."
	}

	if f.plain {
		return "Error: " + userMsg
	}

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	return errorStyle.Render("✗ " + userMsg)
}
