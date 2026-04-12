# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.1] - 2026-04-12

### Fixed

- `Ordinal(math.MinInt, gender)` used to recurse into `Ordinal(-n, gender)`,
  but `-math.MinInt == math.MinInt` in two's-complement, so the program
  crashed with `fatal error: stack overflow`. The negative branch now
  lifts the magnitude through `uint64` (`uint64(math.MaxInt)+1` for the
  `MinInt` case), identical to the pattern already used by
  `convertIntToWords` since v0.4.0. (#19)
- `DecimalToWords` and `DecimalToWordsPrecision` called `strings.TrimSpace`
  for sign detection but kept the untrimmed input for the main parsing
  path, so leading whitespace crashed `strconv.Atoi` with a confusing
  "invalid syntax" error while the sign-detection path silently accepted
  the same whitespace. Both functions now apply `TrimSpace` uniformly,
  making `" 123.45"`, `"\t1.5"`, and `"\n-0.50"` valid inputs. (#20)
- `MoneyFromString("-0.xx", …)` silently dropped the minus sign because
  `strconv.Atoi("-0") == 0` and `Money` carries no sign information. The
  function now preserves the sign through the same guard pattern PR #16
  added for `DecimalToWords` / `DecimalValueToWords` /
  `DecimalToWordsPrecision`. Consistency:
  `DecimalToWords("-0.50")` → `"минус ноль целых пятьдесят сотых"` and
  `MoneyFromString("-0.50", CurrencyRUB)` → `"минус ноль рублей пятьдесят
  копеек"`. (#21)
- `DecimalValueToWords` silently produced **wrong output** for any
  `decimal.Decimal` whose whole part overflowed `int64` (e.g.
  `10^29.99` rendered as a random `int64`-wrapped phrase). The previous
  `whole > math.MaxInt64` guard was a tautology on 64-bit platforms
  where `math.MaxInt64 == math.MaxInt`. The check now uses
  `d.Truncate(0).BigInt().IsInt64()` **before** calling `IntPart`, so
  out-of-range inputs return `ErrNumberTooLarge` instead of garbage.
  (#22)

### Changed

- `convertPositiveUint64ToWords` used O(triads²) allocations due to
  `parts = append([]string{triadWords}, parts...)` prepend. Replaced
  with preallocated linear append + in-place reverse: for
  `math.MaxInt` (7 triads) allocs/op drop from 52 to 40 (−23%) and
  time/op from ~1500 ns to ~1300 ns (−14%) on Apple M1 Pro. The
  improvement scales with triad count — duodecillion-scale inputs
  (13 triads) previously did ~169 inner copies; they now do 13
  appends plus one reverse. No public API or output changes. (#23)

### Documentation

- `examples/examples_test.go` now has runnable godoc `Example*`
  functions for the v0.4.0-introduced APIs that were missing preview
  on pkg.go.dev: `ExampleOrdinal`, `ExampleMoney`, `ExampleMoneyFromString`,
  `ExampleDecimalToWordsPrecision`, plus tenths and higher-precision
  variants. (#24)

### CI

- `test: add readme consistency checks to prevent documentation drift`
  introduced a package-level test that walks `go/ast` and asserts every
  exported top-level declaration is mentioned in both `README.md` and
  `README_EN.md`, plus a `## ` / `### ` header parity check between RU
  and EN. No workflow changes — runs as part of the existing `test`
  job. (#18)
- `ci: add govulncheck step to the lint job` now runs
  `golang.org/x/vuln/cmd/govulncheck@v1.1.4` on every push and pull
  request. Current master is clean (0 reachable vulnerabilities;
  6 non-reachable vulns in transitively-required modules and 1 in an
  imported-but-not-called package are correctly ignored). (#25)

## [0.4.0] - 2026-04-11

### Added

- `Ordinal(n int, gender Gender) string` — Russian ordinal numbers with full
  support for all three grammatical genders (`первый / первая / первое`),
  compound forms (`двадцать первый`, `сорок второй`), and round numbers
  (`тысячный`, `сороковой`, `миллионный`).
- `Money(whole, cents int, c Currency) string` and
  `MoneyFromString(amount string, c Currency) (string, error)` — one-call
  currency formatting with declension and gender handled end-to-end.
- `type Currency` with the fields `WholeOne/Two/Five`, `WholeGender`,
  `FracOne/Two/Five`, `FracGender`, plus the exported presets
  `CurrencyRUB`, `CurrencyUSD`, and `CurrencyEUR`.
- `DecimalToWordsPrecision(s string, precision int) (string, error)` —
  decimal-to-words with arbitrary fractional precision from 1 to 9 digits
  (десятых → миллиардных).
- README (RU + EN): full rewrite with hero snippet, unified API section,
  and four realistic usage scenarios. RU ↔ EN are now strict structural
  mirrors (same 8 `##` headers, same Go snippets).

### Fixed

- Grammar: `DecimalToWords`, `DecimalValueToWords`, and `DecimalToWordsPrecision`
  now use the feminine gender for the whole part, producing grammatically
  correct Russian forms — e.g. `"одна целая"` / `"две целых"` instead of the
  previous `"один целых"` / `"два целых"`. This also affects compound numbers
  such as `21`, `432`, and `6453345242432`, whose last triad is now declined in
  the feminine gender before `"целая/целых"`.
- `Decline` (and internal `getDeclension`) now correctly handle negative
  numbers. Previously, Go's sign-preserving `%` operator caused negative inputs
  to fall through to the `five` form regardless of the actual last digit — for
  example `Decline(-1, ...)` returned `five` instead of `one`. The function
  now reduces modulo 100 first and then negates the remainder, which is also
  safe for `math.MinInt` (where `-n` on the full-magnitude value would
  overflow).
- `IntToWords` / `IntToWordsGender` now handle `math.MinInt` correctly.
  Previously the `"минус " + convertIntToWords(-n, dict)` path overflowed at
  `-math.MinInt`; conversion now routes the positive magnitude through a
  `uint64` helper so the smallest int produces a valid phrase without
  panicking.
- `DecimalToWords`, `DecimalValueToWords`, and `DecimalToWordsPrecision` now
  preserve the minus sign for values like `-0.50` where the whole part is
  zero and the fractional part is non-zero. Previously the sign was lost
  because the zero whole part carries no sign information, so `"-0.50"` and
  `"0.50"` rendered identically. All three entry points now prefix
  `"минус "` in this case, keeping string- and `decimal.Decimal`-based APIs
  consistent.

### BREAKING CHANGES

- Output of `DecimalToWords`, `DecimalValueToWords`, and `DecimalToWordsPrecision`
  has changed for inputs whose whole part ends in `1` or `2` (excluding teens
  `11`–`19`). Golden-string tests or snapshots referencing these functions must
  be updated. Example migrations:
  - `"один целых девяносто девять сотых"` → `"одна целая девяносто девять сотых"`
  - `"минус сорок два целых пятнадцать сотых"` → `"минус сорок две целых пятнадцать сотых"`
  - `"...четыреста тридцать два целых сорок две сотых"` → `"...четыреста тридцать две целых сорок две сотых"`
- `Decline` now returns different (correct) forms for negative inputs; callers
  relying on the previous (incorrect) behavior must be updated.

## [0.3.0] - 2026-04-08

### Added

- `doc.go` with comprehensive package documentation for pkg.go.dev
- Use cases table and "Why go-propisyu" section in both READMEs
- Hero code snippet in READMEs for instant understanding
- More `Example*` functions for richer pkg.go.dev documentation page
- `FUNDING.yml` for GitHub Sponsors
- Issue template config with links to pkg.go.dev and discussions
- GitHub topics, description, and homepage for better discoverability

### Changed

- Improved Go doc comments on all exported functions
- Updated SECURITY.md with actual version information
- Lowered minimum Go version to 1.22 for wider compatibility
- Upgraded CI: separate lint job, golangci-lint-action v7, codecov v6

### Fixed

- CI pipeline failure on Go 1.24 caused by gosec requiring Go 1.25

## [0.2.0] - 2025-11-06

### Added

- Decimal number conversion via `DecimalToWords` and `DecimalValueToWords`
- `GenderNeuter` support for neuter grammatical gender
- Support for `decimal.Decimal` input type

## [0.1.0] - 2025-11-06

### Added

- Initial release: `IntToWords`, `IntToWordsGender`, `Decline`
- Support for numbers up to duodecillions (10³⁹)
- Masculine and feminine grammatical gender support
- Automatic noun declension via `Decline`

[Unreleased]: https://github.com/rekurt/go-propisyu/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/rekurt/go-propisyu/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/rekurt/go-propisyu/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/rekurt/go-propisyu/releases/tag/v0.1.0
