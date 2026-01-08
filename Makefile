test:
	go test ./... -race -coverprofile=coverage.out

lint:
	golangci-lint run

benchmark:
	go test ./... -bench=. -benchmem

coverage:
	go test ./... -coverprofile=coverage.out

clean:
	rm -f coverage.out

all: test lint benchmark coverage