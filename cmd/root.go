package cmd

import (
	"io"
	"os"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/giantswarm/oncall-scheduler/cmd/report"
	"github.com/giantswarm/oncall-scheduler/cmd/schedule"
	"github.com/spf13/cobra"
)

const (
	name        = "oncall-scheduler"
	description = "Tool to schedule oncall shifts"
)

type Config struct {
	Logger micrologger.Logger
	Stderr io.Writer
	Stdout io.Writer
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}
	if config.Stderr == nil {
		config.Stderr = os.Stderr
	}
	if config.Stdout == nil {
		config.Stdout = os.Stdout
	}

	var err error

	var reportCmd *cobra.Command
	{
		c := report.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		reportCmd, err = report.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var scheduleCmd *cobra.Command
	{
		c := schedule.Config{
			Logger: config.Logger,
			Stderr: config.Stderr,
			Stdout: config.Stdout,
		}

		scheduleCmd, err = schedule.New(c)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	f := &flag{}

	r := &runner{
		flag:   f,
		logger: config.Logger,
		stderr: config.Stderr,
		stdout: config.Stdout,
	}

	c := &cobra.Command{
		Use:          name,
		Short:        description,
		Long:         description,
		RunE:         r.Run,
		SilenceUsage: true,
	}

	f.Init(c)

	c.AddCommand(reportCmd)
	c.AddCommand(scheduleCmd)

	return c, nil
}
