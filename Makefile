.PHONY: build
build:
	go build -o build/gist main.go

.PHONY: test
test:
	go test
