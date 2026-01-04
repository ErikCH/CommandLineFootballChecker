package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"nfl-scores/client"
	"nfl-scores/formatter"
	"nfl-scores/models"
	"nfl-scores/service"
	"nfl-scores/ui"

	tea "github.com/charmbracelet/bubbletea"
)

const helpText = `NFL Scores CLI - Display current NFL game scores

Usage:
  nfl-scores [options]

Options:
  -h, --help      Show this help message
  --plain         Disable colors and icons (for basic terminals)
  --watch         Watch a live game with play-by-play updates
  --replay        Replay a completed game play-by-play
  --stats         Show detailed game statistics (box score)
  --game ID       Specify game ID directly
  --dates RANGE   Date range for historical games (format: YYYYMMDD-YYYYMMDD)
  --mascot        Show animated mascot with team colors

Examples:
  nfl-scores                          Display current NFL scores
  nfl-scores --dates 20241201-20241208  Show games from Dec 1-8, 2024
  nfl-scores --stats                  View stats for a game
  nfl-scores --stats --game ID        Stats for specific game
  nfl-scores --replay                 Select and replay a completed game
  nfl-scores --watch --mascot         Watch live game with mascot
  nfl-scores -h                       Show help
`

func main() {
	// Parse flags
	help := flag.Bool("h", false, "Show help")
	flag.BoolVar(help, "help", false, "Show help")
	plain := flag.Bool("plain", false, "Disable colors and icons")
	watch := flag.Bool("watch", false, "Watch a live game")
	replay := flag.Bool("replay", false, "Replay a completed game")
	showStats := flag.Bool("stats", false, "Show game statistics")
	gameID := flag.String("game", "", "Game ID to watch or replay")
	dates := flag.String("dates", "", "Date range (YYYYMMDD-YYYYMMDD)")
	mascot := flag.Bool("mascot", false, "Show animated mascot")
	flag.Parse()

	if *help {
		fmt.Print(helpText)
		os.Exit(0)
	}

	// Initialize components
	espnClient := client.NewESPNClient()
	scoreService := service.NewScoreService(espnClient)
	termFormatter := formatter.NewTerminalFormatter(80, *plain)

	if *showStats {
		runStatsMode(scoreService, termFormatter, *gameID, *dates)
		return
	}

	if *replay {
		runReplayMode(scoreService, termFormatter, *gameID, *dates, *plain, *mascot)
		return
	}

	if *watch {
		runWatchMode(scoreService, termFormatter, *gameID, *plain, *mascot)
		return
	}

	// Default: show all scores
	games, err := scoreService.GetScoresByDates(*dates)
	if err != nil {
		fmt.Fprintln(os.Stderr, termFormatter.FormatError(err))
		os.Exit(1)
	}

	fmt.Print(termFormatter.FormatScoreboard(games))
	os.Exit(0)
}

func runWatchMode(svc *service.ScoreService, f *formatter.TerminalFormatter, gameID string, plain bool, mascot bool) {
	// If no game ID provided, let user select from live games
	if gameID == "" {
		games, err := svc.GetLiveGames()
		if err != nil {
			fmt.Fprintln(os.Stderr, f.FormatError(err))
			os.Exit(1)
		}

		if len(games) == 0 {
			fmt.Println("No live games available right now. Try again during game time.")
			os.Exit(0)
		}

		fmt.Print(f.FormatGameSelection(games))

		// Read user selection
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(games) {
			fmt.Println("Invalid selection.")
			os.Exit(1)
		}

		gameID = games[num-1].ID
	}

	// Run Bubble Tea UI with mouse support
	var model ui.Model
	if mascot {
		model = ui.NewModelWithMascot(gameID, svc, plain)
	} else {
		model = ui.NewModel(gameID, svc, plain)
	}
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
		os.Exit(1)
	}
}

func runReplayMode(svc *service.ScoreService, f *formatter.TerminalFormatter, gameID string, dates string, plain bool, mascot bool) {
	// If no game ID provided, let user select from completed games
	if gameID == "" {
		games, err := svc.GetScoresByDates(dates)
		if err != nil {
			fmt.Fprintln(os.Stderr, f.FormatError(err))
			os.Exit(1)
		}

		// Filter to completed games only
		var completed []models.Game
		for _, g := range games {
			if g.Status == models.StatusFinal {
				completed = append(completed, g)
			}
		}

		if len(completed) == 0 {
			fmt.Println("No completed games available for replay.")
			os.Exit(0)
		}

		dateInfo := ""
		if dates != "" {
			dateInfo = fmt.Sprintf(" (%s)", dates)
		}
		fmt.Printf("\nSelect a completed game to replay%s:\n", dateInfo)
		fmt.Println()
		for i, g := range completed {
			fmt.Printf("  [%d] %s %d @ %s %d (Final)\n",
				i+1, g.AwayTeam.Abbreviation, g.AwayTeam.Score,
				g.HomeTeam.Abbreviation, g.HomeTeam.Score)
		}
		fmt.Print("\nEnter number: ")

		// Read user selection
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(completed) {
			fmt.Println("Invalid selection.")
			os.Exit(1)
		}

		gameID = completed[num-1].ID
	}

	// Run replay UI
	model := ui.NewReplayModel(gameID, svc, plain, mascot)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running replay: %v\n", err)
		os.Exit(1)
	}
}

func runStatsMode(svc *service.ScoreService, f *formatter.TerminalFormatter, gameID string, dates string) {
	// If no game ID provided, let user select a game
	if gameID == "" {
		games, err := svc.GetScoresByDates(dates)
		if err != nil {
			fmt.Fprintln(os.Stderr, f.FormatError(err))
			os.Exit(1)
		}

		if len(games) == 0 {
			fmt.Println("No games available.")
			os.Exit(0)
		}

		dateInfo := ""
		if dates != "" {
			dateInfo = fmt.Sprintf(" (%s)", dates)
		}
		fmt.Printf("\nSelect a game to view stats%s:\n", dateInfo)
		fmt.Println()
		for i, g := range games {
			status := g.StatusText
			if status == "" {
				status = g.Status.String()
			}
			fmt.Printf("  [%d] %s %d @ %s %d (%s)\n",
				i+1, g.AwayTeam.Abbreviation, g.AwayTeam.Score,
				g.HomeTeam.Abbreviation, g.HomeTeam.Score, status)
		}
		fmt.Print("\nEnter number: ")

		// Read user selection
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err != nil || num < 1 || num > len(games) {
			fmt.Println("Invalid selection.")
			os.Exit(1)
		}

		gameID = games[num-1].ID
	}

	// Fetch and display stats
	stats, err := svc.GetGameStats(gameID)
	if err != nil {
		fmt.Fprintln(os.Stderr, f.FormatError(err))
		os.Exit(1)
	}

	fmt.Print(f.FormatGameStats(stats))
}
