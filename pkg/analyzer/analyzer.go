// Package analyzer implements the tagguard go/analysis analyzer.
// It can be used as a standalone linter or integrated into golangci-lint.
package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

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
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		structType := n.(*ast.StructType)

		for _, field := range structType.Fields.List {
			if field.Tag == nil {
				continue
			}

			// Strip surrounding backticks
			raw := strings.Trim(field.Tag.Value, "`")
			parsedTags := tags.Parse(raw)

			for _, tag := range parsedTags {
				// Rule 1: Unknown or typo'd tag key
				if suggestion, known := rules.CheckTagKey(tag.Key); !known {
					if suggestion != "" {
						pass.Reportf(field.Tag.Pos(),
							"unknown tag key %q (did you mean %q?)", tag.Key, suggestion)
					} else {
						pass.Reportf(field.Tag.Pos(),
							"unknown tag key %q", tag.Key)
					}
				}

				// Rule 2: Invalid validate / binding rules
				if tag.Key == "validate" || tag.Key == "binding" {
					for _, issue := range rules.CheckValidateRules(tag.Raw) {
						pass.Reportf(field.Tag.Pos(),
							"in %q tag: %s", tag.Key, issue)
					}
				}
			}

			// Rule 3: Naming inconsistency across serialization tags
			for _, issue := range rules.CheckNamingConsistency(parsedTags) {
				pass.Reportf(field.Tag.Pos(), issue)
			}
		}
	})

	return nil, nil
}
