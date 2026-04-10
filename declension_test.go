package propisyu

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type declineCase struct {
	description string
	forms       [3]string
	num         int
}

const (
	formOneIdx = iota
	formTwoIdx
	formFiveIdx
)

func (c declineCase) expected() string {
	switch {
	case c.num%100 >= 11 && c.num%100 <= 19:
		return c.forms[formFiveIdx]
	case c.num%10 == 1:
		return c.forms[formOneIdx]
	case c.num%10 >= 2 && c.num%10 <= 4:
		return c.forms[formTwoIdx]
	default:
		return c.forms[formFiveIdx]
	}
}

func TestIntToWordsBasic(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		want string
		num  int
	}{
		{name: "zero", num: 0, want: "ноль"},
		{name: "single digit", num: 7, want: "семь"},
		{name: "teens exact", num: 15, want: "пятнадцать"},
		{name: "tens exact", num: 40, want: "сорок"},
		{name: "two digits", num: 42, want: "сорок два"},
		{name: "hundreds", num: 305, want: "триста пять"},
		{name: "hundreds with tens", num: 512, want: "пятьсот двенадцать"},
		{name: "thousand boundary", num: 1000, want: "одна тысяча"},
		{name: "thousand with remainder", num: 2001, want: "две тысячи один"},
		{name: "complex", num: 987654, want: "девятьсот восемьдесят семь тысяч шестьсот пятьдесят четыре"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, IntToWords(tc.num))
		})
	}
}

func TestIntToWordsLargeNumbers(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		want string
		num  int
	}{
		{
			name: "one million",
			num:  1_000_000,
			want: "один миллион",
		},
		{
			name: "multi million",
			num:  21_304_015,
			want: "двадцать один миллион триста четыре тысячи пятнадцать",
		},
		{
			name: "one billion",
			num:  1_000_000_000,
			want: "один миллиард",
		},
		{
			name: "max int32",
			num:  2_147_483_647,
			want: "два миллиарда сто сорок семь миллионов четыреста восемьдесят три тысячи шестьсот сорок семь",
		},
		{
			name: "trillion scale",
			num:  6_453_345_242_432,
			want: "шесть триллионов четыреста пятьдесят три миллиарда триста сорок пять миллионов двести сорок две тысячи четыреста тридцать два",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, IntToWords(tc.num))
		})
	}
}

func TestIntToWordsNegativeAndSpacing(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		want string
		num  int
	}{
		{name: "negative", num: -512, want: "минус пятьсот двенадцать"},
		{name: "teens under thousand", num: 1_115, want: "одна тысяча сто пятнадцать"},
		{name: "feminine thousand one", num: 11_001, want: "одиннадцать тысяч один"},
		{name: "feminine thousand two", num: 22_002, want: "двадцать две тысячи два"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := IntToWords(tc.num)
			assert.Equal(t, tc.want, got)
			assert.NotContains(t, got, "  ", "result should not contain double spaces")
			assert.False(t, strings.HasSuffix(got, " "), "result should not end with a space")
		})
	}
}

func TestIntToWordsGender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		want   string
		num    int
		gender Gender
	}{
		{name: "masculine default", want: "сорок два", num: 42, gender: GenderMasculine},
		{name: "feminine override", want: "сорок две", num: 42, gender: GenderFeminine},
		{name: "neuter override", want: "одно", num: 1, gender: GenderNeuter},
		{name: "invalid fallback", want: "пять", num: 5, gender: Gender(99)},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, IntToWordsGender(tt.num, tt.gender))
		})
	}
}

func TestDeclineCommonCurrencies(t *testing.T) {
	t.Parallel()

	cases := []declineCase{
		{description: "singular", forms: [3]string{"рубль", "рубля", "рублей"}, num: 1},
		{description: "few", forms: [3]string{"рубль", "рубля", "рублей"}, num: 2},
		{description: "many", forms: [3]string{"рубль", "рубля", "рублей"}, num: 5},
		{description: "teens override", forms: [3]string{"рубль", "рубля", "рублей"}, num: 11},
		{description: "ends with one", forms: [3]string{"рубль", "рубля", "рублей"}, num: 21},
		{description: "ends with four", forms: [3]string{"рубль", "рубля", "рублей"}, num: 104},
		{description: "teens override multi hundred", forms: [3]string{"доллар", "доллара", "долларов"}, num: 111},
		{description: "invariant plural", forms: [3]string{"евро", "евро", "евро"}, num: 1234},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected(), Decline(tc.num, tc.forms[formOneIdx], tc.forms[formTwoIdx], tc.forms[formFiveIdx]))
		})
	}
}

func TestGetDeclensionEdgeCases(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		want string
		num  int
	}{
		{name: "zero uses five form", num: 0, want: "товаров"},
		{name: "ending with 3", num: 93, want: "товара"},
		{name: "teens upper bound", num: 219, want: "товаров"},
		{name: "teens lower bound", num: 111, want: "товаров"},
		{name: "ends with 2", num: 1002, want: "товара"},
		{name: "ends with 7", num: 1007, want: "товаров"},
		{name: "negative one uses one form", num: -1, want: "товар"},
		{name: "negative two uses two form", num: -2, want: "товара"},
		{name: "negative eleven uses five form", num: -11, want: "товаров"},
		{name: "negative twenty one uses one form", num: -21, want: "товар"},
		{name: "negative forty two uses two form", num: -42, want: "товара"},
		{name: "negative hundred one uses one form", num: -101, want: "товар"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, getDeclension(tc.num, "товар", "товара", "товаров"))
		})
	}
}

func TestDecimalToWords(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		decimal string
		want    string
		wantErr bool
	}{
		{
			name:    "simple decimal",
			decimal: "123.45",
			want:    "сто двадцать три целых сорок пять сотых",
			wantErr: false,
		},
		{
			name:    "large number from README",
			decimal: "6453345242432.42",
			want:    "шесть триллионов четыреста пятьдесят три миллиарда триста сорок пять миллионов двести сорок две тысячи четыреста тридцать две целых сорок две сотых",
			wantErr: false,
		},
		{
			name:    "no fractional part",
			decimal: "100",
			want:    "сто целых ноль сотых",
			wantErr: false,
		},
		{
			name:    "one digit fraction",
			decimal: "50.5",
			want:    "пятьдесят целых пятьдесят сотых",
			wantErr: false,
		},
		{
			name:    "zero whole",
			decimal: "0.99",
			want:    "ноль целых девяносто девять сотых",
			wantErr: false,
		},
		{
			name:    "truncate extra decimals",
			decimal: "1.999",
			want:    "одна целая девяносто девять сотых",
			wantErr: false,
		},
		{
			name:    "one hundredth",
			decimal: "5.01",
			want:    "пять целых одна сотая",
			wantErr: false,
		},
		{
			name:    "two hundredths",
			decimal: "10.02",
			want:    "десять целых две сотых",
			wantErr: false,
		},
		{
			name:    "one and feminine celaya",
			decimal: "1.5",
			want:    "одна целая пятьдесят сотых",
			wantErr: false,
		},
		{
			name:    "two feminine celyh",
			decimal: "2.5",
			want:    "две целых пятьдесят сотых",
			wantErr: false,
		},
		{
			name:    "compound twenty one feminine",
			decimal: "21.15",
			want:    "двадцать одна целая пятнадцать сотых",
			wantErr: false,
		},
		{
			name:    "invalid whole number",
			decimal: "abc.45",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid fraction",
			decimal: "123.xyz",
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := DecimalToWords(tc.decimal)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestDecimalToWordsPrecision(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		decimal   string
		want      string
		precision int
		wantErr   bool
	}{
		{
			name:      "tenths",
			decimal:   "3.5",
			precision: 1,
			want:      "три целых пять десятых",
		},
		{
			name:      "hundredths",
			decimal:   "3.14",
			precision: 2,
			want:      "три целых четырнадцать сотых",
		},
		{
			name:      "thousandths",
			decimal:   "3.145",
			precision: 3,
			want:      "три целых сто сорок пять тысячных",
		},
		{
			name:      "one and five tenths feminine",
			decimal:   "1.5",
			precision: 1,
			want:      "одна целая пять десятых",
		},
		{
			name:      "two and five tenths feminine",
			decimal:   "2.5",
			precision: 1,
			want:      "две целых пять десятых",
		},
		{
			name:      "compound twenty one and fifteen hundredths",
			decimal:   "21.15",
			precision: 2,
			want:      "двадцать одна целая пятнадцать сотых",
		},
		{
			name:      "invalid precision too low",
			decimal:   "1.5",
			precision: 0,
			wantErr:   true,
		},
		{
			name:      "invalid precision too high",
			decimal:   "1.5",
			precision: 10,
			wantErr:   true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := DecimalToWordsPrecision(tc.decimal, tc.precision)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestDecimalValueToWords(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		decimal decimal.Decimal
		want    string
		wantErr bool
	}{
		{
			name:    "simple decimal from float",
			decimal: decimal.NewFromFloat(123.45),
			want:    "сто двадцать три целых сорок пять сотых",
			wantErr: false,
		},
		{
			name:    "large number from README",
			decimal: decimal.RequireFromString("6453345242432.42"),
			want:    "шесть триллионов четыреста пятьдесят три миллиарда триста сорок пять миллионов двести сорок две тысячи четыреста тридцать две целых сорок две сотых",
			wantErr: false,
		},
		{
			name:    "from string",
			decimal: decimal.RequireFromString("100.00"),
			want:    "сто целых ноль сотых",
			wantErr: false,
		},
		{
			name:    "one digit fraction",
			decimal: decimal.NewFromFloat(50.5),
			want:    "пятьдесят целых пятьдесят сотых",
			wantErr: false,
		},
		{
			name:    "zero whole",
			decimal: decimal.NewFromFloat(0.99),
			want:    "ноль целых девяносто девять сотых",
			wantErr: false,
		},
		{
			name:    "truncate extra decimals",
			decimal: decimal.RequireFromString("1.999"),
			want:    "одна целая девяносто девять сотых",
			wantErr: false,
		},
		{
			name:    "truncate not round",
			decimal: decimal.RequireFromString("1.995"),
			want:    "одна целая девяносто девять сотых",
			wantErr: false,
		},
		{
			name:    "one hundredth",
			decimal: decimal.NewFromFloat(5.01),
			want:    "пять целых одна сотая",
			wantErr: false,
		},
		{
			name:    "two hundredths",
			decimal: decimal.NewFromFloat(10.02),
			want:    "десять целых две сотых",
			wantErr: false,
		},
		{
			name:    "negative number",
			decimal: decimal.NewFromFloat(-42.15),
			want:    "минус сорок две целых пятнадцать сотых",
			wantErr: false,
		},
		{
			name:    "negative one and half",
			decimal: decimal.NewFromFloat(-1.5),
			want:    "минус одна целая пятьдесят сотых",
			wantErr: false,
		},
		{
			name:    "negative compound twenty one",
			decimal: decimal.NewFromFloat(-21.15),
			want:    "минус двадцать одна целая пятнадцать сотых",
			wantErr: false,
		},
		{
			name:    "very precise number rounded",
			decimal: decimal.RequireFromString("3.141592653589793"),
			want:    "три целых четырнадцать сотых",
			wantErr: false,
		},
		{
			name:    "zero",
			decimal: decimal.Zero,
			want:    "ноль целых ноль сотых",
			wantErr: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := DecimalValueToWords(tc.decimal)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
