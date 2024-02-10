mocks:
	# Generate mockery mocks for use with testify, which are checked into git
	# To get the tool: brew install mockery
	mockery --all --inpackage --case underscore --dir ./ --output -./
.PHONY: mocks

demo:
	# Run the ncurses demo with some sensible defaults
	# CGO_CFLAGS is needed to ignore warnings from goncurses
	CGO_CFLAGS="-w" go run demo/demo.go -adult -players=4 -input=reward -delay=200
.PHONY: demo

format:
	# Format the source tree using gofumpt
	# To get the tool: go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .
.PHONY: format

test:
	# Run the test suite with caching disabled
	# CGO_CFLAGS is needed to ignore warnings from goncurses
	CGO_CFLAGS="-w" go test -race -count=1 ./...
.PHONY: test
