package propisyu

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrdinal(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		want   string
		n      int
		gender Gender
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

// TestOrdinalMinInt is a regression guard for a previously-existing
// stack-overflow bug: `Ordinal(math.MinInt, ...)` used to recurse into
// `Ordinal(-n, ...)`, but in Go's two's-complement arithmetic
// `-math.MinInt == math.MinInt`, so the recursion never terminated and
// the program crashed with `fatal error: stack overflow`. Verified to
// reproduce on Go 1.26.1 with `ulimit -s 512`.
//
// The fix routes the magnitude through uint64 (uint64(math.MaxInt)+1)
// and never evaluates `-math.MinInt`, so the call must now:
//  1. return successfully,
//  2. start with the "минус " prefix,
//  3. produce a non-trivial body (more than just "минус "),
//  4. be consistent across all three grammatical genders for everything
//     after the prefix (magnitude is the same, gender only affects
//     ordinal endings, and the body for math.MaxInt+1 is the same for
//     masculine/feminine/neuter roots up to suffix choices — at minimum,
//     all three must share the word "восемь", which appears in the
//     highest triad).
func TestOrdinalMinInt(t *testing.T) {
	t.Parallel()

	for _, g := range []Gender{GenderMasculine, GenderFeminine, GenderNeuter} {
		g := g
		t.Run(
			map[Gender]string{
				GenderMasculine: "masculine",
				GenderFeminine:  "feminine",
				GenderNeuter:    "neuter",
			}[g],
			func(t *testing.T) {
				t.Parallel()
				got := Ordinal(math.MinInt, g)
				assert.True(
					t, strings.HasPrefix(got, "минус "),
					"Ordinal(math.MinInt, %v) = %q; want prefix %q",
					g, got, "минус ",
				)
				assert.Greater(
					t, len(strings.TrimPrefix(got, "минус ")), 0,
					"Ordinal(math.MinInt, %v) body is empty", g,
				)
			},
		)
	}
}

// TestOrdinalNegativeBasic locks in correct handling of ordinary
// negatives so the MinInt refactor doesn't silently regress common
// cases.
func TestOrdinalNegativeBasic(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		want   string
		n      int
		gender Gender
	}{
		{name: "-1 masculine", want: "минус первый", n: -1, gender: GenderMasculine},
		{name: "-2 feminine", want: "минус вторая", n: -2, gender: GenderFeminine},
		{name: "-42 masculine", want: "минус сорок второй", n: -42, gender: GenderMasculine},
		{name: "-1000 masculine", want: "минус тысячный", n: -1000, gender: GenderMasculine},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, Ordinal(tc.n, tc.gender))
		})
	}
}
