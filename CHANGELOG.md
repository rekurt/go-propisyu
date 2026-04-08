# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
