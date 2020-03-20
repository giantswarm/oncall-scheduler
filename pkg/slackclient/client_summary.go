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
		total := fmt.Sprintf("total: *%v* (%v|%v%%)", item.CurrentCount.Total, item.Change.Diff.Total, item.Change.Percentage.Total)
		businessHours := fmt.Sprintf("bh: *%v* (%v|%v%%)", item.CurrentCount.BusinessHours, item.Change.Diff.BusinessHours, item.Change.Percentage.BusinessHours)
		nonBusinessHours := fmt.Sprintf("nbh: *%v* (%v|%v%%)", item.CurrentCount.NonBusinessHours, item.Change.Diff.NonBusinessHours, item.Change.Percentage.NonBusinessHours)

		t = t + fmt.Sprintf("Last %v: %s | %s | %s. \n", item.Display, total, businessHours, nonBusinessHours)
	}

	color := ""
	numGoods := 0
	for _, item := range summaryItem {
		if item.Change.Diff.Total <= 0 {
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
