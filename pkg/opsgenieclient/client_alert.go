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
)

type AlertSummary []AlertSummaryItem

type AlertSummaryItem struct {
	Count         int
	PreviousCount int
	Change        int
	Display       string
}

type Period struct {
	NumDays int
	Display string
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

func (c *Client) GetAlertSummary() (AlertSummary, error) {
	periods := []Period{
		{
			NumDays: 1,
			Display: "24 hours",
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
	query := "createdAt < %v AND createdAt > %v AND tag: stable"

	alertSummary := AlertSummary{}

	for _, period := range periods {
		count, err := c.CountAlerts(
			fmt.Sprintf(
				query,
				c.getUnixTime(time.Now(), 0),
				c.getUnixTime(time.Now(), period.NumDays),
			),
		)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		previousCount, err := c.CountAlerts(
			fmt.Sprintf(
				query,
				c.getUnixTime(time.Now(), period.NumDays),
				c.getUnixTime(time.Now(), period.NumDays*2),
			),
		)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		change := c.calculatePercentageChange(previousCount, count)

		summaryItem := AlertSummaryItem{
			Count:         count,
			PreviousCount: previousCount,
			Change:        change,
			Display:       period.Display,
		}

		alertSummary = append(alertSummary, summaryItem)
	}

	return alertSummary, nil
}

// getUnixTime returns the UNIX time in milliseconds,
// shifted into the past by dayShift days if dayShift != 0.
func (c *Client) getUnixTime(when time.Time, dayShift int) int64 {
	return when.AddDate(0, 0, -dayShift).UnixNano() / int64(time.Millisecond)
}

func (c *Client) calculatePercentageChange(a, b int) int {
	if a == 0 {
		return 0
	}

	if a < b {
		return int((float64(b-a) / float64(a)) * 100)
	} else {
		return -int((float64(a-b) / float64(a)) * 100)
	}
}
