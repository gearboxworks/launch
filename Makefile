all:
	@cat Makefile
	@echo git tag -a v1.4 -m "Better error handling"
	@echo git push origin v1.4

build:
	@goreleaser --snapshot --skip-publish --rm-dist

release:
	@goreleaser --rm-dist

sync:
	@rsync -HvaxP dist/launch_darwin_amd64/launch mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Darwin/launch
	@rsync -HvaxP dist/launch_linux_amd64/launch mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Linux/launch
	@rsync -HvaxP dist/launch_windows_amd64/launch.exe mick@macpro:~/Documents/GitHub/containers/docker-template/bin/Windows/launch.exe

push:
	@echo "Pushing to: $(shell git branch)"
	@git config core.hooksPath .git-hooks
	@git add .
	@git commit .
	@git push

