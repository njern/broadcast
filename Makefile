.PHONY: build lint test fmt modernize

build:
	go build .

lint:
	golangci-lint-v2 run ./...

test:
	go test ./...

fmt:
	gofumpt -l -w .

modernize:
	go run golang.org/x/tools/go/analysis/passes/modernize/cmd/modernize@latest -fix ./...
