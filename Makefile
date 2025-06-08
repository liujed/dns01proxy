BUILDS := $(shell ls -d internal/builds/*/)

.PHONY: default
default: clean build

# Removes compiled artifacts, but leaves generated artifacts.
.PHONY: clean
clean:
	rm -rf dist/

# Builds compiled artifacts.
.PHONY: build
build: $(BUILDS)

.PHONY: $(BUILDS)
$(BUILDS):
	$(MAKE) -C $@

# Removes compiled and generated artifacts.
.PHONY: pristine
pristine: clean
	rm -rf internal/builds/*/ .goreleaser.yml
	go mod tidy

# Regenerates build targets.
.PHONY: gen-builds
gen-builds: pristine
	go generate ./internal/gen/builds
	go mod tidy
	$(MAKE) -C internal/builds

# Regenerates goreleaser config.
.PHONY: gen-release
gen-release:
	go generate ./internal/gen/release

# Regenerates all generated artifacts.
.PHONY: gen-all
gen-all: gen-builds
	$(MAKE) gen-release
