package models

import "time"

// GameStatus represents the current state of a game
type GameStatus int

const (
	StatusScheduled GameStatus = iota
	StatusInProgress
	StatusFinal
)

// String returns a human-readable status string
func (s GameStatus) String() string {
	switch s {
	case StatusScheduled:
		return "Scheduled"
	case StatusInProgress:
		return "In Progress"
	case StatusFinal:
		return "Final"
	default:
		return "Unknown"
	}
}

// Team represents an NFL team with score
type Team struct {
	Name         string
	Abbreviation string
	Score        int
}

// Game represents a single NFL game
type Game struct {
	ID         string
	HomeTeam   Team
	AwayTeam   Team
	Status     GameStatus
	StatusText string
	StartTime  time.Time
}
