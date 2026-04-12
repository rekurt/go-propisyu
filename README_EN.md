# go-propisyu · Russian Number-to-Words for Go

**English version · [Русская версия](README.md)**

[![CI](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml/badge.svg)](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/rekurt/go-propisyu.svg)](https://pkg.go.dev/github.com/rekurt/go-propisyu)
[![Go Report Card](https://goreportcard.com/badge/github.com/rekurt/go-propisyu)](https://goreportcard.com/report/github.com/rekurt/go-propisyu)
[![codecov](https://codecov.io/gh/rekurt/go-propisyu/branch/master/graph/badge.svg)](https://codecov.io/gh/rekurt/go-propisyu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Russian number-to-words: amounts with currency, decimals, ordinal numbers
in all three grammatical genders, and noun declension for any word.

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

- **Full Go `int` range** — from `math.MinInt` to `math.MaxInt`
  (±9.2 × 10¹⁸ on 64-bit). Edge values (including `MinInt`) do not
  crash `-n` — the magnitude is safely computed via `uint64`.
- **Three grammatical genders** — `GenderMasculine`, `GenderFeminine`,
  `GenderNeuter` — for cardinals (`IntToWordsGender`) and ordinals
  (`Ordinal`).
- **Ordinal numbers** — `Ordinal(n, gender)`: compound forms
  ("сорок второй"), round numbers ("тысячный", "миллионный"), all
  three genders.
- **Decimals with arbitrary precision** — `DecimalToWordsPrecision`
  supports 1–9 fractional digits (tenths through billionths).
  `DecimalValueToWords` accepts `shopspring/decimal`.
- **Currencies** — `Money` + presets `CurrencyRUB`, `CurrencyUSD`,
  `CurrencyEUR`; easy to create your own `Currency`.
- **Declension** — `Decline` picks the correct noun form by number,
  with proper handling of 11–14 and negatives.
- **Zero deps in the core** — `shopspring/decimal` is only needed
  for `DecimalValueToWords`.

## Installation

```bash
go get github.com/rekurt/go-propisyu
```

`shopspring/decimal` is only needed for `DecimalValueToWords`:

```bash
go get github.com/shopspring/decimal
```

## Usage

### Integers

| Function | Description |
|---|---|
| `IntToWords(n int) string` | Integer to words, masculine gender by default |
| `IntToWordsGender(n int, gender Gender) string` | Same with an explicit gender |

```go
propisyu.IntToWords(42)                                // сорок два
propisyu.IntToWords(1000)                              // одна тысяча
propisyu.IntToWords(-321)                              // минус триста двадцать один

propisyu.IntToWordsGender(1, propisyu.GenderMasculine) // один
propisyu.IntToWordsGender(1, propisyu.GenderFeminine)  // одна
propisyu.IntToWordsGender(1, propisyu.GenderNeuter)    // одно
```

Gender constants: `GenderMasculine`, `GenderFeminine`, `GenderNeuter`
(type `Gender`).

### Currency & Money

| Function / type | Description |
|---|---|
| `Money(whole, cents int, c Currency) string` | Amount in words |
| `MoneyFromString(amount string, c Currency) (string, error)` | Parses `"1234.56"` and returns the result |
| `CurrencyRUB`, `CurrencyUSD`, `CurrencyEUR` | Built-in presets |

```go
propisyu.Money(1234, 56, propisyu.CurrencyRUB)
// одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек

propisyu.Money(1, 1, propisyu.CurrencyRUB)
// один рубль одна копейка

propisyu.Money(100, 99, propisyu.CurrencyEUR)
// сто евро девяносто девять центов
```

Custom `Currency` preset — for any unit (tokens, points, loyalty):

```go
myTokens := propisyu.Currency{
    WholeOne: "токен", WholeTwo: "токена", WholeFive: "токенов",
    WholeGender: propisyu.GenderMasculine,
    FracOne: "юнит", FracTwo: "юнита", FracFive: "юнитов",
    FracGender: propisyu.GenderMasculine,
}
propisyu.Money(42, 5, myTokens)
// сорок два токена пять юнитов
```

### Decimals

| Function | Description |
|---|---|
| `DecimalToWords(s string) (string, error)` | String with fixed `.xx` precision |
| `DecimalValueToWords(d decimal.Decimal) (string, error)` | `shopspring/decimal` directly |
| `DecimalToWordsPrecision(s string, precision int) (string, error)` | Arbitrary 1–9 digit precision |

```go
propisyu.DecimalToWords("123.45")
// сто двадцать три целых сорок пять сотых

propisyu.DecimalToWordsPrecision("3.14159", 5)
// три целых четырнадцать тысяч сто пятьдесят девять стотысячных

propisyu.DecimalToWords("-0.50")
// минус ноль целых пятьдесят сотых
```

The fractional part is truncated, not rounded. The whole part is always
rendered in the feminine gender ("одна целая", "две целых"). The minus
sign is preserved for `-0.xx`.

### Ordinals

| Function | Description |
|---|---|
| `Ordinal(n int, gender Gender) string` | Ordinal number in the chosen gender |

```go
propisyu.Ordinal(1, propisyu.GenderMasculine)     // первый
propisyu.Ordinal(1, propisyu.GenderFeminine)      // первая
propisyu.Ordinal(42, propisyu.GenderMasculine)    // сорок второй
propisyu.Ordinal(1000, propisyu.GenderFeminine)   // тысячная
propisyu.Ordinal(1_000_000, propisyu.GenderMasculine) // миллионный
```

### Declension

| Function | Description |
|---|---|
| `Decline(n int, one, two, five string) string` | Picks a noun form based on the number |

```go
propisyu.Decline(1,  "рубль", "рубля", "рублей") // рубль
propisyu.Decline(5,  "день",  "дня",   "дней")   // дней
propisyu.Decline(21, "день",  "дня",   "дней")   // день
propisyu.Decline(11, "рубль", "рубля", "рублей") // рублей
```

| Last digit of `n` | `n % 100 ∈ 11…19` | Form |
|---|---|---|
| 1 | no | `one` |
| 2, 3, 4 | no | `two` |
| 0, 5–9 | — | `five` |
| any | yes | `five` |

Negative numbers are handled by magnitude.

### Errors

`ErrNumberTooLarge` is returned from `DecimalValueToWords` when the
whole part of a `decimal.Decimal` does not fit in Go `int`.

## Limitations

- Integers are bounded by Go `int` — on 64-bit platforms this is
  ±9.2 × 10¹⁸.
- `DecimalToWords` and `DecimalValueToWords` work with two fractional
  digits (anything beyond is truncated). For higher precision use
  `DecimalToWordsPrecision` (1–9 digits).

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT. See [LICENSE](LICENSE).
