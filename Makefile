.PHONY: build
build:
	go build -o build/gist -v

.PHONY: test
test:
	go test

.PHONY: lint
lint:
	golint -set_exit_status ./...

.PHONY: verify
verify: clean lint test

.PHONY: clean
clean:
	go clean
	rm -rf build/

.PHONY: fmt
fmt:
	go fmt

.PHONY: prepare-example
prepare-example:
	@echo :prepare-example
	@mkdir example
	@cp testdata/profile.yml example/profile.yml
	@cat example/profile.yml
	@echo

.PHONY: profile-example
profile-example: prepare-example
	@echo :profile-example
	build/gist -f example/profile.yml profile --profile default -t aa00bb11cc22
	@cat example/profile.yml
	@echo

.PHONY: cleanup-example
cleanup-example:
	@echo :cleanup-example
	@rm -rf example/
	@echo

.PHONY: example
example: clean build profile-example cleanup-example
