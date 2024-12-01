.PHONY: build test lint run migrup migrdown migrcreate swago
build: ## Собираем бинарник
	CGO_ENABLED=0 go build -o bin/app -mod=readonly

test: ## Запускаем тесты
	go test -v ./... -race -cover

lint: ## Запускаем линтер (конфиг .golangci.yaml)
	golangci-lint run

run: ## Запускаем сервис (без билда)
	go run ./cmd/app/main.go

swago: ## Запускаем swag (генерим сваггер)
	# https://github.com/swaggo/swag
	swag init --parseDependency --parseInternal --propertyStrategy pascalcase --parseDepth 3 -g cmd/app/main.go

migrup: ## Накатываем миграции
	go run ./cmd/migration/migrator.go up

migrdown: ## Откатываем миграции
	go run ./cmd/migration/migrator.go down

migrcreate: ## Создаёт новую миграцию (параметр ARG1 обязателен)
	@if [ -z "$(ARG1)" ]; then \
		echo "Ошибка: Не задан параметр ARG1! Пример: make migrcreate ARG1=my_migration"; \
	else \
		go run ./cmd/migration/migrator.go -dir migrations create $(ARG1) sql; \
	fi
help:  ## Показываем список команд
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
