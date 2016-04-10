SOURCE := $(shell find -name '*.go')
BIN := main

.PHONY: clean remake

$(BIN): $(SOURCE)
	go build -o $(BIN)

clean:
	rm $(BIN)

remake: clean $(BIN)
