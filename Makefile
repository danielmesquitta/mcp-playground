default: run

.PHONY: install
install:
	@go mod download

.PHONY: update
update:
	@go mod tidy && go get -u ./...
