package pigment

import (
	"os"
	"strings"
	"testing"

	lex "github.com/1buran/pigment/lexers"
	"github.com/muesli/termenv"
)

var (
	output = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.TrueColor))

	strColor1 = output.Color("34")
	strColor2 = output.Color("122")
)

type testPigment struct {
	i *int // call counter, blinking emulation
}

func (tpg testPigment) Inc()        { *tpg.i++ }
func (tpg testPigment) IsOdd() bool { return *tpg.i%2 == 0 }
func (tpg testPigment) Style(k, v string, t lex.Token) (bool, termenv.Style) {
	var blinkColor termenv.Color

	tpg.Inc()

	switch t {
	case lex.STRING:
		if tpg.IsOdd() {
			blinkColor = strColor2
		} else {
			blinkColor = strColor1
		}
		return true, termenv.Style{}.Foreground(blinkColor)
	}
	return false, termenv.Style{}
}

func (tpg testPigment) Format(k, v string, t lex.Token) (bool, string) {
	return false, v
}

func TestPigmentize(t *testing.T) {
	t.Parallel()

	var n int
	var tpg testPigment = testPigment{&n}

	s := `{ "Priority": "low", "text": "Hello\t world!", "msg": "Wake up, Neo...",
            "date": "2024-07-19", "time": "16:55:03", "ts": "2024-07-19T16:11:00+00:00",
            "level": "Warn", "status": "Failed", "isValid": "none", "isAlert": "red",
            "count": "10230", "notifyLevel": "Information", "errorMsg": "user not found",
            "array": [ "string", "yes", "true", "false", "null"] }`

	t.Run("JSON", func(t *testing.T) {
		r := Pigmentize(JSON, tpg, s)

		t.Log(r)

		if !strings.Contains(r, "\x1b[38;5;34m\"Priority\"\x1b[0m") {
			t.Error("Priority has wrong color or missed")
		}

		if !strings.Contains(r, "\x1b[38;5;122m\"low\"\x1b[0m") {
			t.Error("Low has wrong color or missed")
		}

		if !strings.Contains(r, "\x1b[38;5;34m\"text\"\x1b[0m") {
			t.Error("Text has wrong color or missed")
		}

		if !strings.Contains(r, "\x1b[38;5;122m\"none\"\x1b[0m") {
			t.Error("None has wrong color or missed")
		}
	})

	t.Run("unknown", func(t *testing.T) {
		r := Pigmentize(100, tpg, s)

		t.Log(r)

		if s != r {
			t.Error("Original string should not be modified")
		}
	})
}
