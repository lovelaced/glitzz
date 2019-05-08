all: dependencies test build check-gofmt

ci: dependencies test-ci build check-gofmt

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

check-gofmt:
	 ./.travis.gofmt.sh

test-verbose:
	go test -v ./...

test-short:
	go test -short ./...

clean:
	rm -f ./build/glitzz

dependencies:
	go get -t ./...

.PHONY: all ci build run doc test test-ci test-verbose test-short clean dependencies check-gofmt
