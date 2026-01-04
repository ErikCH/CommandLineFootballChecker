# 🏈 NFL Scores CLI

A terminal-based application for viewing live NFL game scores, play-by-play updates, game replays, and detailed statistics.

Built with [Kiro](https://kiro.dev)

## Features

- **Live Scoreboard** - View all current NFL game scores
- **Live Game Tracking** - Watch games in real-time with play-by-play updates and visual field position
- **Game Replay** - Replay completed games play-by-play with manual or auto-play controls
- **Game Statistics** - View detailed box scores with team and player stats
- **Animated Mascot Mode** - Fun dancing mascot with team colors, fireworks on scores, and victory celebrations
- **Historical Games** - Access games from any date range
- **Plain Text Mode** - Works on basic terminals without color support

## Installation

```bash
# Clone the repository
git clone https://github.com/erikch/nfl-scores.git
cd nfl-scores

# Build
go build -o nfl-scores .

# Run
./nfl-scores
```

## Usage

```bash
# Display current NFL scores
./nfl-scores

# Watch a live game with play-by-play
./nfl-scores --watch

# Watch with animated mascot
./nfl-scores --watch --mascot

# Replay a completed game
./nfl-scores --replay

# View game statistics/box score
./nfl-scores --stats

# Access historical games
./nfl-scores --dates 20241201-20241208
./nfl-scores --replay --dates 20241201-20241208
./nfl-scores --stats --dates 20241201-20241208

# Plain text mode (no colors/icons)
./nfl-scores --plain

# Specify a game directly
./nfl-scores --watch --game 401671793
```

## Replay Controls

| Key            | Action               |
| -------------- | -------------------- |
| `←` / `→`      | Previous / Next play |
| `Space`        | Toggle auto-play     |
| `+` / `-`      | Speed up / Slow down |
| `Home` / `End` | Jump to start / end  |
| `q`            | Quit                 |

## Statistics

The `--stats` flag shows detailed box scores including:

- **Team Stats**: Total yards, passing/rushing yards, first downs, 3rd down efficiency, turnovers, time of possession
- **Passing**: Completions, yards, TDs, INTs, sacks, QBR, passer rating
- **Rushing**: Carries, yards, average, TDs, long
- **Receiving**: Receptions, yards, average, TDs, long

## Data Source

All game data is fetched from the ESPN public API.

## Tech Stack

- Go 1.25+
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

MIT
# CommandLineFootballChecker
