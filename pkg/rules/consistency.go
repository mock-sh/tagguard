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
