PACKAGES=$(shell go list ./... | egrep -v 'go-apologies\/demo')

mocks:
	mockery --all --inpackage --case underscore --dir ./ --output -./
.PHONY: mocks

demo:
	CGO_CFLAGS="-w" go run demo/demo.go -adult -players=4 -input=reward -delay=200
.PHONY: demo

format:
	goimports -w .
.PHONY: format

test:
	CGO_FLAGS="-w" go test -race -v $(PACKAGES)
.PHONY: test
