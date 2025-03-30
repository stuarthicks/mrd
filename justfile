all: fmt build test

fmt:
  go mod tidy
  go fmt ./...

build:
  go build ./...

test:
  go test -race ./...

clean:
  git clean -f .

snapshot:
  goreleaser release --clean --snapshot

release:
  goreleaser release --clean
