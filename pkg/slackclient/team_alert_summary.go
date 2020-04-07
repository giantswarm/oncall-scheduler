package slackclient

import (
	"fmt"
	"sort"
	"strings"

	"github.com/giantswarm/microerror"
	"github.com/nlopes/slack"

	"github.com/giantswarm/oncall-scheduler/pkg/opsgenieclient"
)

func (c *Client) PostAlertSummaries(thread_ts string, summary opsgenieclient.AlertSummary) error {
	teamNames := []string{}
	for teamName := range summary {
		teamNames = append(teamNames, teamName)
	}
	sort.Strings(teamNames)

	for _, teamName := range teamNames {
		teamAlerts := summary[teamName]
		if len(teamAlerts) > 0 {
			attachment := c.buildAlertSummaryAttachment(teamAlerts)

			_, _, err := c.client.PostMessage(
				c.channel,
				slack.MsgOptionAsUser(true),
				slack.MsgOptionText(fmt.Sprintf("*%s*: non-business hours alerts from the last on-call:", c.formatTeamName(teamName)), false),
				slack.MsgOptionAttachments(attachment),
				slack.MsgOptionTS(thread_ts),
			)
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	return nil
}

func (c *Client) buildAlertSummaryAttachment(alerts []opsgenieclient.Alert) slack.Attachment {
	var alertAttachment []string
	{
		for _, alert := range alerts {
			var alertStatus string
			{
				if alert.Status == "open" {
					alertStatus = ":alert-red:"
				} else if alert.Status == "closed" {
					alertStatus = ":alert-green:"
				} else {
					alertStatus = ":question:"
				}
			}
			alertLink := fmt.Sprintf(opsgenieclient.OpsGenieAlertUrl, alert.ID)

			text := fmt.Sprintf("* %v <%v|%v>", alertStatus, alertLink, alert.Message)
			alertAttachment = append(alertAttachment, text)
		}
	}

	attachment := slack.Attachment{
		Text: strings.Join(alertAttachment, "\n"),
	}

	return attachment
}
