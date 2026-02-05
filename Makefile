up:
	@cd configs/local && docker compose up -d

up-build:
	@cd configs/local && docker compose up --build --force-recreate -d

down:
	@cd configs/local && docker compose down

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags -a -installsuffix cgo cmd/main.go

restart:
	@cd configs/local && docker compose restart

check:
	@echo "CHECK LINT ANALYSIS"
	@echo "=========================================================="
	@golangci-lint run --config .golangci.yml --timeout 5m

unit_test:
	go test -v -failfast -cover -count=1 -coverprofile=coverage.out $(shell go list ./... | grep -v /test)

it_test:
	go test -v -failfast -cover -count=1 -coverprofile=coverage.out ./test

swag:
	@if ! command -v swag >/dev/null; then \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@swag fmt && swag init -g server.go -o docs --dir internal/adapters/http,pkg/errors

clean:
	@echo "🧹 Limpando arquivos temporários..."
	@rm -rf bin/ coverage.out .scannerwork/
	@docker system prune -f

.PHONY: check test swag up up-build down restart dev-run status clean help \
	infra-local infra-prod infra-destroy-local infra-destroy-prod