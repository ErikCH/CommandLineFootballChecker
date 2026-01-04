package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"nfl-scores/models"
)

const (
	baseURL         = "https://site.api.espn.com"
	scoreboardPath  = "/apis/site/v2/sports/football/nfl/scoreboard"
	summaryPath     = "/apis/site/v2/sports/football/nfl/summary"
	defaultTimeout  = 10 * time.Second
	maxResponseSize = 10 * 1024 * 1024 // 10MB max response size
)

// Validation patterns
var (
	dateRangePattern  = regexp.MustCompile(`^\d{8}-\d{8}$`)
	singleDatePattern = regexp.MustCompile(`^\d{8}$`)
	gameIDPattern     = regexp.MustCompile(`^\d+$`)
)

// ESPNClient handles communication with the ESPN API
type ESPNClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewESPNClient creates a new ESPN API client with default configuration
func NewESPNClient() *ESPNClient {
	return &ESPNClient{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL: baseURL,
	}
}

// FetchScoreboard retrieves the current NFL scoreboard data
func (c *ESPNClient) FetchScoreboard() (*models.ScoreboardResponse, error) {
	return c.FetchScoreboardByDates("")
}

// FetchScoreboardByDates retrieves NFL scoreboard for a date range (format: YYYYMMDD-YYYYMMDD)
func (c *ESPNClient) FetchScoreboardByDates(dates string) (*models.ScoreboardResponse, error) {
	url := c.baseURL + scoreboardPath
	if dates != "" {
		if !dateRangePattern.MatchString(dates) && !singleDatePattern.MatchString(dates) {
			return nil, fmt.Errorf("invalid date format: expected YYYYMMDD or YYYYMMDD-YYYYMMDD")
		}
		url = fmt.Sprintf("%s?dates=%s", url, dates)
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scoreboard: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var scoreboard models.ScoreboardResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&scoreboard); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &scoreboard, nil
}

// FetchGameSummary retrieves detailed game data including plays
func (c *ESPNClient) FetchGameSummary(gameID string) (*models.SummaryResponse, error) {
	if !gameIDPattern.MatchString(gameID) {
		return nil, fmt.Errorf("invalid game ID: expected numeric value")
	}

	url := fmt.Sprintf("%s%s?event=%s", c.baseURL, summaryPath, gameID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch game summary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var summary models.SummaryResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&summary); err != nil {
		return nil, fmt.Errorf("failed to parse game summary: %w", err)
	}

	return &summary, nil
}

// FetchGameReplay retrieves full game data for replay mode
func (c *ESPNClient) FetchGameReplay(gameID string) (*models.GameReplay, error) {
	summary, err := c.FetchGameSummary(gameID)
	if err != nil {
		return nil, err
	}
	return summary.ToGameReplay(), nil
}

// FetchGameStats retrieves game statistics
func (c *ESPNClient) FetchGameStats(gameID string) (*models.GameStats, error) {
	summary, err := c.FetchGameSummary(gameID)
	if err != nil {
		return nil, err
	}
	return summary.ToGameStats(), nil
}
