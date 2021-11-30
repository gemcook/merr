.PHONY: test
test:
	go test -cover -coverprofile=./coverage.out ./...
	go tool cover -func=coverage.out | grep "total:"

# Prepare golangci
.PHONY: bootstrap-golangci
bootstrap-golangci:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s

# If bin/golangci-lint does not exist,
# run make bootstrap-golangci and then run lint.
.PHONY: lint
lint:
	@if [ ! -e bin/golangci-lint ]; then $(MAKE) bootstrap-golangci ; fi
	bin/golangci-lint run ./...

.PHONY: clean
clean:
	rm -rf ./bin
