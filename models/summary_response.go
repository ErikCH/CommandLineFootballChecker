package models

import (
	"strconv"
	"time"
)

// SummaryResponse represents the ESPN game summary API response
type SummaryResponse struct {
	Header   SummaryHeader    `json:"header"`
	Drives   Drives           `json:"drives"`
	Boxscore BoxscoreResponse `json:"boxscore"`
}

type SummaryHeader struct {
	ID           string               `json:"id"`
	Competitions []SummaryCompetition `json:"competitions"`
}

type SummaryCompetition struct {
	ID          string              `json:"id"`
	Date        string              `json:"date"`
	Status      SummaryStatus       `json:"status"`
	Competitors []SummaryCompetitor `json:"competitors"`
}

type SummaryStatus struct {
	Type SummaryStatusType `json:"type"`
}

type SummaryStatusType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Completed   bool   `json:"completed"`
	ShortDetail string `json:"shortDetail"`
}

type SummaryCompetitor struct {
	ID       string          `json:"id"`
	HomeAway string          `json:"homeAway"`
	Score    string          `json:"score"`
	Team     SummaryTeamInfo `json:"team"`
}

type SummaryTeamInfo struct {
	ID           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
}

type Drives struct {
	Current  *CurrentDrive `json:"current"`
	Previous []DriveInfo   `json:"previous"`
}

type CurrentDrive struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Team        DriveTeam  `json:"team"`
	Plays       []PlayInfo `json:"plays"`
	Start       DriveStart `json:"start"`
}

type DriveTeam struct {
	Abbreviation string `json:"abbreviation"`
	DisplayName  string `json:"displayName"`
}

type DriveStart struct {
	Text string `json:"text"`
}

type DriveInfo struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Team        DriveTeam  `json:"team"`
	Plays       []PlayInfo `json:"plays"`
}

type PlayInfo struct {
	ID          string       `json:"id"`
	Text        string       `json:"text"`
	Type        PlayType     `json:"type"`
	Clock       PlayClock    `json:"clock"`
	Period      PlayPeriod   `json:"period"`
	HomeScore   int          `json:"homeScore"`
	AwayScore   int          `json:"awayScore"`
	ScoringPlay bool         `json:"scoringPlay"`
	Start       PlayPosition `json:"start"`
	End         PlayPosition `json:"end"`
}

type PlayType struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type PlayClock struct {
	DisplayValue string `json:"displayValue"`
}

type PlayPeriod struct {
	Number int `json:"number"`
}

type PlayPosition struct {
	Down             int    `json:"down"`
	Distance         int    `json:"distance"`
	YardLine         int    `json:"yardLine"`
	YardsToEndzone   int    `json:"yardsToEndzone"`
	DownDistanceText string `json:"downDistanceText"`
	PossessionText   string `json:"possessionText"`
}

// BoxscoreResponse represents the boxscore section of the ESPN API
type BoxscoreResponse struct {
	Teams   []BoxscoreTeam   `json:"teams"`
	Players []BoxscorePlayer `json:"players"`
}

type BoxscoreTeam struct {
	Team       BoxscoreTeamInfo   `json:"team"`
	Statistics []BoxscoreTeamStat `json:"statistics"`
}

type BoxscoreTeamInfo struct {
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
	DisplayName  string `json:"displayName"`
}

type BoxscoreTeamStat struct {
	Name         string `json:"name"`
	DisplayValue string `json:"displayValue"`
}

type BoxscorePlayer struct {
	Team       BoxscoreTeamInfo      `json:"team"`
	Statistics []BoxscorePlayerGroup `json:"statistics"`
}

type BoxscorePlayerGroup struct {
	Name     string                `json:"name"`
	Labels   []string              `json:"labels"`
	Athletes []BoxscoreAthleteData `json:"athletes"`
}

type BoxscoreAthleteData struct {
	Athlete BoxscoreAthlete `json:"athlete"`
	Stats   []string        `json:"stats"`
}

type BoxscoreAthlete struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Position    string `json:"position"`
}

// ToGameSummary converts the API response to our internal model
func (r *SummaryResponse) ToGameSummary() *GameSummary {
	if len(r.Header.Competitions) == 0 {
		return nil
	}

	comp := r.Header.Competitions[0]

	game := Game{
		ID:         r.Header.ID,
		StatusText: comp.Status.Type.ShortDetail,
		Status:     mapStatus(comp.Status.Type.State),
	}

	// Parse start time
	if t, err := time.Parse(time.RFC3339, comp.Date); err == nil {
		game.StartTime = t
	}

	// Extract teams
	for _, c := range comp.Competitors {
		score, _ := strconv.Atoi(c.Score)
		team := Team{
			Name:         c.Team.DisplayName,
			Abbreviation: c.Team.Abbreviation,
			Score:        score,
		}
		if c.HomeAway == "home" {
			game.HomeTeam = team
		} else {
			game.AwayTeam = team
		}
	}

	summary := &GameSummary{
		Game:        game,
		RecentPlays: make([]Play, 0),
	}

	// Get current situation and plays
	if r.Drives.Current != nil && len(r.Drives.Current.Plays) > 0 {
		lastPlay := r.Drives.Current.Plays[len(r.Drives.Current.Plays)-1]

		summary.CurrentPlay = &Play{
			ID:             lastPlay.ID,
			Text:           lastPlay.Text,
			Type:           lastPlay.Type.Text,
			Clock:          lastPlay.Clock.DisplayValue,
			Period:         lastPlay.Period.Number,
			HomeScore:      lastPlay.HomeScore,
			AwayScore:      lastPlay.AwayScore,
			ScoringPlay:    lastPlay.ScoringPlay,
			Down:           lastPlay.End.DownDistanceText,
			Possession:     r.Drives.Current.Team.Abbreviation,
			YardsToEndzone: lastPlay.End.YardsToEndzone,
		}

		summary.Situation = lastPlay.End.DownDistanceText
		summary.YardsToEndzone = lastPlay.End.YardsToEndzone

		// Get recent plays from current drive
		for i := len(r.Drives.Current.Plays) - 1; i >= 0 && len(summary.RecentPlays) < 5; i-- {
			p := r.Drives.Current.Plays[i]
			summary.RecentPlays = append(summary.RecentPlays, Play{
				ID:          p.ID,
				Text:        p.Text,
				Type:        p.Type.Text,
				Clock:       p.Clock.DisplayValue,
				Period:      p.Period.Number,
				HomeScore:   p.HomeScore,
				AwayScore:   p.AwayScore,
				ScoringPlay: p.ScoringPlay,
			})
		}
	}

	return summary
}

// ToGameReplay converts the API response to a full game replay with all plays
func (r *SummaryResponse) ToGameReplay() *GameReplay {
	if len(r.Header.Competitions) == 0 {
		return nil
	}

	comp := r.Header.Competitions[0]

	game := Game{
		ID:         r.Header.ID,
		StatusText: comp.Status.Type.ShortDetail,
		Status:     mapStatus(comp.Status.Type.State),
	}

	// Extract teams
	for _, c := range comp.Competitors {
		score, _ := strconv.Atoi(c.Score)
		team := Team{
			Name:         c.Team.DisplayName,
			Abbreviation: c.Team.Abbreviation,
			Score:        score,
		}
		if c.HomeAway == "home" {
			game.HomeTeam = team
		} else {
			game.AwayTeam = team
		}
	}

	replay := &GameReplay{
		Game:   game,
		Plays:  make([]ReplayPlay, 0),
		Drives: make([]ReplayDrive, 0),
	}

	// Collect all plays from all drives in order
	for _, drive := range r.Drives.Previous {
		rd := ReplayDrive{
			ID:          drive.ID,
			Description: drive.Description,
			Team:        drive.Team.Abbreviation,
			StartIndex:  len(replay.Plays),
		}

		for _, p := range drive.Plays {
			play := ReplayPlay{
				ID:             p.ID,
				Text:           p.Text,
				Type:           p.Type.Text,
				Clock:          p.Clock.DisplayValue,
				Period:         p.Period.Number,
				HomeScore:      p.HomeScore,
				AwayScore:      p.AwayScore,
				ScoringPlay:    p.ScoringPlay,
				Possession:     drive.Team.Abbreviation,
				YardsToEndzone: p.End.YardsToEndzone,
				Down:           p.End.DownDistanceText,
				DriveID:        drive.ID,
			}
			replay.Plays = append(replay.Plays, play)
		}

		rd.EndIndex = len(replay.Plays) - 1
		replay.Drives = append(replay.Drives, rd)
	}

	return replay
}

// ToGameStats converts the API response to game statistics
func (r *SummaryResponse) ToGameStats() *GameStats {
	if len(r.Header.Competitions) == 0 {
		return nil
	}

	comp := r.Header.Competitions[0]

	game := Game{
		ID:         r.Header.ID,
		StatusText: comp.Status.Type.ShortDetail,
		Status:     mapStatus(comp.Status.Type.State),
	}

	// Extract teams
	for _, c := range comp.Competitors {
		score, _ := strconv.Atoi(c.Score)
		team := Team{
			Name:         c.Team.DisplayName,
			Abbreviation: c.Team.Abbreviation,
			Score:        score,
		}
		if c.HomeAway == "home" {
			game.HomeTeam = team
		} else {
			game.AwayTeam = team
		}
	}

	stats := &GameStats{
		Game: game,
		HomeStats: TeamStats{
			TeamName: game.HomeTeam.Name,
			TeamAbbr: game.HomeTeam.Abbreviation,
			Totals:   make(map[string]string),
		},
		AwayStats: TeamStats{
			TeamName: game.AwayTeam.Name,
			TeamAbbr: game.AwayTeam.Abbreviation,
			Totals:   make(map[string]string),
		},
	}

	// Parse team totals
	for _, teamBox := range r.Boxscore.Teams {
		var teamStats *TeamStats
		if teamBox.Team.Abbreviation == game.HomeTeam.Abbreviation {
			teamStats = &stats.HomeStats
		} else {
			teamStats = &stats.AwayStats
		}

		for _, stat := range teamBox.Statistics {
			teamStats.Totals[stat.Name] = stat.DisplayValue
		}
	}

	// Parse player stats
	for _, playerBox := range r.Boxscore.Players {
		var teamStats *TeamStats
		if playerBox.Team.Abbreviation == game.HomeTeam.Abbreviation {
			teamStats = &stats.HomeStats
		} else {
			teamStats = &stats.AwayStats
		}

		for _, statGroup := range playerBox.Statistics {
			cat := PlayerStatCategory{
				Category: statGroup.Name,
				Labels:   statGroup.Labels,
			}
			for _, athlete := range statGroup.Athletes {
				cat.Players = append(cat.Players, PlayerStatLine{
					Name:     athlete.Athlete.DisplayName,
					Position: athlete.Athlete.Position,
					Stats:    athlete.Stats,
				})
			}
			teamStats.PlayerStats = append(teamStats.PlayerStats, cat)
		}
	}

	return stats
}
