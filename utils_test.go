package o2go

import (
	"reflect"
	"testing"
)

func TestParseParams(t *testing.T) {
	reserved := map[string]struct{}{
		"reserved1": {},
		"reserved2": {},
	}
	params := map[string]string{
		"reserved1": "val1",
		"custom1":   "customval1",
		"custom2":   "customval2",
	}

	got := make(map[string]string)
	parseParams(reserved, params, func(key, value string) {
		got[key] = value
	})

	want := map[string]string{
		"custom1": "customval1",
		"custom2": "customval2",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseParams() = %v, want %v", got, want)
	}
}

func TestBaseReservedParams(t *testing.T) {
	tests := []struct {
		name          string
		extraReserved []string
		wantKeys      []string
	}{
		{
			name:          "base only",
			extraReserved: nil,
			wantKeys:      []string{"client_id", "client_secret", "redirect_uri", "code", "refresh_token"},
		},
		{
			name:          "with extra",
			extraReserved: []string{"grant_type", "scope"},
			wantKeys:      []string{"client_id", "client_secret", "redirect_uri", "code", "refresh_token", "grant_type", "scope"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := baseReservedParams(tt.extraReserved)
			if len(got) != len(tt.wantKeys) {
				t.Errorf("baseReservedParams() returned %d keys, want %d", len(got), len(tt.wantKeys))
			}
			for _, key := range tt.wantKeys {
				if _, ok := got[key]; !ok {
					t.Errorf("baseReservedParams() missing key %s", key)
				}
			}
		})
	}
}
