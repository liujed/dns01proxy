.PHONY: default
default: clean build

.PHONY: clean
clean:
	rm -rf dist/ internal/builds/*/
	go mod tidy

.PHONY: gen-builds
gen-builds: clean
	go generate ./internal/gen/builds
	go mod tidy

.PHONY: builds
builds: clean gen-builds
	$(MAKE) -C internal/builds

.PHONY: build
build: clean builds gen-release

.PHONY: gen-release
gen-release: builds
	go generate ./internal/gen/release
