package filter

import (
	"testing"

	"github.com/jaisonv/telegram-cad-bot/internal/cad"
)

func TestNormalizeStreet(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"S. Broadway", "s broadway"},
		{"Main   Street", "main street"},
		{"Oak Ave.", "oak ave"},
		{"N.  Main  St.", "n main st"},
		{"1st, 2nd Ave", "1st 2nd ave"},
		{"  Extra  Spaces  ", "extra spaces"},
		{"UPPERCASE", "uppercase"},
	}

	for _, tt := range tests {
		result := normalizeStreet(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeStreet(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestMatchStreet(t *testing.T) {
	tests := []struct {
		name            string
		callAddress     string
		monitoredStreet string
		shouldMatch     bool
	}{
		{
			name:            "Period in monitored street",
			callAddress:     "123 S Broadway Ave",
			monitoredStreet: "S. Broadway",
			shouldMatch:     true,
		},
		{
			name:            "Period in address",
			callAddress:     "123 S. Broadway Ave",
			monitoredStreet: "S Broadway",
			shouldMatch:     true,
		},
		{
			name:            "Extra spaces in monitored street",
			callAddress:     "456 Main Street",
			monitoredStreet: "Main   Street",
			shouldMatch:     true,
		},
		{
			name:            "Different case",
			callAddress:     "789 OAK AVENUE",
			monitoredStreet: "oak avenue",
			shouldMatch:     true,
		},
		{
			name:            "Abbreviation with period",
			callAddress:     "100 N Main St",
			monitoredStreet: "N. Main",
			shouldMatch:     true,
		},
		{
			name:            "No match",
			callAddress:     "123 Elm Street",
			monitoredStreet: "Oak Avenue",
			shouldMatch:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := &cad.CADCall{
				Address: tt.callAddress,
			}
			monitoredStreets := []string{tt.monitoredStreet}
			result := MatchStreet(call, monitoredStreets)

			if result != tt.shouldMatch {
				t.Errorf("MatchStreet(%q, %q) = %v, want %v",
					tt.callAddress, tt.monitoredStreet, result, tt.shouldMatch)
			}
		})
	}
}
