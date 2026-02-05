.PHONY: test bench cover vet clean

test:
	go test -race ./...

bench:
	go test -bench=. -benchmem ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | tail -1

vet:
	go vet ./...

clean:
	rm -f coverage.out
