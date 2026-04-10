package rules

import (
	"fmt"
	"strings"
)

// knownValidateRules lists all valid rule names from go-playground/validator v10.
// Source: https://pkg.go.dev/github.com/go-playground/validator/v10
var knownValidateRules = []string{
	// Existence
	"required",
	"required_if",
	"required_unless",
	"required_with",
	"required_with_all",
	"required_without",
	"required_without_all",
	"excluded_if",
	"excluded_unless",
	"excluded_with",
	"excluded_with_all",
	"excluded_without",
	"excluded_without_all",
	"omitempty",
	"omitnil",
	"isdefault",

	// Comparisons
	"eq",
	"eq_ignore_case",
	"gt",
	"gte",
	"lt",
	"lte",
	"ne",
	"ne_ignore_case",

	// Other
	"len",
	"max",
	"min",
	"oneof",

	// Strings
	"alpha",
	"alphanum",
	"alphanumunicode",
	"alphaunicode",
	"ascii",
	"boolean",
	"contains",
	"containsany",
	"containsrune",
	"endsnotwith",
	"endswith",
	"excludes",
	"excludesall",
	"excludesrune",
	"lowercase",
	"multibyte",
	"number",
	"numeric",
	"printascii",
	"startsnotwith",
	"startswith",
	"uppercase",

	// Format
	"base32",
	"base64",
	"base64url",
	"base64rawurl",
	"bic",
	"bcp47_language_tag",
	"btc_addr",
	"btc_addr_bech32",
	"credit_card",
	"mongodb",
	"mongodb_connection_string",
	"cron",
	"spicedb",
	"datetime",
	"e164",
	"email",
	"eth_addr",
	"hexadecimal",
	"hexcolor",
	"hsl",
	"hsla",
	"html",
	"html_encoded",
	"http_url",
	"url",
	"uri",
	"url_encoded",
	"urn_rfc2141",
	"ip",
	"ip4_addr",
	"ip6_addr",
	"ip_addr",
	"ipv4",
	"ipv6",
	"cidr",
	"cidrv4",
	"cidrv6",
	"tcp_addr",
	"tcp4_addr",
	"tcp6_addr",
	"udp_addr",
	"udp4_addr",
	"udp6_addr",
	"unix_addr",
	"mac",
	"hostname",
	"hostname_port",
	"hostname_rfc1123",
	"fqdn",
	"jwt",
	"latitude",
	"longitude",
	"postcode_iso3166_alpha2",
	"postcode_iso3166_alpha2_field",
	"rgb",
	"rgba",
	"semver",
	"ssn",
	"timezone",
	"uuid",
	"uuid3",
	"uuid4",
	"uuid5",
	"uuid_rfc4122",
	"uuid3_rfc4122",
	"uuid4_rfc4122",
	"uuid5_rfc4122",
	"md5",
	"sha256",
	"sha384",
	"sha512",
	"ripemd128",
	"ripemd160",
	"tiger128",
	"tiger160",
	"tiger192",
	"issn",
	"luhn_checksum",
	"mongodb_connection_string",
	"cve",

	// Comparisons across fields
	"eqfield",
	"eqcsfield",
	"gtfield",
	"gtcsfield",
	"gtefield",
	"gtecsfield",
	"ltfield",
	"ltcsfield",
	"ltefield",
	"ltecsfield",
	"nefield",
	"necsfield",

	// Slices / maps / arrays
	"dive",
	"keys",
	"endkeys",
	"unique",

	// Special
	"-",
}

var knownValidateRulesSet map[string]bool

func init() {
	knownValidateRulesSet = make(map[string]bool, len(knownValidateRules))
	for _, r := range knownValidateRules {
		knownValidateRulesSet[r] = true
	}
}

// CheckValidateRules validates all rules in a validate tag value string.
// e.g. "required,min=2,max=50,email"
// Returns a list of human-readable issue strings.
func CheckValidateRules(rawValue string) []string {
	var issues []string

	// Handle pipe-separated rules (cross-field OR conditions)
	// e.g. "required_if=Field value|required_unless=Other value"
	parts := strings.Split(rawValue, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "-" {
			continue
		}

		// Rules can have parameters: min=2, oneof=red blue, datetime=2006-01-02
		ruleName := part
		if idx := strings.IndexAny(part, "="); idx >= 0 {
			ruleName = part[:idx]
		}

		// Also strip pipe-separated OR alternatives: required_if=X|required_unless=Y
		if idx := strings.Index(ruleName, "|"); idx >= 0 {
			ruleName = ruleName[:idx]
		}

		ruleName = strings.TrimSpace(ruleName)
		if ruleName == "" {
			continue
		}

		if !knownValidateRulesSet[ruleName] {
			threshold := 2
			if len(ruleName) <= 4 {
				threshold = 1
			}
			if suggestion, ok := closestMatch(ruleName, knownValidateRules, threshold); ok {
				issues = append(issues, fmt.Sprintf("unknown validate rule %q (did you mean %q?)", ruleName, suggestion))
			} else {
				issues = append(issues, fmt.Sprintf("unknown validate rule %q", ruleName))
			}
		}
	}

	return issues
}
