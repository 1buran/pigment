# pigment - tiny toy colorizer

[![codecov](https://codecov.io/gh/1buran/pigment/graph/badge.svg?token=3F7HTBT028)](https://codecov.io/gh/1buran/pigment)
[![Go Reference](https://pkg.go.dev/badge/github.com/1buran/pigment.svg)](https://pkg.go.dev/github.com/1buran/pigment)
[![goreportcard](https://goreportcard.com/badge/github.com/1buran/pigment)](https://goreportcard.com/report/github.com/1buran/pigment)

![Main demo](https://i.imgur.com/ojdMg7W.png)

## Introduction

This is very very tiny toy colorizer. It was created for personal needs: i needed something tiny
for colorize JSON data and make something like 'hard' formatting - not just re-indent
the fields, i needed remove some escape sequences, applied pretty printing
for some special json field values etc.

Pigment may be useful for internal / personal tools, scripts whatever when you don't want
add more heavy, but of course, more mature and featured, libs. However, it can offer
**more styling freedom** than others: custom fg/bg color, bold, underline, etc everything what
[termenv](https://github.com/muesli/termenv) supports.

## Features

Currently implemented:

- JSON lexer
- colorize JSON with [termenv](https://github.com/muesli/termenv)
- apply custom formatters

Both features based on node field name or field value, lexer token. That means you may code
any business logic of colorizing / formatting: literal or regex matching of passed json
field name or value or lexer token etc. You may find examples of usage in tests.

## Contributing

New lexers are welcome! Please take a look at the json lexer, as start point for writing
a new one. Please do not forget write tests along with the code.

## Usage

Add to project:

```
go get github.com/1buran/pigment
```

As i said this lib is pretty flexible, for example you may want colorize every even string
by green and every odd by the bright green color, this is funny,
it's looks like a python (snake), here is the output of test with this trick:

![python](https://i.imgur.com/8s7QMyd.png)

You may find the code in `pigment_test.go` file, but for clarity (how to use this module)
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

### How does it work

There is an interface `Pigmentizer`, this is core of this module: it defines two functions,
which will used for colorize input string of data.

The first one is `Style(k, v string, t Token) (bool, termenv.Style)` function. It uses
to check whether the part of input string should be styled. You may use literal matching
or matching by regexp or token matching for create all you needed conditions. If the part
is matched rules, then the function return `true` and `termenv.Style` which should be applied
to this part of processed data, otherwise it returns `false, termenv.Style{}`.

The second is `Format(k, v string, t Token) (bool, string)` function. Its purpose the same,
but in context of string content: you may use it for override some values of json or
correcting the formatting or what else you needed.

Here are some examples (full code in tests). A `Pigmentizer.Style` function:

```go
func (tpg testPigment) Style(k, v string, t Token) (bool, termenv.Style) {
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

```

it does:

- apply `nullColor` for every `null` value found
- apply `numColor` for every number found
- apply `trueColor` for every `true` boolean
- apply `falseColor` for every `false` boolean
- for every string found:
  - apply `msgColor` for value of json field `msg`
  - apply `lemonColor` for value of json fields: `date`, `ts`, `dateTime`, `timestamp`
  - apply `wrnColor` for field values matched these words: `Warn`, `Warning`
  - apply `infColor` for field values matched these words: `Info`, `Information`
  - apply `errColor` for field values matched these words: `Error`, `Failed`, `High`
  - apply `strColor` for other values (default color for all strings)

again, you may write any other logic e.g. colorize all date time meaning fields,
that matched by regex: `(?i).*(date|time|ts).*)` etc.

A `Pigmentizer.Format` function:

```go
func (tpg testPigment) Format(k, v string, t Token) (bool, string) {
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
```

it does:

- replace the original value of field `Priority` with `High`
- add prefix `AA` and sufix `BB` to the value of field `text`
- convert all `null` values to upper case
