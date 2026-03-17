# Contributing to go-propisyu

Thank you for your interest in contributing!

## Getting Started

```bash
git clone https://github.com/rekurt/go-propisyu.git
cd go-propisyu
go mod tidy
```

## Running Tests

```bash
make test          # Run all tests
make coverage      # Run tests and open HTML coverage report
```

## Running Linter

```bash
make lint          # Run golangci-lint (auto-installs if missing)
make secure        # Run gosec security scanner
```

## Code Style

This project uses [golangci-lint](https://golangci-lint.run/) with default settings. Run `make lint` before submitting a PR to ensure your code passes all checks.

## Opening a Pull Request

1. Fork the repository and create a feature branch from `master`
2. Write or update tests for any changed behavior
3. Run `make test` and `make lint` — both must pass
4. Use [Conventional Commits](https://www.conventionalcommits.org/) for your commit messages:
   - `feat: add support for ordinal numbers`
   - `fix: correct declension for 11-19`
   - `docs: update README examples`
5. Open a PR against the `master` branch

## Adding a New Language Form

If you need to add a new grammatical form or extend declension rules:

1. Edit `declension.go` — the `Decline` function and its helpers live there
2. Edit `num2word.go` — the `dictionary` struct holds word forms
3. Add test cases in `declension_test.go` covering the new form
4. Update `README.md` and `README_EN.md` if the public API changes

## Reporting Bugs

Please use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md) when opening issues.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
