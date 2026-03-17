package propisyu_test

import (
	"fmt"
	"log"

	propisyu "github.com/rekurt/go-propisyu"
)

func ExampleIntToWords() {
	fmt.Println(propisyu.IntToWords(321))
	// Output: триста двадцать один
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

func ExampleDecimalToWords() {
	result, err := propisyu.DecimalToWords("123.45")
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
