package lexers

import (
	"github.com/muesli/termenv"
)

const (
	INIT = iota
	WS
	OBJECT
	VALUE
	ARRAY
	STRING
	NUMBER
	TRUE
	FALSE
	NULL
)

type Token int

func (v Token) String() string {
	return []string{"Init", "Whitespace", "Object", "Value", "Array", "String", "Number",
		"True", "False", "Null"}[v]
}

type Pigmentizer interface {
	// Styler detects whether the current key value is need to be styled.
	// Here could be anything you wanted: a regual expression matching
	// or literal matching or something else. You may compared by:
	// k - key (JSON field)
	// v - value (JSON field value)
	// t - lexer token
	//
	// Take a look at the test for examples.
	Style(k, v string, t Token) (bool, termenv.Style)

	// Formatter checks whether any custom format to key value should be applied.
	// Here could be anything related to value changing: apply custom format
	// or replace/add/delete something special words or by use regular expression.
	// You may use the same params:
	// k - key (JSON field)
	// v - value (JSON field value)
	// t - lexer token
	//
	// Take a look at the test for examples.
	Format(k, v string, t Token) (bool, string)
}
