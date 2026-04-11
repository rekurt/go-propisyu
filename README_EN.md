# go-propisyu · Russian Number-to-Words for Go

**English version · [Русская версия](README.md)**

[![CI](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml/badge.svg)](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/rekurt/go-propisyu.svg)](https://pkg.go.dev/github.com/rekurt/go-propisyu)
[![Go Report Card](https://goreportcard.com/badge/github.com/rekurt/go-propisyu)](https://goreportcard.com/report/github.com/rekurt/go-propisyu)
[![codecov](https://codecov.io/gh/rekurt/go-propisyu/branch/master/graph/badge.svg)](https://codecov.io/gh/rekurt/go-propisyu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

One library that covers **everything** about spelling Russian numbers:
amounts with currency, decimals of arbitrary precision, ordinal numbers
in all three grammatical genders, and manual declension for any noun.
One line of Go instead of a mountain of hand-written conditions — ready
for invoices, fiscal receipts, financial documents, 1C integrations,
and voice assistants.

```go
propisyu.Money(1234, 56, propisyu.CurrencyRUB)
// одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек

propisyu.Ordinal(1, propisyu.GenderFeminine)
// первая

propisyu.DecimalToWordsPrecision("3.141", 3)
// три целых сто сорок одна тысячная

propisyu.Decline(5, "заказ", "заказа", "заказов")
// заказов
```

## Features

- **Full Go `int` range** — from `math.MinInt` to `math.MaxInt`, up to
  duodecillions (10³⁹). Edge values (including `MinInt`) do not crash
  `-n`; the library safely lifts the magnitude through `uint64`.
- **Three grammatical genders** — `GenderMasculine`, `GenderFeminine`,
  `GenderNeuter` — applied to both cardinals and ordinals.
- **Ordinal numbers** — `Ordinal(n, gender)` produces "первый / первая /
  первое", correctly handles compound forms ("сорок второй"), round
  numbers ("тысячный", "сороковой", "миллионный"), and all three genders.
- **Decimals with arbitrary precision** — `DecimalToWordsPrecision`
  supports 1–9 fractional digits (tenths → billionths).
  `DecimalValueToWords` accepts `shopspring/decimal` directly.
- **Ready-made currencies** — `Money` + presets `CurrencyRUB`,
  `CurrencyUSD`, `CurrencyEUR`. The pair `(1234, 56)` becomes "одна
  тысяча двести тридцать четыре рубля пятьдесят шесть копеек" in one
  call.
- **Your own nouns via `Decline`** — correct Russian forms for any word,
  with proper handling of 11–14, 21, and negative numbers.
- **Zero external deps in the core** — `shopspring/decimal` is only
  needed if you use `DecimalValueToWords` yourself. CI, linter, tests,
  semantic versioning, release via goreleaser.

## Use Cases

| Domain | Example |
|---|---|
| Fintech & Banking | Amount in words on payment orders and statements |
| Accounting | Invoices, acts, and waybills |
| Fiscal Receipts | POS / OFD — amount in words (Russian 54-FZ) |
| Voice Assistants | TTS pronunciation of amounts and quantities |
| Chatbots | Natural-language responses with amounts |
| Document Generation | Contracts, powers of attorney, acts |

## Installation

```bash
go get github.com/rekurt/go-propisyu
```

`shopspring/decimal` is only needed for `DecimalValueToWords`:

```bash
go get github.com/shopspring/decimal
```

## Usage

### Amount in Words for an Invoice

```go
res, _ := propisyu.MoneyFromString("1234.56", propisyu.CurrencyRUB)
fmt.Println(res)
// одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек

fmt.Println(propisyu.Money(1, 1, propisyu.CurrencyRUB))
// один рубль одна копейка

fmt.Println(propisyu.Money(100, 99, propisyu.CurrencyEUR))
// сто евро девяносто девять центов
```

`Money` picks the correct gender and declension for both the whole and
fractional parts — you never have to juggle "рубль / рубля / рублей" by
hand.

### Decimals with More Than Two Digits

```go
res, _ := propisyu.DecimalToWordsPrecision("3.14159", 5)
fmt.Println(res)
// три целых четырнадцать тысяч сто пятьдесят девять стотысячных

res, _ = propisyu.DecimalToWordsPrecision("3.5", 1)
fmt.Println(res)
// три целых пять десятых
```

The `precision` parameter is the number of fractional digits (1–9) and
selects how the fractional part is declined: tenths (десятых), hundredths
(сотых), thousandths (тысячных), ten-thousandths, hundred-thousandths,
millionths, ten-millionths, hundred-millionths, or billionths.

### Ordinal Numbers for Documents

```go
fmt.Println(propisyu.Ordinal(1, propisyu.GenderMasculine))     // первый
fmt.Println(propisyu.Ordinal(1, propisyu.GenderFeminine))      // первая
fmt.Println(propisyu.Ordinal(1, propisyu.GenderNeuter))        // первое

fmt.Println(propisyu.Ordinal(42, propisyu.GenderMasculine))    // сорок второй
fmt.Println(propisyu.Ordinal(1000, propisyu.GenderFeminine))   // тысячная
fmt.Println(propisyu.Ordinal(40, propisyu.GenderMasculine))    // сороковой
```

Works for both round numbers ("тысячный", "миллионный") and compound
forms ("сорок второй", "двадцать первый").

### Custom Declension via `Decline`

```go
fmt.Println(propisyu.Decline(5,  "товар", "товара", "товаров")) // товаров
fmt.Println(propisyu.Decline(5,  "день",  "дня",    "дней"))    // дней
fmt.Println(propisyu.Decline(21, "день",  "дня",    "дней"))    // день
fmt.Println(propisyu.Decline(11, "рубль", "рубля",  "рублей"))  // рублей
```

The "1 / 2–4 / 5–20" rule with the 11–14 and negative-number exceptions
is built in — just pass your noun in three forms and get the correct
one back.

## API

### Integers

| Function | Description |
|---|---|
| `IntToWords(n int) string` | Integer to words, masculine gender by default |
| `IntToWordsGender(n int, gender Gender) string` | Same with an explicit gender |

```go
propisyu.IntToWords(42)                                // сорок два
propisyu.IntToWords(1000)                              // одна тысяча
propisyu.IntToWords(-321)                              // минус триста двадцать один

propisyu.IntToWordsGender(2, propisyu.GenderMasculine) // два
propisyu.IntToWordsGender(2, propisyu.GenderFeminine)  // две
propisyu.IntToWordsGender(1, propisyu.GenderNeuter)    // одно
```

Gender constants: `GenderMasculine`, `GenderFeminine`, `GenderNeuter`.

### Decimals

| Function | Description |
|---|---|
| `DecimalToWords(s string) (string, error)` | String with fixed `.xx` precision → words |
| `DecimalValueToWords(d decimal.Decimal) (string, error)` | `shopspring/decimal` → words |
| `DecimalToWordsPrecision(s string, precision int) (string, error)` | String with arbitrary 1–9 digit precision |

```go
propisyu.DecimalToWords("123.45")
// сто двадцать три целых сорок пять сотых

propisyu.DecimalValueToWords(decimal.NewFromFloat(3.14159))
// три целых четырнадцать сотых

propisyu.DecimalToWords("-0.50")
// минус ноль целых пятьдесят сотых
```

Key things to know:

- The fractional part is **truncated, not rounded**: `1.999` → "одна
  целая девяносто девять сотых".
- The whole part is always rendered in the **feminine gender** ("одна
  целая", "две целых").
- The minus sign is preserved even for `-0.xx`, where the whole part is
  zero.

### Ordinals

| Function | Description |
|---|---|
| `Ordinal(n int, gender Gender) string` | Ordinal number in the chosen gender |

```go
propisyu.Ordinal(21,        propisyu.GenderMasculine) // двадцать первый
propisyu.Ordinal(1000,      propisyu.GenderFeminine)  // тысячная
propisyu.Ordinal(1_000_000, propisyu.GenderMasculine) // миллионный
```

### Currency & Money

| Function / type | Description |
|---|---|
| `type Currency struct { ... }` | Currency descriptor: three forms for the whole part, three for the fraction, gender for each |
| `CurrencyRUB`, `CurrencyUSD`, `CurrencyEUR` | Built-in presets |
| `Money(whole, cents int, c Currency) string` | Amount in words from pre-parsed fields |
| `MoneyFromString(amount string, c Currency) (string, error)` | Parses `"1234.56"` and returns the result |

```go
propisyu.Money(1234, 56, propisyu.CurrencyRUB)
// одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек

propisyu.Money(42, 0, propisyu.CurrencyUSD)
// сорок два доллара ноль центов
```

`Currency` fields: `WholeOne`, `WholeTwo`, `WholeFive`, `WholeGender`,
`FracOne`, `FracTwo`, `FracFive`, `FracGender`. Build your own preset for
any currency (tokens, points, loyalty units) without touching the
library code.

### Declension

| Function | Description |
|---|---|
| `Decline(n int, one, two, five string) string` | Picks a noun form based on the number |

The rule matches Russian grammar:

| Last digit of `n` | `n % 100 ∈ 11…19` | Form | Example |
|---|---|---|---|
| 1 | no | `one` | рубль, день |
| 2, 3, 4 | no | `two` | рубля, дня |
| 0, 5–9 | — | `five` | рублей, дней |
| any | yes | `five` | 11 → рублей, 19 → рублей |

Negative numbers are handled by magnitude: `Decline(-1, …)` → the `one`
form.

### Errors

`ErrNumberTooLarge` is returned from `DecimalValueToWords` when the
whole part of a `decimal.Decimal` does not fit in Go `int`.

## Limitations

- Integers are bounded by Go `int` — on 64-bit platforms this is
  ±9.2·10¹⁸.
- `DecimalToWords` and `DecimalValueToWords` only work with two
  fractional digits (anything beyond is truncated). For higher precision
  use `DecimalToWordsPrecision` (1–9 digits).
- Classifier and mixed counting systems (e.g. Yakut-style agreement)
  are not supported — the library targets Russian grammar strictly.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT. See [LICENSE](LICENSE).
