package shift

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/oncall-scheduler/pkg/opsgenieclient"
	"github.com/spf13/cobra"
)

type runner struct {
	flag   *flag
	logger micrologger.Logger
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	var opsgenie *opsgenieclient.Client
	{
		c := opsgenieclient.Config{
			Logger: r.logger,

			APIKey: r.flag.OpsGenieAPIKey,
		}

		opsgenie, err = opsgenieclient.New(c)
		if err != nil {
			return microerror.Mask(err)
		}
	}

	// var slack *slackclient.Client
	// {
	// 	c := slackclient.Config{
	// 		Logger: r.logger,
	//
	// 		Channel: r.flag.SlackChannel,
	// 		Token:   r.flag.SlackToken,
	// 	}
	//
	// 	slack, err = slackclient.New(c)
	// 	if err != nil {
	// 		return microerror.Mask(err)
	// 	}
	// }

	summary, err := opsgenie.GetShiftSummary(
		"ops_schedule",
		time.Date(2019, time.June, 1, 00, 0, 0, 0, time.UTC),
	)
	if err != nil {
		return microerror.Mask(err)
	}

	fmt.Println(summary)

	return nil
}
