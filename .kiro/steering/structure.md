# Project Structure

```
nfl-scores/
├── main.go              # Entry point, CLI flag parsing
├── client/
│   └── espn.go          # ESPN API client (HTTP requests)
├── service/
│   └── scores.go        # Business logic layer
├── models/
│   ├── game.go          # Domain models (Game, Team, GameStatus)
│   ├── play.go          # Play, GameSummary, GameReplay models
│   ├── stats.go         # GameStats, TeamStats, PlayerStatLine models
│   ├── response.go      # ESPN scoreboard API response mapping
│   └── summary_response.go  # ESPN summary API response mapping (includes boxscore)
├── formatter/
│   ├── terminal.go      # Scoreboard output formatting
│   ├── live.go          # Live game formatting
│   └── stats.go         # Statistics/box score formatting
└── ui/
    ├── live.go          # Bubble Tea TUI model for live games
    ├── replay.go        # Bubble Tea TUI model for game replay
    ├── field.go         # ASCII football field renderer
    ├── mascot.go        # Animated mascot with team colors
    └── victory.go       # Victory celebration screen
```

## Architecture Layers

1. **client** - External API communication only
2. **models** - Data structures and API response mapping (`ToXxx` converters)
3. **service** - Orchestrates client calls, filters data
4. **formatter** - Terminal output rendering (plain + styled modes)
5. **ui** - Interactive TUI using Bubble Tea

## Conventions

- API response structs live in `models/` with `ToXxx()` methods for conversion
- Each package has a single responsibility
- Plain/styled rendering handled via `plain bool` flag throughout
- Bubble Tea pattern: `Model`, `Init()`, `Update()`, `View()`
