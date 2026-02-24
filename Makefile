BIN := tui-timer

.PHONY: build run clean

build:
	go build -o $(BIN) ./cmd/tui-timer

run: build
	./$(BIN)

clean:
	rm -f $(BIN)
