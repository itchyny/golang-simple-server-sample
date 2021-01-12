BIN := golang-simple-server-sample
export GO111MODULE=on

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BIN) .

.PHONY: clean
clean:
	rm -rf $(BIN)
	go clean
