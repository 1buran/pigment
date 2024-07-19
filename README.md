# pigment - tiny toy colorizer

[![codecov](https://codecov.io/gh/1buran/pigment/graph/badge.svg?token=3F7HTBT028)](https://codecov.io/gh/1buran/pigment)
[![goreportcard](https://goreportcard.com/badge/github.com/1buran/pigment)](https://goreportcard.com/report/github.com/1buran/pigment)

![Main demo](https://i.imgur.com/ojdMg7W.png)

## Introduction

This is very very tiny toy colorizer. It was created for personal needs: i needed something tiny
for colorize JSON data and make something like 'hard' formatting - not just re-indent
the fields, i needed remove some escape sequences, applied pretty printing
for some special json field values etc.

## Features

Currently implemented:

- colorize JSON with [termenv](https://github.com/muesli/termenv)
- apply custom formatters

Both features based on node field name or field value, lexer token. That means you may code
any business logic of colorizing / formatting: literal or regex matrching of passed json
field name or value or lexer token etc. You may find examples of usage in tests.
