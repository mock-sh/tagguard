// Package analyzer implements the tagguard go/analysis analyzer.
// It can be used as a standalone linter or integrated into golangci-lint.
package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/mokshg/tagguard/pkg/config"
	"github.com/mokshg/tagguard/pkg/rules"
	"github.com/mokshg/tagguard/pkg/tags"
)

// Analyzer is the main go/analysis.Analyzer for tagguard.
// Register this with your analysis driver or golangci-lint plugin.
var Analyzer = &analysis.Analyzer{
	Name:     "tagguard",
	Doc:      "checks struct tags for unknown keys, typos, invalid validate rules, and naming inconsistencies",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Load config by searching upward from the first file's directory
	cfg := &config.Config{}
	if len(pass.Files) > 0 {
		dir := pass.Fset.File(pass.Files[0].Pos()).Name()
		// dir is the file path — get its directory
		if idx := strings.LastIndex(dir, "/"); idx >= 0 {
			dir = dir[:idx]
		}
		if loaded, err := config.Load(dir); err == nil {
			cfg = loaded
		}
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		structType := n.(*ast.StructType)

		// Collect all field tag entries for struct-wide consistency check
		var allFieldEntries []rules.FieldTagEntry
		// Map field tag position → field name (for reporting struct-wide issues)
		type fieldMeta struct {
			pos  token.Pos
			name string
		}
		var fieldMetas []fieldMeta

		for _, field := range structType.Fields.List {
			if field.Tag == nil {
				continue
			}

			raw := strings.Trim(field.Tag.Value, "`")
			parsedTags := tags.Parse(raw)

			// Determine field name for reporting
			fieldName := "_"
			if len(field.Names) > 0 {
				fieldName = field.Names[0].Name
			}

			// Rule 1: Unknown or typo'd tag key
			if !cfg.IsDisabled("unknown-key") {
				for _, tag := range parsedTags {
					if suggestion, known := rules.CheckTagKey(tag.Key, cfg.ExtraKnownKeys); !known {
						if suggestion != "" {
							pass.Reportf(field.Tag.Pos(),
								"unknown tag key %q (did you mean %q?)", tag.Key, suggestion)
						} else {
							pass.Reportf(field.Tag.Pos(),
								"unknown tag key %q", tag.Key)
						}
					}
				}
			}

			// Rule 2: Invalid validate / binding rules
			if !cfg.IsDisabled("validate-rules") {
				for _, tag := range parsedTags {
					if tag.Key == "validate" || tag.Key == "binding" {
						for _, issue := range rules.CheckValidateRules(tag.Raw) {
							pass.Reportf(field.Tag.Pos(),
								"in %q tag: %s", tag.Key, issue)
						}
					}
				}
			}

			if !cfg.IsDisabled("naming-consistency") {
				// Rule 3: Per-field naming inconsistency (e.g. json camelCase vs db snake_case)
				for _, issue := range rules.CheckNamingConsistency(parsedTags) {
					pass.Reportf(field.Tag.Pos(), issue)
				}

				// Rule 4: Enforced project-wide naming style from config
				if cfg.NamingStyle != "" {
					for _, tag := range parsedTags {
						if issue := rules.CheckEnforcedNamingStyle(tag.Key, tag.Value, cfg.NamingStyle); issue != "" {
							pass.Reportf(field.Tag.Pos(), issue)
						}
					}
				}

				// Collect for struct-wide check (Rule 5)
				entries := rules.CollectFieldStyles(fieldName, parsedTags)
				allFieldEntries = append(allFieldEntries, entries...)
				if len(entries) > 0 {
					fieldMetas = append(fieldMetas, fieldMeta{pos: field.Tag.Pos(), name: fieldName})
				}
			}
		}

		// Rule 5: Struct-wide naming consistency
		if !cfg.IsDisabled("naming-consistency") && len(allFieldEntries) > 0 {
			structIssues := rules.CheckStructNamingConsistency(allFieldEntries)

			// Report issues on the right field's tag position
			for _, field := range structType.Fields.List {
				if field.Tag == nil {
					continue
				}
				fieldName := "_"
				if len(field.Names) > 0 {
					fieldName = field.Names[0].Name
				}
				for _, issue := range structIssues[fieldName] {
					pass.Reportf(field.Tag.Pos(), issue)
				}
			}
		}
	})

	return nil, nil
}
