package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Victory banner ASCII art frames
var victoryBanners = []string{
	`
  ██╗    ██╗██╗███╗   ██╗███╗   ██╗███████╗██████╗ ██╗
  ██║    ██║██║████╗  ██║████╗  ██║██╔════╝██╔══██╗██║
  ██║ █╗ ██║██║██╔██╗ ██║██╔██╗ ██║█████╗  ██████╔╝██║
  ██║███╗██║██║██║╚██╗██║██║╚██╗██║██╔══╝  ██╔══██╗╚═╝
  ╚███╔███╔╝██║██║ ╚████║██║ ╚████║███████╗██║  ██║██╗
   ╚══╝╚══╝ ╚═╝╚═╝  ╚═══╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝
`,
	`
  ░██╗░░░██╗░██╗███╗░░██╗███╗░░██╗███████╗██████╗░██╗
  ░██║░░░██║░██║████╗░██║████╗░██║██╔════╝██╔══██╗██║
  ░██║░█░██║░██║██╔██╗██║██╔██╗██║█████╗░░██████╔╝██║
  ░██║█████║░██║██║╚████║██║╚████║██╔══╝░░██╔══██╗╚═╝
  ░╚███╔███╔╝░██║██║░╚███║██║░╚███║███████╗██║░░██║██╗
  ░░╚══╝╚══╝░░╚═╝╚═╝░░╚══╝╚═╝░░╚══╝╚══════╝╚═╝░░╚═╝╚═╝
`,
}

// Confetti characters
var confettiChars = []string{"★", "✦", "◆", "●", "▲", "♦", "✶", "◉", "✸"}

// Trophy ASCII art
var trophyArt = `
       ___________
      '._==_==_=_.'
      .-\:      /-.
     | (|:.     |) |
      '-|:.     |-'
        \::.    /
         '::. .'
           ) (
         _.' '._
        '-------'
`

// RenderVictoryScreen renders the full victory celebration
func RenderVictoryScreen(winnerName, winnerAbbr string, winnerScore, loserScore int, frame int, width, height int) string {
	var sb strings.Builder

	// Get team color
	color := "255"
	if c, ok := teamColors[winnerAbbr]; ok {
		color = c
	}

	// Styles
	bannerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	// Alternate colors for flashing effect
	colors := []string{color, "226", "255", "196", "46"}
	flashColor := colors[frame%len(colors)]

	teamStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(flashColor)).
		Bold(true).
		Background(lipgloss.Color("235")).
		Padding(1, 3)

	scoreStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	trophyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("220")).
		Bold(true)

	// Generate confetti line
	confettiLine := renderConfettiLine(width, frame)
	confettiStyle := lipgloss.NewStyle().Bold(true)

	// Top confetti
	sb.WriteString(confettiStyle.Render(confettiLine) + "\n")
	sb.WriteString(confettiStyle.Render(confettiLine) + "\n\n")

	// Winner banner (alternating frames)
	bannerIdx := frame % len(victoryBanners)
	sb.WriteString(bannerStyle.Render(victoryBanners[bannerIdx]) + "\n\n")

	// Trophy
	sb.WriteString(trophyStyle.Render(trophyArt) + "\n")

	// Team name
	teamDisplay := fmt.Sprintf("🏈  %s  🏈", winnerName)
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, teamStyle.Render(teamDisplay)) + "\n\n")

	// Final score
	scoreDisplay := fmt.Sprintf("FINAL SCORE: %d - %d", winnerScore, loserScore)
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, scoreStyle.Render(scoreDisplay)) + "\n\n")

	// Dancing mascots
	mascot1 := RenderMascotWithState(winnerAbbr, frame, MascotCelebrating, false)
	mascot2 := RenderMascotWithState(winnerAbbr, frame+1, MascotCelebrating, false)
	mascot3 := RenderMascotWithState(winnerAbbr, frame+2, MascotCelebrating, false)
	mascotsRow := lipgloss.JoinHorizontal(lipgloss.Top, "    ", mascot1, "        ", mascot2, "        ", mascot3, "    ")
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, mascotsRow) + "\n\n")

	// More fireworks
	fw1 := RenderFireworks(frame, false)
	fw2 := RenderFireworks(frame+2, false)
	fw3 := RenderFireworks(frame+4, false)
	fw4 := RenderFireworks(frame+1, false)
	fwRow := lipgloss.JoinHorizontal(lipgloss.Top, fw1, "  ", fw2, "  ", fw3, "  ", fw4)
	sb.WriteString(lipgloss.PlaceHorizontal(width, lipgloss.Center, fwRow) + "\n")

	// Bottom confetti
	sb.WriteString("\n" + confettiStyle.Render(confettiLine) + "\n")

	// Exit hint
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	sb.WriteString("\n" + lipgloss.PlaceHorizontal(width, lipgloss.Center, hintStyle.Render("Press q to exit")) + "\n")

	return sb.String()
}

// renderConfettiLine creates a line of random confetti
func renderConfettiLine(width int, frame int) string {
	var sb strings.Builder
	colors := []string{"196", "226", "21", "201", "46", "208", "51", "213"}

	for i := 0; i < width; i++ {
		if (i+frame)%3 == 0 {
			charIdx := (i + frame) % len(confettiChars)
			colorIdx := (i + frame/2) % len(colors)
			style := lipgloss.NewStyle().Foreground(lipgloss.Color(colors[colorIdx]))
			sb.WriteString(style.Render(confettiChars[charIdx]))
		} else {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}
