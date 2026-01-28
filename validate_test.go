package o2go

import (
	"net/http"
	"testing"
)

type mockProvider struct {
	authURL  string
	tokenURL string
}

func (m *mockProvider) AuthURL() string  { return m.authURL }
func (m *mockProvider) TokenURL() string { return m.tokenURL }

func TestValidate(t *testing.T) {
	validConfig := &Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://example.com/callback",
	}
	validProvider := &mockProvider{
		authURL:  "https://example.com/auth",
		tokenURL: "https://example.com/token",
	}

	tests := []struct {
		name    string
		o       *OAuth2
		wantErr error
	}{
		{
			name:    "nil instance",
			o:       nil,
			wantErr: ErrNilInstance,
		},
		{
			name: "nil config",
			o: &OAuth2{
				Provider:   validProvider,
				HTTPClient: http.DefaultClient,
			},
			wantErr: ErrNilConfig,
		},
		{
			name: "nil provider",
			o: &OAuth2{
				Config:     validConfig,
				HTTPClient: http.DefaultClient,
			},
			wantErr: ErrNilProvider,
		},
		{
			name: "empty auth url",
			o: &OAuth2{
				Config: validConfig,
				Provider: &mockProvider{
					authURL:  "",
					tokenURL: "https://example.com/token",
				},
				HTTPClient: http.DefaultClient,
			},
			wantErr: ErrEmptyAuthURL,
		},
		{
			name: "empty token url",
			o: &OAuth2{
				Config: validConfig,
				Provider: &mockProvider{
					authURL:  "https://example.com/auth",
					tokenURL: "",
				},
				HTTPClient: http.DefaultClient,
			},
			wantErr: ErrEmptyTokenURL,
		},
		{
			name: "empty client id",
			o: &OAuth2{
				Config: &Config{
					ClientID:     "",
					ClientSecret: "secret",
					RedirectURL:  "url",
				},
				Provider:   validProvider,
				HTTPClient: http.DefaultClient,
			},
			wantErr: ErrEmptyClientID,
		},
		{
			name: "empty client secret",
			o: &OAuth2{
				Config: &Config{
					ClientID:     "id",
					ClientSecret: "",
					RedirectURL:  "url",
				},
				Provider:   validProvider,
				HTTPClient: http.DefaultClient,
			},
			wantErr: ErrEmptyClientSecret,
		},
		{
			name: "empty redirect url",
			o: &OAuth2{
				Config: &Config{
					ClientID:     "id",
					ClientSecret: "secret",
					RedirectURL:  "",
				},
				Provider:   validProvider,
				HTTPClient: http.DefaultClient,
			},
			wantErr: ErrEmptyRedirectURL,
		},
		{
			name: "happy path",
			o: &OAuth2{
				Config:     validConfig,
				Provider:   validProvider,
				HTTPClient: http.DefaultClient,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.validate()
			if err != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
