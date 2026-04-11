# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
