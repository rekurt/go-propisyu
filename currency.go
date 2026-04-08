package propisyu

import (
	"fmt"
	"strconv"
	"strings"
)

// Currency describes a currency for formatting amounts in words.
type Currency struct {
	WholeOne  string // "рубль"
	WholeTwo  string // "рубля"
	WholeFive string // "рублей"
	FracOne   string // "копейка"
	FracTwo   string // "копейки"
	FracFive  string // "копеек"

	WholeGender Gender // GenderMasculine
	FracGender  Gender // GenderFeminine
}

//nolint:gochecknoglobals // exported preset Currency values; const not supported for structs
var (
	CurrencyRUB = Currency{
		WholeOne: "рубль", WholeTwo: "рубля", WholeFive: "рублей", WholeGender: GenderMasculine,
		FracOne: "копейка", FracTwo: "копейки", FracFive: "копеек", FracGender: GenderFeminine,
	}
	CurrencyUSD = Currency{
		WholeOne: "доллар", WholeTwo: "доллара", WholeFive: "долларов", WholeGender: GenderMasculine,
		FracOne: "цент", FracTwo: "цента", FracFive: "центов", FracGender: GenderMasculine,
	}
	CurrencyEUR = Currency{
		WholeOne: "евро", WholeTwo: "евро", WholeFive: "евро", WholeGender: GenderNeuter,
		FracOne: "цент", FracTwo: "цента", FracFive: "центов", FracGender: GenderMasculine,
	}
)

// Money formats an amount as words with currency.
// whole is the integer part, cents is the fractional part (0-99).
// Example: Money(1234, 56, &CurrencyRUB) returns
// "одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек"
func Money(whole int, cents int, c *Currency) string {
	wholeWords := IntToWordsGender(whole, c.WholeGender)
	wholeDecl := Decline(whole, c.WholeOne, c.WholeTwo, c.WholeFive)
	centsWords := IntToWordsGender(cents, c.FracGender)
	centsDecl := Decline(cents, c.FracOne, c.FracTwo, c.FracFive)
	return wholeWords + " " + wholeDecl + " " + centsWords + " " + centsDecl
}

// MoneyFromString parses "1234.56" and formats with currency.
func MoneyFromString(amount string, c *Currency) (string, error) {
	parts := strings.SplitN(amount, ".", 2)
	whole, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid amount: %w", err)
	}
	cents := 0
	if len(parts) == 2 {
		frac := parts[1]
		if len(frac) > 2 {
			frac = frac[:2]
		}
		for len(frac) < 2 {
			frac += "0"
		}
		cents, err = strconv.Atoi(frac)
		if err != nil {
			return "", fmt.Errorf("invalid fractional part: %w", err)
		}
	}
	return Money(whole, cents, c), nil
}
