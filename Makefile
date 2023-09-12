build: 
	@go build -o bin/gobank

run: build
	@./bin/gobank

dev:
	@air

test:
	@go test -v ./...