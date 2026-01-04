package formatter

import (
	"fmt"
	"strings"

	"nfl-scores/models"

	"github.com/charmbracelet/lipgloss"
)

// FormatLiveGame renders a live game with play-by-play
func (f *TerminalFormatter) FormatLiveGame(summary *models.GameSummary) string {
	if f.plain {
		return f.formatLivePlain(summary)
	}
	return f.formatLiveStyled(summary)
}

func (f *TerminalFormatter) formatLivePlain(summary *models.GameSummary) string {
	var sb strings.Builder
	g := summary.Game

	line := strings.Repeat("=", 76)
	sb.WriteString("\n" + line + "\n")
	sb.WriteString(fmt.Sprintf("  %s %d  @  %s %d\n",
		g.AwayTeam.Name, g.AwayTeam.Score,
		g.HomeTeam.Name, g.HomeTeam.Score))
	sb.WriteString(fmt.Sprintf("  %s\n", g.StatusText))
	sb.WriteString(line + "\n\n")

	if summary.Situation != "" {
		sb.WriteString(fmt.Sprintf("  Situation: %s\n", summary.Situation))
		if summary.CurrentPlay != nil {
			sb.WriteString(fmt.Sprintf("  Possession: %s\n\n", summary.CurrentPlay.Possession))
		}
	}

	sb.WriteString("  RECENT PLAYS:\n")
	sb.WriteString("  " + strings.Repeat("-", 72) + "\n")

	for _, play := range summary.RecentPlays {
		playText := truncate(play.Text, 68)
		// Replace newlines with spaces
		playText = strings.ReplaceAll(playText, "\n", " ")
		sb.WriteString(fmt.Sprintf("  Q%d %s | %s\n", play.Period, play.Clock, playText))
	}

	sb.WriteString("\n" + line + "\n")
	return sb.String()
}

func (f *TerminalFormatter) formatLiveStyled(summary *models.GameSummary) string {
	var sb strings.Builder
	g := summary.Game

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	teamStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255"))

	scoreStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226"))

	liveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true).
		Blink(true)

	situationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	playStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	clockStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	scoringStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("40")).
		Bold(true)

	border := borderStyle.Render(strings.Repeat("━", 76))

	sb.WriteString("\n" + border + "\n")

	// Score header
	scoreHeader := fmt.Sprintf("  %s %s  %s  %s %s",
		teamStyle.Render(g.AwayTeam.Abbreviation),
		scoreStyle.Render(fmt.Sprintf("%d", g.AwayTeam.Score)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("@"),
		teamStyle.Render(g.HomeTeam.Abbreviation),
		scoreStyle.Render(fmt.Sprintf("%d", g.HomeTeam.Score)),
	)
	sb.WriteString(scoreHeader + "  " + liveStyle.Render(iconLive+" LIVE") + "\n")
	sb.WriteString("  " + clockStyle.Render(g.StatusText) + "\n")
	sb.WriteString(border + "\n\n")

	// Current situation
	if summary.Situation != "" {
		sb.WriteString("  " + headerStyle.Render("󰈍 SITUATION") + "\n")
		sb.WriteString("  " + situationStyle.Render(summary.Situation))
		if summary.CurrentPlay != nil {
			sb.WriteString("  " + clockStyle.Render("Ball: "+summary.CurrentPlay.Possession))
		}
		sb.WriteString("\n\n")
	}

	// Recent plays
	sb.WriteString("  " + headerStyle.Render(" RECENT PLAYS") + "\n")
	sb.WriteString("  " + borderStyle.Render(strings.Repeat("─", 72)) + "\n")

	for _, play := range summary.RecentPlays {
		playText := truncate(play.Text, 60)
		playText = strings.ReplaceAll(playText, "\n", " ")

		timeInfo := clockStyle.Render(fmt.Sprintf("Q%d %5s", play.Period, play.Clock))

		var playLine string
		if play.ScoringPlay {
			playLine = scoringStyle.Render("󰸞 " + playText)
		} else {
			playLine = playStyle.Render(playText)
		}

		sb.WriteString(fmt.Sprintf("  %s │ %s\n", timeInfo, playLine))
	}

	sb.WriteString("\n" + border + "\n")
	sb.WriteString(clockStyle.Render("  Press Ctrl+C to exit • Refreshing every 10s") + "\n")

	return sb.String()
}

// FormatGameSelection renders a list of live games for selection
func (f *TerminalFormatter) FormatGameSelection(games []models.Game) string {
	if len(games) == 0 {
		if f.plain {
			return "No live games available."
		}
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Render("No live games available.")
	}

	var sb strings.Builder

	if f.plain {
		sb.WriteString("\nSelect a live game to track:\n\n")
		for i, g := range games {
			sb.WriteString(fmt.Sprintf("  [%d] %s %d @ %s %d (%s)\n",
				i+1, g.AwayTeam.Abbreviation, g.AwayTeam.Score,
				g.HomeTeam.Abbreviation, g.HomeTeam.Score, g.StatusText))
		}
		sb.WriteString("\nEnter number: ")
	} else {
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39"))

		numStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

		teamStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

		scoreStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("40"))

		statusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

		sb.WriteString("\n" + headerStyle.Render(iconFootball+" Select a live game to track:") + "\n\n")

		for i, g := range games {
			num := numStyle.Render(fmt.Sprintf("[%d]", i+1))
			away := teamStyle.Render(g.AwayTeam.Abbreviation)
			home := teamStyle.Render(g.HomeTeam.Abbreviation)
			awayScore := scoreStyle.Render(fmt.Sprintf("%d", g.AwayTeam.Score))
			homeScore := scoreStyle.Render(fmt.Sprintf("%d", g.HomeTeam.Score))
			status := statusStyle.Render(fmt.Sprintf("(%s)", g.StatusText))

			sb.WriteString(fmt.Sprintf("  %s %s %s @ %s %s %s\n",
				num, away, awayScore, home, homeScore, status))
		}

		sb.WriteString("\n" + statusStyle.Render("Enter number: "))
	}

	return sb.String()
}
