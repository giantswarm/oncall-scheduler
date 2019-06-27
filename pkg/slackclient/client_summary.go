package slackclient

import (
	"fmt"
	"math"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/oncall-scheduler/pkg/opsgenieclient"
	"github.com/nlopes/slack"
)

func (c *Client) PostAlertSummary(summary opsgenieclient.AlertSummary) error {
	attachments := []slack.Attachment{}

	for _, summaryItem := range summary {
		attachment := c.buildAttachment(summaryItem)
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

// buildAttachment builds a Slack attachment for a given AlertSummaryItem.
// See https://api.slack.com/docs/formatting for formatting docs.
func (c *Client) buildAttachment(summaryItem opsgenieclient.AlertSummaryItem) slack.Attachment {
	attachment := slack.Attachment{
		Text: fmt.Sprintf("%v alerts over the last %v", summaryItem.Count, summaryItem.Display),
	}

	absChange := int(math.Abs(float64(summaryItem.Change)))

	switch change := summaryItem.Change; {
	case change < 0:
		attachment.Color = "good"
		attachment.Footer = fmt.Sprintf("Decrease of %v%%, %v fewer alerts!", absChange, summaryItem.PreviousCount-summaryItem.Count)
	case change == 0:
		attachment.Color = "info"
		attachment.Footer = "No change~"
	default:
		attachment.Color = "danger"
		attachment.Footer = fmt.Sprintf("Increase of %v%%, %v more alerts", absChange, summaryItem.Count-summaryItem.PreviousCount)
	}

	return attachment
}
