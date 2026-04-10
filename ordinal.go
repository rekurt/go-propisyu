package propisyu

import "strings"

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
		return "минус " + Ordinal(-n, gender)
	}

	lastTriad := n % 1000
	prefix := n / 1000

	if lastTriad > 0 {
		// Convert all preceding triads as cardinal, then make last triad ordinal
		var parts []string
		if prefix > 0 {
			cardinalPrefix := convertHighTriads(prefix)
			parts = append(parts, cardinalPrefix)
		}
		parts = append(parts, ordinalTriad(lastTriad, gender))
		return strings.Join(parts, " ")
	}

	// lastTriad == 0: round order number (like 1000, 2000000, etc.)
	return roundOrdinal(prefix, gender)
}

// convertHighTriads converts the higher triads (everything above the last triad)
// into cardinal words. The input is the number with the last triad removed (n/1000),
// so it represents multiples of thousands.
func convertHighTriads(prefix int) string {
	// prefix represents the number of thousands
	// We need to convert prefix * 1000 as cardinal, but only the prefix part
	// using the full number conversion with appropriate thousand forms
	dict := newDictionary(GenderMasculine)
	return convertIntToWords(prefix*1000, dict)
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
