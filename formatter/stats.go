package formatter

import (
	"fmt"
	"strings"

	"nfl-scores/models"

	"github.com/charmbracelet/lipgloss"
)

// FormatGameStats renders game statistics
func (f *TerminalFormatter) FormatGameStats(stats *models.GameStats) string {
	if f.plain {
		return f.formatStatsPlain(stats)
	}
	return f.formatStatsStyled(stats)
}

func (f *TerminalFormatter) formatStatsPlain(stats *models.GameStats) string {
	var sb strings.Builder
	g := stats.Game

	line := strings.Repeat("=", 76)
	sb.WriteString("\n" + line + "\n")
	fmt.Fprintf(&sb, "  %s %d  @  %s %d  -  %s\n",
		g.AwayTeam.Name, g.AwayTeam.Score,
		g.HomeTeam.Name, g.HomeTeam.Score,
		g.StatusText)
	sb.WriteString(line + "\n\n")

	// Team comparison
	sb.WriteString("  TEAM STATS\n")
	sb.WriteString("  " + strings.Repeat("-", 72) + "\n")
	fmt.Fprintf(&sb, "  %-20s %12s %12s\n", "", stats.AwayStats.TeamAbbr, stats.HomeStats.TeamAbbr)

	teamStatLabels := []struct{ key, label string }{
		{"totalYards", "Total Yards"},
		{"netPassingYards", "Passing Yards"},
		{"rushingYards", "Rushing Yards"},
		{"firstDowns", "First Downs"},
		{"thirdDownEff", "3rd Down"},
		{"turnovers", "Turnovers"},
		{"possession", "Possession"},
	}

	for _, s := range teamStatLabels {
		away := stats.AwayStats.Totals[s.key]
		home := stats.HomeStats.Totals[s.key]
		if away != "" || home != "" {
			fmt.Fprintf(&sb, "  %-20s %12s %12s\n", s.label, away, home)
		}
	}

	// Player stats by category
	categories := []string{"passing", "rushing", "receiving"}
	for _, cat := range categories {
		fmt.Fprintf(&sb, "\n  %s\n", strings.ToUpper(cat))
		sb.WriteString("  " + strings.Repeat("-", 72) + "\n")
		f.formatPlayerCategoryPlain(&sb, stats.AwayStats, cat)
		f.formatPlayerCategoryPlain(&sb, stats.HomeStats, cat)
	}

	sb.WriteString("\n" + line + "\n")
	return sb.String()
}

func (f *TerminalFormatter) formatPlayerCategoryPlain(sb *strings.Builder, team models.TeamStats, category string) {
	for _, cat := range team.PlayerStats {
		if cat.Category == category && len(cat.Players) > 0 {
			fmt.Fprintf(sb, "  %s:\n", team.TeamAbbr)

			// Show more columns for passing (to include RTG)
			numCols := 5
			if category == "passing" {
				numCols = 8
			}

			// Header
			fmt.Fprintf(sb, "    %-18s", "Player")
			for i := 0; i < min(numCols, len(cat.Labels)); i++ {
				fmt.Fprintf(sb, " %7s", cat.Labels[i])
			}
			sb.WriteString("\n")
			// Players
			for _, p := range cat.Players {
				fmt.Fprintf(sb, "    %-18s", truncateName(p.Name, 18))
				for i := 0; i < min(numCols, len(p.Stats)); i++ {
					fmt.Fprintf(sb, " %7s", p.Stats[i])
				}
				sb.WriteString("\n")
			}
		}
	}
}

func (f *TerminalFormatter) formatStatsStyled(stats *models.GameStats) string {
	var sb strings.Builder
	g := stats.Game

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	teamStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226"))

	scoreStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226"))

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	playerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	border := borderStyle.Render(strings.Repeat("━", 76))

	sb.WriteString("\n" + border + "\n")
	fmt.Fprintf(&sb, "  %s %s  @  %s %s  -  %s\n",
		teamStyle.Render(g.AwayTeam.Abbreviation),
		scoreStyle.Render(fmt.Sprintf("%d", g.AwayTeam.Score)),
		teamStyle.Render(g.HomeTeam.Abbreviation),
		scoreStyle.Render(fmt.Sprintf("%d", g.HomeTeam.Score)),
		labelStyle.Render(g.StatusText))
	sb.WriteString(border + "\n\n")

	// Team comparison header
	sb.WriteString("  " + headerStyle.Render("📊 TEAM STATS") + "\n")
	sb.WriteString("  " + borderStyle.Render(strings.Repeat("─", 72)) + "\n")

	// Column headers - use same widths as data rows
	sb.WriteString("  " + padRight("", 20) + " ")
	sb.WriteString(teamStyle.Render(padLeft(stats.AwayStats.TeamAbbr, 12)) + " ")
	sb.WriteString(teamStyle.Render(padLeft(stats.HomeStats.TeamAbbr, 12)) + "\n")

	teamStatLabels := []struct{ key, label string }{
		{"totalYards", "Total Yards"},
		{"netPassingYards", "Passing Yards"},
		{"rushingYards", "Rushing Yards"},
		{"firstDowns", "First Downs"},
		{"thirdDownEff", "3rd Down Eff"},
		{"turnovers", "Turnovers"},
		{"possession", "Time of Poss"},
	}

	for _, s := range teamStatLabels {
		away := stats.AwayStats.Totals[s.key]
		home := stats.HomeStats.Totals[s.key]
		if away != "" || home != "" {
			sb.WriteString("  " + labelStyle.Render(padRight(s.label, 20)) + " ")
			sb.WriteString(valueStyle.Render(padLeft(away, 12)) + " ")
			sb.WriteString(valueStyle.Render(padLeft(home, 12)) + "\n")
		}
	}

	// Player stats
	catInfo := []struct {
		name  string
		icon  string
		title string
	}{
		{"passing", "🏈", "PASSING"},
		{"rushing", "🏃", "RUSHING"},
		{"receiving", "🎯", "RECEIVING"},
	}

	for _, ci := range catInfo {
		sb.WriteString("\n  " + headerStyle.Render(ci.icon+" "+ci.title) + "\n")
		sb.WriteString("  " + borderStyle.Render(strings.Repeat("─", 72)) + "\n")
		f.formatPlayerCategoryStyled(&sb, stats.AwayStats, ci.name, teamStyle, playerStyle, labelStyle)
		f.formatPlayerCategoryStyled(&sb, stats.HomeStats, ci.name, teamStyle, playerStyle, labelStyle)
	}

	sb.WriteString("\n" + border + "\n")
	return sb.String()
}

func (f *TerminalFormatter) formatPlayerCategoryStyled(sb *strings.Builder, team models.TeamStats, category string,
	teamStyle, playerStyle, labelStyle lipgloss.Style) {
	for _, cat := range team.PlayerStats {
		if cat.Category == category && len(cat.Players) > 0 {
			// Show more columns for passing (to include RTG)
			numCols := 5
			if category == "passing" {
				numCols = 8 // C/ATT, YDS, AVG, TD, INT, SACKS, QBR, RTG
			}

			// Team header with column labels
			sb.WriteString("  " + teamStyle.Render(padRight(team.TeamAbbr, 18)))
			for i := 0; i < min(numCols, len(cat.Labels)); i++ {
				sb.WriteString(" " + labelStyle.Render(padLeft(cat.Labels[i], 7)))
			}
			sb.WriteString("\n")

			// Player rows
			for _, p := range cat.Players {
				sb.WriteString("    " + playerStyle.Render(padRight(truncateName(p.Name, 16), 16)))
				for i := 0; i < min(numCols, len(p.Stats)); i++ {
					sb.WriteString(" " + valueStyle.Render(padLeft(p.Stats[i], 7)))
				}
				sb.WriteString("\n")
			}
		}
	}
}

var valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))

func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func truncateName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	return name[:maxLen-1] + "…"
}
