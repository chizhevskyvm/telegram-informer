GO := go
LINTER := golangci-lint
PKG := ./...

.DEFAULT_GOAL := help

## Форматирование кода (goimports + gofmt)
fmt:
	goimports -w .
	$(GO) fmt $(PKG)

## Линтер
lint:
	$(LINTER) run ./...

## Автофикс линтера
lint-fix:
	$(LINTER) run --fix ./...

## Помощь
help:
	@echo "make fmt       - Форматировать код (goimports + gofmt)"
	@echo "make lint      - Запустить линтер golangci-lint"
	@echo "make lint-fix  - Автофикс линтера"