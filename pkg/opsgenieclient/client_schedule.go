package opsgenieclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
)

const (
	OpsGenieAPIScheduleUrlFormat = "https://api.opsgenie.com/v2/schedules/%s/timeline?identifierType=name&date=%v&interval=1&intervalUnit=days"
)

type ShiftSummary []ShiftSummaryItem

type ShiftSummaryItem struct {
	Email string
	Count int
}

func (c *Client) GetShiftSummary(schedule string, t time.Time) (ShiftSummary, error) {
	type OpsGenieScheduleResponseRecipient struct {
		Name string `json:"name"`
	}

	type OpsGenieScheduleResponsePeriod struct {
		StartDate time.Time                         `json:"startDate"`
		EndDate   time.Time                         `json:"endDate"`
		Recipient OpsGenieScheduleResponseRecipient `json:"recipient"`
	}

	type OpsGenieScheduleResponseRotation struct {
		Periods []OpsGenieScheduleResponsePeriod `json:"periods"`
	}

	type OpsGenieScheduleResponseFinalTimeline struct {
		Rotations []OpsGenieScheduleResponseRotation `json:"rotations"`
	}

	type OpsGenieScheduleResponseData struct {
		FinalTimeline OpsGenieScheduleResponseFinalTimeline `json:"finalTimeline"`
	}

	type OpsGenieScheduleResponse struct {
		Data OpsGenieScheduleResponseData `json:"data"`
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(OpsGenieAPIScheduleUrlFormat, schedule, t.Format("2006-01-02T15:04:05Z")),
		nil,
	)
	if err != nil {
		return ShiftSummary{}, nil
	}

	req.Header.Add("Authorization", fmt.Sprintf("GenieKey %v", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return ShiftSummary{}, nil
	}
	defer resp.Body.Close()

	opsgenieResponse := &OpsGenieScheduleResponse{}
	if err := json.NewDecoder(resp.Body).Decode(opsgenieResponse); err != nil {
		return ShiftSummary{}, microerror.Mask(err)
	}

	fmt.Println(opsgenieResponse)

	return ShiftSummary{}, nil
}
