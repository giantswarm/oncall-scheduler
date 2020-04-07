package opsgenieclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
)

const (
	OpsGenieAlertUrl     = "https://giantswarm.app.opsgenie.com/alert/detail/%s/details"
	OpsGenieAPIAlertsUrl = "https://api.opsgenie.com/v2/alerts"
)

// GetTeamAlerts returns list of OpsGenie alerts matched query.
func (c *Client) QueryAlerts(query string) ([]Alert, error) {
	type OpsgenieAlertsResponse struct {
		Data []Alert `json:"data"`
	}

	req, err := http.NewRequest(http.MethodGet, OpsGenieAPIAlertsUrl, nil)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("GenieKey %v", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	defer resp.Body.Close()

	opsgenieResponse := &OpsgenieAlertsResponse{}
	if err := json.NewDecoder(resp.Body).Decode(opsgenieResponse); err != nil {
		return nil, microerror.Mask(err)
	}

	return opsgenieResponse.Data, nil
}

func (c *Client) GetNonBusinessAlertSummary() (AlertSummary, error) {
	excludeAlertsRouterTeam := true
	teams, err := c.GetTeams(excludeAlertsRouterTeam)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	summary := AlertSummary{}

	for _, team := range teams {
		teamNonBusinessHoursAlertsQuery := fmt.Sprintf(
			queryFormatNonBusinessHours,
			c.getUnixTime(time.Now(), 0),
			c.getUnixTime(time.Now(), 1),
			team,
		)

		teamAlerts, err := c.QueryAlerts(teamNonBusinessHoursAlertsQuery)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		summary[team] = teamAlerts
	}

	return summary, nil
}
