all: build

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

test-verbose:
	go test -v ./...

test-short:
	go test -short ./...

bench:
	go test -v -run=XXX -bench=. ./...

clean:
	rm -f ./build/glitzz

.PHONY: all build run doc test test-verbose test-short bench clean
