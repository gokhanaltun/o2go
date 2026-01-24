package o2go

import (
	"fmt"
	"net/url"
)

func (o *OAuth2) AuthCodeURL(state string) (string, error) {
	if err := o.validate(); err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("client_id", o.Config.ClientID)
	params.Set("redirect_uri", o.Config.RedirectURL)
	params.Set("response_type", "code")

	parseParams(baseReservedParams(nil), o.Config.AuthURLParams, func(key, value string) {
		params.Set(key, value)
	})

	if state != "" {
		params.Set("state", state)
	}

	return fmt.Sprintf("%s?%s", o.Provider.AuthURL(), params.Encode()), nil
}
