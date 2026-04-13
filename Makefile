.PHONY: build start

build:
	go build -o bin/go-stock-scanner .

start: build
	./bin/go-stock-scanner scan