package schedule

import (
	"context"
	"fmt"
	"io"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
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
	fmt.Println("fucking schedule it yourself")

	// Fetch oncallers for team `teamName`
	// Calculate days that do not have overrides (they are unscheduled) between now and `horizon`
	//
	// For each 'unscheduled' day:
	//     Copy oncallers
	//
	//     Fetch AFK events on, or covering the day
	//     Remove any oncallers that have AFKs on this day
	//     If the day is a Saturday or a Sunday, remove any oncallers that were oncall the previous Saturday or Sunday
	//     Remove any oncallers that have an AFK lasting more than 5 days finishing yesterday
	//
	//     Fetch OpsGenie schedule `schedule-name` from `pastHorizon` to day - 1 day
	//     Calculate the oncaller that has not been oncall the longest
	//
	//     Add override for the determined oncaller to `schedule-name`

	return nil
}
