# Design Document: NFL Scores CLI

## Overview

This document describes the technical design for the NFL Scores CLI application—a Go-based command-line tool that fetches and displays NFL game scores. The application uses the ESPN public API to retrieve real-time game data and presents it in a professionally formatted terminal output.

The design prioritizes simplicity, reliability, and clean code organization following Go best practices.

## Architecture

The application follows a layered architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────┐
│                      CLI Layer                          │
│              (main.go - argument parsing)               │
├─────────────────────────────────────────────────────────┤
│                   Service Layer                         │
│           (score fetching & orchestration)              │
├─────────────────────────────────────────────────────────┤
│                   Client Layer                          │
│              (ESPN API HTTP client)                     │
├─────────────────────────────────────────────────────────┤
│                  Formatter Layer                        │
│            (terminal output formatting)                 │
└─────────────────────────────────────────────────────────┘
```

### Data Flow

1. User invokes CLI → CLI parses arguments
2. CLI calls Score Service → Service calls ESPN API Client
3. API Client fetches data → Returns parsed Game structs
4. Service passes Games to Formatter → Formatter renders output
5. CLI displays output → Returns exit code

## Components and Interfaces

### 1. API Client (`client/espn.go`)

Responsible for HTTP communication with the ESPN API.

```go
package client

type ESPNClient struct {
    httpClient *http.Client
    baseURL    string
}

// NewESPNClient creates a new ESPN API client with default configuration
func NewESPNClient() *ESPNClient

// FetchScoreboard retrieves the current NFL scoreboard data
func (c *ESPNClient) FetchScoreboard() (*ScoreboardResponse, error)
```

**ESPN API Endpoint:**

- Base URL: `https://site.api.espn.com`
- Scoreboard: `/apis/site/v2/sports/football/nfl/scoreboard`

### 2. Models (`models/game.go`)

Data structures representing NFL game information.

```go
package models

type Game struct {
    ID         string
    HomeTeam   Team
    AwayTeam   Team
    HomeScore  int
    AwayScore  int
    Status     GameStatus
    StatusText string
    StartTime  time.Time
}

type Team struct {
    Name         string
    Abbreviation string
    Score        int
}

type GameStatus int

const (
    StatusScheduled GameStatus = iota
    StatusInProgress
    StatusFinal
)
```

### 3. Score Service (`service/scores.go`)

Orchestrates fetching and processing of score data.

```go
package service

type ScoreService struct {
    client *client.ESPNClient
}

// NewScoreService creates a new score service instance
func NewScoreService(c *client.ESPNClient) *ScoreService

// GetCurrentScores retrieves and processes current NFL scores
func (s *ScoreService) GetCurrentScores() ([]models.Game, error)
```

### 4. Formatter (`formatter/terminal.go`)

Handles terminal output formatting.

```go
package formatter

type TerminalFormatter struct {
    width int
}

// NewTerminalFormatter creates a formatter with specified width
func NewTerminalFormatter(width int) *TerminalFormatter

// FormatScoreboard renders games as formatted terminal output
func (f *TerminalFormatter) FormatScoreboard(games []models.Game) string

// FormatError renders error messages for terminal display
func (f *TerminalFormatter) FormatError(err error) string
```

### 5. CLI (`main.go`)

Entry point handling argument parsing and orchestration.

```go
package main

func main() {
    // Parse flags
    // Initialize components
    // Fetch and display scores
    // Handle errors and exit codes
}
```

## Data Models

### ESPN API Response Structure

The ESPN scoreboard endpoint returns JSON with the following relevant structure:

```json
{
  "events": [
    {
      "id": "401547417",
      "name": "Team A at Team B",
      "status": {
        "type": {
          "id": "3",
          "name": "STATUS_FINAL",
          "state": "post",
          "completed": true
        }
      },
      "competitions": [
        {
          "competitors": [
            {
              "homeAway": "home",
              "team": {
                "displayName": "Team B",
                "abbreviation": "TB"
              },
              "score": "24"
            },
            {
              "homeAway": "away",
              "team": {
                "displayName": "Team A",
                "abbreviation": "TA"
              },
              "score": "17"
            }
          ]
        }
      ]
    }
  ]
}
```

### Internal Game Model Mapping

| ESPN Field                        | Internal Field      | Transformation         |
| --------------------------------- | ------------------- | ---------------------- |
| `events[].id`                     | `Game.ID`           | Direct copy            |
| `competitors[].team.displayName`  | `Team.Name`         | Direct copy            |
| `competitors[].team.abbreviation` | `Team.Abbreviation` | Direct copy            |
| `competitors[].score`             | `Team.Score`        | Parse string to int    |
| `status.type.state`               | `Game.Status`       | Map to GameStatus enum |
| `status.type.shortDetail`         | `Game.StatusText`   | Direct copy            |

## Correctness Properties

_A property is a characteristic or behavior that should hold true across all valid executions of a system—essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees._

### Property 1: JSON Parsing Preserves Game Data

_For any_ valid Game structure, when serialized to ESPN-like JSON format and then parsed back, the resulting Game structure SHALL be equivalent to the original.

**Validates: Requirements 1.2**

### Property 2: Formatted Output Contains All Game Information

_For any_ valid Game structure, the formatted output string SHALL contain:

- The home team name
- The away team name
- The home team score
- The away team score
- A status indicator (Scheduled, In Progress, or Final)

**Validates: Requirements 2.1, 2.2, 2.5**

### Property 3: Output Width Constraint

_For any_ list of Game structures (including empty lists), every line in the formatted output SHALL be at most 80 characters in length.

**Validates: Requirements 2.3, 2.4**

### Property 4: Malformed JSON Produces Parse Error

_For any_ malformed or invalid JSON input, the parser SHALL return an error rather than a partial or incorrect Game structure.

**Validates: Requirements 4.2**

### Property 5: Exit Code and Output Stream Consistency

_For any_ execution of the CLI:

- Successful execution SHALL return exit code 0 and write output to stdout
- Failed execution SHALL return exit code 1 and write errors to stderr

**Validates: Requirements 3.3, 4.3, 4.4**

## Error Handling

### Error Categories

| Error Type         | Source      | User Message                                                                                  | Exit Code |
| ------------------ | ----------- | --------------------------------------------------------------------------------------------- | --------- |
| Network Timeout    | HTTP Client | "Unable to connect to NFL data service. Please check your internet connection and try again." | 1         |
| Connection Refused | HTTP Client | "NFL data service is unavailable. Please try again later."                                    | 1         |
| Invalid JSON       | Parser      | "Received invalid data from NFL service. Please try again."                                   | 1         |
| No Games           | Service     | "No NFL games are currently scheduled."                                                       | 0         |
| Unknown            | Any         | "An unexpected error occurred. Please try again."                                             | 1         |

### Error Handling Strategy

1. **Wrap errors with context**: Use `fmt.Errorf` with `%w` to preserve error chain
2. **Categorize at boundaries**: Convert low-level errors to user-friendly messages at CLI layer
3. **Log vs Display**: Internal details logged (if verbose), user messages displayed
4. **Graceful degradation**: Always exit cleanly, never panic

```go
// Error wrapping example
func (c *ESPNClient) FetchScoreboard() (*ScoreboardResponse, error) {
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch scoreboard: %w", err)
    }
    // ...
}
```

## Testing Strategy

### Unit Tests

Unit tests verify specific examples and edge cases:

- **Client tests**: Mock HTTP responses, verify correct URL construction
- **Parser tests**: Test with sample ESPN JSON responses
- **Formatter tests**: Verify output format with known game data
- **Error handling tests**: Verify correct error messages for each error type

### Property-Based Tests

Property-based tests verify universal properties using the `rapid` library for Go:

- **Minimum 100 iterations** per property test
- **Tag format**: `Feature: nfl-scores-cli, Property N: [description]`
- Each correctness property maps to one property-based test

### Test Organization

```
├── client/
│   ├── espn.go
│   └── espn_test.go
├── models/
│   ├── game.go
│   └── game_test.go
├── formatter/
│   ├── terminal.go
│   └── terminal_test.go
├── service/
│   ├── scores.go
│   └── scores_test.go
└── main_test.go
```

### Testing Library

- **Property-based testing**: `pgregory.net/rapid` - Go's most popular PBT library
- **Standard testing**: Go's built-in `testing` package
- **Assertions**: `github.com/stretchr/testify/assert` for cleaner assertions

## Project Structure

```
nfl-scores/
├── main.go              # CLI entry point
├── client/
│   ├── espn.go          # ESPN API client
│   └── espn_test.go     # Client tests
├── models/
│   ├── game.go          # Game data structures
│   ├── response.go      # ESPN API response types
│   └── game_test.go     # Model tests
├── formatter/
│   ├── terminal.go      # Terminal output formatter
│   └── terminal_test.go # Formatter tests
├── service/
│   ├── scores.go        # Score service
│   └── scores_test.go   # Service tests
├── go.mod
└── go.sum
```
