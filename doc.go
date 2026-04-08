// Package propisyu converts integers and decimal numbers into Russian words
// with correct grammatical gender and noun declension.
//
// This package is useful for generating invoices, receipts, checks, financial
// documents, voice assistant prompts, chatbots, and any system that must
// spell out numbers in fluent Russian (числа прописью).
//
// # Features
//
//   - Numbers up to duodecillions (10^39)
//   - Three grammatical genders: masculine, feminine, neuter
//   - Automatic noun declension via [Decline]
//   - Decimal number support via strings or [github.com/shopspring/decimal]
//   - Zero external dependencies for core integer functions
//
// # Integer conversion
//
//	propisyu.IntToWords(42)
//	// "сорок два"
//
//	propisyu.IntToWordsGender(2, propisyu.GenderFeminine)
//	// "две"
//
// # Decimal conversion
//
//	result, _ := propisyu.DecimalToWords("123.45")
//	// "сто двадцать три целых сорок пять сотых"
//
// # Noun declension
//
//	propisyu.Decline(5, "рубль", "рубля", "рублей")
//	// "рублей"
//
//	propisyu.Decline(21, "день", "дня", "дней")
//	// "день"
//
// See the project README for more examples:
// https://github.com/rekurt/go-propisyu
package propisyu
