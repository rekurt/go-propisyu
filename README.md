# go-propisyu · Russian Number-to-Words Converter

**Русская версия · [English version](README_EN.md)**

`go-propisyu` — Go-библиотека, которая переводит числа в русские слова и подбирает правильные склонения. Решение подходит для счетов, фискальных чеков, голосовых ассистентов, генерации документов и любых сервисов, где важно грамотно проговаривать суммы и количества.

## Ключевые особенности

- 🔢 Поддержка огромных чисел вплоть до дуодециллионов (10³⁹)
- 🧠 Грамматические роды: мужской, женский и средний для корректных окончаний
- 💬 Функция `Decline` для автоматического склонения существительных
- 💰 Работа с десятичными числами через строки или `decimal.Decimal`
- ✅ Нулевые внешние зависимости для базовых функций

## Установка

```bash
go get github.com/rekurt/go-propisyu
```

Для работы с `decimal.Decimal`:
```bash
go get github.com/shopspring/decimal
```

## Публичные функции

| Функция | Описание |
|---------|----------|
| `IntToWords(n int) string` | Конвертирует целое число в слова (мужской род) |
| `IntToWordsGender(n int, gender Gender) string` | Конвертирует целое число в слова с указанием рода |
| `DecimalToWords(decimalStr string) (string, error)` | Конвертирует десятичное число из строки в слова |
| `DecimalValueToWords(d decimal.Decimal) (string, error)` | Конвертирует `decimal.Decimal` значение в слова |
| `Decline(n int, one, two, five string) string` | Выбирает правильную форму склонения существительного |

### Константы рода

```go
const (
    GenderMasculine Gender = 1  // Мужской род: "один", "два"
    GenderFeminine  Gender = 2  // Женский род: "одна", "две"
    GenderNeuter    Gender = 3  // Средний род: "одно", "два"
)
```

### Ошибки

- `ErrNumberTooLarge` - число слишком велико для конвертации (не помещается в `int`)

## Быстрый старт

### Целые числа

```go
package main

import (
	"fmt"

	"github.com/rekurt/go-propisyu"
)

func main() {
	// Базовая конвертация (мужской род по умолчанию)
	fmt.Println(propisyu.IntToWords(321))
	// триста двадцать один

	// Конвертация с указанием рода
	fmt.Println(propisyu.IntToWordsGender(2, propisyu.GenderFeminine))
	// две

	fmt.Println(propisyu.IntToWordsGender(2, propisyu.GenderMasculine))
	// два

	// Автоматическое склонение существительных
	fmt.Println(propisyu.Decline(1, "рубль", "рубля", "рублей"))   // рубль
	fmt.Println(propisyu.Decline(2, "рубль", "рубля", "рублей"))   // рубля
	fmt.Println(propisyu.Decline(5, "рубль", "рубля", "рублей"))   // рублей
	fmt.Println(propisyu.Decline(21, "рубль", "рубля", "рублей"))  // рубль
}
```

### Десятичные числа

#### Способ 1: Используя строку

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

#### Способ 2: Используя decimal.Decimal

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

## Примеры использования

### Сумма прописью для чека

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

### Большие числа

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

### Склонение с разными существительными

```go
package main

import (
	"fmt"

	"github.com/rekurt/go-propisyu"
)

func main() {
	count := 5

	// Валюты
	fmt.Println(count, propisyu.Decline(count, "доллар", "доллара", "долларов"))
	// 5 долларов

	// Единицы измерения
	fmt.Println(count, propisyu.Decline(count, "метр", "метра", "метров"))
	// 5 метров

	// Предметы
	fmt.Println(count, propisyu.Decline(count, "товар", "товара", "товаров"))
	// 5 товаров

	// Для 21
	count = 21
	fmt.Println(count, propisyu.Decline(count, "день", "дня", "дней"))
	// 21 день
}
```

## API

### Функции для целых чисел

#### `IntToWords(n int) string`
Конвертирует целое число в слова (мужской род по умолчанию).

```go
propisyu.IntToWords(42)    // "сорок два"
propisyu.IntToWords(1000)  // "одна тысяча"
```

#### `IntToWordsGender(n int, gender Gender) string`
Конвертирует целое число в слова с указанием рода.

Доступные роды:
- `GenderMasculine` (1) - мужской род
- `GenderFeminine` (2) - женский род
- `GenderNeuter` (3) - средний род

```go
propisyu.IntToWordsGender(2, propisyu.GenderMasculine)  // "два"
propisyu.IntToWordsGender(2, propisyu.GenderFeminine)   // "две"
propisyu.IntToWordsGender(1, propisyu.GenderNeuter)     // "одно"
```

### Функции для десятичных чисел

#### `DecimalToWords(decimalStr string) (string, error)`
Конвертирует десятичное число из строки в слова. Дробная часть обрезается до 2 знаков.

```go
result, err := propisyu.DecimalToWords("3.14")
// "три целых четырнадцать сотых"
```

#### `DecimalValueToWords(d decimal.Decimal) (string, error)`
Конвертирует `decimal.Decimal` значение в слова. Дробная часть обрезается (не округляется!) до 2 знаков.

```go
d := decimal.NewFromFloat(3.14159)
result, err := propisyu.DecimalValueToWords(d)
// "три целых четырнадцать сотых" (не округляет 3.14159 до 3.14, а обрезает)
```

**Важно:** Функция обрезает (truncate), а не округляет:
- `1.999` → "один целых девяносто девять сотых"
- `1.995` → "один целых девяносто девять сотых"

### Функция склонения

#### `Decline(n int, one, two, five string) string`
Выбирает правильную форму существительного в зависимости от числа.

Параметры:
- `n` - число
- `one` - форма для 1, 21, 31... (рубль, день, товар)
- `two` - форма для 2-4, 22-24... (рубля, дня, товара)
- `five` - форма для 0, 5-20, 25-30... (рублей, дней, товаров)

```go
propisyu.Decline(1, "рубль", "рубля", "рублей")   // "рубль"
propisyu.Decline(2, "рубль", "рубля", "рублей")   // "рубля"
propisyu.Decline(5, "рубль", "рубля", "рублей")   // "рублей"
propisyu.Decline(11, "рубль", "рубля", "рублей")  // "рублей"
propisyu.Decline(21, "рубль", "рубля", "рублей")  // "рубль"
```

## Ограничения

- Целые числа: поддерживаются значения в диапазоне `int` (обычно -2³¹ до 2³¹-1 или -2⁶³ до 2⁶³-1)
- Десятичные числа: поддерживаются только 2 знака после запятой (остальное обрезается)
- `DecimalValueToWords` вернет ошибку `ErrNumberTooLarge`, если число не помещается в `int`

## Тесты

```bash
go test ./...              # Запустить все тесты
go test -v ./...           # С подробным выводом
go test -cover ./...       # С покрытием кода
```

## Лицензия

MIT

---

**SEO keywords:** Go Russian number to words, конвертер чисел в слова, Russian declension library, golang number spelling, propisyu Go package, decimal to words, число прописью Go.
