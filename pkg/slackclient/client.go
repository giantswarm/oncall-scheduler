package slackclient

import (
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"github.com/nlopes/slack"
)

type Config struct {
	Logger micrologger.Logger

	Channel string
	Token   string
}

type Client struct {
	client *slack.Client
	logger micrologger.Logger

	channel string
}

func New(config Config) (*Client, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.Channel == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Channel must not be empty", config)
	}
	if config.Token == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.Token must not be empty", config)
	}

	slackClient := slack.New(config.Token)

	c := &Client{
		client: slackClient,
		logger: config.Logger,

		channel: config.Channel,
	}

	return c, nil
}
