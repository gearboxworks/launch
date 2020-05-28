################################################################################
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

all:
	@echo "build		- Build for local testing."
	@echo "release		- Build for published release."
	@echo "push		- Push repo to GitHub."
	@echo ""
	@$(BUILDTOOL) get all

build:
	#@make pkgreflect
	@$(BUILDTOOL) build

release:
	#@make pkgreflect
	@$(BUILDTOOL) release

push:
	#@make pkgreflect
	@$(BUILDTOOL) push

pkgreflect:
	#@$(BUILDTOOL) pkgreflect jtc/helpers

sync:
	@echo "sync		- Used only by MickMake"
	@rsync -HvaxP dist/$(BINARY)_darwin_amd64/$(BINARY) mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Darwin/$(BINARY)
	@rsync -HvaxP dist/$(BINARY)_linux_amd64/$(BINARY) mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Linux/$(BINARY)
	@rsync -HvaxP dist/$(BINARY)_windows_amd64/$(BINARY).exe mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Windows/$(BINARY).exe

