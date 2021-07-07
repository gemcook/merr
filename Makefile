.PHONY: test
test:
	go test -cover -coverprofile=./coverage.out ./...
	go tool cover -func=coverage.out | grep "total:"

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run