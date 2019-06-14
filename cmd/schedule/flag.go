package schedule

import (
	"time"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagHorizon      = "horizon"
	flagPastHorizon  = "past-horizon"
	flagScheduleName = "schedule-name"
	flagTeamName     = "team-name"
)

const (
	defaultHorizon      = time.Hour * 24 * 7
	defaultPastHorizon  = time.Hour * 24 * 7 * 4
	defaultScheduleName = "ops_schedule"
	defaultTeamName     = "ops_team"
)

type flag struct {
	Horizon      time.Duration
	PastHorizon  time.Duration
	ScheduleName string
	TeamName     string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().DurationVar(&f.Horizon, flagHorizon, defaultHorizon, "How far ahead to ensure oncallers are scheduled")
	cmd.Flags().DurationVar(&f.PastHorizon, flagPastHorizon, defaultPastHorizon, "How far back to look for scheduling")
	cmd.Flags().StringVar(&f.ScheduleName, flagScheduleName, defaultScheduleName, "OpsGenie schedule to ensure")
	cmd.Flags().StringVar(&f.TeamName, flagTeamName, defaultTeamName, "OpsGenie team to draw oncallers from")

}

func (f *flag) Validate() error {
	if f.Horizon == 0 {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagHorizon)
	}
	if f.PastHorizon == 0 {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagPastHorizon)
	}
	if f.ScheduleName == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagScheduleName)
	}
	if f.TeamName == "" {
		return microerror.Maskf(invalidFlagError, "--%s must not be empty", flagTeamName)
	}

	return nil
}
