package opsgenieclient

import (
	"net/http"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type Config struct {
	Logger micrologger.Logger

	APIKey string
}

type Client struct {
	client *http.Client
	logger micrologger.Logger

	apiKey string
}

func New(config Config) (*Client, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	if config.APIKey == "" {
		return nil, microerror.Maskf(invalidConfigError, "%T.APIKey must not be empty", config)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	c := &Client{
		client: httpClient,
		logger: config.Logger,

		apiKey: config.APIKey,
	}

	return c, nil
}
