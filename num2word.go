// Package propisyu provides helpers for converting integers and decimal numbers
// to Russian words with correct grammatical gender and noun declension.
//
// Supported range: integers up to duodecillions (10³⁹), decimals up to 2 decimal places.
//
// Quick start:
//
//	propisyu.IntToWords(321)                                    // "триста двадцать один"
//	propisyu.IntToWordsGender(2, propisyu.GenderFeminine)       // "две"
//	propisyu.Decline(5, "рубль", "рубля", "рублей")            // "рублей"
package propisyu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	ErrNumberTooLarge = errors.New("number is too large to convert")
)

type Gender int

const (
	// GenderMasculine applies masculine endings (e.g. "один").
	GenderMasculine Gender = 1
	// GenderFeminine applies feminine endings (e.g. "одна").
	GenderFeminine Gender = 2
	// GenderNeuter applies neuter endings (e.g. "одно").
	GenderNeuter Gender = 3
)

type dictionary struct {
	ones       [][4]string
	tens       []string
	teens      []string
	orders     [][]string
	baseGender int
}

func newDictionary(baseGender Gender) *dictionary {
	return &dictionary{
		ones: [][4]string{
			{0: "", 1: "один", 2: "одна", 3: "одно"},
			{0: "", 1: "два", 2: "две", 3: "два"},
			{0: "три", 1: "три", 2: "три", 3: "три"},
			{0: "четыре", 1: "четыре", 2: "четыре", 3: "четыре"},
			{0: "пять", 1: "пять", 2: "пять", 3: "пять"},
			{0: "шесть", 1: "шесть", 2: "шесть", 3: "шесть"},
			{0: "семь", 1: "семь", 2: "семь", 3: "семь"},
			{0: "восемь", 1: "восемь", 2: "восемь", 3: "восемь"},
			{0: "девять", 1: "девять", 2: "девять", 3: "девять"},
		},
		tens: []string{
			"", "десять", "двадцать", "тридцать", "сорок", "пятьдесят",
			"шестьдесят", "семьдесят", "восемьдесят", "девяносто",
		},
		teens: []string{
			"десять", "одиннадцать", "двенадцать", "тринадцать", "четырнадцать",
			"пятнадцать", "шестнадцать", "семнадцать", "восемнадцать", "девятнадцать",
		},
		orders: [][]string{
			{"", "", ""}, // thousands, millions, etc.
			{"тысяча", "тысячи", "тысяч"},                       // 10^3
			{"миллион", "миллиона", "миллионов"},                // 10^6
			{"миллиард", "миллиарда", "миллиардов"},             // 10^9
			{"триллион", "триллиона", "триллионов"},             // 10^12
			{"квадриллион", "квадриллиона", "квадриллионов"},    // 10^15
			{"квинтиллион", "квинтиллиона", "квинтиллионов"},    // 10^18
			{"секстиллион", "секстиллиона", "секстиллионов"},    // 10^21
			{"септиллион", "септиллиона", "септиллионов"},       // 10^24
			{"октиллион", "октиллиона", "октиллионов"},          // 10^27
			{"нониллион", "нониллиона", "нониллионов"},          // 10^30
			{"дециллион", "дециллиона", "дециллионов"},          // 10^33
			{"ундециллион", "ундециллиона", "ундециллионов"},    // 10^36
			{"дуодециллион", "дуодециллиона", "дуодециллионов"}, // 10^39
		},
		baseGender: clampGender(baseGender),
	}
}

func IntToWords(n int) string {
	return convertIntToWords(n, newDictionary(GenderMasculine))
}

func IntToWordsGender(n int, gender Gender) string {
	return convertIntToWords(n, newDictionary(gender))
}

// DecimalToWords converts a decimal number string to Russian words.
// The input should be a string like "123.45" or "6453345242432.42".
// Returns the number in Russian with proper declensions, e.g.
// "сто двадцать три целых сорок пять сотых".
// Only the first two decimal places are used; additional digits are truncated.
func DecimalToWords(decimalStr string) (string, error) {
	parts := strings.SplitN(decimalStr, ".", 2)

	whole, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid whole number part: %w", err)
	}

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

	hundredths, err := strconv.Atoi(fraction)
	if err != nil {
		return "", fmt.Errorf("invalid fractional part: %w", err)
	}

	result := fmt.Sprintf(
		"%s целых %s %s",
		IntToWords(whole),
		IntToWordsGender(hundredths, GenderFeminine),
		Decline(hundredths, "сотая", "сотых", "сотых"),
	)

	return result, nil
}

// DecimalValueToWords converts a decimal.Decimal value to Russian words.
// The input should be a decimal.Decimal value like decimal.NewFromFloat(123.45).
// Returns the number in Russian with proper declensions, e.g.
// "сто двадцать три целых сорок пять сотых".
// Only the first two decimal places are used; additional digits are truncated.
func DecimalValueToWords(d decimal.Decimal) (string, error) {
	whole := d.IntPart()

	if whole > int64(^uint(0)>>1) || whole < -int64(^uint(0)>>1)-1 {
		return "", fmt.Errorf("%w: %v", ErrNumberTooLarge, whole)
	}

	fractionalPart := d.Sub(decimal.NewFromInt(whole))

	hundredths := fractionalPart.Mul(decimal.NewFromInt(100)).Abs().Truncate(0).IntPart()

	result := fmt.Sprintf(
		"%s целых %s %s",
		IntToWords(int(whole)),
		IntToWordsGender(int(hundredths), GenderFeminine),
		Decline(int(hundredths), "сотая", "сотых", "сотых"),
	)

	return result, nil
}

func convertIntToWords(n int, dict *dictionary) string {
	if n == 0 {
		return "ноль"
	}

	if n < 0 {
		return "минус " + convertIntToWords(-n, dict)
	}

	var parts []string
	order := 0

	for n > 0 {
		triad := n % 1000
		n /= 1000

		if triad != 0 {
			triadWords := dict.triadToWords(triad, order)
			if order > 0 && order < len(dict.orders) {
				forms := dict.orders[order]
				triadWords += " " + getDeclension(triad%100, forms[0], forms[1], forms[2])
			}
			parts = append([]string{triadWords}, parts...)
		}
		order++
	}

	return strings.Join(parts, " ")
}

func (d *dictionary) triadToWords(n, order int) string {
	var s []string

	h := n / 100
	t := (n % 100) / 10
	o := n % 10

	if h > 0 {
		s = append(s, []string{
			"сто", "двести", "триста", "четыреста", "пятьсот", "шестьсот", "семьсот", "восемьсот", "девятьсот",
		}[h-1])
	}

	if t == 1 {
		s = append(s, d.teens[o])
	} else {
		if t > 0 {
			s = append(s, d.tens[t])
		}
		if o > 0 {
			form := 1
			if order == 0 {
				form = d.baseGender
			}
			if form < 1 || form > 3 {
				form = 1
			}
			if order == 1 { // тысяча
				form = 2 // feminine
			}
			s = append(s, d.ones[o-1][form])
		}
	}

	if len(s) == 0 {
		return ""
	}
	return strings.Join(s, " ")
}

func clampGender(g Gender) int {
	switch g {
	case GenderMasculine, GenderFeminine, GenderNeuter:
		return int(g)
	default:
		return int(GenderMasculine)
	}
}
