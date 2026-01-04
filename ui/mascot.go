package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MascotState represents the mascot's emotional state
type MascotState int

const (
	MascotNormal MascotState = iota
	MascotCelebrating
	MascotSad
)

// Mascot animation frames - a dancing football player
var mascotFrames = []string{
	`   O
  /|\
  / \`,
	`   O
  \|/
  / \`,
	`   O
  /|/
   |`,
	`   O
  \|\
   |`,
}

// Sad mascot frames - dejected poses
var sadMascotFrames = []string{
	`   O
  /|
  / \`,
	`   o
  /|
  / \`,
}

// Celebrating mascot frames - jumping with arms up
var celebratingFrames = []string{
	`  \O/
   |
  / \`,
	`  \O/
   |
  /|\`,
	`   O
  \|/
   ^`,
	`  \O/
   |
  < >`,
}

// Firework frames
var fireworkFrames = [][]string{
	{
		`    *    `,
		`         `,
		`         `,
	},
	{
		`   ***   `,
		`   *•*   `,
		`         `,
	},
	{
		`  * * *  `,
		`  *•*•*  `,
		`  * * *  `,
	},
	{
		` *  *  * `,
		`  * • *  `,
		` *  *  * `,
	},
	{
		`*   *   *`,
		`    •    `,
		`*   *   *`,
	},
}

// NFL team colors (primary color)
var teamColors = map[string]string{
	"ARI": "161", // Cardinal Red
	"ATL": "196", // Falcon Red
	"BAL": "55",  // Ravens Purple
	"BUF": "21",  // Bills Blue
	"CAR": "39",  // Panthers Blue
	"CHI": "202", // Bears Orange
	"CIN": "208", // Bengals Orange
	"CLE": "208", // Browns Orange
	"DAL": "21",  // Cowboys Blue
	"DEN": "208", // Broncos Orange
	"DET": "39",  // Lions Blue
	"GB":  "28",  // Packers Green
	"HOU": "124", // Texans Red
	"IND": "21",  // Colts Blue
	"JAX": "30",  // Jaguars Teal
	"KC":  "196", // Chiefs Red
	"LAC": "39",  // Chargers Blue
	"LAR": "21",  // Rams Blue
	"LV":  "247", // Raiders Silver
	"MIA": "37",  // Dolphins Aqua
	"MIN": "55",  // Vikings Purple
	"NE":  "21",  // Patriots Blue
	"NO":  "220", // Saints Gold
	"NYG": "21",  // Giants Blue
	"NYJ": "28",  // Jets Green
	"PHI": "30",  // Eagles Green
	"PIT": "220", // Steelers Gold
	"SEA": "28",  // Seahawks Green
	"SF":  "196", // 49ers Red
	"TB":  "196", // Bucs Red
	"TEN": "39",  // Titans Blue
	"WAS": "124", // Commanders Burgundy
}

// RenderMascot returns an animated mascot with team colors
func RenderMascot(teamAbbr string, frame int, plain bool) string {
	return RenderMascotWithState(teamAbbr, frame, MascotNormal, plain)
}

// RenderMascotWithState returns mascot with specific emotional state
func RenderMascotWithState(teamAbbr string, frame int, state MascotState, plain bool) string {
	var frames []string
	switch state {
	case MascotCelebrating:
		frames = celebratingFrames
	case MascotSad:
		frames = sadMascotFrames
	default:
		frames = mascotFrames
	}

	frameIdx := frame % len(frames)
	mascot := frames[frameIdx]

	if plain {
		return teamAbbr + "\n" + mascot
	}

	color := "255" // default white
	if c, ok := teamColors[teamAbbr]; ok {
		color = c
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	return labelStyle.Render(teamAbbr) + "\n" + style.Render(mascot)
}

// RenderFireworks returns animated fireworks
func RenderFireworks(frame int, plain bool) string {
	frameIdx := frame % len(fireworkFrames)
	fw := fireworkFrames[frameIdx]

	if plain {
		return strings.Join(fw, "\n")
	}

	// Cycle through colors for fireworks
	colors := []string{"196", "226", "21", "201", "46", "208"}
	color := colors[frame%len(colors)]

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)

	var lines []string
	for _, line := range fw {
		lines = append(lines, style.Render(line))
	}
	return strings.Join(lines, "\n")
}
