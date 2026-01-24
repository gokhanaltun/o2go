package o2go

import "errors"

var (
	ErrNilInstance       = errors.New("oauth2 instance is nil")
	ErrNilConfig         = errors.New("config is nil")
	ErrNilProvider       = errors.New("provider is nil")
	ErrNilHTTPClient     = errors.New("http client is nil")
	ErrEmptyAuthURL      = errors.New("provider auth URL is empty")
	ErrEmptyTokenURL     = errors.New("provider token URL is empty")
	ErrEmptyClientID     = errors.New("client ID is empty")
	ErrEmptyClientSecret = errors.New("client secret is empty")
	ErrEmptyRedirectURL  = errors.New("redirect URL is empty")
	ErrEmptyCode         = errors.New("code parameter cannot be empty")
	ErrEmptyRefreshToken = errors.New("refresh token parameter cannot be empty")
)
