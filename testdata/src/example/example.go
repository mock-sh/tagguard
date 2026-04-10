// Package example contains structs used to manually test tagguard.
// Run: tagguard ./testdata/src/example/
package example

// GoodStruct should produce zero warnings.
type GoodStruct struct {
	ID    int    `json:"id" db:"id" validate:"required,gt=0"`
	Name  string `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" db:"email" validate:"required,email"`
	Age   int    `json:"age,omitempty" db:"age" validate:"omitempty,gte=0,lte=150"`
}

// TypoInKey should warn: `jsno` is not a known tag key (did you mean `json`?)
type TypoInKey struct {
	Name string `jsno:"name"` // want `unknown tag key "jsno"`
}

// TypoInValidateRule should warn: `requred` is not a known validate rule.
type TypoInValidateRule struct {
	Name string `json:"name" validate:"requred,min=2"` // want `unknown validate rule "requred"`
}

// InconsistentNaming should warn: json uses camelCase, db uses snake_case.
type InconsistentNaming struct {
	UserID string `json:"userId" db:"user_id"` // want `inconsistent naming`
}

// UnknownKey should warn about a completely unknown tag key.
type UnknownKey struct {
	Data string `foobar:"data"` // want `unknown tag key "foobar"`
}

// MultipleIssues demonstrates multiple issues on one struct.
type MultipleIssues struct {
	Name  string `jsno:"name" validate:"requred"`        // typo in key + typo in rule
	Email string `json:"email" validate:"emai,omitempty"` // typo in rule
}
