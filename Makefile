build: *.go
	@go build -o promgrep .

run: build
	./promgrep

GOTESTSUM_PATH ?= $(shell which gotestsum)

test:
	$(if $(GOTESTSUM_PATH), gotestsum --, go test) ./...

watch:
	gotestsum --watch --format=standard-verbose -- ./...
