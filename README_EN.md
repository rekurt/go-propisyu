# go-propisyu · Go Library for Russian Number Spelling

**English version · [Русская версия](README.md)**

`go-propisyu` converts integers into Russian words and applies the correct grammatical declensions. The package is ideal for invoices, receipts, voice prompts, document generators, fintech dashboards, and any interface that must spell numbers in fluent Russian.

## Highlights

- 🔢 Handles gigantic numbers up to duodecillions (10³⁹) with no external dependencies.
- 🧠 Supports masculine, feminine, and neuter genders via `IntToWordsGender`.
- 💬 Provides `Decline` to choose noun forms for currencies, measurements, and custom labels.

## Installation

```bash
go get github.com/rekurt/go-propisyu
```

## Quick Start

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

## Example: 6,453,345,242,432.42

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
	// шесть триллионов четыреста пятьдесят три миллиарда триста сорок пять миллионов
	// двести сорок две тысячи четыреста тридцать два целых сорок две сотых
}
```
---

**SEO keywords:** go-propisyu, Go Russian number to words, Russian number declension, golang number spelling library, Russian currency declension Go.
