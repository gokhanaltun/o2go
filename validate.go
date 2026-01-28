package o2go

import "strings"

func (o *OAuth2) validate() error {
	if o == nil {
		return ErrNilInstance
	}

	if o.Config == nil {
		return ErrNilConfig
	}

	if o.Provider == nil {
		return ErrNilProvider
	}

	if strings.TrimSpace(o.Provider.AuthURL()) == "" {
		return ErrEmptyAuthURL
	}

	if strings.TrimSpace(o.Provider.TokenURL()) == "" {
		return ErrEmptyTokenURL
	}

	if strings.TrimSpace(o.Config.ClientID) == "" {
		return ErrEmptyClientID
	}

	if strings.TrimSpace(o.Config.ClientSecret) == "" {
		return ErrEmptyClientSecret
	}

	if strings.TrimSpace(o.Config.RedirectURL) == "" {
		return ErrEmptyRedirectURL
	}

	return nil
}
