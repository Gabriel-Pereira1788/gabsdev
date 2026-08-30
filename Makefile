.PHONY: dev run build generate fmt tidy test clean

# Live reload dev server (templ generate + go build on change)
dev:
	air

# Generate templ files then run the server once (no reload)
run: generate
	go run ./cmd/server/main.go

# Compile the binary into ./bin/server
build: generate
	go build -o ./bin/server ./cmd/server/main.go

# Regenerate *_templ.go from .templ sources
generate:
	templ generate

fmt:
	templ fmt .
	gofmt -w .

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf ./bin ./tmp
