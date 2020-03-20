package opsgenieclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
)

const (
	OpsGenieAPICountUrl = "https://api.opsgenie.com/v2/alerts/count"
	OpsGenieAPITeamUrl  = "https://api.opsgenie.com/v2/teams"

	queryFormatBusinessHours    = "createdAt < %v AND createdAt > %v AND (tag: stable or [Pingdom]) AND teams: %v AND NOT teams: ops_team"
	queryFormatNonBusinessHours = "createdAt < %v AND createdAt > %v AND (tag: stable or [Pingdom]) AND teams: %v AND teams: ops_team"
)

var (
	blocklist = []string{
		"se",
		"sre_team",
		"ops_team",
	}
)

type Summary map[string]AlertSummary

type AlertSummary []AlertSummaryItem

type AlertSummaryItem struct {
	CurrentCount  Count
	PreviousCount Count
	Change        Count
	Display       string
}

type Count struct {
	BusinessHours    int
	NonBusinessHours int
	Total            int
}

type Period struct {
	NumDays int
	Display string
}

func (c *Client) GetTeams() ([]string, error) {
	type OpsgenieTeamResponseData struct {
		Name string `json:name`
	}

	type OpsgenieTeamResponse struct {
		Data []OpsgenieTeamResponseData `json:"data"`
	}

	req, err := http.NewRequest(http.MethodGet, OpsGenieAPITeamUrl, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("GenieKey %v", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer resp.Body.Close()

	opsgenieResponse := &OpsgenieTeamResponse{}
	if err := json.NewDecoder(resp.Body).Decode(opsgenieResponse); err != nil {
		return nil, microerror.Mask(err)
	}

	teamNames := []string{}
	for _, team := range opsgenieResponse.Data {
		if c.contains(blocklist, team.Name) {
			continue
		}

		teamNames = append(teamNames, team.Name)
	}

	return teamNames, nil
}

func (c *Client) CountAlerts(query string) (int, error) {
	type OpsgenieCountResponseData struct {
		Count int `json:"count"`
	}

	type OpsgenieCountResponse struct {
		Data OpsgenieCountResponseData `json:"data"`
	}

	req, err := http.NewRequest(http.MethodGet, OpsGenieAPICountUrl, nil)
	if err != nil {
		return 0, microerror.Mask(err)
	}

	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("GenieKey %v", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, microerror.Mask(err)
	}
	defer resp.Body.Close()

	opsgenieResponse := &OpsgenieCountResponse{}
	if err := json.NewDecoder(resp.Body).Decode(opsgenieResponse); err != nil {
		return 0, microerror.Mask(err)
	}

	return opsgenieResponse.Data.Count, nil
}

func (c *Client) GetAlertSummary(team string, periods []Period) (AlertSummary, error) {
	alertSummary := AlertSummary{}

	for _, period := range periods {
		var currentCount Count
		var previousCount Count
		var change Count

		currentQueryBusinessHours := fmt.Sprintf(
			queryFormatBusinessHours,
			c.getUnixTime(time.Now(), 0),
			c.getUnixTime(time.Now(), period.NumDays),
			team,
		)
		currentCount.BusinessHours, _ = c.CountAlerts(currentQueryBusinessHours)

		currentNonBusinessHours := fmt.Sprintf(
			queryFormatNonBusinessHours,
			c.getUnixTime(time.Now(), 0),
			c.getUnixTime(time.Now(), period.NumDays),
			team,
		)
		currentCount.NonBusinessHours, _ = c.CountAlerts(currentNonBusinessHours)

		currentCount.Total = currentCount.BusinessHours + currentCount.NonBusinessHours

		previousQueryBusinessHours := fmt.Sprintf(
			queryFormatBusinessHours,
			c.getUnixTime(time.Now(), period.NumDays),
			c.getUnixTime(time.Now(), period.NumDays*2),
			team,
		)
		previousCount.BusinessHours, _ = c.CountAlerts(previousQueryBusinessHours)

		previousNonBusinessHours := fmt.Sprintf(
			queryFormatNonBusinessHours,
			c.getUnixTime(time.Now(), period.NumDays),
			c.getUnixTime(time.Now(), period.NumDays*2),
			team,
		)
		previousCount.NonBusinessHours, _ = c.CountAlerts(previousNonBusinessHours)

		previousCount.Total = previousCount.BusinessHours + previousCount.NonBusinessHours

		change.BusinessHours = c.calculatePercentageChange(previousCount.BusinessHours, currentCount.BusinessHours)
		change.NonBusinessHours = c.calculatePercentageChange(previousCount.NonBusinessHours, currentCount.NonBusinessHours)
		change.Total = c.calculatePercentageChange(previousCount.Total, currentCount.Total)

		summaryItem := AlertSummaryItem{
			CurrentCount:  currentCount,
			PreviousCount: previousCount,
			Change:        change,
			Display:       period.Display,
		}

		alertSummary = append(alertSummary, summaryItem)
	}

	return alertSummary, nil

}

func (c *Client) GetSummary() (Summary, error) {
	teams, err := c.GetTeams()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	periods := []Period{
		{
			NumDays: 1,
			Display: "day",
		},
		{
			NumDays: 7,
			Display: "week",
		},
		{
			NumDays: 30,
			Display: "month",
		},
	}

	summary := Summary{}

	for _, team := range teams {
		alertSummary, err := c.GetAlertSummary(team, periods)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		summary[team] = alertSummary
	}

	return summary, nil
}

// getUnixTime returns the UNIX time in milliseconds,
// shifted into the past by dayShift days if dayShift != 0.
func (c *Client) getUnixTime(when time.Time, dayShift int) int64 {
	return when.AddDate(0, 0, -dayShift).UnixNano() / int64(time.Millisecond)
}

func (c *Client) calculatePercentageChange(a, b int) int {
	if a == 0 && b == 0 {
		return 0
	}
	if a == 0 && b > 0 {
		return 100
	}

	if a < b {
		return int((float64(b-a) / float64(a)) * 100)
	} else {
		return -int((float64(a-b) / float64(a)) * 100)
	}
}

func (c *Client) contains(a []string, b string) bool {
	for _, x := range a {
		if x == b {
			return true
		}
	}

	return false
}
