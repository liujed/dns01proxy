.PHONY: default
default: clean build

.PHONY: clean
clean:
	rm -rf dist/ internal/builds/*/
	go mod tidy

.PHONY: generate
generate: clean
	go generate ./...
	go mod tidy

.PHONY: build
build: clean generate
	$(MAKE) -C internal/builds
