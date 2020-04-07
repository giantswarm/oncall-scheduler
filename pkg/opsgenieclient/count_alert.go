package opsgenieclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
)

const (
	OpsGenieAPIAlertsCountUrl = "https://api.opsgenie.com/v2/alerts/count"

	queryFormatBusinessHours    = "createdAt < %v AND createdAt > %v AND (tag: stable or [Pingdom]) AND teams: %v AND NOT teams: ops_team"
	queryFormatNonBusinessHours = "createdAt < %v AND createdAt > %v AND (tag: stable or [Pingdom]) AND teams: %v AND teams: ops_team"
)

func (c *Client) CountAlerts(query string) (int, error) {
	type OpsgenieCountResponseData struct {
		Count int `json:"count"`
	}

	type OpsgenieCountResponse struct {
		Data OpsgenieCountResponseData `json:"data"`
	}

	req, err := http.NewRequest(http.MethodGet, OpsGenieAPIAlertsCountUrl, nil)
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

func (c *Client) GetAlertCountSummary(team string, periods []Period) (AlertCountSummary, error) {
	alertSummary := AlertCountSummary{}

	for _, period := range periods {
		var currentCount Count
		var previousCount Count
		var change Change

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

		change.Diff.BusinessHours, change.Percentage.BusinessHours = c.calculateChange(previousCount.BusinessHours, currentCount.BusinessHours)
		change.Diff.NonBusinessHours, change.Percentage.NonBusinessHours = c.calculateChange(previousCount.NonBusinessHours, currentCount.NonBusinessHours)
		change.Diff.Total, change.Percentage.Total = c.calculateChange(previousCount.Total, currentCount.Total)

		summaryItem := AlertCountSummaryItem{
			CurrentCount:  currentCount,
			PreviousCount: previousCount,
			Change:        change,
			Display:       period.Display,
		}

		alertSummary = append(alertSummary, summaryItem)
	}

	return alertSummary, nil

}

func (c *Client) GetCountSummary() (CountSummary, error) {
	excludeAlertsRouterTeam := false
	teams, err := c.GetTeams(excludeAlertsRouterTeam)
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

	summary := CountSummary{}

	for _, team := range teams {
		alertSummary, err := c.GetAlertCountSummary(team, periods)
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

func (c *Client) calculateChange(a, b int) (int, int) {
	var diff, percentage int

	if a == 0 && b == 0 {
		return 0, 0
	}
	if a == 0 && b > 0 {
		return b, 100
	}

	if a < b {
		diff = int(b - a)
		percentage = int((float64(b-a) / float64(a)) * 100)
	} else {
		diff = -int(a - b)
		percentage = -int((float64(a-b) / float64(a)) * 100)
	}

	return diff, percentage
}

func (c *Client) contains(a []string, b string) bool {
	for _, x := range a {
		if x == b {
			return true
		}
	}

	return false
}
