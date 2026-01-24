package o2go

import "net/http"

type OAuth2 struct {
	Provider   Provider
	Config     *Config
	HTTPClient *http.Client
}

type Option func(*OAuth2)

func WithHttpClient(client *http.Client) Option {
	return func(o *OAuth2) {
		o.HTTPClient = client
	}
}

func New(provider Provider, cfg *Config, opts ...Option) (*OAuth2, error) {
	o := &OAuth2{
		Provider:   provider,
		Config:     cfg,
		HTTPClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(o)
	}

	if err := o.validate(); err != nil {
		return nil, err
	}

	return o, nil
}
