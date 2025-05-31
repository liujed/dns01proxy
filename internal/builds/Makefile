SUBDIRS := $(shell ls -d */)

.PHONY: all
all: $(SUBDIRS)
	go mod tidy

.PHONY: $(SUBDIRS)
$(SUBDIRS):
	# Attempt to compile the build. Remove it if compilation fails.
	$(MAKE) -C $@ || rm -rf $@
