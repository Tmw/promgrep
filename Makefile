build: *.go
	@go build -o promgrep .

run: build
	./promgrep
