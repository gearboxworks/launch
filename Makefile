all:
	@cat Makefile

build:
	@goreleaser --snapshot --skip-publish --rm-dist

sync:
	@rsync -HvaxP dist/gb-launch_darwin_amd64/gb-launch mick@macpro:~/Documents/GitHub/containers/docker-template/bin/gb-launch-Darwin
	@rsync -HvaxP dist/gb-launch_linux_amd64/gb-launch mick@macpro:~/Documents/GitHub/containers/docker-template/bin/gb-launch-Linux

push:
	@echo "Pushing to: $(shell git branch)"
	@git config core.hooksPath .git-hooks
	@git add .
	@git commit .
	@git push

