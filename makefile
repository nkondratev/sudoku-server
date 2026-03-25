TARGET=sudoku-server

.PHONY: build run test clean

run: build
	@./$(TARGET)

build: main.go
	@go build .

test: main.go
	@go test -v ./...

clean: $(TARGET)
	@rm $(TARGET)

