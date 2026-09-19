.PHONY: build test test-v run-web clean install lint

BINARY_NAME=openstories

build:
	@echo "🔨 Compilando $(BINARY_NAME)..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/openstories

test:
	@echo "🧪 Executando testes..."
	go test ./...

test-v:
	@echo "🧪 Executando testes detalhados..."
	go test -v ./...

run-web: build
	@echo "🌐 Iniciando servidor web..."
	./$(BINARY_NAME) serve --port 8080

install: build
	@echo "📦 Instalando binário no PATH..."
	go install ./cmd/openstories

clean:
	@echo "🧹 Limpando binários..."
	rm -f $(BINARY_NAME)
