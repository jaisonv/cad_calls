package filter

import (
	"strings"

	"github.com/jaisonv/telegram-cad-bot/internal/cad"
)

// MatchStreet checks if a CAD call address matches any of the monitored streets
// Uses case-insensitive partial matching
func MatchStreet(call *cad.CADCall, monitoredStreets []string) bool {
	if call.Address == "" {
		return false
	}

	callAddress := strings.ToLower(call.Address)

	for _, street := range monitoredStreets {
		monitoredStreet := strings.ToLower(strings.TrimSpace(street))
		if monitoredStreet == "" {
			continue
		}

		// Check if the monitored street name appears in the call address
		if strings.Contains(callAddress, monitoredStreet) {
			return true
		}

		// Also check for common abbreviations
		// e.g., "St" vs "Street", "Ave" vs "Avenue", etc.
		if matchesWithAbbreviations(callAddress, monitoredStreet) {
			return true
		}
	}

	return false
}

// matchesWithAbbreviations handles common street abbreviations
func matchesWithAbbreviations(address, street string) bool {
	abbreviations := map[string][]string{
		"street":    {"st", "str"},
		"avenue":    {"ave", "av"},
		"road":      {"rd"},
		"drive":     {"dr"},
		"boulevard": {"blvd", "boul"},
		"lane":      {"ln"},
		"court":     {"ct"},
		"place":     {"pl"},
		"circle":    {"cir"},
		"way":       {"wy"},
	}

	// Check if expanding abbreviations creates a match
	for full, abbrevs := range abbreviations {
		for _, abbrev := range abbrevs {
			// If street contains abbreviated form, try full form in address
			if strings.Contains(street, " "+abbrev) || strings.HasSuffix(street, " "+abbrev) {
				expandedStreet := strings.ReplaceAll(street, " "+abbrev, " "+full)
				if strings.Contains(address, expandedStreet) {
					return true
				}
			}

			// If address contains abbreviated form, try full form in street
			if strings.Contains(address, " "+abbrev) || strings.Contains(address, abbrev+" ") {
				expandedAddress := strings.ReplaceAll(address, " "+abbrev, " "+full)
				if strings.Contains(expandedAddress, street) {
					return true
				}
			}
		}
	}

	return false
}

// FilterCalls returns only the calls that match the monitored streets
func FilterCalls(calls []cad.CADCall, monitoredStreets []string) []cad.CADCall {
	var matched []cad.CADCall

	for _, call := range calls {
		if MatchStreet(&call, monitoredStreets) {
			matched = append(matched, call)
		}
	}

	return matched
}
