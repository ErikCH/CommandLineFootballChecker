package ui

import (
	"fmt"
	"strings"
	"time"

	"nfl-scores/models"
	"nfl-scores/service"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Messages
type tickMsg time.Time
type gameDataMsg *models.GameSummary
type errorMsg error
type mascotTickMsg time.Time

// Model holds the UI state
type Model struct {
	gameID         string
	service        *service.ScoreService
	summary        *models.GameSummary
	prevSummary    *models.GameSummary
	spinner        spinner.Model
	loading        bool
	err            error
	plain          bool
	lastPlayID     string
	flashScore     bool
	flashPlay      bool
	width          int
	height         int
	selectedPlay   int // -1 means no play selected, >= 0 is index
	expandedPlay   int // -1 means no play expanded
	playsStartY    int // Y coordinate where plays list starts
	showMascot     bool
	mascotFrame    int
	mascotState    MascotState
	lastPossession string
	showFireworks  bool
}

// NewModel creates a new live game UI model
func NewModel(gameID string, svc *service.ScoreService, plain bool) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		gameID:       gameID,
		service:      svc,
		spinner:      s,
		loading:      true,
		plain:        plain,
		width:        80,
		height:       24,
		selectedPlay: -1,
		expandedPlay: -1,
		showMascot:   false,
	}
}

// NewModelWithMascot creates a new live game UI model with mascot animation
func NewModelWithMascot(gameID string, svc *service.ScoreService, plain bool) Model {
	m := NewModel(gameID, svc, plain)
	m.showMascot = true
	return m
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.spinner.Tick,
		fetchGameDataCmd(m.gameID, m.service),
		tickCmd(),
	}
	if m.showMascot {
		cmds = append(cmds, mascotTickCmd())
	}
	return tea.Batch(cmds...)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			if m.expandedPlay >= 0 {
				// Close expanded play first
				m.expandedPlay = -1
				return m, nil
			}
			return m, tea.Quit
		case "up", "k":
			if m.summary != nil && len(m.summary.RecentPlays) > 0 {
				if m.selectedPlay > 0 {
					m.selectedPlay--
				} else {
					m.selectedPlay = len(m.summary.RecentPlays) - 1
				}
			}
		case "down", "j":
			if m.summary != nil && len(m.summary.RecentPlays) > 0 {
				if m.selectedPlay < len(m.summary.RecentPlays)-1 {
					m.selectedPlay++
				} else {
					m.selectedPlay = 0
				}
			}
		case "enter", " ":
			if m.selectedPlay >= 0 && m.selectedPlay < len(m.summary.RecentPlays) {
				if m.expandedPlay == m.selectedPlay {
					m.expandedPlay = -1 // Toggle off
				} else {
					m.expandedPlay = m.selectedPlay // Expand
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			// Plays start after: header(4) + field(9) + situation(3) + plays header(2) = ~18 rows
			// But this varies, so we use a range check
			playsStartRow := 18
			playRow := msg.Y - playsStartRow

			if playRow >= 0 && m.summary != nil && playRow < len(m.summary.RecentPlays) {
				m.selectedPlay = playRow
				if m.expandedPlay == playRow {
					m.expandedPlay = -1
				} else {
					m.expandedPlay = playRow
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tickMsg:
		return m, tea.Batch(fetchGameDataCmd(m.gameID, m.service), tickCmd())

	case mascotTickMsg:
		m.mascotFrame++
		return m, mascotTickCmd()

	case gameDataMsg:
		m.loading = false
		m.prevSummary = m.summary
		m.summary = msg

		// Track current play ID and possession changes
		if m.summary != nil && m.summary.CurrentPlay != nil {
			currentPossession := m.summary.CurrentPlay.Possession

			// Check for changes to trigger animations
			if m.prevSummary != nil {
				// Score changed - celebrate with fireworks!
				if m.summary.Game.HomeTeam.Score != m.prevSummary.Game.HomeTeam.Score ||
					m.summary.Game.AwayTeam.Score != m.prevSummary.Game.AwayTeam.Score {
					m.flashScore = true
					m.showFireworks = true
					m.mascotState = MascotCelebrating
					return m, clearCelebration()
				}

				// Turnover detection - possession changed without a score
				if m.lastPossession != "" && currentPossession != m.lastPossession {
					m.mascotState = MascotSad
					return m, clearSadMascot()
				}

				// New play
				if m.summary.CurrentPlay.ID != m.lastPlayID {
					m.flashPlay = true
					m.lastPlayID = m.summary.CurrentPlay.ID
					return m, clearFlashPlay()
				}
			}
			m.lastPlayID = m.summary.CurrentPlay.ID
			m.lastPossession = currentPossession
		}

	case clearCelebrationMsg:
		m.flashScore = false
		m.showFireworks = false
		m.mascotState = MascotNormal

	case clearSadMascotMsg:
		m.mascotState = MascotNormal

	case clearFlashScoreMsg:
		m.flashScore = false

	case clearFlashPlayMsg:
		m.flashPlay = false

	case errorMsg:
		m.err = msg
		m.loading = false
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.loading && m.summary == nil {
		return fmt.Sprintf("\n\n   %s Loading game data...\n", m.spinner.View())
	}

	if m.err != nil {
		return fmt.Sprintf("\n\n   Error: %v\n\n   Press q to quit.\n", m.err)
	}

	if m.summary == nil {
		return "\n\n   No game data available.\n\n   Press q to quit.\n"
	}

	// Check if game is final - show victory screen in mascot mode
	if m.showMascot && m.summary.Game.Status == models.StatusFinal {
		return m.renderVictory()
	}

	if m.plain {
		return m.renderPlain()
	}
	return m.renderStyled()
}

func (m Model) renderVictory() string {
	g := m.summary.Game

	// Determine winner
	var winnerName, winnerAbbr string
	var winnerScore, loserScore int

	if g.HomeTeam.Score > g.AwayTeam.Score {
		winnerName = g.HomeTeam.Name
		winnerAbbr = g.HomeTeam.Abbreviation
		winnerScore = g.HomeTeam.Score
		loserScore = g.AwayTeam.Score
	} else {
		winnerName = g.AwayTeam.Name
		winnerAbbr = g.AwayTeam.Abbreviation
		winnerScore = g.AwayTeam.Score
		loserScore = g.HomeTeam.Score
	}

	return RenderVictoryScreen(winnerName, winnerAbbr, winnerScore, loserScore, m.mascotFrame, m.width, m.height)
}

func (m Model) renderPlain() string {
	var sb strings.Builder
	g := m.summary.Game

	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 70) + "\n")
	sb.WriteString(fmt.Sprintf("  %s %d  @  %s %d   [%s]\n",
		g.AwayTeam.Abbreviation, g.AwayTeam.Score,
		g.HomeTeam.Abbreviation, g.HomeTeam.Score,
		g.StatusText))
	sb.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Field - use actual yards to endzone
	yardsToEndzone := m.summary.YardsToEndzone
	if yardsToEndzone == 0 {
		yardsToEndzone = 50
	}
	possession := ""
	if m.summary.CurrentPlay != nil {
		possession = m.summary.CurrentPlay.Possession
	}
	sb.WriteString(RenderField(yardsToEndzone, possession, true))

	// Situation
	if m.summary.Situation != "" {
		sb.WriteString(fmt.Sprintf("\n  SITUATION: %s\n", m.summary.Situation))
	}

	// Recent plays
	sb.WriteString("\n  RECENT PLAYS:\n")
	sb.WriteString("  " + strings.Repeat("-", 66) + "\n")
	for _, play := range m.summary.RecentPlays {
		text := strings.ReplaceAll(play.Text, "\n", " ")
		if len(text) > 55 {
			text = text[:52] + "..."
		}
		sb.WriteString(fmt.Sprintf("  Q%d %5s │ %s\n", play.Period, play.Clock, text))
	}

	sb.WriteString("\n  Press q to quit\n")
	return sb.String()
}

func (m Model) renderStyled() string {
	var sb strings.Builder
	g := m.summary.Game

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("235")).
		Padding(0, 2).
		MarginBottom(1)

	teamStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255"))

	scoreStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226"))

	if m.flashScore {
		scoreStyle = scoreStyle.
			Background(lipgloss.Color("196")).
			Foreground(lipgloss.Color("255"))
	}

	liveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Bold(true)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	situationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	playStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	newPlayStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("40")).
		Bold(true)

	clockStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	scoringStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true).
		Background(lipgloss.Color("22"))

	border := borderStyle.Render(strings.Repeat("━", 62))

	// Score header bar - always at top (no leading newline)
	awayScore := scoreStyle.Render(fmt.Sprintf("%d", g.AwayTeam.Score))
	homeScore := scoreStyle.Render(fmt.Sprintf("%d", g.HomeTeam.Score))

	liveIndicator := ""
	switch g.Status {
	case models.StatusInProgress:
		liveIndicator = liveStyle.Render(" ● LIVE")
	case models.StatusFinal:
		liveIndicator = lipgloss.NewStyle().Foreground(lipgloss.Color("40")).Render(" ✓ FINAL")
	}

	// Score line at very top
	scoreLine := fmt.Sprintf("  %s %s  @  %s %s   %s %s",
		teamStyle.Render(g.AwayTeam.Abbreviation),
		awayScore,
		teamStyle.Render(g.HomeTeam.Abbreviation),
		homeScore,
		statusStyle.Render(g.StatusText),
		liveIndicator,
	)

	sb.WriteString(border + "\n")
	sb.WriteString(scoreLine + "\n")
	sb.WriteString(border + "\n")

	// Show fireworks above field when celebrating
	if m.showMascot && m.showFireworks {
		fw1 := RenderFireworks(m.mascotFrame, false)
		fw2 := RenderFireworks(m.mascotFrame+2, false)
		fw3 := RenderFireworks(m.mascotFrame+4, false)
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, "  ", fw1, "    ", fw2, "    ", fw3) + "\n")
	}

	// Football field - use actual yards to endzone
	yardsToEndzone := m.summary.YardsToEndzone
	if yardsToEndzone == 0 {
		yardsToEndzone = 50 // Default to midfield if unknown
	}
	possession := ""
	if m.summary.CurrentPlay != nil {
		possession = m.summary.CurrentPlay.Possession
	}
	sb.WriteString("\n")

	// Render field with optional mascot
	fieldStr := RenderField(yardsToEndzone, possession, false)
	if m.showMascot && possession != "" {
		mascotStr := RenderMascotWithState(possession, m.mascotFrame, m.mascotState, false)
		// Join field and mascot side by side
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, fieldStr, "  ", mascotStr))
	} else {
		sb.WriteString(fieldStr)
	}

	// Current situation
	if m.summary.Situation != "" {
		sb.WriteString("\n  " + headerStyle.Render("󰈍 SITUATION") + "\n")
		sb.WriteString("  " + situationStyle.Render(m.summary.Situation) + "\n")
	}

	// Recent plays
	sb.WriteString("\n  " + headerStyle.Render(" PLAYS (↑↓ to select, Enter to expand)") + "\n")
	sb.WriteString("  " + borderStyle.Render(strings.Repeat("─", 58)) + "\n")

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("236")).
		Foreground(lipgloss.Color("255"))

	expandedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		PaddingLeft(12)

	for i, play := range m.summary.RecentPlays {
		text := strings.ReplaceAll(play.Text, "\n", " ")
		isSelected := i == m.selectedPlay
		isExpanded := i == m.expandedPlay

		timeInfo := clockStyle.Render(fmt.Sprintf("Q%d %5s", play.Period, play.Clock))

		var playText string
		displayText := text
		if !isExpanded && len(displayText) > 50 {
			displayText = displayText[:47] + "..."
		}

		if play.ScoringPlay {
			playText = scoringStyle.Render("󰸞 " + displayText)
		} else if i == 0 && m.flashPlay {
			playText = newPlayStyle.Render("► " + displayText)
		} else {
			playText = playStyle.Render(displayText)
		}

		line := fmt.Sprintf("  %s │ %s", timeInfo, playText)

		if isSelected {
			// Highlight selected row
			line = selectedStyle.Render(line)
		}

		sb.WriteString(line + "\n")

		// Show full text if expanded
		if isExpanded && len(text) > 50 {
			fullText := strings.ReplaceAll(play.Text, "\n", "\n            ")
			sb.WriteString(expandedStyle.Render("└─ "+fullText) + "\n")
		}
	}

	sb.WriteString("\n" + border + "\n")
	sb.WriteString(statusStyle.Render("  Press q to quit • Auto-refreshing every 10s") + "\n")

	return sb.String()
}

// Commands
func fetchGameDataCmd(gameID string, svc *service.ScoreService) tea.Cmd {
	return func() tea.Msg {
		summary, err := svc.GetGameSummary(gameID)
		if err != nil {
			return errorMsg(err)
		}
		return gameDataMsg(summary)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type clearFlashScoreMsg struct{}
type clearFlashPlayMsg struct{}
type clearCelebrationMsg struct{}
type clearSadMascotMsg struct{}

func clearFlashScore() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return clearFlashScoreMsg{}
	})
}

func clearFlashPlay() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return clearFlashPlayMsg{}
	})
}

func clearCelebration() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return clearCelebrationMsg{}
	})
}

func clearSadMascot() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return clearSadMascotMsg{}
	})
}

func mascotTickCmd() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return mascotTickMsg(t)
	})
}
