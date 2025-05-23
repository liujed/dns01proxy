.PHONY: default
default: clean build

.PHONY: clean
clean:
	rm -rf dist/

.PHONY: build
build: clean
	go build -ldflags '-w -s' \
		-tags=nobadger,nomysql,nopgx \
		-trimpath \
		-o dist/dns01proxy \
		./cmd/dns01proxy
