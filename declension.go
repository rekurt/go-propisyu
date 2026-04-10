package propisyu

// Decline returns the correct Russian noun declension form for the given number.
// It follows standard Russian pluralization rules:
//   - one: used for 1, 21, 31, ... (e.g. "рубль", "день", "товар")
//   - two: used for 2-4, 22-24, ... (e.g. "рубля", "дня", "товара")
//   - five: used for 0, 5-20, 25-30, ... (e.g. "рублей", "дней", "товаров")
//
// Example:
//
//	propisyu.Decline(1, "рубль", "рубля", "рублей")  // "рубль"
//	propisyu.Decline(5, "рубль", "рубля", "рублей")  // "рублей"
//	propisyu.Decline(21, "рубль", "рубля", "рублей") // "рубль"
func Decline(n int, one, two, five string) string {
	return getDeclension(n, one, two, five)
}
func getDeclension(n int, one, two, five string) string {
	var nAbs uint64
	if n < 0 {
		nAbs = uint64(-(n + 1))
		nAbs++
	} else {
		nAbs = uint64(n)
	}

	nMod100 := int(nAbs % 100)
	if nMod100 >= 11 && nMod100 <= 19 {
		return five
	}
	switch nMod100 % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return two
	default:
		return five
	}
}
