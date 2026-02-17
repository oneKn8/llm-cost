BINARY=llm-cost

build:
	go build -o $(BINARY) .

install:
	go install .

test:
	go test ./...

clean:
	rm -f $(BINARY)

.PHONY: build install test clean
