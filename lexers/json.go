package pigment

import (
	"log"
	"strings"

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

// Styler detects whether the current key value is need to be styled.
// Here could be anything you wanted: a regual expression matching
// or literal matching or something else. You may compared by:
// k - key (JSON field)
// v - value (JSON field value)
// t - lexer token
//
// Take a look at the test for examples.
type Styler func(k, v string, t Token) (bool, termenv.Style)

// Formatter checks whether any custom format to key value should be applied.
// Here could be anything related to value changing: apply custom format
// or replace/add/delete something special words or by use regular expression.
// You may use the same params:
// k - key (JSON field)
// v - value (JSON field value)
// t - lexer token
//
// Take a look at the test for examples.
type Formatter func(k, v string, t Token) (bool, string)

// Tiny toy lexer which is purpose is getting from JSON meaningful parts: tokens,
// fields and values and bring the ability to highlight something or format/override.
func JSONLexer(s string, styler Styler, formatter Formatter) string {
	var (
		v               Token
		buf, out        strings.Builder
		key, lastString string
		lastChar        rune
	)

	// Apply custom format, style for key value based on current key, token value.
	// The passed fmtStyler is responsible for formatting & styling,
	// user may add there any wanted logic. This function is shortcut for invoke
	// formatter & styler functions.
	applyCustomFormatStyle := func(key, bufString string, token Token) string {
		customFormat, formatted := formatter(key, bufString, token)

		if customFormat {
			bufString = formatted
		}

		customStyle, style := styler(key, bufString, token)
		if customStyle {
			bufString = style.Styled(bufString)
		}
		return bufString
	}

	for _, c := range s {
		switch v {
		case INIT, VALUE, OBJECT, ARRAY:
			switch c {
			case '{':
				v = OBJECT
				out.WriteRune(c)
			case '"':
				v = STRING
				buf.WriteRune(c)
			case '[':
				v = ARRAY
				out.WriteRune(c)
			case 't':
				v = TRUE
				buf.WriteRune(c)
			case 'f':
				v = FALSE
				buf.WriteRune(c)
			case 'n':
				v = NULL
				buf.WriteRune(c)
			case 'e', 'E', '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				v = NUMBER
				buf.WriteRune(c)
			case ' ', '\t', '\n', ',', ':', '}', ']':
				if c == ':' {
					key = lastString
				} else if c == ',' || c == '}' || c == ']' {
					key = ""
				}
				out.WriteRune(c)
			}
		case STRING:
			buf.WriteRune(c)
			if c != '"' {
				lastChar = c
				continue
			} else {
				if lastChar == '\\' {
					lastChar = c
					continue
				}
				lastString = buf.String()
				bufString := applyCustomFormatStyle(key, lastString, STRING)
				out.WriteString(bufString)
				buf.Reset()
				v = VALUE
			}
		case NUMBER:
			switch c {
			case 'e', 'E', '-', '+', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				buf.WriteRune(c)
			case ',', '}', ']':
				bufString := applyCustomFormatStyle(key, buf.String(), NUMBER)
				out.WriteString(bufString)
				out.WriteRune(c)
				buf.Reset()
				v = VALUE
			}
		case NULL:
			switch c {
			case 'n', 'u', 'l':
				buf.WriteRune(c)
			case ',', '}', ']': // stop parsing
				bufString := buf.String()
				if bufString != "null" {
					log.Fatal("malformed json data, expected null, got: ", bufString)
				}
				bufString = applyCustomFormatStyle(key, bufString, NULL)
				out.WriteString(bufString)
				out.WriteRune(c)
				buf.Reset()
				v = VALUE
			case ' ', '\t', '\n':
				if buf.String() == "null" {
					continue
				} else {
					log.Fatalf("during parsing null unexpected whitespace found")
				}
			default:
				log.Fatalf("during parsing null unexpected rune found: %q", c)
			}
		case FALSE:
			switch c {
			case 'f', 'a', 'l', 's', 'e':
				buf.WriteRune(c)
			case ',', '}', ']':
				bufString := buf.String()
				if bufString != "false" {
					log.Fatal("malformed json data, expected false, got: ", bufString)
				}
				bufString = applyCustomFormatStyle(key, bufString, FALSE)
				out.WriteString(bufString)
				out.WriteRune(c)
				buf.Reset()
				v = VALUE
			case ' ', '\t', '\n':
				if buf.String() == "false" {
					continue
				} else {
					log.Fatalf("during parsing false unexpected whitespace found")
				}
			default:
				log.Fatalf("during parsing false unexpected rune found: %q", c)
			}
		case TRUE:
			switch c {
			case 't', 'r', 'u', 'e':
				buf.WriteRune(c)
			case ',', '}', ']':
				bufString := buf.String()
				if bufString != "true" {
					log.Fatal("malformed json data, expected true, got: ", bufString)
				}
				bufString = applyCustomFormatStyle(key, bufString, TRUE)
				out.WriteString(bufString)
				out.WriteRune(c)
				buf.Reset()
				v = VALUE
			case ' ', '\t', '\n':
				if buf.String() == "true" {
					continue
				} else {
					log.Fatalf("during parsing true unexpected whitespace found")
				}
			default:
				log.Fatalf("during parsing true unexpected rune found: %q", c)
			}
		}
	}
	return out.String()
}
