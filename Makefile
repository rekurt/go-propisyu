lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; GOBIN=$$(go env GOPATH)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2; }
	$$(go env GOPATH)/bin/golangci-lint run --fix
secure:
	@command -v gosec >/dev/null 2>&1 || { echo "Installing gosec..."; GOBIN=$$(go env GOPATH)/bin go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	$$(go env GOPATH)/bin/gosec -exclude-dir=protocol ./...
tidy:
	go mod tidy
test:
	go test ./...
