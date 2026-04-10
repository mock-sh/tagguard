package tags_test

import (
	"testing"

	"github.com/mokshg/tagguard/pkg/tags"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTags []tags.Tag
	}{
		{
			name:  "single tag",
			input: `json:"name"`,
			wantTags: []tags.Tag{
				{Key: "json", Value: "name", Options: nil, Raw: "name"},
			},
		},
		{
			name:  "tag with options",
			input: `json:"name,omitempty"`,
			wantTags: []tags.Tag{
				{Key: "json", Value: "name", Options: []string{"omitempty"}, Raw: "name,omitempty"},
			},
		},
		{
			name:  "multiple tags",
			input: `json:"name,omitempty" db:"user_name" validate:"required,min=2"`,
			wantTags: []tags.Tag{
				{Key: "json", Value: "name", Options: []string{"omitempty"}, Raw: "name,omitempty"},
				{Key: "db", Value: "user_name", Options: nil, Raw: "user_name"},
				{Key: "validate", Value: "required", Options: []string{"min=2"}, Raw: "required,min=2"},
			},
		},
		{
			name:     "empty string",
			input:    "",
			wantTags: nil,
		},
		{
			name:  "dash value",
			input: `json:"-"`,
			wantTags: []tags.Tag{
				{Key: "json", Value: "-", Options: nil, Raw: "-"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tags.Parse(tt.input)

			if len(got) != len(tt.wantTags) {
				t.Fatalf("Parse(%q): got %d tags, want %d\ngot: %+v", tt.input, len(got), len(tt.wantTags), got)
			}

			for i, want := range tt.wantTags {
				g := got[i]
				if g.Key != want.Key {
					t.Errorf("tag[%d].Key = %q, want %q", i, g.Key, want.Key)
				}
				if g.Value != want.Value {
					t.Errorf("tag[%d].Value = %q, want %q", i, g.Value, want.Value)
				}
				if g.Raw != want.Raw {
					t.Errorf("tag[%d].Raw = %q, want %q", i, g.Raw, want.Raw)
				}
			}
		})
	}
}
