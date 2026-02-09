.PHONY: build test clean

build:
	mkdir -p internal/transpile/lang
	cp lang/*.json internal/transpile/lang/
	go build -o bin/pgo ./cmd/pgo

test:
	mkdir -p internal/transpile/lang
	cp lang/*.json internal/transpile/lang/
	go test ./...

clean:
	rm -rf bin .pgo_gen
