.PHONY: build
build:
	go build -o build/gist main.go

.PHONY: test
test:
	go test

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: verify
verify: lint test
