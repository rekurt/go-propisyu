package propisyu_test

import (
	"fmt"
	"log"

	propisyu "github.com/rekurt/go-propisyu"
	"github.com/shopspring/decimal"
)

func ExampleIntToWords() {
	fmt.Println(propisyu.IntToWords(321))
	// Output: триста двадцать один
}

func ExampleIntToWords_zero() {
	fmt.Println(propisyu.IntToWords(0))
	// Output: ноль
}

func ExampleIntToWords_negative() {
	fmt.Println(propisyu.IntToWords(-42))
	// Output: минус сорок два
}

func ExampleIntToWords_million() {
	fmt.Println(propisyu.IntToWords(1_000_000))
	// Output: один миллион
}

func ExampleIntToWordsGender() {
	fmt.Println(propisyu.IntToWordsGender(2, propisyu.GenderFeminine))
	fmt.Println(propisyu.IntToWordsGender(2, propisyu.GenderMasculine))
	fmt.Println(propisyu.IntToWordsGender(1, propisyu.GenderNeuter))
	// Output:
	// две
	// два
	// одно
}

func ExampleDecline() {
	fmt.Println(propisyu.Decline(1, "рубль", "рубля", "рублей"))
	fmt.Println(propisyu.Decline(2, "рубль", "рубля", "рублей"))
	fmt.Println(propisyu.Decline(5, "рубль", "рубля", "рублей"))
	fmt.Println(propisyu.Decline(21, "рубль", "рубля", "рублей"))
	// Output:
	// рубль
	// рубля
	// рублей
	// рубль
}

func ExampleDecline_units() {
	for _, n := range []int{1, 2, 5, 11, 21} {
		fmt.Printf("%d %s\n", n, propisyu.Decline(n, "день", "дня", "дней"))
	}
	// Output:
	// 1 день
	// 2 дня
	// 5 дней
	// 11 дней
	// 21 день
}

func ExampleDecimalToWords() {
	result, err := propisyu.DecimalToWords("123.45")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// Output: сто двадцать три целых сорок пять сотых
}

func ExampleDecimalValueToWords() {
	d := decimal.NewFromFloat(123.45)
	result, err := propisyu.DecimalValueToWords(d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// Output: сто двадцать три целых сорок пять сотых
}

func ExampleIntToWords_receipt() {
	amount := 1234
	words := propisyu.IntToWords(amount)
	decl := propisyu.Decline(amount, "рубль", "рубля", "рублей")
	fmt.Printf("%s %s 00 копеек\n", words, decl)
	// Output: одна тысяча двести тридцать четыре рубля 00 копеек
}

func ExampleIntToWords_invoice() {
	amount := 42
	words := propisyu.IntToWordsGender(amount, propisyu.GenderFeminine)
	decl := propisyu.Decline(amount, "штука", "штуки", "штук")
	fmt.Printf("Количество: %s (%s)\n", words, decl)
	// Output: Количество: сорок две (штуки)
}
