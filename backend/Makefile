.PHONY: run build clean test

run:
	go run main.go

build:
	go build -o llmscope-backend main.go

clean:
	rm -f llmscope-backend

test:
	go test ./...

deps:
	go mod download
	go mod tidy
