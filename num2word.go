package propisyu

import (
	"errors"
	"fmt"
	"math"
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

// IntToWords converts an integer to its Russian word representation
// using masculine gender by default.
//
//	propisyu.IntToWords(42)   // "сорок два"
//	propisyu.IntToWords(0)    // "ноль"
//	propisyu.IntToWords(-5)   // "минус пять"
//	propisyu.IntToWords(1000) // "одна тысяча"
func IntToWords(n int) string {
	return convertIntToWords(n, newDictionary(GenderMasculine))
}

// IntToWordsGender converts an integer to its Russian word representation
// with the specified grammatical gender. Gender affects the forms of
// "один"/"одна"/"одно" and "два"/"две".
//
//	propisyu.IntToWordsGender(1, propisyu.GenderMasculine) // "один"
//	propisyu.IntToWordsGender(1, propisyu.GenderFeminine)  // "одна"
//	propisyu.IntToWordsGender(1, propisyu.GenderNeuter)    // "одно"
func IntToWordsGender(n int, gender Gender) string {
	return convertIntToWords(n, newDictionary(gender))
}

// DecimalToWords converts a decimal number string to Russian words.
// The input should be a string like "123.45" or "6453345242432.42".
// Returns the number in Russian with proper declensions, e.g.
// "сто двадцать три целых сорок пять сотых".
// Only the first two decimal places are used; additional digits are truncated.
//
// Leading and trailing whitespace in the input is ignored, so
// " 123.45\n" is accepted just like "123.45".
func DecimalToWords(decimalStr string) (string, error) {
	decimalStr = strings.TrimSpace(decimalStr)
	isNegative := strings.HasPrefix(decimalStr, "-")

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
		"%s %s %s %s",
		IntToWordsGender(whole, GenderFeminine),
		Decline(whole, "целая", "целых", "целых"),
		IntToWordsGender(hundredths, GenderFeminine),
		Decline(hundredths, "сотая", "сотых", "сотых"),
	)

	if isNegative && whole == 0 && hundredths > 0 {
		result = "минус " + result
	}

	return result, nil
}

// DecimalValueToWords converts a decimal.Decimal value to Russian words.
// The input should be a decimal.Decimal value like decimal.NewFromFloat(123.45).
// Returns the number in Russian with proper declensions, e.g.
// "сто двадцать три целых сорок пять сотых".
// Only the first two decimal places are used; additional digits are truncated.
//
// Returns ErrNumberTooLarge when the whole part does not fit in Go int.
// shopspring/decimal can hold arbitrary-precision values, but this package
// renders through the int-based pipeline used by IntToWordsGender.
func DecimalValueToWords(d decimal.Decimal) (string, error) {
	// Check for int64 fit BEFORE calling IntPart(). decimal.Decimal.IntPart()
	// silently wraps on overflow (e.g. 10^29 → 7886392056514347007), so the
	// previous `whole > math.MaxInt64` guard never fired for values that had
	// already been truncated by IntPart — we have to look at the original
	// arbitrary-precision integer representation.
	truncBig := d.Truncate(0).BigInt()
	if !truncBig.IsInt64() {
		return "", fmt.Errorf("%w: %s", ErrNumberTooLarge, truncBig.String())
	}
	whole := truncBig.Int64()
	// On 32-bit platforms int is narrower than int64. The core conversion
	// pipeline takes int, so re-check the int fit explicitly instead of
	// relying on the silent int64 → int narrowing below.
	if whole > int64(math.MaxInt) || whole < int64(math.MinInt) {
		return "", fmt.Errorf("%w: %d", ErrNumberTooLarge, whole)
	}

	fractionalPart := d.Sub(decimal.NewFromInt(whole))

	hundredths := fractionalPart.Mul(decimal.NewFromInt(100)).Abs().Truncate(0).IntPart()

	result := fmt.Sprintf(
		"%s %s %s %s",
		IntToWordsGender(int(whole), GenderFeminine),
		Decline(int(whole), "целая", "целых", "целых"),
		IntToWordsGender(int(hundredths), GenderFeminine),
		Decline(int(hundredths), "сотая", "сотых", "сотых"),
	)

	// When the whole part is zero, IntPart() drops the sign, so "-0.50"
	// would otherwise render the same as "0.50". Preserve the minus to
	// stay consistent with DecimalToWords(`-0.xx`).
	if d.IsNegative() && whole == 0 && hundredths > 0 {
		result = "минус " + result
	}

	return result, nil
}

func convertIntToWords(n int, dict *dictionary) string {
	if n == 0 {
		return "ноль"
	}

	if n < 0 {
		// For math.MinInt, -n would overflow int (|math.MinInt| = math.MaxInt+1
		// does not fit in a signed int). Route the magnitude through uint64
		// instead. math.MaxInt is a positive compile-time constant, so the
		// conversion to uint64 is statically safe.
		if n == math.MinInt {
			return "минус " + convertPositiveUint64ToWords(uint64(math.MaxInt)+1, dict)
		}
		return "минус " + convertIntToWords(-n, dict)
	}

	return convertPositiveUint64ToWords(uint64(n), dict)
}

func convertPositiveUint64ToWords(n uint64, dict *dictionary) string {
	if n == 0 {
		return "ноль"
	}

	// uint64 max is ~1.8e19, which is 7 triads. dict.orders goes up to
	// 13 (duodecillion). Preallocate enough for the worst case so the
	// append loop never has to grow, and build the slice in natural
	// (least-significant-first) order — we reverse in place at the end
	// rather than prepending each iteration (O(n²) → O(n)).
	parts := make([]string, 0, len(dict.orders))
	order := 0

	for n > 0 {
		// n%1000 is in [0, 999]; the narrowing cast is statically safe.
		// gosec's G115 is purely type-based and cannot see the value range,
		// so suppress it with gosec's native #nosec directive (inline
		// `//nolint:gosec` is not reliably propagated by golangci-lint v1).
		triad := int(n % 1000) //#nosec G115 -- bounded to [0, 999]
		n /= 1000

		if triad != 0 {
			triadWords := dict.triadToWords(triad, order)
			if order > 0 && order < len(dict.orders) {
				forms := dict.orders[order]
				triadWords += " " + getDeclension(triad%100, forms[0], forms[1], forms[2])
			}
			parts = append(parts, triadWords)
		}
		order++
	}

	// Reverse in place so the most significant triad comes first.
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
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

// DecimalToWordsPrecision converts a decimal string with specified precision (1-9).
// Precision 1 = десятые, 2 = сотые, 3 = тысячные, etc.
// The fractional part is always in feminine gender.
//
//	DecimalToWordsPrecision("3.14", 2)  // "три целых четырнадцать сотых"
//	DecimalToWordsPrecision("3.5", 1)   // "три целых пять десятых"
//
// Leading and trailing whitespace in the input is ignored.
func DecimalToWordsPrecision(decimalStr string, precision int) (string, error) {
	if precision < 1 || precision > 9 {
		return "", fmt.Errorf("precision must be between 1 and 9, got %d", precision)
	}

	fractionUnits := [][3]string{
		{"десятая", "десятых", "десятых"},                            // 1
		{"сотая", "сотых", "сотых"},                                  // 2
		{"тысячная", "тысячных", "тысячных"},                         // 3
		{"десятитысячная", "десятитысячных", "десятитысячных"},       // 4
		{"стотысячная", "стотысячных", "стотысячных"},                // 5
		{"миллионная", "миллионных", "миллионных"},                   // 6
		{"десятимиллионная", "десятимиллионных", "десятимиллионных"}, // 7
		{"стомиллионная", "стомиллионных", "стомиллионных"},          // 8
		{"миллиардная", "миллиардных", "миллиардных"},                // 9
	}

	decimalStr = strings.TrimSpace(decimalStr)
	isNegative := strings.HasPrefix(decimalStr, "-")

	parts := strings.SplitN(decimalStr, ".", 2)

	whole, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid whole number part: %w", err)
	}

	fracStr := "0"
	if len(parts) == 2 {
		fracStr = parts[1]
	}

	// Truncate or pad to the specified precision
	if len(fracStr) > precision {
		fracStr = fracStr[:precision]
	}
	for len(fracStr) < precision {
		fracStr += "0"
	}

	fracVal, err := strconv.Atoi(fracStr)
	if err != nil {
		return "", fmt.Errorf("invalid fractional part: %w", err)
	}

	units := fractionUnits[precision-1]

	result := fmt.Sprintf(
		"%s %s %s %s",
		IntToWordsGender(whole, GenderFeminine),
		Decline(whole, "целая", "целых", "целых"),
		IntToWordsGender(fracVal, GenderFeminine),
		Decline(fracVal, units[0], units[1], units[2]),
	)

	// Same rationale as DecimalToWords: for inputs like "-0.5" the minus
	// sign would be lost because strconv.Atoi("0") → 0 and IntToWordsGender
	// drops it. Preserve it explicitly.
	if isNegative && whole == 0 && fracVal > 0 {
		result = "минус " + result
	}

	return result, nil
}

func clampGender(g Gender) int {
	switch g {
	case GenderMasculine, GenderFeminine, GenderNeuter:
		return int(g)
	default:
		return int(GenderMasculine)
	}
}
