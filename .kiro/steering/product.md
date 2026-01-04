# NFL Scores CLI

A terminal-based application for viewing live NFL game scores, play-by-play updates, game replays, and detailed statistics.

## Core Features

- Display current NFL scoreboard with all games
- Live game tracking with real-time play-by-play updates
- Visual football field showing ball position
- Game replay mode for completed games with manual/auto-play controls
- Detailed game statistics and box scores
- Animated mascot mode with team colors, fireworks, and victory celebrations
- Historical game access via date ranges
- Interactive TUI with keyboard/mouse navigation
- Plain text mode for basic terminals

## CLI Flags

- `--watch` - Live game tracking with play-by-play
- `--replay` - Replay completed games
- `--stats` - View game statistics/box score
- `--mascot` - Show animated mascot with team colors
- `--dates YYYYMMDD-YYYYMMDD` - Access historical games
- `--game ID` - Specify game directly
- `--plain` - Disable colors and icons

## Data Source

All game data is fetched from the ESPN public API (`site.api.espn.com`).
