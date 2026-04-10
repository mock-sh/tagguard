package rules

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mokshg/tagguard/pkg/tags"
)

// namingStyle represents a detected naming convention.
type namingStyle string

const (
	styleSnakeCase  namingStyle = "snake_case"
	styleCamelCase  namingStyle = "camelCase"
	stylePascalCase namingStyle = "PascalCase"
	styleKebabCase  namingStyle = "kebab-case"
	styleUnknown    namingStyle = "unknown"
	styleDash       namingStyle = "-" // means "omit this field"
)

// serializationTags are tags whose values represent field names in output formats.
// We check these for naming consistency with each other.
var serializationTags = map[string]bool{
	"json":  true,
	"yaml":  true,
	"xml":   true,
	"toml":  true,
	"bson":  true,
	"db":    true,
	"form":  true,
	"query": true,
	"uri":   true,
	"param": true,
}

// detectStyle infers the naming convention of a field name value.
func detectStyle(value string) namingStyle {
	if value == "" || value == "-" {
		return styleDash
	}
	if strings.Contains(value, "_") {
		return styleSnakeCase
	}
	if strings.Contains(value, "-") {
		return styleKebabCase
	}
	if len(value) > 0 && unicode.IsUpper(rune(value[0])) {
		return stylePascalCase
	}
	// Check for camelCase: lowercase start but contains an uppercase
	for _, r := range value[1:] {
		if unicode.IsUpper(r) {
			return styleCamelCase
		}
	}
	// All lowercase — could be snake_case without underscores or just lowercase
	return styleSnakeCase
}

// FieldTagEntry holds a tag value and its detected style for a single field.
type FieldTagEntry struct {
	FieldName string
	TagKey    string
	Value     string
	Style     namingStyle
}

// CollectFieldStyles extracts style entries from a field's tags.
// Used to build up per-struct style data for struct-wide consistency checks.
func CollectFieldStyles(fieldName string, fieldTags []tags.Tag) []FieldTagEntry {
	var entries []FieldTagEntry
	for _, t := range fieldTags {
		if !serializationTags[t.Key] {
			continue
		}
		if t.Value == "" || t.Value == "-" {
			continue
		}
		s := detectStyle(t.Value)
		if s == styleUnknown || s == styleDash {
			continue
		}
		entries = append(entries, FieldTagEntry{
			FieldName: fieldName,
			TagKey:    t.Key,
			Value:     t.Value,
			Style:     s,
		})
	}
	return entries
}

// CheckStructNamingConsistency checks that all fields in a struct use the same
// naming style for each tag key. For example, if most json tags use snake_case
// but one uses camelCase, that field is flagged.
//
// entries should contain the collected FieldTagEntry values for all fields in one struct.
func CheckStructNamingConsistency(entries []FieldTagEntry) map[string][]string {
	// Group entries by tag key
	byKey := map[string][]FieldTagEntry{}
	for _, e := range entries {
		byKey[e.TagKey] = append(byKey[e.TagKey], e)
	}

	issues := map[string][]string{} // fieldName → issues

	for tagKey, keyEntries := range byKey {
		if len(keyEntries) < 2 {
			continue
		}

		// Count styles for this tag key across the struct
		counts := map[namingStyle]int{}
		for _, e := range keyEntries {
			counts[e.Style]++
		}

		// Find dominant style
		dominant := styleUnknown
		maxCount := 0
		for style, count := range counts {
			if count > maxCount {
				maxCount = count
				dominant = style
			}
		}

		for _, e := range keyEntries {
			if e.Style != dominant {
				msg := fmt.Sprintf(
					"struct-wide inconsistency: %s tag on field %q uses %s (%q) but most fields use %s",
					tagKey, e.FieldName, e.Style, e.Value, dominant,
				)
				issues[e.FieldName] = append(issues[e.FieldName], msg)
			}
		}
	}

	return issues
}

// CheckEnforcedNamingStyle checks a single tag value against a required naming style.
// Returns an issue string if the value doesn't match, or "" if it's fine.
func CheckEnforcedNamingStyle(tagKey, value, requiredStyle string) string {
	if value == "" || value == "-" {
		return ""
	}
	if !serializationTags[tagKey] {
		return ""
	}
	detected := string(detectStyle(value))
	if detected != requiredStyle {
		return fmt.Sprintf(
			"%s tag value %q uses %s but project requires %s",
			tagKey, value, detected, requiredStyle,
		)
	}
	return ""
}

// CheckNamingConsistency checks that serialization tags on one field use a consistent
// naming style. For example, mixing json:"userId" and db:"user_id" is flagged.
func CheckNamingConsistency(fieldTags []tags.Tag) []string {
	type styleEntry struct {
		key   string
		value string
		style namingStyle
	}

	var entries []styleEntry
	for _, t := range fieldTags {
		if !serializationTags[t.Key] {
			continue
		}
		if t.Value == "" || t.Value == "-" {
			continue
		}
		s := detectStyle(t.Value)
		if s == styleUnknown || s == styleDash {
			continue
		}
		entries = append(entries, styleEntry{key: t.Key, value: t.Value, style: s})
	}

	if len(entries) < 2 {
		return nil
	}

	// Find the dominant style (most common)
	counts := map[namingStyle]int{}
	for _, e := range entries {
		counts[e.style]++
	}
	dominant := styleUnknown
	maxCount := 0
	for style, count := range counts {
		if count > maxCount {
			maxCount = count
			dominant = style
		}
	}

	var issues []string
	for _, e := range entries {
		if e.style != dominant {
			issues = append(issues, fmt.Sprintf(
				"inconsistent naming: %s tag uses %s (%q) but other tags use %s",
				e.key, e.style, e.value, dominant,
			))
		}
	}
	return issues
}
