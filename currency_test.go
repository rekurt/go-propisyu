package propisyu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoney(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		whole int
		cents int
		cur   Currency
		want  string
	}{
		{
			name:  "RUB standard",
			whole: 1234,
			cents: 56,
			cur:   CurrencyRUB(),
			want:  "одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек",
		},
		{
			name:  "RUB singular",
			whole: 1,
			cents: 1,
			cur:   CurrencyRUB(),
			want:  "один рубль одна копейка",
		},
		{
			name:  "USD zero cents",
			whole: 5,
			cents: 0,
			cur:   CurrencyUSD(),
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
		cur     Currency
		want    string
		wantErr bool
	}{
		{
			name:   "RUB from string",
			amount: "1234.56",
			cur:    CurrencyRUB(),
			want:   "одна тысяча двести тридцать четыре рубля пятьдесят шесть копеек",
		},
		{
			name:   "no fractional part",
			amount: "100",
			cur:    CurrencyRUB(),
			want:   "сто рублей ноль копеек",
		},
		{
			name:    "invalid amount",
			amount:  "abc.45",
			cur:     CurrencyRUB(),
			wantErr: true,
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
