package slackclient

import (
	"fmt"
	"sort"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/oncall-scheduler/pkg/opsgenieclient"
	"github.com/nlopes/slack"
)

func (c *Client) PostSummary(summary opsgenieclient.Summary) error {
	teamNames := []string{}
	for teamName := range summary {
		teamNames = append(teamNames, teamName)
	}
	sort.Strings(teamNames)

	attachments := []slack.Attachment{}

	for _, teamName := range teamNames {
		alertSummary := summary[teamName]

		attachment := c.buildAlertSummaryAttachment(teamName, alertSummary)
		attachments = append(attachments, attachment)
	}

	_, _, err := c.client.PostMessage(
		c.channel,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionText("Today's OpsGenie Alert Summary!", false),
		slack.MsgOptionAttachments(attachments...),
	)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (c *Client) buildAlertSummaryAttachment(team string, summaryItem opsgenieclient.AlertSummary) slack.Attachment {
	name := c.formatTeamName(team)

	t := ""
	for _, item := range summaryItem {
		totalDiff := diff(item.CurrentCount.Total, item.PreviousCount.Total)
		total := fmt.Sprintf("total: %v (%v|%v%%)", item.CurrentCount.Total, totalDiff, item.Change.Total)
		businessHoursDiff := diff(item.CurrentCount.BusinessHours, item.PreviousCount.BusinessHours)
		businessHours := fmt.Sprintf("bh: %v (%v|%v%%)", item.CurrentCount.BusinessHours, businessHoursDiff, item.Change.BusinessHours)
		nonBusinessHoursDiff := diff(item.CurrentCount.NonBusinessHours, item.PreviousCount.NonBusinessHours)
		nonBusinessHours := fmt.Sprintf("nbh: %v (%v|%v%%)", item.CurrentCount.NonBusinessHours, nonBusinessHoursDiff, item.Change.NonBusinessHours)

		t = t + fmt.Sprintf("Last %v: %s | %s | %s. \n", item.Display, total, businessHours, nonBusinessHours)
	}

	color := ""
	numGoods := 0
	for _, item := range summaryItem {
		if item.Change.Total <= 0 {
			numGoods++
		}
	}
	switch numGoods {
	case 3:
		color = "#28a745"
	case 2:
		color = "#007bff"
	case 1:
		color = "#ffc107"
	case 0:
		color = "#dc3545"
	}

	attachment := slack.Attachment{
		Title: name,
		Text:  t,
		Color: color,
	}

	return attachment
}

func (c *Client) formatTeamName(name string) string {
	if strings.HasSuffix(name, "_team") {
		name = strings.Replace(name, "_team", "", 1)
	}

	name = strings.Title(name)

	return name
}

func diff(num1, num2 int) string {
	var diff string

	switch change := num1 - num2; {
	case change < 0:
		diff = fmt.Sprintf("-%d", num2-num1)
	case change == 0:
		diff = "0"
	default:
		diff = fmt.Sprintf("+%d", num1-num2)
	}

	return diff
}
