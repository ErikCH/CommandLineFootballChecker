package ui

import (
	"fmt"
	"strings"
	"time"

	"nfl-scores/models"
	"nfl-scores/service"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Replay messages
type replayTickMsg time.Time
type replayDataMsg *models.GameReplay
type replayErrorMsg error

// ReplayModel holds the replay UI state
type ReplayModel struct {
	gameID      string
	service     *service.ScoreService
	replay      *models.GameReplay
	loading     bool
	err         error
	plain       bool
	playIndex   int  // Current play index
	autoPlay    bool // Auto-advance mode
	autoSpeed   int  // Seconds between plays (1-5)
	width       int
	height      int
	showMascot  bool
	mascotFrame int
}

// NewReplayModel creates a new replay UI model
func NewReplayModel(gameID string, svc *service.ScoreService, plain, mascot bool) ReplayModel {
	return ReplayModel{
		gameID:     gameID,
		service:    svc,
		loading:    true,
		plain:      plain,
		playIndex:  0,
		autoPlay:   false,
		autoSpeed:  2,
		width:      80,
		height:     24,
		showMascot: mascot,
	}
}

// Init initializes the replay model
func (m ReplayModel) Init() tea.Cmd {
	return tea.Batch(
		fetchReplayDataCmd(m.gameID, m.service),
		replayMascotTickCmd(),
	)
}

// Update handles messages
func (m ReplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "right", "l", "n":
			// Next play
			if m.replay != nil && m.playIndex < len(m.replay.Plays)-1 {
				m.playIndex++
			}
		case "left", "h", "p":
			// Previous play
			if m.playIndex > 0 {
				m.playIndex--
			}
		case " ":
			// Toggle auto-play
			m.autoPlay = !m.autoPlay
			if m.autoPlay {
				return m, replayAutoTickCmd(m.autoSpeed)
			}
		case "+", "=":
			// Speed up
			if m.autoSpeed > 1 {
				m.autoSpeed--
			}
		case "-", "_":
			// Slow down
			if m.autoSpeed < 5 {
				m.autoSpeed++
			}
		case "home", "0":
			// Go to start
			m.playIndex = 0
		case "end", "$":
			// Go to end
			if m.replay != nil {
				m.playIndex = len(m.replay.Plays) - 1
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case replayTickMsg:
		m.mascotFrame++
		return m, replayMascotTickCmd()

	case replayDataMsg:
		m.loading = false
		m.replay = msg

	case replayErrorMsg:
		m.err = msg
		m.loading = false

	case replayAutoTickMsg:
		if m.autoPlay && m.replay != nil && m.playIndex < len(m.replay.Plays)-1 {
			m.playIndex++
			return m, replayAutoTickCmd(m.autoSpeed)
		}
		m.autoPlay = false
	}

	return m, nil
}

// View renders the replay UI
func (m ReplayModel) View() string {
	if m.loading {
		return "\n\n   Loading game data...\n"
	}

	if m.err != nil {
		return fmt.Sprintf("\n\n   Error: %v\n\n   Press q to quit.\n", m.err)
	}

	if m.replay == nil || len(m.replay.Plays) == 0 {
		return "\n\n   No play data available for this game.\n\n   Press q to quit.\n"
	}

	if m.plain {
		return m.renderPlain()
	}
	return m.renderStyled()
}

func (m ReplayModel) renderPlain() string {
	var sb strings.Builder
	play := m.replay.Plays[m.playIndex]
	g := m.replay.Game

	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 70) + "\n")
	sb.WriteString(fmt.Sprintf("  %s %d  @  %s %d   [REPLAY]\n",
		g.AwayTeam.Abbreviation, play.AwayScore,
		g.HomeTeam.Abbreviation, play.HomeScore))
	sb.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Progress
	sb.WriteString(fmt.Sprintf("  Play %d of %d\n", m.playIndex+1, len(m.replay.Plays)))
	sb.WriteString(fmt.Sprintf("  Q%d %s\n\n", play.Period, play.Clock))

	// Field
	yardsToEndzone := play.YardsToEndzone
	if yardsToEndzone == 0 {
		yardsToEndzone = 50
	}
	sb.WriteString(RenderField(yardsToEndzone, play.Possession, true))

	// Situation
	if play.Down != "" {
		sb.WriteString(fmt.Sprintf("\n  SITUATION: %s\n", play.Down))
	}

	// Play description
	sb.WriteString("\n  PLAY:\n")
	sb.WriteString("  " + strings.Repeat("-", 66) + "\n")
	sb.WriteString(fmt.Sprintf("  %s\n", play.Text))

	// Controls
	sb.WriteString("\n  " + strings.Repeat("-", 66) + "\n")
	autoStatus := "OFF"
	if m.autoPlay {
		autoStatus = fmt.Sprintf("ON (%ds)", m.autoSpeed)
	}
	sb.WriteString(fmt.Sprintf("  ←/→: prev/next | SPACE: auto-play [%s] | +/-: speed | q: quit\n", autoStatus))

	return sb.String()
}

func (m ReplayModel) renderStyled() string {
	var sb strings.Builder
	play := m.replay.Plays[m.playIndex]
	g := m.replay.Game

	// Styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Background(lipgloss.Color("235")).
		Padding(0, 2)

	teamStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("255"))

	scoreStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226"))

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	situationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("214")).
		Bold(true)

	playTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39"))

	replayBadge := lipgloss.NewStyle().
		Foreground(lipgloss.Color("201")).
		Bold(true).
		Render("▶ REPLAY")

	border := borderStyle.Render(strings.Repeat("━", 62))

	// Score header
	scoreLine := fmt.Sprintf("  %s %s  @  %s %s   Q%d %s  %s",
		teamStyle.Render(g.AwayTeam.Abbreviation),
		scoreStyle.Render(fmt.Sprintf("%d", play.AwayScore)),
		teamStyle.Render(g.HomeTeam.Abbreviation),
		scoreStyle.Render(fmt.Sprintf("%d", play.HomeScore)),
		play.Period,
		play.Clock,
		replayBadge,
	)

	sb.WriteString(border + "\n")
	sb.WriteString(scoreLine + "\n")
	sb.WriteString(border + "\n")

	// Progress bar
	progress := float64(m.playIndex+1) / float64(len(m.replay.Plays))
	barWidth := 50
	filled := int(progress * float64(barWidth))
	progressBar := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(strings.Repeat("█", filled))
	progressBar += lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(strings.Repeat("░", barWidth-filled))
	sb.WriteString(fmt.Sprintf("\n  Play %d/%d  [%s]\n", m.playIndex+1, len(m.replay.Plays), progressBar))

	// Field
	yardsToEndzone := play.YardsToEndzone
	if yardsToEndzone == 0 {
		yardsToEndzone = 50
	}

	sb.WriteString("\n")
	fieldStr := RenderField(yardsToEndzone, play.Possession, false)
	if m.showMascot && play.Possession != "" {
		state := MascotNormal
		if play.ScoringPlay {
			state = MascotCelebrating
		}
		mascotStr := RenderMascotWithState(play.Possession, m.mascotFrame, state, false)
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, fieldStr, "  ", mascotStr))
	} else {
		sb.WriteString(fieldStr)
	}

	// Situation
	if play.Down != "" {
		sb.WriteString("\n  " + headerStyle.Render("󰈍 SITUATION") + "\n")
		sb.WriteString("  " + situationStyle.Render(play.Down) + "\n")
	}

	// Play description
	sb.WriteString("\n  " + headerStyle.Render(" PLAY") + "\n")
	sb.WriteString("  " + borderStyle.Render(strings.Repeat("─", 58)) + "\n")

	// Wrap play text
	playText := play.Text
	if play.ScoringPlay {
		playText = "🏈 " + playText + " 🎉"
		playTextStyle = playTextStyle.Foreground(lipgloss.Color("226")).Bold(true)
	}
	sb.WriteString("  " + playTextStyle.Render(playText) + "\n")

	// Controls
	sb.WriteString("\n" + border + "\n")
	autoStatus := statusStyle.Render("OFF")
	if m.autoPlay {
		autoStatus = lipgloss.NewStyle().Foreground(lipgloss.Color("40")).Bold(true).Render(fmt.Sprintf("ON (%ds)", m.autoSpeed))
	}
	controls := fmt.Sprintf("  ←/→: prev/next • SPACE: auto [%s] • +/-: speed • HOME/END: jump • q: quit", autoStatus)
	sb.WriteString(statusStyle.Render(controls) + "\n")

	return sb.String()
}

// Commands
func fetchReplayDataCmd(gameID string, svc *service.ScoreService) tea.Cmd {
	return func() tea.Msg {
		replay, err := svc.GetGameReplay(gameID)
		if err != nil {
			return replayErrorMsg(err)
		}
		return replayDataMsg(replay)
	}
}

type replayAutoTickMsg time.Time

func replayAutoTickCmd(seconds int) tea.Cmd {
	return tea.Tick(time.Duration(seconds)*time.Second, func(t time.Time) tea.Msg {
		return replayAutoTickMsg(t)
	})
}

func replayMascotTickCmd() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return replayTickMsg(t)
	})
}
