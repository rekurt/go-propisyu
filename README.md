# go-propisyu · Числа прописью на Go

**Русская версия · [English version](README_EN.md)**

[![CI](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml/badge.svg)](https://github.com/rekurt/go-propisyu/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/rekurt/go-propisyu.svg)](https://pkg.go.dev/github.com/rekurt/go-propisyu)
[![Go Report Card](https://goreportcard.com/badge/github.com/rekurt/go-propisyu)](https://goreportcard.com/report/github.com/rekurt/go-propisyu)
[![codecov](https://codecov.io/gh/rekurt/go-propisyu/branch/master/graph/badge.svg)](https://codecov.io/gh/rekurt/go-propisyu)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Одна библиотека, которая закрывает всё про **«числа прописью»** на русском:
суммы с валютой, десятичные числа произвольной точности, порядковые номера
во всех трёх родах, ручные склонения для любого существительного. Одна
строка Go вместо километра ручных условий — для счетов, фискальных чеков,
документов, 1С-интеграций и голосовых ассистентов.

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

- **Весь диапазон Go `int`** — от `math.MinInt` до `math.MaxInt`, до
  дуодециллионов (10³⁹) включительно. Крайние значения (в т.ч. `MinInt`)
  не ломают `-n`, библиотека безопасно поднимает магнитуду через `uint64`.
- **Три рода** — `GenderMasculine`, `GenderFeminine`, `GenderNeuter` —
  применяются и к кардиналам, и к порядковым.
- **Порядковые числительные** — `Ordinal(n, gender)` даёт «первый / первая
  / первое», корректно для составных («сорок второй»), для round-чисел
  («тысячный», «сороковой», «миллионный») и для всех трёх родов.
- **Десятичные с произвольной точностью** — `DecimalToWordsPrecision`
  поддерживает 1–9 знаков после запятой: десятых → миллиардных.
  `DecimalValueToWords` принимает `shopspring/decimal` напрямую.
- **Готовые валюты** — `Money` + пресеты `CurrencyRUB`, `CurrencyUSD`,
  `CurrencyEUR`. Пара `(1234, 56)` превращается в «одна тысяча двести
  тридцать четыре рубля пятьдесят шесть копеек» одной строкой.
- **Свои существительные через `Decline`** — правильные русские формы
  для любого слова, корректно обрабатываются 11–14, 21 и отрицательные
  числа.
- **Нулевые внешние зависимости в core** — `shopspring/decimal` нужен
  только если вы сами используете `DecimalValueToWords`. CI, линтер,
  тесты, семантическое версионирование, release через goreleaser.

## Где применяется

| Сфера | Пример |
|---|---|
| Финтех и банкинг | Сумма прописью в платёжных поручениях и выписках |
| Бухгалтерия и 1С | Счета-фактуры, акты, накладные |
| Фискальные чеки | Касса / ОФД — сумма словами по 54-ФЗ |
| Голосовые ассистенты | TTS-озвучка сумм и количеств |
| Чат-боты | Ответы с суммами на естественном языке |
| Генерация документов | Договоры, доверенности, акты |

## Установка

```bash
go get github.com/rekurt/go-propisyu
```

`shopspring/decimal` нужен только для `DecimalValueToWords`:

```bash
go get github.com/shopspring/decimal
```

## Использование

### Сумма прописью для счёта

```go
res, _ := propisyu.MoneyFromString("1234.56", propisyu.CurrencyRUB)
fmt.Println(res)
// одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек

fmt.Println(propisyu.Money(1, 1, propisyu.CurrencyRUB))
// один рубль одна копейка

fmt.Println(propisyu.Money(100, 99, propisyu.CurrencyEUR))
// сто евро девяносто девять центов
```

`Money` сам выбирает род и склонение для целой и дробной частей — не надо
вручную жонглировать «рубль / рубля / рублей».

### Десятичные с точностью больше двух знаков

```go
res, _ := propisyu.DecimalToWordsPrecision("3.14159", 5)
fmt.Println(res)
// три целых четырнадцать тысяч сто пятьдесят девять стотысячных

res, _ = propisyu.DecimalToWordsPrecision("3.5", 1)
fmt.Println(res)
// три целых пять десятых
```

Параметр `precision` — это число знаков после запятой (1–9). От него
зависит, во что склоняется дробная часть: десятых, сотых, тысячных,
десятитысячных, стотысячных, миллионных, десятимиллионных,
стомиллионных, миллиардных.

### Порядковые номера для документов

```go
fmt.Println(propisyu.Ordinal(1, propisyu.GenderMasculine))     // первый
fmt.Println(propisyu.Ordinal(1, propisyu.GenderFeminine))      // первая
fmt.Println(propisyu.Ordinal(1, propisyu.GenderNeuter))        // первое

fmt.Println(propisyu.Ordinal(42, propisyu.GenderMasculine))    // сорок второй
fmt.Println(propisyu.Ordinal(1000, propisyu.GenderFeminine))   // тысячная
fmt.Println(propisyu.Ordinal(40, propisyu.GenderMasculine))    // сороковой
```

Работает и для round-чисел («тысячный», «миллионный»), и для составных
(«сорок второй», «двадцать первый»).

### Свои склонения через `Decline`

```go
fmt.Println(propisyu.Decline(5,  "товар", "товара", "товаров")) // товаров
fmt.Println(propisyu.Decline(5,  "день",  "дня",    "дней"))    // дней
fmt.Println(propisyu.Decline(21, "день",  "дня",    "дней"))    // день
fmt.Println(propisyu.Decline(11, "рубль", "рубля",  "рублей"))  // рублей
```

Правило «1 / 2–4 / 5–20» с исключениями для 11–14 и для отрицательных
чисел встроено — передавайте своё существительное в трёх формах и
получайте правильную.

## API

### Целые числа

| Функция | Описание |
|---|---|
| `IntToWords(n int) string` | Целое число в слова, мужской род по умолчанию |
| `IntToWordsGender(n int, gender Gender) string` | То же с явным родом |

```go
propisyu.IntToWords(42)                                // сорок два
propisyu.IntToWords(1000)                              // одна тысяча
propisyu.IntToWords(-321)                              // минус триста двадцать один

propisyu.IntToWordsGender(2, propisyu.GenderMasculine) // два
propisyu.IntToWordsGender(2, propisyu.GenderFeminine)  // две
propisyu.IntToWordsGender(1, propisyu.GenderNeuter)    // одно
```

Константы рода: `GenderMasculine`, `GenderFeminine`, `GenderNeuter`.

### Десятичные числа

| Функция | Описание |
|---|---|
| `DecimalToWords(s string) (string, error)` | Строка с фиксированной точностью `.xx` → слова |
| `DecimalValueToWords(d decimal.Decimal) (string, error)` | `shopspring/decimal` → слова |
| `DecimalToWordsPrecision(s string, precision int) (string, error)` | Строка с произвольной точностью 1–9 знаков |

```go
propisyu.DecimalToWords("123.45")
// сто двадцать три целых сорок пять сотых

propisyu.DecimalValueToWords(decimal.NewFromFloat(3.14159))
// три целых четырнадцать сотых

propisyu.DecimalToWords("-0.50")
// минус ноль целых пятьдесят сотых
```

Важные особенности:

- Дробная часть **обрезается, а не округляется** (truncate): `1.999` →
  `одна целая девяносто девять сотых`.
- Целая часть всегда идёт в **женском роде** (`одна целая`, `две целых`).
- Знак «минус» сохраняется даже для `-0.xx`, где целая часть равна нулю.

### Порядковые числительные

| Функция | Описание |
|---|---|
| `Ordinal(n int, gender Gender) string` | Порядковое число в указанном роде |

```go
propisyu.Ordinal(21,        propisyu.GenderMasculine) // двадцать первый
propisyu.Ordinal(1000,      propisyu.GenderFeminine)  // тысячная
propisyu.Ordinal(1_000_000, propisyu.GenderMasculine) // миллионный
```

### Валюты и деньги

| Функция / тип | Описание |
|---|---|
| `type Currency struct { ... }` | Описание валюты: три формы целой, три формы дробной, род каждой части |
| `CurrencyRUB`, `CurrencyUSD`, `CurrencyEUR` | Готовые пресеты |
| `Money(whole, cents int, c Currency) string` | Сумма прописью по разобранным полям |
| `MoneyFromString(amount string, c Currency) (string, error)` | Парсит `"1234.56"` и сразу отдаёт результат |

```go
propisyu.Money(1234, 56, propisyu.CurrencyRUB)
// одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек

propisyu.Money(42, 0, propisyu.CurrencyUSD)
// сорок два доллара ноль центов
```

Поля `Currency`: `WholeOne`, `WholeTwo`, `WholeFive`, `WholeGender`,
`FracOne`, `FracTwo`, `FracFive`, `FracGender`. Соберите свой пресет для
любой валюты (токены, баллы, условные единицы) без правок в коде
библиотеки.

### Склонения

| Функция | Описание |
|---|---|
| `Decline(n int, one, two, five string) string` | Выбирает форму существительного по числу |

Правило совпадает с русской грамматикой:

| Последняя цифра `n` | `n % 100 ∈ 11…19` | Форма | Пример |
|---|---|---|---|
| 1 | нет | `one` | рубль, день |
| 2, 3, 4 | нет | `two` | рубля, дня |
| 0, 5–9 | — | `five` | рублей, дней |
| любая | да | `five` | 11 → рублей, 19 → рублей |

Отрицательные числа обрабатываются по модулю: `Decline(-1, …)` → форма
`one`.

### Ошибки

`ErrNumberTooLarge` — возвращается из `DecimalValueToWords`, если целая
часть `decimal.Decimal` не помещается в Go `int`.

## Ограничения

- Целые числа ограничены Go `int` — на 64-битных платформах это ±9.2·10¹⁸.
- `DecimalToWords` и `DecimalValueToWords` работают только с двумя знаками
  после запятой (остальное обрезается). Для большей точности используйте
  `DecimalToWordsPrecision` (1–9 знаков).
- Классификаторы и смешанные системы счёта (например, якутское спряжение)
  не поддерживаются — библиотека строго под русскую грамматику.

## Contributing

Contributions are welcome! См. [CONTRIBUTING.md](CONTRIBUTING.md).

## Лицензия

MIT. См. [LICENSE](LICENSE).
