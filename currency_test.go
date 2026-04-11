package propisyu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoney(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		want  string
		cur   Currency
		whole int
		cents int
	}{
		{
			name:  "RUB standard",
			whole: 1234,
			cents: 56,
			cur:   CurrencyRUB,
			want:  "одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек",
		},
		{
			name:  "RUB singular",
			whole: 1,
			cents: 1,
			cur:   CurrencyRUB,
			want:  "один рубль одна копейка",
		},
		{
			name:  "USD zero cents",
			whole: 5,
			cents: 0,
			cur:   CurrencyUSD,
			want:  "пять долларов ноль центов",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := Money(tc.whole, tc.cents, tc.cur)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMoneyFromString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		amount  string
		want    string
		cur     Currency
		wantErr bool
	}{
		{
			name:   "RUB from string",
			amount: "1234.56",
			cur:    CurrencyRUB,
			want:   "одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек",
		},
		{
			name:   "no fractional part",
			amount: "100",
			cur:    CurrencyRUB,
			want:   "сто рублей ноль копеек",
		},
		{
			name:    "invalid amount",
			amount:  "abc.45",
			cur:     CurrencyRUB,
			wantErr: true,
		},
		// Regression guards: the minus sign must be preserved when the whole
		// part rounds to zero. Without the fix MoneyFromString(-0.xx) used
		// to drop the sign because strconv.Atoi("-0") == 0 and Money carries
		// no sign information. Matches the same guard pattern PR #16 added
		// for DecimalToWords / DecimalValueToWords / DecimalToWordsPrecision.
		{
			name:   "negative zero whole RUB 50 cents",
			amount: "-0.50",
			cur:    CurrencyRUB,
			want:   "минус ноль рублей пятьдесят копеек",
		},
		{
			name:   "negative zero whole RUB one cent",
			amount: "-0.01",
			cur:    CurrencyRUB,
			want:   "минус ноль рублей одна копейка",
		},
		{
			name:   "negative zero whole USD 5 cents",
			amount: "-0.05",
			cur:    CurrencyUSD,
			want:   "минус ноль долларов пять центов",
		},
		// Regular negatives keep working — IntToWordsGender already carries
		// the sign in the whole part, so the guard must not double-prefix.
		{
			name:   "negative whole does not double minus",
			amount: "-1234.56",
			cur:    CurrencyRUB,
			want:   "минус одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек",
		},
		// Zero with zero cents is a true zero — no minus even when the
		// input string has a leading "-" (which would never legitimately
		// appear for zero, but the guard must still not fire).
		{
			name:   "negative-prefixed zero stays plain zero",
			amount: "-0.00",
			cur:    CurrencyRUB,
			want:   "ноль рублей ноль копеек",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := MoneyFromString(tc.amount, tc.cur)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}
