APP_NAME := crm-api-core
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)

.PHONY: help setup install-tools mod lint test format build clean

help:
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: install-tools mod ## Prepara o ambiente local instalando todas as dependências e ferramentas
	@echo "==> Configurando hooks do pre-commit..."
	pre-commit install
	@echo "✅ Ambiente configurado com sucesso! Você está pronto para codar."

install-tools: ## Instala as ferramentas globais necessárias (linter, formatador)
	@echo "==> Instalando ferramentas de desenvolvimento Go..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "⚠️  Nota: Certifique-se de ter o 'pre-commit' instalado na sua máquina (via apt, brew ou pip)."

mod: ## Baixa e atualiza as dependências do Go (go.mod)
	@echo "==> Baixando dependências Go..."
	go mod tidy
	go mod download

lint: ## Roda o golangci-lint em todo o projeto
	@echo "==> Rodando linter..."
	golangci-lint run ./...

format: ## Formata o código e organiza os imports
	@echo "==> Formatando código..."
	go fmt ./...
	goimports -w .

test: ## Executa os testes unitários com identificação de race conditions
	@echo "==> Rodando testes..."
	go test -v -race ./...

build: ## Compila o binário da aplicação
	@echo "==> Compilando o projeto..."
	go build -o bin/$(APP_NAME) main.go
	@echo "✅ Build concluído em bin/$(APP_NAME)"

clean: ## Limpa os arquivos compilados e cache
	@echo "==> Limpando o projeto..."
	rm -rf bin/
	go clean -testcache
