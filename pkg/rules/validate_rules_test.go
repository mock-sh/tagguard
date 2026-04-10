package rules_test

import (
	"strings"
	"testing"

	"github.com/mokshg/tagguard/pkg/rules"
)

func TestCheckValidateRules(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantIssue string // substring expected in first issue, or "" for no issues
	}{
		{
			name:      "valid rules",
			input:     "required,min=2,max=50,email",
			wantIssue: "",
		},
		{
			name:      "omitempty is valid",
			input:     "omitempty,gte=0",
			wantIssue: "",
		},
		{
			name:      "typo required",
			input:     "requred",
			wantIssue: `did you mean "required"`,
		},
		{
			name:      "typo email",
			input:     "emai",
			wantIssue: `did you mean "email"`,
		},
		{
			name:      "completely unknown rule",
			input:     "superstrongvalidation",
			wantIssue: `unknown validate rule`,
		},
		{
			name:      "dash is valid (skip field)",
			input:     "-",
			wantIssue: "",
		},
		{
			name:      "mixed valid and invalid",
			input:     "required,emial,min=2",
			wantIssue: `did you mean "email"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := rules.CheckValidateRules(tt.input)

			if tt.wantIssue == "" {
				if len(issues) > 0 {
					t.Errorf("CheckValidateRules(%q): unexpected issues: %v", tt.input, issues)
				}
				return
			}

			if len(issues) == 0 {
				t.Fatalf("CheckValidateRules(%q): expected issue containing %q, got none", tt.input, tt.wantIssue)
			}

			found := false
			for _, issue := range issues {
				if strings.Contains(issue, tt.wantIssue) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("CheckValidateRules(%q): expected issue containing %q\ngot: %v", tt.input, tt.wantIssue, issues)
			}
		})
	}
}

func TestCheckTagKey(t *testing.T) {
	tests := []struct {
		key        string
		wantKnown  bool
		wantSuggest string
	}{
		{"json", true, ""},
		{"yaml", true, ""},
		{"validate", true, ""},
		{"jsno", false, "json"},
		{"ymal", false, "yaml"},
		{"foobar", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			suggestion, known := rules.CheckTagKey(tt.key, nil)
			if known != tt.wantKnown {
				t.Errorf("CheckTagKey(%q): known=%v, want %v", tt.key, known, tt.wantKnown)
			}
			if suggestion != tt.wantSuggest {
				t.Errorf("CheckTagKey(%q): suggestion=%q, want %q", tt.key, suggestion, tt.wantSuggest)
			}
		})
	}
}
