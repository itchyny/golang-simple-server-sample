.PHONY: build
build:
	go build -ldflags="-s -w" .

.PHONY: clean
clean:
	go clean
