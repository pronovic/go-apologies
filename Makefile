mocks:
	mockery --all --inpackage --case underscore --dir ./ --output -./

demo:
	CGO_CFLAGS="-w" go run demo/demo.go -adult -players=4 -input=reward -delay=200

.PHONY: demo mocks
