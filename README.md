# go-propisyu · Числа прописью на Go

**Русская версия · [English version](README_EN.md)**

[![CI](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml/badge.svg)](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/rekurt/go-propisyu.svg)](https://pkg.go.dev/github.com/rekurt/go-propisyu)
[![Go Report Card](https://goreportcard.com/badge/github.com/rekurt/go-propisyu)](https://goreportcard.com/report/github.com/rekurt/go-propisyu)
[![codecov](https://codecov.io/gh/rekurt/go-propisyu/branch/master/graph/badge.svg)](https://codecov.io/gh/rekurt/go-propisyu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Числа прописью на русском: суммы с валютой, десятичные дроби, порядковые
числительные во всех трёх родах, склонения для любого существительного.

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

## Возможности

- **Весь диапазон Go `int`** — от `math.MinInt` до `math.MaxInt`
  (±9,2 × 10¹⁸ на 64-бит). Крайние значения (в т.ч. `MinInt`) не ломают
  `-n` — магнитуда безопасно вычисляется через `uint64`.
- **Три рода** — `GenderMasculine`, `GenderFeminine`, `GenderNeuter` —
  для количественных (`IntToWordsGender`) и порядковых (`Ordinal`).
- **Порядковые числительные** — `Ordinal(n, gender)`: составные
  («сорок второй»), круглые («тысячный», «миллионный»), все три рода.
- **Десятичные с произвольной точностью** — `DecimalToWordsPrecision`
  поддерживает 1–9 знаков после запятой (от десятых до миллиардных).
  `DecimalValueToWords` принимает `shopspring/decimal`.
- **Валюты** — `Money` + пресеты `CurrencyRUB`, `CurrencyUSD`,
  `CurrencyEUR`; легко создать свой `Currency`.
- **Склонения** — `Decline` выбирает правильную форму существительного
  по числу, с обработкой 11–14 и отрицательных.
- **Нулевые зависимости в core** — `shopspring/decimal` нужен только
  для `DecimalValueToWords`.

## Установка

```bash
go get github.com/rekurt/go-propisyu
```

`shopspring/decimal` нужен только для `DecimalValueToWords`:

```bash
go get github.com/shopspring/decimal
```

## Использование

### Целые числа

| Функция | Описание |
|---|---|
| `IntToWords(n int) string` | Число в слова, мужской род по умолчанию |
| `IntToWordsGender(n int, gender Gender) string` | То же с явным родом |

```go
propisyu.IntToWords(42)                                // сорок два
propisyu.IntToWords(1000)                              // одна тысяча
propisyu.IntToWords(-321)                              // минус триста двадцать один

propisyu.IntToWordsGender(1, propisyu.GenderMasculine) // один
propisyu.IntToWordsGender(1, propisyu.GenderFeminine)  // одна
propisyu.IntToWordsGender(1, propisyu.GenderNeuter)    // одно
```

Константы рода: `GenderMasculine`, `GenderFeminine`, `GenderNeuter`
(тип `Gender`).

### Валюты и деньги

| Функция / тип | Описание |
|---|---|
| `Money(whole, cents int, c Currency) string` | Сумма прописью |
| `MoneyFromString(amount string, c Currency) (string, error)` | Парсит `"1234.56"` и отдаёт результат |
| `CurrencyRUB`, `CurrencyUSD`, `CurrencyEUR` | Готовые пресеты |

```go
propisyu.Money(1234, 56, propisyu.CurrencyRUB)
// одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек

propisyu.Money(1, 1, propisyu.CurrencyRUB)
// один рубль одна копейка

propisyu.Money(100, 99, propisyu.CurrencyEUR)
// сто евро девяносто девять центов
```

Свой пресет `Currency` — для любой единицы (токены, баллы, валюты):

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

### Десятичные числа

| Функция | Описание |
|---|---|
| `DecimalToWords(s string) (string, error)` | Строка с фиксированной точностью `.xx` |
| `DecimalValueToWords(d decimal.Decimal) (string, error)` | `shopspring/decimal` напрямую |
| `DecimalToWordsPrecision(s string, precision int) (string, error)` | Произвольная точность 1–9 знаков |

```go
propisyu.DecimalToWords("123.45")
// сто двадцать три целых сорок пять сотых

propisyu.DecimalToWordsPrecision("3.14159", 5)
// три целых четырнадцать тысяч сто пятьдесят девять стотысячных

propisyu.DecimalToWords("-0.50")
// минус ноль целых пятьдесят сотых
```

Дробная часть обрезается, а не округляется. Целая часть идёт в женском
роде («одна целая», «две целых»). Знак минус сохраняется для `-0.xx`.

### Порядковые числительные

| Функция | Описание |
|---|---|
| `Ordinal(n int, gender Gender) string` | Порядковое числительное в указанном роде |

```go
propisyu.Ordinal(1, propisyu.GenderMasculine)     // первый
propisyu.Ordinal(1, propisyu.GenderFeminine)      // первая
propisyu.Ordinal(42, propisyu.GenderMasculine)    // сорок второй
propisyu.Ordinal(1000, propisyu.GenderFeminine)   // тысячная
propisyu.Ordinal(1_000_000, propisyu.GenderMasculine) // миллионный
```

### Склонения

| Функция | Описание |
|---|---|
| `Decline(n int, one, two, five string) string` | Выбирает форму существительного по числу |

```go
propisyu.Decline(1,  "рубль", "рубля", "рублей") // рубль
propisyu.Decline(5,  "день",  "дня",   "дней")   // дней
propisyu.Decline(21, "день",  "дня",   "дней")   // день
propisyu.Decline(11, "рубль", "рубля", "рублей") // рублей
```

| Последняя цифра `n` | `n % 100 ∈ 11…19` | Форма |
|---|---|---|
| 1 | нет | `one` |
| 2, 3, 4 | нет | `two` |
| 0, 5–9 | — | `five` |
| любая | да | `five` |

Отрицательные числа обрабатываются по модулю.

### Ошибки

`ErrNumberTooLarge` — возвращается из `DecimalValueToWords`, если целая
часть `decimal.Decimal` не помещается в Go `int`.

## Ограничения

- Целые числа ограничены Go `int` — на 64-битных платформах это
  ±9,2 × 10¹⁸.
- `DecimalToWords` и `DecimalValueToWords` работают с двумя знаками
  после запятой (остальное обрезается). Для большей точности —
  `DecimalToWordsPrecision` (1–9 знаков).

## Contributing

Contributions are welcome! См. [CONTRIBUTING.md](CONTRIBUTING.md).

## Лицензия

MIT. См. [LICENSE](LICENSE).
