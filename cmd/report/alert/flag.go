package alert

import (
	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagOpsGenieAPIKey = "opsgenie-api-key"
	flagSlackChannel   = "slack-channel"
	flagSlackToken     = "slack-token"
)

const (
	defaultOpsGenieAPIKey = ""
	defaultSlackChannel   = "test_bot" // TODO: Change to #ops.
	defaultSlackToken     = ""
)

type flag struct {
	OpsGenieAPIKey string
	SlackChannel   string
	SlackToken     string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.OpsGenieAPIKey, flagOpsGenieAPIKey, defaultOpsGenieAPIKey, "OpsGenie API Key to authenticate with")
	cmd.Flags().StringVar(&f.SlackChannel, flagSlackChannel, defaultSlackChannel, "Slack channel to post to")
	cmd.Flags().StringVar(&f.SlackToken, flagSlackToken, defaultSlackToken, "Slack token to authenticate with")
}

func (f *flag) Validate() error {
	if f.OpsGenieAPIKey == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagOpsGenieAPIKey)
	}
	if f.SlackChannel == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagSlackChannel)
	}
	if f.SlackToken == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagSlackToken)
	}

	return nil
}
