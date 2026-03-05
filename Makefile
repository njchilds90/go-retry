
# Makefile for njchilds90/go-retry

.PHONY: build
build:
	@go build -o go-retry ./...

.PHONY: test
test:
	@go test -race ./...

.PHONY: lint
lint:
	@golint ./...

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: bench
bench:
	@go test -bench=. ./...

.PHONY: coverage
coverage:
	@go test -coverpkg=./... -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

.PHONY: clean
clean:
	@go clean
	@rm -f coverage.out
