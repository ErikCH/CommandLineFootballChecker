package service

import (
	"nfl-scores/client"
	"nfl-scores/models"
)

// ScoreService orchestrates fetching and processing of score data
type ScoreService struct {
	client *client.ESPNClient
}

// NewScoreService creates a new score service instance
func NewScoreService(c *client.ESPNClient) *ScoreService {
	return &ScoreService{
		client: c,
	}
}

// GetCurrentScores retrieves and processes current NFL scores
func (s *ScoreService) GetCurrentScores() ([]models.Game, error) {
	return s.GetScoresByDates("")
}

// GetScoresByDates retrieves NFL scores for a date range (format: YYYYMMDD-YYYYMMDD)
func (s *ScoreService) GetScoresByDates(dates string) ([]models.Game, error) {
	response, err := s.client.FetchScoreboardByDates(dates)
	if err != nil {
		return nil, err
	}

	return response.ToGames(), nil
}

// GetLiveGames returns only games that are currently in progress
func (s *ScoreService) GetLiveGames() ([]models.Game, error) {
	games, err := s.GetCurrentScores()
	if err != nil {
		return nil, err
	}

	live := make([]models.Game, 0)
	for _, g := range games {
		if g.Status == models.StatusInProgress {
			live = append(live, g)
		}
	}
	return live, nil
}

// GetGameSummary retrieves detailed game info with play-by-play
func (s *ScoreService) GetGameSummary(gameID string) (*models.GameSummary, error) {
	response, err := s.client.FetchGameSummary(gameID)
	if err != nil {
		return nil, err
	}

	return response.ToGameSummary(), nil
}

// GetGameReplay retrieves full game data for replay mode
func (s *ScoreService) GetGameReplay(gameID string) (*models.GameReplay, error) {
	return s.client.FetchGameReplay(gameID)
}

// GetGameStats retrieves game statistics
func (s *ScoreService) GetGameStats(gameID string) (*models.GameStats, error) {
	return s.client.FetchGameStats(gameID)
}
