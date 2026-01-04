# Implementation Plan: NFL Scores CLI

## Overview

Streamlined implementation plan for a Go CLI that fetches and displays NFL scores from the ESPN API. Tasks are consolidated for rapid MVP delivery.

## Tasks

- [x] 1. Set up project and create all data models

  - Create directory structure: `client/`, `models/`, `formatter/`, `service/`
  - Update `go.mod` with module name `nfl-scores`
  - Create `models/game.go` with Game, Team structs and GameStatus enum
  - Create `models/response.go` with ESPN API response types and `ToGames()` converter
  - _Requirements: 1.1, 1.2, 2.1, 2.2_

- [x] 2. Implement ESPN client and score service

  - Create `client/espn.go` with ESPNClient struct
  - Implement `FetchScoreboard()` to GET `https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard`
  - Create `service/scores.go` with ScoreService that wraps the client
  - Implement `GetCurrentScores()` returning `[]Game`
  - _Requirements: 1.1, 1.2, 1.3, 4.1, 4.2_

- [x] 3. Implement formatter and CLI entry point

  - Create `formatter/terminal.go` with TerminalFormatter
  - Implement `FormatScoreboard()` displaying games in aligned format (≤80 chars)
  - Implement `FormatError()` for user-friendly error messages
  - Create `main.go` with flag parsing (-h/--help), component wiring, stdout/stderr handling
  - Exit code 0 on success, 1 on error
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 3.1, 3.2, 3.3, 3.4, 4.3, 4.4_

- [x] 4. Final verification

  - Run `go build` to verify compilation
  - Test CLI with `./nfl-scores` and `./nfl-scores -h`
  - Ensure all tests pass, ask the user if questions arise.

- [ ]\* 5. Add property-based tests (optional)
  - [ ]\* 5.1 Property test: JSON parsing round-trip
    - **Property 1: JSON Parsing Preserves Game Data**
    - **Validates: Requirements 1.2**
  - [ ]\* 5.2 Property test: Formatted output contains all game info
    - **Property 2: Formatted Output Contains All Game Information**
    - **Validates: Requirements 2.1, 2.2, 2.5**
  - [ ]\* 5.3 Property test: Output width ≤80 characters
    - **Property 3: Output Width Constraint**
    - **Validates: Requirements 2.3, 2.4**

## Notes

- Tasks 1-4 deliver a working MVP
- Task 5 (marked `*`) is optional for comprehensive testing later
- ESPN API is free, no authentication required
