package o2go

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type mockRoundTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTrip(req)
}

func TestExchangeAuthCode(t *testing.T) {
	config := &Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://example.com/callback",
		TokenParams: map[string]string{
			"custom_param": "custom_value",
		},
	}
	provider := &mockProvider{
		authURL:  "https://example.com/auth",
		tokenURL: "https://example.com/token",
	}

	client := &http.Client{
		Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				req.ParseForm()
				if req.FormValue("code") != "auth-code" {
					t.Errorf("expected code auth-code, got %s", req.FormValue("code"))
				}
				if req.FormValue("custom_param") != "custom_value" {
					t.Errorf("expected custom_param custom_value, got %s", req.FormValue("custom_param"))
				}
				if req.FormValue("grant_type") != "authorization_code" {
					t.Errorf("expected grant_type authorization_code, got %s", req.FormValue("grant_type"))
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "secret-token"}`)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}

	o, _ := New(provider, config, WithHttpClient(client))

	resp, err := o.ExchangeAuthCode(context.Background(), "auth-code")
	if err != nil {
		t.Fatalf("ExchangeAuthCode() error = %v", err)
	}

	if resp.Data["access_token"] != "secret-token" {
		t.Errorf("expected access_token secret-token, got %v", resp.Data["access_token"])
	}
}

func TestExchangeAuthCode_EmptyCode(t *testing.T) {
	o := &OAuth2{}
	_, err := o.ExchangeAuthCode(context.Background(), "   ")
	if err != ErrEmptyCode {
		t.Errorf("expected ErrEmptyCode, got %v", err)
	}
}

func TestExchangeRefreshToken(t *testing.T) {
	config := &Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://example.com/callback",
	}
	provider := &mockProvider{
		tokenURL: "https://example.com/token",
		authURL:  "https://example.com/auth",
	}

	client := &http.Client{
		Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				req.ParseForm()
				if req.FormValue("refresh_token") != "ref-token" {
					t.Errorf("expected refresh_token ref-token, got %s", req.FormValue("refresh_token"))
				}
				if req.FormValue("grant_type") != "refresh_token" {
					t.Errorf("expected grant_type refresh_token, got %s", req.FormValue("grant_type"))
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "new-token"}`)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}

	o := &OAuth2{
		Config:     config,
		Provider:   provider,
		HTTPClient: client,
	}

	resp, err := o.ExchangeRefreshToken(context.Background(), "ref-token")
	if err != nil {
		t.Fatalf("ExchangeRefreshToken() error = %v", err)
	}

	if resp.Data["access_token"] != "new-token" {
		t.Errorf("expected access_token new-token, got %v", resp.Data["access_token"])
	}
}

func TestExchangeResponse_Decode(t *testing.T) {
	type Token struct {
		AccessToken string `json:"access_token"`
	}

	resp := &ExchangeResponse{
		Raw: []byte(`{"access_token": "decoded-token"}`),
	}

	var token Token
	err := resp.Decode(&token)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if token.AccessToken != "decoded-token" {
		t.Errorf("expected decoded-token, got %s", token.AccessToken)
	}
}

func TestExchange_UnmarshalError(t *testing.T) {
	provider := &mockProvider{
		tokenURL: "https://example.com/token",
		authURL:  "https://example.com/auth",
	}
	config := &Config{ClientID: "id", ClientSecret: "secret", RedirectURL: "url"}
	
	client := &http.Client{
		Transport: &mockRoundTripper{
			roundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
					Header:     make(http.Header),
				}, nil
			},
		},
	}

	o := &OAuth2{
		Config:     config,
		Provider:   provider,
		HTTPClient: client,
	}

	_, err := o.ExchangeAuthCode(context.Background(), "code")
	if err == nil {
		t.Fatal("expected unmarshal error, got nil")
	}
}
