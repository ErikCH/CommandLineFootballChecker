package models

import (
	"strconv"
	"time"
)

// ScoreboardResponse represents the ESPN API scoreboard response
type ScoreboardResponse struct {
	Events []Event `json:"events"`
}

// Event represents a single game event
type Event struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Date         string        `json:"date"`
	Status       EventStatus   `json:"status"`
	Competitions []Competition `json:"competitions"`
}

// EventStatus contains game status information
type EventStatus struct {
	Type StatusType `json:"type"`
}

// StatusType contains detailed status info
type StatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	ShortDetail string `json:"shortDetail"`
}

// Competition represents a single competition within an event
type Competition struct {
	Competitors []Competitor `json:"competitors"`
}

// Competitor represents a team in the competition
type Competitor struct {
	HomeAway string   `json:"homeAway"`
	Team     TeamInfo `json:"team"`
	Score    string   `json:"score"`
}

// TeamInfo contains team details
type TeamInfo struct {
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
}

// ToGames converts the ESPN response to a slice of Game structs
func (r *ScoreboardResponse) ToGames() []Game {
	games := make([]Game, 0, len(r.Events))

	for _, event := range r.Events {
		if len(event.Competitions) == 0 {
			continue
		}

		game := Game{
			ID:         event.ID,
			StatusText: event.Status.Type.ShortDetail,
			Status:     mapStatus(event.Status.Type.State),
		}

		// Parse start time
		if t, err := time.Parse(time.RFC3339, event.Date); err == nil {
			game.StartTime = t
		}

		// Extract teams
		for _, comp := range event.Competitions[0].Competitors {
			score, _ := strconv.Atoi(comp.Score)
			team := Team{
				Name:         comp.Team.DisplayName,
				Abbreviation: comp.Team.Abbreviation,
				Score:        score,
			}

			if comp.HomeAway == "home" {
				game.HomeTeam = team
			} else {
				game.AwayTeam = team
			}
		}

		games = append(games, game)
	}

	return games
}

// mapStatus converts ESPN status state to GameStatus
func mapStatus(state string) GameStatus {
	switch state {
	case "in":
		return StatusInProgress
	case "post":
		return StatusFinal
	default:
		return StatusScheduled
	}
}
