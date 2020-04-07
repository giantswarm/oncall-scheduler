package opsgenieclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/giantswarm/microerror"
)

const (
	OpsGenieAPITeamUrl = "https://api.opsgenie.com/v2/teams"

	alertsRouterTeam = "alerts_router_team"
)

var (
	blocklist = []string{
		"se",
		"sre_team",
		"ops_team",
	}
)

func (c *Client) GetTeams(excludeAlertsRouterTeam bool) ([]string, error) {
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

		if team.Name == alertsRouterTeam && excludeAlertsRouterTeam {
			continue
		}

		teamNames = append(teamNames, team.Name)
	}

	return teamNames, nil
}
