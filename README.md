# pigment - tiny toy colorizer

[![codecov](https://codecov.io/gh/1buran/pigment/graph/badge.svg?token=3F7HTBT028)](https://codecov.io/gh/1buran/pigment)
[![goreportcard](https://goreportcard.com/badge/github.com/1buran/pigment)](https://goreportcard.com/report/github.com/1buran/pigment)

![Main demo](https://i.imgur.com/ojdMg7W.png)

## Introduction

This is very very tiny toy colorizer. It was created for personal needs: i needed something tiny
for colorize JSON data and make something like 'hard' formatting - not just re-indent
the fields, i needed remove some escape sequences, applied pretty printing
for some special json field values etc.

Pigment may be useful for internal / personal tools, scripts whatever when you don't want
add more heavy, but of course, more mature and featured, libs.

## Features

Currently implemented:

- colorize JSON with [termenv](https://github.com/muesli/termenv)
- apply custom formatters

Both features based on node field name or field value, lexer token. That means you may code
any business logic of colorizing / formatting: literal or regex matrching of passed json
field name or value or lexer token etc. You may find examples of usage in tests.

## Contributing

New lexers are welcome! Please take a look at the json lexer, as start point for writing
a new one. Please do not forget write tests along with the code.

## Usage

As i said this lib is pretty flexible, for example you may want colorize every even string
by green and every odd by the bright green color, this is funny,
it's looks like a python (snake), here is the output of test with this trick:

![python](https://i.imgur.com/8s7QMyd.png)

You may found the code in `pigment_test.go` file, but for clarity (how to use this module)
I am attaching the code of a simple script that does the same thing:

```go
package main

import (
	"fmt"
	"os"

	"github.com/1buran/pigment"
	"github.com/1buran/pigment/lexers"
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
func (tpg testPigment) Style(k, v string, t lexers.Token) (bool, termenv.Style) {
	var blinkColor termenv.Color

	tpg.Inc()

	switch t {
	case lexers.STRING:
		if tpg.IsOdd() {
			blinkColor = strColor2
		} else {
			blinkColor = strColor1
		}
		return true, termenv.Style{}.Foreground(blinkColor)
	}
	return false, termenv.Style{}
}

func (tpg testPigment) Format(k, v string, t lexers.Token) (bool, string) {
	return false, v
}

func main() {
	var n int
	var tpg testPigment = testPigment{&n}

	s := `{ "Priority": "low", "text": "Hello\t world!", "msg": "Wake up, Neo...",
            "date": "2024-07-19", "time": "16:55:03", "ts": "2024-07-19T16:11:00+00:00",
            "level": "Warn", "status": "Failed", "isValid": "none", "isAlert": "red",
            "count": "10230", "notifyLevel": "Information", "errorMsg": "user not found",
            "array": [ "string", "yes", "true", "false", "null"] }`

	r := pigment.Pigmentize(pigment.JSON, tpg, s)
	fmt.Println(r)
}
```
