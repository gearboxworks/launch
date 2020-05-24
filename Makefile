VERSION := $(shell tools/getVersion.sh)
COMMENT := $(shell tools/getComment.sh)

all:
	@echo "Current launch version is:	v$(VERSION)"
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
	@echo "Current launch version is v$(VERSION)"
	@git add .
	@git commit -a -m "Latest push"
	-@git push
	@git tag -a v$(VERSION) -m '"Release v$(VERSION)"'
	@git push origin v$(VERSION)
	@goreleaser --rm-dist

sync:
	@rsync -HvaxP dist/launch_darwin_amd64/launch mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Darwin/launch
	@rsync -HvaxP dist/launch_linux_amd64/launch mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Linux/launch
	@rsync -HvaxP dist/launch_windows_amd64/launch.exe mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Windows/launch.exe

push:
	@echo "Pushing to: $(shell git branch)"
	@git config core.hooksPath .git-hooks
	@git add .
	@git commit -m '"$(COMMENT)"' .
	@git push

