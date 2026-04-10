// tagguard is a linter for Go struct tags.
// It detects:
//   - Unknown or typo'd tag keys (e.g. `jsno:"name"` → did you mean `json`?)
//   - Invalid validate rules (e.g. `validate:"requred"` → did you mean `required`?)
//   - Naming inconsistencies across serialization tags (e.g. json uses camelCase but db uses snake_case)
//
// Usage:
//
//	tagguard ./...
//	tagguard ./path/to/package
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/mokshg/tagguard/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
