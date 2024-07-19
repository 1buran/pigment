package pigment

import (
	lex "github.com/1buran/pigment/lexers"
)

const (
	JSON = iota
)

type SyntaxType int

func Pigmentize(t SyntaxType, pg lex.Pigmentizer, s string) string {
	switch t {
	case JSON:
		return lex.JSONLexer(s, pg)
	}
	return s
}
