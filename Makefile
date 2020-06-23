################################################################################
SHELL=/bin/bash
ifeq (, $(shell which buildtool))
$(warning "Installing buildtool...")
$(warning "go get github.com/gearboxworks/buildtool")
$(shell go get github.com/gearboxworks/buildtool)
endif
BUILDTOOL := $(shell which buildtool)
ifeq (, $(BUILDTOOL))
$(error "No buildtool found...")
endif
################################################################################

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

all:
	@:

%:
	@:

################################################################################

help:
	@$(BUILDTOOL) $@ $(args)

build:
	@$(BUILDTOOL) $@ $(args)

clone:
	@$(BUILDTOOL) $@ $(args)

commit:
	@$(BUILDTOOL) $@ $(args)

get:
	@$(BUILDTOOL) $@ $(args)

ghr:
	@$(BUILDTOOL) $@ $(args)

go:
	@$(BUILDTOOL) $@ $(args)

pkgreflect:
	@$(BUILDTOOL) $@ $(args)

pull:
	@$(BUILDTOOL) $@ $(args)

push:
	@$(BUILDTOOL) $@ $(args)

release:
	@$(BUILDTOOL) $@ $(args)

selfupdate:
	@$(BUILDTOOL) $@ $(args)

set:
	@$(BUILDTOOL) $@ $(args)

sync:
	@$(BUILDTOOL) $@ $(args)

version:
	@$(BUILDTOOL) $@ $(args)

vfsgen:
	@$(BUILDTOOL) $@ $(args)

