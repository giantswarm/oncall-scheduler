package slackclient

import (
	"fmt"
	"math"
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
		absChange := int(math.Abs(float64(item.Change)))

		switch change := item.Change; {
		case change < 0:
			t = t + fmt.Sprintf("%v alerts over the last %v (%v fewer alerts, decrease of %v%%)\n", item.Count, item.Display, item.PreviousCount-item.Count, absChange)
		case change == 0:
			t = t + fmt.Sprintf("%v alerts over the last %v, same as previous\n", item.Count, item.Display)
		default:
			t = t + fmt.Sprintf("%v alerts over the last %v (%v more alerts, increase of %v%%)\n", item.Count, item.Display, item.Count-item.PreviousCount, absChange)
		}
	}

	color := ""
	numGoods := 0
	for _, item := range summaryItem {
		if item.Change <= 0 {
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
