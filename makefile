TARGET=sudoku-server

.PHONY: build run test clean

run: build
	@./$(TARGET)

build: main.go
	@go build .

test: main.go
	@go clean -testcache
	@go test  ./...

clean: $(TARGET)
	@rm $(TARGET)

