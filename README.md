# tagguard

A Go linter for struct tags. Catches the bugs that the compiler misses.

```
$ tagguard ./...

user.go:12:14: unknown tag key "jsno" (did you mean "json"?)
user.go:18:14: in "validate" tag: unknown validate rule "requred" (did you mean "required"?)
user.go:24:16: inconsistent naming: db tag uses snake_case ("user_id") but other tags use camelCase
```

---

## Why

Go struct tags are plain strings — the compiler has no idea what's inside them.
A typo silently fails at runtime:

```go
type User struct {
    Name  string `jsno:"name"`          // ← compiles fine, JSON silently breaks
    Email string `validate:"requred"`   // ← compiles fine, validation never runs
}
```

`tagguard` catches these before they ship.

---

## What It Detects

### 1. Unknown or Typo'd Tag Keys

Detects tag keys that don't match any known Go ecosystem tag, with typo suggestions using [Damerau-Levenshtein](https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance) distance.

```go
Name string `jsno:"name"`    // unknown tag key "jsno" (did you mean "json"?)
Name string `ymal:"name"`    // unknown tag key "ymal" (did you mean "yaml"?)
Name string `foobar:"name"`  // unknown tag key "foobar"
```

Known keys include: `json`, `yaml`, `xml`, `toml`, `bson`, `db`, `gorm`, `validate`, `binding`, `form`, `query`, `uri`, `env`, `mapstructure`, and [30+ more](docs/known-keys.md).

### 2. Invalid Validate Rules

Detects unknown or typo'd rules in `validate` and `binding` tags (from [go-playground/validator](https://github.com/go-playground/validator)).

```go
Age   int    `validate:"requred,min=0"`       // unknown validate rule "requred" (did you mean "required"?)
Email string `validate:"emai,omitempty"`       // unknown validate rule "emai" (did you mean "email"?)
Score int    `validate:"superstrongvalidation"` // unknown validate rule "superstrongvalidation"
```

### 3. Naming Inconsistencies

Detects when serialization tags on the same field use different naming conventions.

```go
// inconsistent naming: db tag uses snake_case ("user_id") but other tags use camelCase
UserID string `json:"userId" db:"user_id"`
```

Checks these tags for consistency: `json`, `yaml`, `xml`, `toml`, `bson`, `db`, `form`, `query`, `uri`, `param`.

---

## Installation

```bash
go install github.com/mokshg/tagguard/cmd/tagguard@latest
```

Or build from source:

```bash
git clone https://github.com/mokshg/tagguard
cd tagguard
go build ./cmd/tagguard
```

---

## Usage

```bash
# Lint a single package
tagguard ./...

# Lint a specific directory
tagguard ./internal/models/...

# Use as part of a Go analysis pipeline
tagguard -json ./... | jq .
```

tagguard uses the standard `go/analysis` framework, so it accepts all the same flags as `go vet`.

---

## golangci-lint Integration

Add tagguard as a custom linter in your `.golangci.yml`:

```yaml
linters-settings:
  custom:
    tagguard:
      path: ./bin/tagguard
      description: Struct tag linter
      original-url: github.com/mokshg/tagguard
```

---

## Configuration (Planned)

Future versions will support a `.tagguard.yaml` config file:

```yaml
# Add your project-specific tag keys as known/allowed
extra-known-keys:
  - mytag
  - internal

# Disable specific rules
disable:
  - naming-consistency

# Enforce a specific naming style across all tags
naming-style: snake_case
```

---

## How It Works

tagguard uses Go's `go/analysis` framework to walk the AST of your Go source files.
For each struct field with a tag:

1. **Parses** the raw tag string into individual key-value pairs
2. **Checks each key** against a list of 40+ known ecosystem tag keys, using fuzzy matching for typo suggestions
3. **Validates validate/binding rules** against the full go-playground/validator rule set
4. **Compares naming styles** across serialization tags on the same field

---

## Limitations

- Custom tag keys (e.g. internal company tags) will be flagged as unknown until added to the known list or config
- Validate rules from custom validators won't be recognized
- Naming consistency only checks fields individually — it does not enforce a single style across the whole struct

---

## Contributing

Pull requests welcome. See [CONTRIBUTING.md](docs/contributing.md) for setup instructions.

Common areas to contribute:
- Add more known tag keys
- Add more known validate rules
- Implement the config file
- Add golangci-lint plugin support
- Add auto-fix support (`-fix` flag)

---

## License

MIT
