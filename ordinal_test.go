package propisyu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrdinal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		n      int
		gender Gender
		want   string
	}{
		{name: "1 masculine", n: 1, gender: GenderMasculine, want: "первый"},
		{name: "1 feminine", n: 1, gender: GenderFeminine, want: "первая"},
		{name: "1 neuter", n: 1, gender: GenderNeuter, want: "первое"},
		{name: "2 masculine", n: 2, gender: GenderMasculine, want: "второй"},
		{name: "3 feminine", n: 3, gender: GenderFeminine, want: "третья"},
		{name: "7 masculine", n: 7, gender: GenderMasculine, want: "седьмой"},
		{name: "10 masculine", n: 10, gender: GenderMasculine, want: "десятый"},
		{name: "11 masculine", n: 11, gender: GenderMasculine, want: "одиннадцатый"},
		{name: "20 masculine", n: 20, gender: GenderMasculine, want: "двадцатый"},
		{name: "21 masculine", n: 21, gender: GenderMasculine, want: "двадцать первый"},
		{name: "42 masculine", n: 42, gender: GenderMasculine, want: "сорок второй"},
		{name: "42 feminine", n: 42, gender: GenderFeminine, want: "сорок вторая"},
		{name: "100 masculine", n: 100, gender: GenderMasculine, want: "сотый"},
		{name: "101 masculine", n: 101, gender: GenderMasculine, want: "сто первый"},
		{name: "200 masculine", n: 200, gender: GenderMasculine, want: "двухсотый"},
		{name: "500 feminine", n: 500, gender: GenderFeminine, want: "пятисотая"},
		{name: "1000 masculine", n: 1000, gender: GenderMasculine, want: "тысячный"},
		{name: "2000 masculine", n: 2000, gender: GenderMasculine, want: "двухтысячный"},
		{name: "1000000 masculine", n: 1000000, gender: GenderMasculine, want: "миллионный"},
		{name: "2026 masculine", n: 2026, gender: GenderMasculine, want: "две тысячи двадцать шестой"},
		{name: "zero masculine", n: 0, gender: GenderMasculine, want: "нулевой"},
		{name: "negative masculine", n: -1, gender: GenderMasculine, want: "минус первый"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Ordinal(tc.n, tc.gender)
			assert.Equal(t, tc.want, got)
		})
	}
}
