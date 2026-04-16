TARGET=app

.PHONY: build run test clean

run: build
	@./$(TARGET)

build: main.go
	@go build . -o $(TARGET)

test: main.go
	@go clean -testcache
	@go test  ./...

clean: $(TARGET)
	@rm $(TARGET)

release: main.go
	@go build -tags release $(TARGET)
