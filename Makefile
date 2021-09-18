BIN := golang-simple-server-sample

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags="-s -w" -o $(BIN) .

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean
