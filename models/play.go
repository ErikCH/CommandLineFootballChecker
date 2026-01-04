package models

// Play represents a single play in a game
type Play struct {
	ID             string
	Text           string
	Type           string
	Clock          string
	Period         int
	HomeScore      int
	AwayScore      int
	ScoringPlay    bool
	Down           string
	Possession     string
	YardsToEndzone int
}

// GameSummary contains detailed game info with plays
type GameSummary struct {
	Game           Game
	CurrentPlay    *Play
	RecentPlays    []Play
	Situation      string // e.g., "1st & 10 at CAR 25"
	YardsToEndzone int
}

// ReplayPlay represents a play with full context for replay
type ReplayPlay struct {
	ID             string
	Text           string
	Type           string
	Clock          string
	Period         int
	HomeScore      int
	AwayScore      int
	ScoringPlay    bool
	Possession     string
	YardsToEndzone int
	Down           string
	DriveID        string
}

// ReplayDrive represents a drive in the replay
type ReplayDrive struct {
	ID          string
	Description string
	Team        string
	StartIndex  int // Index of first play in Plays slice
	EndIndex    int // Index of last play in Plays slice
}

// GameReplay contains full game data for replay mode
type GameReplay struct {
	Game   Game
	Plays  []ReplayPlay
	Drives []ReplayDrive
}
