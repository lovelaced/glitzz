all: dependencies test build

ci: dependencies test-ci build

build:
	mkdir -p build
	go build -o ./build/glitzz ./cmd/glitzz

run:
	./build/glitzz

doc:
	@echo "http://localhost:6060/pkg/github.com/lovelaced/glitzz/"
	godoc -http=:6060

test:
	go test ./...

test-ci:
	go test -coverprofile=coverage.txt -covermode=atomic ./...

test-verbose:
	go test -v ./...

test-short:
	go test -short ./...

clean:
	rm -f ./build/glitzz

dependencies:
	go get -t ./...

.PHONY: all build run doc test test-verbose test-short clean dependencies
