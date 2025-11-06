package propisyu

func Decline(n int, one, two, five string) string {
	return getDeclension(n, one, two, five)
}
func getDeclension(n int, one, two, five string) string {
	n %= 100
	if n >= 11 && n <= 19 {
		return five
	}
	switch n % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return two
	default:
		return five
	}
}
