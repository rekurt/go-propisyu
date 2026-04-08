.PHONY: help lint secure tidy test coverage build

help: ## Показать справку по командам
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

lint: ## Запустить golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; GOBIN=$$(go env GOPATH)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.1.6; }
	$$(go env GOPATH)/bin/golangci-lint run --fix

secure: ## Запустить gosec
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; GOBIN=$$(go env GOPATH)/bin go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	$$(go env GOPATH)/bin/gosec -exclude-dir=protocol ./...

tidy: ## go mod tidy
	go mod tidy

test: ## Запустить все тесты
	go test ./...

coverage: ## Тесты с покрытием (открыть HTML отчёт)
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

build: ## Собрать бинарник (для проверки компиляции)
	go build ./...
