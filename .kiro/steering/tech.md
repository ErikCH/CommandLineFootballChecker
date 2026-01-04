# Tech Stack

## Language & Runtime

- Go 1.25+
- Module name: `nfl-scores`

## Key Dependencies

- `github.com/charmbracelet/bubbletea` - Terminal UI framework (Elm architecture)
- `github.com/charmbracelet/bubbles` - UI components (spinner, etc.)
- `github.com/charmbracelet/lipgloss` - Terminal styling and layout

## External APIs

- ESPN Scoreboard: `https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard`
- ESPN Game Summary: `https://site.api.espn.com/apis/site/v2/sports/football/nfl/summary`

## Common Commands

```bash
# Build
go build -o nfl-scores .

# Run
go run .

# Run with flags
go run . --plain         # No colors/icons
go run . --watch         # Live game tracking
go run . --watch --mascot # Live with animated mascot
go run . --replay        # Replay completed games
go run . --stats         # View game statistics
go run . --game ID       # Watch specific game
go run . --dates YYYYMMDD-YYYYMMDD  # Historical games

# Dependencies
go mod download
go mod tidy
```

## Code Style

- Use `fmt.Errorf` with `%w` for error wrapping
- Constructor functions follow `NewXxx` pattern
- HTTP clients should have configurable timeouts
- Prefer lipgloss styles over raw ANSI codes
- Pad strings before applying lipgloss styles for proper column alignment
