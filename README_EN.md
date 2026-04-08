# go-propisyu · Russian Number-to-Words for Go

**English version · [Русская версия](README.md)**

[![CI](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml/badge.svg)](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/rekurt/go-propisyu.svg)](https://pkg.go.dev/github.com/rekurt/go-propisyu)
[![Go Report Card](https://goreportcard.com/badge/github.com/rekurt/go-propisyu)](https://goreportcard.com/report/github.com/rekurt/go-propisyu)
[![codecov](https://codecov.io/gh/rekurt/go-propisyu/branch/master/graph/badge.svg)](https://codecov.io/gh/rekurt/go-propisyu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/rekurt/go-propisyu)](go.mod)

`go-propisyu` is a Go library that converts integers and decimals into Russian words with correct grammatical gender and noun declension. Perfect for invoices, fiscal receipts, accounting documents, bank statements, voice assistants (TTS), chatbots, document generation, and any service that needs to spell out numbers in fluent Russian.

```go
propisyu.IntToWords(42)                              // "сорок два"
propisyu.Decline(5, "рубль", "рубля", "рублей")      // "рублей"
propisyu.DecimalToWords("1234.56")                    // "одна тысяча двести тридцать четыре целых пятьдесят шесть сотых"
```

## Use Cases

| Domain | Example |
|--------|---------|
| Fintech & Banking | Amount in words on payment orders and statements |
| Accounting | Invoice, act, and waybill generation |
| Fiscal Receipts | POS / OFD — amount in words (Russian 54-FZ) |
| Voice Assistants | TTS pronunciation of amounts and quantities |
| Chatbots | Natural-language responses with amounts |
| Document Generation | Contract, power-of-attorney, and act templates |

## Highlights

- Handles numbers up to duodecillions (10³⁹)
- Supports masculine, feminine, and neuter grammatical genders
- `Decline` helper for automatic Russian noun declension
- Decimal support via plain strings or `decimal.Decimal`
- Zero external dependencies for core functions
- High test coverage, CI/CD, linter, semantic versioning

## Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Examples](#usage-examples)
- [API](#api)
- [Why go-propisyu](#why-go-propisyu)
- [Limitations](#limitations)
- [Contributing](#contributing)
- [License](#license)

## Installation

```bash
go get github.com/rekurt/go-propisyu
```

For `decimal.Decimal` support:
```bash
go get github.com/shopspring/decimal
```

## Public Functions

| Function | Description |
|----------|-------------|
| `IntToWords(n int) string` | Converts an integer to words (masculine gender) |
| `IntToWordsGender(n int, gender Gender) string` | Converts an integer to words with specified gender |
| `DecimalToWords(decimalStr string) (string, error)` | Converts a decimal number from string to words |
| `DecimalValueToWords(d decimal.Decimal) (string, error)` | Converts a `decimal.Decimal` value to words |
| `Decline(n int, one, two, five string) string` | Chooses the correct noun declension form |

### Gender Constants

```go
const (
    GenderMasculine Gender = 1  // Masculine: "один", "два"
    GenderFeminine  Gender = 2  // Feminine: "одна", "две"
    GenderNeuter    Gender = 3  // Neuter: "одно", "два"
)
```

### Errors

- `ErrNumberTooLarge` - number is too large to convert (doesn't fit in `int`)

## Quick Start

### Integer Numbers

```go
package main

import (
	"fmt"

	"github.com/rekurt/go-propisyu"
)

func main() {
	// Basic conversion (masculine gender by default)
	fmt.Println(propisyu.IntToWords(321))
	// триста двадцать один

	// Conversion with specified gender
	fmt.Println(propisyu.IntToWordsGender(2, propisyu.GenderFeminine))
	// две

	fmt.Println(propisyu.IntToWordsGender(2, propisyu.GenderMasculine))
	// два

	// Automatic noun declension
	fmt.Println(propisyu.Decline(1, "рубль", "рубля", "рублей"))   // рубль
	fmt.Println(propisyu.Decline(2, "рубль", "рубля", "рублей"))   // рубля
	fmt.Println(propisyu.Decline(5, "рубль", "рубля", "рублей"))   // рублей
	fmt.Println(propisyu.Decline(21, "рубль", "рубля", "рублей"))  // рубль
}
```

### Decimal Numbers

#### Option 1: Using String

```go
package main

import (
	"fmt"
	"log"

	"github.com/rekurt/go-propisyu"
)

func main() {
	result, err := propisyu.DecimalToWords("123.45")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// сто двадцать три целых сорок пять сотых
}
```

#### Option 2: Using decimal.Decimal

```go
package main

import (
	"fmt"
	"log"

	"github.com/rekurt/go-propisyu"
	"github.com/shopspring/decimal"
)

func main() {
	d := decimal.NewFromFloat(123.45)
	result, err := propisyu.DecimalValueToWords(d)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// сто двадцать три целых сорок пять сотых
}
```

## Usage Examples

### Amount in Words for Receipt

```go
package main

import (
	"fmt"

	"github.com/rekurt/go-propisyu"
)

func main() {
	amount := 1234
	rubles := propisyu.IntToWords(amount)
	rublesDecl := propisyu.Decline(amount, "рубль", "рубля", "рублей")

	fmt.Printf("%s %s 00 копеек", rubles, rublesDecl)
	// одна тысяча двести тридцать четыре рубля 00 копеек
}
```

### Large Numbers

```go
package main

import (
	"fmt"
	"log"

	"github.com/rekurt/go-propisyu"
)

func main() {
	result, err := propisyu.DecimalToWords("6453345242432.42")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// шесть триллионов четыреста пятьдесят три миллиарда
	// триста сорок пять миллионов двести сорок две тысячи четыреста тридцать два
	// целых сорок две сотых
}
```

### Declension with Different Nouns

```go
package main

import (
	"fmt"

	"github.com/rekurt/go-propisyu"
)

func main() {
	count := 5

	// Currencies
	fmt.Println(count, propisyu.Decline(count, "доллар", "доллара", "долларов"))
	// 5 долларов

	// Units of measurement
	fmt.Println(count, propisyu.Decline(count, "метр", "метра", "метров"))
	// 5 метров

	// Items
	fmt.Println(count, propisyu.Decline(count, "товар", "товара", "товаров"))
	// 5 товаров

	// For 21
	count = 21
	fmt.Println(count, propisyu.Decline(count, "день", "дня", "дней"))
	// 21 день
}
```

## API

### Integer Functions

#### `IntToWords(n int) string`
Converts an integer to words (masculine gender by default).

```go
propisyu.IntToWords(42)    // "сорок два"
propisyu.IntToWords(1000)  // "одна тысяча"
```

#### `IntToWordsGender(n int, gender Gender) string`
Converts an integer to words with specified gender.

Available genders:
- `GenderMasculine` (1) - masculine
- `GenderFeminine` (2) - feminine
- `GenderNeuter` (3) - neuter

```go
propisyu.IntToWordsGender(2, propisyu.GenderMasculine)  // "два"
propisyu.IntToWordsGender(2, propisyu.GenderFeminine)   // "две"
propisyu.IntToWordsGender(1, propisyu.GenderNeuter)     // "одно"
```

### Decimal Functions

#### `DecimalToWords(decimalStr string) (string, error)`
Converts a decimal number from string to words. Fractional part is truncated to 2 decimal places.

```go
result, err := propisyu.DecimalToWords("3.14")
// "три целых четырнадцать сотых"
```

#### `DecimalValueToWords(d decimal.Decimal) (string, error)`
Converts a `decimal.Decimal` value to words. Fractional part is truncated (not rounded!) to 2 decimal places.

```go
d := decimal.NewFromFloat(3.14159)
result, err := propisyu.DecimalValueToWords(d)
// "три целых четырнадцать сотых" (truncates, not rounds)
```

**Important:** The function truncates, not rounds:
- `1.999` → "один целых девяносто девять сотых"
- `1.995` → "один целых девяносто девять сотых"

### Declension Function

#### `Decline(n int, one, two, five string) string`
Chooses the correct noun form based on the number.

Parameters:
- `n` - the number
- `one` - form for 1, 21, 31... (рубль, день, товар)
- `two` - form for 2-4, 22-24... (рубля, дня, товара)
- `five` - form for 0, 5-20, 25-30... (рублей, дней, товаров)

```go
propisyu.Decline(1, "рубль", "рубля", "рублей")   // "рубль"
propisyu.Decline(2, "рубль", "рубля", "рублей")   // "рубля"
propisyu.Decline(5, "рубль", "рубля", "рублей")   // "рублей"
propisyu.Decline(11, "рубль", "рубля", "рублей")  // "рублей"
propisyu.Decline(21, "рубль", "рубля", "рублей")  // "рубль"
```

## Why go-propisyu

- **Pure Go** — not a C wrapper; easy to build and deploy anywhere
- **Correct grammar** — three genders, proper declension for all numeric ranges
- **Zero deps** — `IntToWords` and `Decline` require no third-party packages
- **Production-ready** — CI with linter, tests, semantic versioning, goreleaser
- **Open license** — MIT, free for commercial use

## Limitations

- Integer numbers: supports values in the `int` range (typically -2³¹ to 2³¹-1 or -2⁶³ to 2⁶³-1)
- Decimal numbers: only 2 decimal places are supported (rest is truncated)
- `DecimalValueToWords` returns `ErrNumberTooLarge` error if number doesn't fit in `int`

## Testing

```bash
go test ./...              # Run all tests
go test -v ./...           # With verbose output
go test -cover ./...       # With code coverage
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT
