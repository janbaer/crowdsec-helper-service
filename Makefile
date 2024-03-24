.PHONY: run build test

VERSION=1.0.0

# Install gow first with `go install github.com/mitranim/gow@latest`
run:
	gow run main.go

build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o ./bin/crowdsec-helper-service

test:
	go test ./...
