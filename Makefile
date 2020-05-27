BINARY := $(shell tools/getBinary.sh)
VERSION := $(shell tools/getVersion.sh)
COMMENT := $(shell tools/getComment.sh)

all:
	@echo "Current $(BINARY) version is:	v$(VERSION)"
	@echo "Last commit message is:		'$(COMMENT)'"
	@#echo git tag -a v$(VERSION) -m '"$(COMMENT)"'
	@#echo git push origin v$(VERSION)
	@echo ""
	@echo "build	- Build for local testing."
	@echo "release	- Build for published release."
	@echo "push	- Push repo to GitHub."
	@echo "sync	- Used only by MickMake"

build:
	@goreleaser --snapshot --skip-publish --rm-dist

release:
	@echo "Current $(BINARY) version is v$(VERSION)"
	@git add .
	@git commit -a -m '"Commit before release v$(VERSION)"'
	-@git push
	@git tag -a v$(VERSION) -m '"Release v$(VERSION)"'
	@git push origin v$(VERSION)
	@goreleaser --rm-dist

sync:
	@rsync -HvaxP dist/$(BINARY)_darwin_amd64/$(BINARY) mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Darwin/$(BINARY)
	@rsync -HvaxP dist/$(BINARY)_linux_amd64/$(BINARY) mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Linux/$(BINARY)
	@rsync -HvaxP dist/$(BINARY)_windows_amd64/$(BINARY).exe mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Windows/$(BINARY).exe

push:
	@echo "Pushing to: $(shell git branch)"
	@git config core.hooksPath .git-hooks
	@git add .
	@git commit -m '"$(COMMENT)"' .
	@git push

