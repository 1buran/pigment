# pigment - un pequeño colorizador de juguete

[![codecov](https://codecov.io/gh/1buran/pigment/graph/badge.svg?token=3F7HTBT028)](https://codecov.io/gh/1buran/pigment)
[![Go Reference](https://pkg.go.dev/badge/github.com/1buran/pigment.svg)](https://pkg.go.dev/github.com/1buran/pigment)
[![goreportcard](https://goreportcard.com/badge/github.com/1buran/pigment)](https://goreportcard.com/report/github.com/1buran/pigment)

![Main demo](https://i.imgur.com/ojdMg7W.png)

## Introducción

Este es un colorizador de juguete muy, muy pequeño. Fue creado para necesidades personales: necesitaba algo ligero
para colorear datos JSON y realizar un formateo "estricto": no solo re-indentar
los campos, sino también eliminar algunas secuencias de escape, aplicar un formato legible (pretty printing)
para algunos valores de campo JSON especiales, etc.

Pigment puede ser útil para herramientas internas / personales, scripts o cualquier situación en la que no quieras
añadir librerías más pesadas, aunque estas últimas sean, por supuesto, más maduras y completas. Sin embargo, puede ofrecer
**más libertad de estilo** que otras: color de primer plano/fondo personalizado, negrita, subrayado, etc., todo lo que
[termenv](https://github.com/muesli/termenv) soporte.

## Características

Implementado actualmente:

- Lexer de JSON
- Coloreado de JSON con [termenv](https://github.com/muesli/termenv)
- Aplicación de formateadores personalizados

Ambas características se basan en el nombre del campo del nodo, el valor del campo o el token del lexer. Esto significa que puedes programar
cualquier lógica de negocio de coloreado / formateado: coincidencia literal o por regex del nombre del campo JSON pasado, el valor o el token del lexer, etc. Puedes encontrar ejemplos de uso en las pruebas.

## Contribución

¡Los nuevos lexers son bienvenidos! Por favor, echa un vistazo al lexer de json como punto de partida para escribir
uno nuevo. No olvides escribir pruebas junto con el código.

## Uso

Añadir al proyecto:

```
go get github.com/1buran/pigment
```

Como mencioné, esta librería es bastante flexible; por ejemplo, podrías querer colorear cada cadena par
en verde y cada una impar en color verde brillante; esto es divertido,
parece una pitón (serpiente), aquí tienes el resultado de una prueba con este truco:

![python](https://i.imgur.com/8s7QMyd.png)

Puedes encontrar el código en el archivo `pigment_test.go`, pero para mayor claridad (sobre cómo usar este módulo)
adjunto el código de un script sencillo que hace lo mismo:

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
	i *int // contador de llamadas, emulación de parpadeo
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

### Cómo funciona

Existe una interfaz `Pigmentizer`, que es el núcleo de este módulo: define dos funciones que se utilizarán para colorear la cadena de datos de entrada.

La primera es la función `Style(k, v string, t Token) (bool, termenv.Style)`. Se utiliza para verificar si una parte de la cadena de entrada debe ser estilizada. Puedes usar coincidencias literales, por expresiones regulares o coincidencia de tokens para crear todas las condiciones que necesites. Si la parte coincide con las reglas, la función devuelve `true` y un `termenv.Style` que debe aplicarse a esa parte de los datos procesados; de lo contrario, devuelve `false, termenv.Style{}`.

La segunda es la función `Format(k, v string, t Token) (bool, string)`. Su propósito es el mismo, pero en el contexto del contenido de la cadena: puedes usarla para sobrescribir algunos valores de JSON, corregir el formato o cualquier otra cosa que necesites.

Aquí hay algunos ejemplos (código completo en las pruebas). Una función `Pigmentizer.Style`:

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
		switch k { // resaltar campos JSON especiales (coincidencia literal)
		case `"msg"`:
			return true, termenv.Style{}.Foreground(msgColor)
		case `"date"`, `"ts"`, `"dateTime"`, `"timestamp"`:
			return true, termenv.Style{}.Foreground(lemonColor)
		}

		switch v { // resaltar valores JSON especiales (coincidencia literal)
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

Esto hace lo siguiente:

- aplica `nullColor` para cada valor `null` encontrado
- aplica `numColor` para cada número encontrado
- aplica `trueColor` para cada booleano `true`
- aplica `falseColor` para cada booleano `false`
- para cada cadena encontrada:
  - aplica `msgColor` para el valor del campo json `msg`
  - aplica `lemonColor` para los valores de los campos json: `date`, `ts`, `dateTime`, `timestamp`
  - aplica `wrnColor` para los valores de campo que coincidan con estas palabras: `Warn`, `Warning`
  - aplica `infColor` para los valores de campo que coincidan con estas palabras: `Info`, `Information`
  - aplica `errColor` para los valores de campo que coincidan con estas palabras: `Error`, `Failed`, `High`
  - aplica `strColor` para otros valores (color predeterminado para todas las cadenas)

Nuevamente, puedes escribir cualquier otra lógica, por ejemplo, colorear todos los campos que signifiquen fecha y hora que coincidan con una regex: `(?i).*(date|time|ts).*)` etc.

Una función `Pigmentizer.Format`:

```go
func (tpg testPigment) Format(k, v string, t Token) (bool, string) {
	switch k {
	case `"Priority"`:
		return true, `"High"` // sobrescribir valor del campo
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

Esto hace lo siguiente:

- reemplaza el valor original del campo `Priority` por `High`
- añade el prefijo `AA` y el sufijo `BB` al valor del campo `text`
- convierte todos los valores `null` a mayúsculas
