default:
	@echo "hello world"

#? lint: Run latest golangci-lint linter
lint:
	@echo "--> Running linter"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
.PHONY: lint

#? vulncheck: Run latest govulncheck
vulncheck:
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
.PHONY: vulncheck

#? test: Run tests
test:
	go test -v ./...
.PHONY: test

# Run goimports-reviser to lint and format imports
lint-imports:
	@find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | while read -r file; do \
		goimports-reviser -company-prefixes github.com/0xPellNetwork/pelldvs-libs -rm-unused -format "$$file"; \
	done
