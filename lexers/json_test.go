package pigment

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/muesli/termenv"
)

func TestJSONLexer(t *testing.T) {
	t.Parallel()

	var (
		output = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))

		strColor   = output.Color("75")
		numColor   = output.Color("159")
		nullColor  = output.Color("87")
		trueColor  = output.Color("47")
		falseColor = output.Color("212")

		wrnColor = output.Color("226")
		infColor = output.Color("121")
		errColor = output.Color("212")

		msgColor   = output.Color("183")
		lemonColor = output.Color("192")

		styler Styler = func(k, v string, t Token) (bool, termenv.Style) {
			switch t {
			case NULL:
				return true, termenv.Style{}.Foreground(nullColor).Bold()
			case NUMBER:
				return true, termenv.Style{}.Foreground(numColor)
			case TRUE:
				return true, termenv.Style{}.Foreground(trueColor)
			case FALSE:
				return true, termenv.Style{}.Foreground(falseColor).Bold()
			case STRING:
				switch k { // highlight special JSON fields (literal matching)
				case `"msg"`:
					return true, termenv.Style{}.Foreground(msgColor)
				case `"date"`, `"ts"`, `"dateTime"`, `"timestamp"`:
					return true, termenv.Style{}.Foreground(lemonColor)
				}

				switch v { // highlight special JSON values (literal matching)
				case `"Warn"`, `"Warning"`:
					return true, termenv.Style{}.Foreground(wrnColor)
				case `"Info"`, `"Information"`:
					return true, termenv.Style{}.Foreground(infColor)
				case `"Error"`, `"Failed"`, `"High"`:
					return true, termenv.Style{}.Foreground(errColor)
				default:
					return true, termenv.Style{}.Foreground(strColor)
				}
			}
			return false, termenv.Style{}
		}

		formatter Formatter = func(k, v string, t Token) (bool, string) {
			switch k {
			case `"Priority"`:
				return true, `"High"` // override field value
			case `"text"`:
				r, _ := strconv.Unquote(v)
				return true, `"AA` + r + `BB"`
			}

			switch t {
			case NULL:
				return true, "NULL"
			}
			return false, v
		}
	)

	s := `{ "Priority": "low", "text": "Hello\t world!", "msg": "Wake up, Neo...",
            "date": "2024-07-19", "time": "16:55:03", "ts": "2024-07-19T16:11:00+00:00",
            "level": "Warn", "status": "Failed", "isValid": false, "isAlert": true,
            "count": 10230, "notifyLevel": "Information", "errorMsg": null,
            "array": [ "string", 12, true, false, null] }`

	r := JSONLexer(s, styler, formatter)

	t.Log(r)

	t.Run("strings", func(t *testing.T) {
		if !strings.Contains(r, "\x1b[38;5;75m\"notifyLevel\"\x1b[0m") {
			t.Error("strings were not pigmented")
		}
	})
	t.Run("null", func(t *testing.T) {
		if !strings.Contains(r, "\x1b[38;5;87;1mNULL\x1b[0m") {
			t.Error("null values were not pigmented")
		}
	})
	t.Run("false", func(t *testing.T) {
		if !strings.Contains(r, "\x1b[38;5;212;1mfalse\x1b[0m") {
			t.Error("false values were not pigmented")
		}
	})

	t.Run("true", func(t *testing.T) {
		if !strings.Contains(r, "\x1b[38;5;47mtrue\x1b[0m") {
			t.Error("true values were not pigmented")
		}
	})

	t.Run("numbers", func(t *testing.T) {
		if !strings.Contains(r, "\x1b[38;5;159m12\x1b[0m") {
			t.Error("numbers were not pigmented")
		}
	})

	t.Run("formatter", func(t *testing.T) {
		if strings.Contains(r, "\x1b[38;5;87;1mnull\x1b[0m") {
			t.Error("formatter was not applied: null was not replaced with NULL")
		}
		if !strings.Contains(r, "AAHello\t world!BB") {
			t.Error("formatter was not applied by key:", `"AAHello\t world!BB"`, "not found")
		}
	})
}
