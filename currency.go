package propisyu

import (
	"fmt"
	"strconv"
	"strings"
)

// Currency describes a currency for formatting amounts in words.
type Currency struct { //nolint:govet // fieldalignment: preserve exported field order for API compatibility
	WholeOne    string // "рубль"
	WholeTwo    string // "рубля"
	WholeFive   string // "рублей"
	WholeGender Gender // GenderMasculine
	FracOne     string // "копейка"
	FracTwo     string // "копейки"
	FracFive    string // "копеек"
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
// Example: Money(1234, 56, CurrencyRUB) returns
// "одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек"
func Money(whole, cents int, c Currency) string { //nolint:gocritic // hugeParam: preserve public API compatibility
	wholeWords := IntToWordsGender(whole, c.WholeGender)
	wholeDecl := Decline(whole, c.WholeOne, c.WholeTwo, c.WholeFive)
	centsWords := IntToWordsGender(cents, c.FracGender)
	centsDecl := Decline(cents, c.FracOne, c.FracTwo, c.FracFive)
	return wholeWords + " " + wholeDecl + " " + centsWords + " " + centsDecl
}

// MoneyFromString parses "1234.56" and formats with currency.
//
// The minus sign is preserved even for amounts where the whole part rounds
// to zero and the fractional part is non-zero, so "-0.50" renders as
// "минус ноль рублей пятьдесят копеек". Without this guard the sign would
// be lost, because strconv.Atoi("-0") == 0 and Money(0, 50, …) carries
// no sign information.
func MoneyFromString(amount string, c Currency) (string, error) { //nolint:gocritic // hugeParam: preserve public API compatibility
	isNegative := strings.HasPrefix(amount, "-")

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

	result := Money(whole, cents, c)
	if isNegative && whole == 0 && cents > 0 {
		result = "минус " + result
	}
	return result, nil
}
