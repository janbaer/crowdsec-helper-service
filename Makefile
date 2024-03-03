.PHONY: run build test

# Install gow first with `go install github.com/mitranim/gow@latest`
run:
	gow run main.go

build:
	go build -o ./bin/crowdsec-helper-service

test:
	go test ./...

