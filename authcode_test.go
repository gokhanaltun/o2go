package o2go

import (
	"strings"
	"testing"
)

func TestAuthCodeURL(t *testing.T) {
	config := &Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RedirectURL:  "https://example.com/callback",
		AuthURLParams: map[string]string{
			"scope": "user:email",
			"extra": "value",
		},
	}
	provider := &mockProvider{
		authURL:  "https://example.com/auth",
		tokenURL: "https://example.com/token",
	}

	o, _ := New(provider, config)

	tests := []struct {
		name       string
		state      string
		wantSubstr []string
	}{
		{
			name:  "basic",
			state: "xyz",
			wantSubstr: []string{
				"https://example.com/auth",
				"client_id=client-id",
				"redirect_uri=https%3A%2F%2Fexample.com%2Fcallback",
				"response_type=code",
				"state=xyz",
				"scope=user%3Aemail",
				"extra=value",
			},
		},
		{
			name:  "no state",
			state: "",
			wantSubstr: []string{
				"client_id=client-id",
				"redirect_uri=https%3A%2F%2Fexample.com%2Fcallback",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := o.AuthCodeURL(tt.state)
			if err != nil {
				t.Fatalf("AuthCodeURL() error = %v", err)
			}

			for _, substr := range tt.wantSubstr {
				if !strings.Contains(got, substr) {
					t.Errorf("AuthCodeURL() = %v, want to contain %v", got, substr)
				}
			}

			if tt.state == "" && strings.Contains(got, "state=") {
				t.Errorf("AuthCodeURL() = %v, should not contain state parameter", got)
			}
		})
	}
}

func TestAuthCodeURL_ValidationFailure(t *testing.T) {
	o := &OAuth2{}
	_, err := o.AuthCodeURL("state")
	if err == nil {
		t.Error("AuthCodeURL() expected error due to validation failure, got nil")
	}
}
