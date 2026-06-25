run:
	@go run .

test:
	@go test -v ./...

build:
	@go build .

clean:
	@rm -f lru_cache

.PHONY: run test build clean