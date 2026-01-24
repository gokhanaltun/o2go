package o2go

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type ExchangeResponse struct {
	Raw  []byte
	Data map[string]any
}

func (e *ExchangeResponse) Decode(v any) error {
	return json.Unmarshal(e.Raw, v)
}

func (o *OAuth2) ExchangeAuthCode(ctx context.Context, code string) (*ExchangeResponse, error) {
	if strings.TrimSpace(code) == "" {
		return nil, ErrEmptyCode
	}

	data := url.Values{}
	data.Set("client_id", o.Config.ClientID)
	data.Set("client_secret", o.Config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", o.Config.RedirectURL)
	data.Set("grant_type", "authorization_code")

	return o.exchange(ctx, data, baseReservedParams(nil))
}

func (o *OAuth2) ExchangeRefreshToken(ctx context.Context, refreshToken string) (*ExchangeResponse, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrEmptyRefreshToken
	}

	data := url.Values{}
	data.Set("client_id", o.Config.ClientID)
	data.Set("client_secret", o.Config.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("redirect_uri", o.Config.RedirectURL)
	data.Set("grant_type", "refresh_token")

	return o.exchange(ctx, data, baseReservedParams([]string{"grant_type"}))
}

func (o *OAuth2) exchange(ctx context.Context, form url.Values, reserved map[string]struct{}) (*ExchangeResponse, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	parseParams(reserved, o.Config.TokenParams, func(key, value string) {
		form.Set(key, value)
	})

	rawResp, err := postForm(ctx, o.Provider.TokenURL(), form, o.HTTPClient)
	if err != nil {
		if he, ok := err.(*HTTPError); ok {
			return nil, he
		}
		return nil, &HTTPError{Err: fmt.Errorf("exchange: failed to post token request: %v", err)}
	}

	exchangeResponse := &ExchangeResponse{
		Raw: rawResp,
	}

	respMap := map[string]any{}
	if unmarshalErr := json.Unmarshal(rawResp, &respMap); unmarshalErr != nil {
		return exchangeResponse, &HTTPError{
			Err:  fmt.Errorf("exchange: failed to unmarshal token response: %v", unmarshalErr),
			Body: rawResp,
		}
	}

	exchangeResponse.Data = respMap
	return exchangeResponse, nil
}
