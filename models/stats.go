package models

// GameStats represents processed game statistics
type GameStats struct {
	Game      Game
	HomeStats TeamStats
	AwayStats TeamStats
}

// TeamStats contains all stats for one team
type TeamStats struct {
	TeamName    string
	TeamAbbr    string
	Totals      map[string]string // e.g., "totalYards" -> "350"
	PlayerStats []PlayerStatCategory
}

// PlayerStatCategory groups player stats by category (passing, rushing, etc.)
type PlayerStatCategory struct {
	Category string
	Labels   []string
	Players  []PlayerStatLine
}

// PlayerStatLine represents one player's stats
type PlayerStatLine struct {
	Name     string
	Position string
	Stats    []string
}
