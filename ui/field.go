package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Field dimensions (scaled down)
const (
	fieldWidth  = 60
	fieldHeight = 11
)

// RenderField creates an ASCII football field with ball position
func RenderField(yardsToEndzone int, possession string, plain bool) string {
	if plain {
		return renderFieldPlain(yardsToEndzone, possession)
	}
	return renderFieldStyled(yardsToEndzone, possession)
}

func renderFieldPlain(yardsToEndzone int, possession string) string {
	var sb strings.Builder

	// Calculate ball position (0-100 yards maps to field width)
	ballPos := int(float64(100-yardsToEndzone) / 100.0 * float64(fieldWidth-4))
	if ballPos < 0 {
		ballPos = 0
	}
	if ballPos > fieldWidth-4 {
		ballPos = fieldWidth - 4
	}

	// Top border
	sb.WriteString("  ╔" + strings.Repeat("═", fieldWidth) + "╗\n")

	// End zone labels
	sb.WriteString("  ║" + centerText("END", 5) + "│")
	sb.WriteString(centerText("", fieldWidth-12))
	sb.WriteString("│" + centerText("END", 5) + "║\n")

	// Yard markers
	sb.WriteString("  ║     │")
	markers := "  10   20   30   40   50   40   30   20   10  "
	sb.WriteString(markers)
	sb.WriteString("│     ║\n")

	// Field with ball
	for i := 0; i < 3; i++ {
		sb.WriteString("  ║     │")
		if i == 1 {
			// Ball row
			line := strings.Repeat(" ", ballPos) + "🏈" + strings.Repeat(" ", fieldWidth-12-ballPos-2)
			sb.WriteString(line)
		} else {
			sb.WriteString(strings.Repeat(" ", fieldWidth-12))
		}
		sb.WriteString("│     ║\n")
	}

	// Yard line markers (dashes)
	sb.WriteString("  ║     │")
	for i := 0; i < fieldWidth-12; i++ {
		if i%5 == 2 {
			sb.WriteString("┼")
		} else {
			sb.WriteString("─")
		}
	}
	sb.WriteString("│     ║\n")

	// Bottom info
	sb.WriteString("  ║     │")
	sb.WriteString(centerText(fmt.Sprintf("← %s BALL →", possession), fieldWidth-12))
	sb.WriteString("│     ║\n")

	// Bottom border
	sb.WriteString("  ╚" + strings.Repeat("═", fieldWidth) + "╝\n")

	return sb.String()
}

func renderFieldStyled(yardsToEndzone int, possession string) string {
	var sb strings.Builder

	// Field layout: | END (5) | playing field (46) | END (5) | = 56 inner + 2 borders = 58
	const (
		endZoneWidth = 5
		playingWidth = 46                                                 // Must match yard markers string length
		innerWidth   = endZoneWidth + 1 + playingWidth + 1 + endZoneWidth // 5 + | + 46 + | + 5 = 58
	)

	// Styles
	fieldStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("34")) // Green

	endZoneStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")). // Red
		Bold(true)

	yardStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255"))

	ballStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("208")). // Orange
		Bold(true)

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("34"))

	// Calculate ball position (emoji is 2 chars wide visually)
	ballPos := int(float64(100-yardsToEndzone) / 100.0 * float64(playingWidth-2))
	if ballPos < 0 {
		ballPos = 0
	}
	if ballPos > playingWidth-2 {
		ballPos = playingWidth - 2
	}

	// Top border
	sb.WriteString("  " + borderStyle.Render("╔"+strings.Repeat("═", innerWidth)+"╗") + "\n")

	// End zone row
	sb.WriteString("  " + borderStyle.Render("║"))
	sb.WriteString(endZoneStyle.Render(centerText("END", endZoneWidth)))
	sb.WriteString(borderStyle.Render("│"))
	sb.WriteString(fieldStyle.Render(strings.Repeat(" ", playingWidth)))
	sb.WriteString(borderStyle.Render("│"))
	sb.WriteString(endZoneStyle.Render(centerText("END", endZoneWidth)))
	sb.WriteString(borderStyle.Render("║") + "\n")

	// Yard markers (46 chars: " 10   20   30   40   50   40   30   20   10  ")
	sb.WriteString("  " + borderStyle.Render("║"))
	sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
	sb.WriteString(borderStyle.Render("│"))
	markers := " 10   20   30   40   50   40   30   20   10  "
	sb.WriteString(yardStyle.Render(markers))
	sb.WriteString(borderStyle.Render("│"))
	sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
	sb.WriteString(borderStyle.Render("║") + "\n")

	// Field rows with ball
	for i := 0; i < 3; i++ {
		sb.WriteString("  " + borderStyle.Render("║"))
		sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
		sb.WriteString(borderStyle.Render("│"))

		if i == 1 {
			// Ball row - emoji takes 2 display columns
			before := strings.Repeat("░", ballPos)
			ball := "🏈"
			after := strings.Repeat("░", playingWidth-ballPos-2)
			sb.WriteString(fieldStyle.Render(before) + ballStyle.Render(ball) + fieldStyle.Render(after))
		} else {
			sb.WriteString(fieldStyle.Render(strings.Repeat("░", playingWidth)))
		}

		sb.WriteString(borderStyle.Render("│"))
		sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
		sb.WriteString(borderStyle.Render("║") + "\n")
	}

	// Yard line markers
	sb.WriteString("  " + borderStyle.Render("║"))
	sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
	sb.WriteString(borderStyle.Render("│"))
	var lineMarkers strings.Builder
	for i := 0; i < playingWidth; i++ {
		if i%5 == 2 {
			lineMarkers.WriteString("┼")
		} else {
			lineMarkers.WriteString("─")
		}
	}
	sb.WriteString(fieldStyle.Render(lineMarkers.String()))
	sb.WriteString(borderStyle.Render("│"))
	sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
	sb.WriteString(borderStyle.Render("║") + "\n")

	// Possession indicator
	sb.WriteString("  " + borderStyle.Render("║"))
	sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
	sb.WriteString(borderStyle.Render("│"))
	possText := fmt.Sprintf("◄── %s BALL ──►", possession)
	sb.WriteString(yardStyle.Render(centerText(possText, playingWidth)))
	sb.WriteString(borderStyle.Render("│"))
	sb.WriteString(fieldStyle.Render(strings.Repeat(" ", endZoneWidth)))
	sb.WriteString(borderStyle.Render("║") + "\n")

	// Bottom border
	sb.WriteString("  " + borderStyle.Render("╚"+strings.Repeat("═", innerWidth)+"╝") + "\n")

	return sb.String()
}

func centerText(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	padding := (width - len(s)) / 2
	return strings.Repeat(" ", padding) + s + strings.Repeat(" ", width-len(s)-padding)
}
