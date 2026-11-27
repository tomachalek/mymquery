.PHONY: all build test rtest vet clean

all: test build

build:
	CGO_ENABLED=0 go build -o mymquery .

test:
	go test ./...

rtest:
	go test -race ./...

vet:
	go vet ./...

clean:
	rm -f mymquery
