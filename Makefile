default: build

.PHONY: install
install:
	@go mod download

.PHONY: update
update:
	@go mod tidy && go get -u ./...

.PHONY: build
build:
	@go build -o cep-server main.go
