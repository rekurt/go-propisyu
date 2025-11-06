# go-propisyu · Russian Number-to-Words Converter

**Русская версия · [English version](README_EN.md)**

`go-propisyu` — Go-библиотека, которая переводит числа в русские слова и подбирает правильные склонения. Решение подходит для счетов, фискальных чеков, голосовых ассистентов, генерации документов и любых сервисов, где важно грамотно проговаривать суммы и количества.

## Ключевые особенности

- 🔢 Поддержка огромных чисел вплоть до дуодециллионов (10³⁹) без сторонних зависимостей.
- 🧠 Грамматические роды: мужской, женский и средний для корректных окончаний единиц.
- 💬 Функция `Decline` для валют, единиц измерения и любых существительных.

## Установка

```bash
go get github.com/rekurt/go-propisyu
```

## Быстрый старт

```go
package main

import (
	"fmt"

	"github.com/rekurt/go-propisyu"
)

func main() {
	fmt.Println(propisyu.IntToWords(321)) // триста двадцать один
	fmt.Println(propisyu.Decline(5, "рубль", "рубля", "рублей")) // рублей
}
```

## Пример: 6 453 345 242 432,42

```go
package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rekurt/go-propisyu"
)

func main() {
	const raw = "6453345242432.42"

	parts := strings.SplitN(raw, ".", 2)
	whole, _ := strconv.Atoi(parts[0])

	fraction := "00"
	if len(parts) == 2 {
		fraction = parts[1]
		if len(fraction) > 2 {
			fraction = fraction[:2]
		}
		for len(fraction) < 2 {
			fraction += "0"
		}
	}

	hundredths, _ := strconv.Atoi(fraction)

	result := fmt.Sprintf(
		"%s целых %s %s",
		propisyu.IntToWords(whole),
		propisyu.IntToWordsGender(hundredths, propisyu.GenderFeminine),
		propisyu.Decline(hundredths, "сотая", "сотых", "сотых"),
	)

	fmt.Println(result)
	// Output:
	// шесть триллионов четыреста пятьдесят три миллиарда
	// триста сорок пять миллионов двести сорок две тысячи четыреста тридцать два
	// целых сорок две сотых
}
```
---

**SEO keywords:** Go Russian number to words, конвертер чисел в слова, Russian declension library, golang number spelling, propisyu Go package.
