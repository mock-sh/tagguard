package rules

// KnownTagKeys lists all well-known struct tag keys used across the Go ecosystem.
// Tags not in this list trigger an "unknown tag key" warning (with typo suggestion if close).
var KnownTagKeys = []string{
	// Standard library
	"json",
	"xml",

	// Popular serialization
	"yaml",
	"toml",
	"bson",       // MongoDB
	"msgpack",    // MessagePack
	"cbor",       // CBOR encoding
	"protobuf",   // Protocol Buffers
	"avro",       // Apache Avro

	// Database / ORM
	"db",         // sqlx
	"gorm",       // GORM
	"sql",
	"pg",         // go-pg
	"bigquery",   // Google BigQuery

	// Web frameworks
	"form",       // form binding (gin, echo)
	"query",      // URL query params
	"uri",        // URI params (gin)
	"header",     // HTTP header binding
	"cookie",     // cookie binding
	"binding",    // gin binding (alias for validate)
	"param",      // path param (fiber)
	"path",       // path param

	// Validation
	"validate",   // go-playground/validator
	"valid",      // govalidator

	// Config / env
	"env",        // envconfig, godotenv
	"envconfig",
	"mapstructure", // viper / mapstructure
	"flag",       // flag package

	// Documentation / codegen
	"example",
	"description",
	"default",
	"swaggertype",
	"extensions",

	// Misc
	"redis",
	"dynamodbav", // AWS DynamoDB
	"firestore",  // Google Firestore
	"spanner",    // Google Cloud Spanner
	"datastore",  // Google Cloud Datastore
	"csv",
	"xlsx",
}

// knownKeysSet is a fast lookup set built from KnownTagKeys.
var knownKeysSet map[string]bool

func init() {
	knownKeysSet = make(map[string]bool, len(KnownTagKeys))
	for _, k := range KnownTagKeys {
		knownKeysSet[k] = true
	}
}

// CheckTagKey returns whether the key is known.
// extraKeys is a list of additional keys to treat as known (from config).
// If unknown but close to a known key, returns a suggestion.
//
// Returns:
//   - (suggestion, true)  → key is known, no issue
//   - ("", false)         → key is unknown, no close match
//   - (suggestion, false) → key is unknown, here's what you probably meant
func CheckTagKey(key string, extraKeys []string) (suggestion string, known bool) {
	if knownKeysSet[key] {
		return "", true
	}
	for _, k := range extraKeys {
		if k == key {
			return "", true
		}
	}

	// Typo threshold: allow 1 edit for very short keys, 2 for normal keys.
	// We use Damerau-Levenshtein so transpositions (jsno→json) count as 1 edit.
	threshold := 2
	if len(key) <= 2 {
		threshold = 1
	}

	// Build the full candidate list including extra keys
	allKeys := KnownTagKeys
	if len(extraKeys) > 0 {
		allKeys = append(allKeys, extraKeys...)
	}

	if match, ok := closestMatch(key, allKeys, threshold); ok {
		return match, false
	}

	return "", false
}
