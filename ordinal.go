package propisyu

import (
	"math"
	"strings"
)

// genderIndex converts Gender to an index 0=masculine, 1=feminine, 2=neuter.
func genderIndex(g Gender) int {
	switch g {
	case GenderFeminine:
		return 1
	case GenderNeuter:
		return 2
	default:
		return 0
	}
}

// Ordinal converts a number to its ordinal form in Russian.
//
//	Ordinal(1, GenderMasculine)  // "первый"
//	Ordinal(1, GenderFeminine)   // "первая"
//	Ordinal(42, GenderMasculine) // "сорок второй"
//	Ordinal(1000, GenderMasculine) // "тысячный"
func Ordinal(n int, gender Gender) string {
	gi := genderIndex(gender)

	if n == 0 {
		zeroOrd := [3]string{"нулевой", "нулевая", "нулевое"}
		return zeroOrd[gi]
	}

	if n < 0 {
		var magU uint64
		if n == math.MinInt {
			// |math.MinInt| == math.MaxInt + 1 does not fit in a signed int;
			// computing -n here would overflow back to math.MinInt and
			// recurse forever. Lift the magnitude through uint64 instead,
			// matching the approach used by convertIntToWords.
			magU = uint64(math.MaxInt) + 1
		} else {
			magU = uint64(-n) //#nosec G115 -- n<0 and n!=math.MinInt, so -n fits in int64 and is non-negative
		}
		return "минус " + ordinalFromMagnitude(magU, gender)
	}

	return ordinalFromMagnitude(uint64(n), gender)
}

// ordinalFromMagnitude produces the ordinal phrase for the absolute value
// of the original input. magU is in [1, uint64(math.MaxInt)+1]:
//   - magU % 1000 is always in [0, 999] and fits in int trivially.
//   - magU / 1000 is at most (math.MaxInt + 1) / 1000, which is still
//     below math.MaxInt and therefore fits in a signed int.
//
// These bounds justify the int narrowings below.
func ordinalFromMagnitude(magU uint64, gender Gender) string {
	lastTriad := int(magU % 1000) //#nosec G115 -- bounded to [0, 999]
	prefixU := magU / 1000

	if lastTriad > 0 {
		// Convert all preceding triads as cardinal, then make last triad ordinal
		var parts []string
		if prefixU > 0 {
			parts = append(parts, convertHighTriadsU(prefixU))
		}
		parts = append(parts, ordinalTriad(lastTriad, gender))
		return strings.Join(parts, " ")
	}

	// lastTriad == 0: round order number (like 1000, 2000000, etc.).
	// roundOrdinal takes int (prefix); prefixU here is bounded to
	// math.MaxInt/1000+1 at worst (from math.MinInt magnitude lift), which
	// still fits in a signed int on any supported platform.
	return roundOrdinal(int(prefixU), gender) //#nosec G115 -- bounded above by (math.MaxInt+1)/1000, fits in int
}

// convertHighTriadsU converts the higher triads (everything above the last
// triad) into cardinal words. The input is the magnitude with its last triad
// removed (magU/1000), so it represents multiples of thousands.
//
// This is uint64-native on purpose: when Ordinal is called on math.MinInt,
// prefixU can be as large as (math.MaxInt+1)/1000, and multiplying that by
// 1000 in int would overflow. We stay in uint64 and delegate to
// convertPositiveUint64ToWords, which is safe for the full range.
func convertHighTriadsU(prefixU uint64) string {
	dict := newDictionary(GenderMasculine)
	return convertPositiveUint64ToWords(prefixU*1000, dict)
}

// ordinalTriad converts the last triad (1-999) into ordinal form.
// It uses cardinal forms for hundreds and tens as prefix, and ordinal for the final component.
func ordinalTriad(n int, gender Gender) string { //nolint:funlen // lookup data tables
	h := n / 100
	t := (n % 100) / 10
	o := n % 10
	gi := genderIndex(gender)

	ordinalOnes := [][3]string{
		{"первый", "первая", "первое"},
		{"второй", "вторая", "второе"},
		{"третий", "третья", "третье"},
		{"четвёртый", "четвёртая", "четвёртое"},
		{"пятый", "пятая", "пятое"},
		{"шестой", "шестая", "шестое"},
		{"седьмой", "седьмая", "седьмое"},
		{"восьмой", "восьмая", "восьмое"},
		{"девятый", "девятая", "девятое"},
	}
	ordinalTeens := [][3]string{
		{"десятый", "десятая", "десятое"},
		{"одиннадцатый", "одиннадцатая", "одиннадцатое"},
		{"двенадцатый", "двенадцатая", "двенадцатое"},
		{"тринадцатый", "тринадцатая", "тринадцатое"},
		{"четырнадцатый", "четырнадцатая", "четырнадцатое"},
		{"пятнадцатый", "пятнадцатая", "пятнадцатое"},
		{"шестнадцатый", "шестнадцатая", "шестнадцатое"},
		{"семнадцатый", "семнадцатая", "семнадцатое"},
		{"восемнадцатый", "восемнадцатая", "восемнадцатое"},
		{"девятнадцатый", "девятнадцатая", "девятнадцатое"},
	}
	ordinalTens := [][3]string{
		{"двадцатый", "двадцатая", "двадцатое"},
		{"тридцатый", "тридцатая", "тридцатое"},
		{"сороковой", "сороковая", "сороковое"},
		{"пятидесятый", "пятидесятая", "пятидесятое"},
		{"шестидесятый", "шестидесятая", "шестидесятое"},
		{"семидесятый", "семидесятая", "семидесятое"},
		{"восьмидесятый", "восьмидесятая", "восьмидесятое"},
		{"девяностый", "девяностая", "девяностое"},
	}
	ordinalHundreds := [][3]string{
		{"сотый", "сотая", "сотое"},
		{"двухсотый", "двухсотая", "двухсотое"},
		{"трёхсотый", "трёхсотая", "трёхсотое"},
		{"четырёхсотый", "четырёхсотая", "четырёхсотое"},
		{"пятисотый", "пятисотая", "пятисотое"},
		{"шестисотый", "шестисотая", "шестисотое"},
		{"семисотый", "семисотая", "семисотое"},
		{"восьмисотый", "восьмисотая", "восьмисотое"},
		{"девятисотый", "девятисотая", "девятисотое"},
	}
	hundredsCardinal := []string{
		"сто", "двести", "триста", "четыреста", "пятьсот", "шестьсот", "семьсот", "восемьсот", "девятьсот",
	}
	tensCardinal := []string{
		"", "десять", "двадцать", "тридцать", "сорок", "пятьдесят",
		"шестьдесят", "семьдесят", "восемьдесят", "девяносто",
	}

	var parts []string

	if t == 1 { // teen
		if h > 0 {
			parts = append(parts, hundredsCardinal[h-1])
		}
		parts = append(parts, ordinalTeens[o][gi])
	} else if o > 0 { // has ones
		if h > 0 {
			parts = append(parts, hundredsCardinal[h-1])
		}
		if t > 0 {
			parts = append(parts, tensCardinal[t])
		}
		parts = append(parts, ordinalOnes[o-1][gi])
	} else if t > 0 { // ends with tens
		if h > 0 {
			parts = append(parts, hundredsCardinal[h-1])
		}
		parts = append(parts, ordinalTens[t-2][gi])
	} else { // only hundreds
		parts = append(parts, ordinalHundreds[h-1][gi])
	}

	return strings.Join(parts, " ")
}

// roundOrdinal handles cases where the number is a round order (e.g. 1000, 2000, 5000000).
// prefix is n/1000 where n is the original number.
func roundOrdinal(prefix int, gender Gender) string {
	gi := genderIndex(gender)

	ordinalOrders := [][3]string{
		{"тысячный", "тысячная", "тысячное"},
		{"миллионный", "миллионная", "миллионное"},
		{"миллиардный", "миллиардная", "миллиардное"},
		{"триллионный", "триллионная", "триллионное"},
		{"квадриллионный", "квадриллионная", "квадриллионное"},
		{"квинтиллионный", "квинтиллионная", "квинтиллионное"},
		{"секстиллионный", "секстиллионная", "секстиллионное"},
		{"септиллионный", "септиллионная", "септиллионное"},
		{"октиллионный", "октиллионная", "октиллионное"},
		{"нониллионный", "нониллионная", "нониллионное"},
		{"дециллионный", "дециллионная", "дециллионное"},
		{"ундециллионный", "ундециллионная", "ундециллионное"},
		{"дуодециллионный", "дуодециллионная", "дуодециллионное"},
	}
	onesCompound := []string{
		"одно", "двух", "трёх", "четырёх", "пяти", "шести", "семи", "восьми", "девяти",
	}

	// Find the order level of the round number
	orderLevel := 0
	temp := prefix
	for temp > 0 && temp%1000 == 0 {
		temp /= 1000
		orderLevel++
	}

	actualOrder := orderLevel + 1 // order index (1=тысяча, 2=миллион, etc.)

	if actualOrder < 1 || actualOrder > len(ordinalOrders) {
		// Fallback: just use cardinal
		return IntToWords(prefix*1000) + "-й"
	}

	if temp == 1 {
		// Coefficient is 1: just use the order ordinal (тысячный, миллионный, etc.)
		return ordinalOrders[actualOrder-1][gi]
	}

	if temp >= 1 && temp <= 9 {
		// Simple coefficient: use compound prefix
		return onesCompound[temp-1] + ordinalOrders[actualOrder-1][gi]
	}

	// Coefficient > 9: use cardinal form + space + order ordinal
	cardinalCoeff := IntToWords(temp)
	return cardinalCoeff + " " + ordinalOrders[actualOrder-1][gi]
}
