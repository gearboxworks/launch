#
# Standard version level Makefile used to build a Docker container for Gearbox - https://github.com/gearboxworks/gearbox/
#

.ONESHELL:

ifdef TARGET_VERSION

################################################################################
JSONCMD := ./bin/JsonToConfig -json $(TARGET_VERSION)/gearbox.json
JSONTEST := $(shell $(JSONCMD) -template-string '{{ .Json.name }}')
ifeq ($(JSONTEST),)
$(error "# Gearbox: ERROR - No ./bin/JsonToConfig binary.")

else
define GetFromPkg
$(shell ./bin/JsonToConfig -json $(TARGET_VERSION)/gearbox.json -template-string '$(1)')
endef

endif

################################################################################
# Set global variables from container file.
ORGANIZATION	:= $(call GetFromPkg,organization)
NAME		:= $(call GetFromPkg,name)
VERSION		:= $(call GetFromPkg,version)
MAJORVERSION	:= $(call GetFromPkg,majorversion)
LATEST		:= $(call GetFromPkg,latest)
CLASS		:= $(call GetFromPkg,class)
NETWORK		:= $(call GetFromPkg,network)
PORTS		:= $(call GetFromPkg,ports)
VOLUMES		:= $(call GetFromPkg,volumes)
RESTART		:= $(call GetFromPkg,restart)
ARGS		:= $(call GetFromPkg,args)
STATE		:= $(call GetFromPkg,state)
ENV		:= $(call GetFromPkg,env)
GIT_COMMENT	:= Release $(TARGET_VERSION) commit.

# The critical bit. Determines what Dockerfile to build against.
SKIP		:= no
DOCKERFILE	:= $(TARGET_VERSION)/DockerfileRuntime
BUILD_ARGS	:= 
NAME		:= $(NAME)

IMAGE_NAME	?= $(ORGANIZATION)/$(NAME)
CONTAINER_NAME	?= $(NAME)-$(VERSION)
CONTAINER_JSON	?= '$(call rawJSON)'

else
BASEDIR := $(shell pwd)
VERSIONS := $(subst /,, $(sort $(filter-out ./, $(dir $(wildcard *.*/) ) ) ) )
#$(error TARGET_VERSION is not set)
endif


.PHONY: init build push release clean list logs inspect test shell run start stop rm

################################################################################
# Image related commands.
.DEFAULT:
	@make -k help

default: all

all:
	@make -k help

help:
	@echo "################################################################################"
	@echo "init				- Initialize repository from TEMPLATE."
	@echo "update				- Update repository TEMPLATE."
	@echo "git-release			- Generate a repository release."
	@echo ""
	@echo "clean-[all | <VERSION>]		- Clean runtime container image."
	@echo "build-[all | <VERSION>]		- Generate runtime container image."
	@echo "test-[all | <VERSION>]		- Execute container unit tests."
	@echo "push-[all | <VERSION>]		- Push runtime container image to DockerHub & GitHub."
	@echo ""
	@echo "info-[all | <VERSION>]		- Generate info from runtime container image."
	@echo "logs-[all | <VERSION>]		- Show logs from last build."
	@echo "readme				- Show README.md."
	@echo ""
	@echo "release-[all | <VERSION>]	- Execute the following process:"
	@echo "	clean"
	@echo "	build"
	@echo "	test"
	@echo "	push"


################################################################################
#%-all:
#	@./bin/$*.sh all


################################################################################
readme:
	@cat README.md


################################################################################
init: *.json

.FORCE:
%.json: .FORCE
	@echo "Gearbox: Initialize repository."
	@./bin/create-build.sh "all"
	@./bin/create-version.sh "all"


################################################################################
update:
	@./bin/TemplateUpdate.sh
	@make init

git-release:
	@./bin/TemplateRelease.sh


################################################################################
info:
	@./bin/$@.sh $(VERSION)
info-%:
	@./bin/info.sh "$*"


################################################################################
clean:
	@./bin/$@.sh $(VERSION)
clean-%:
	@./bin/clean.sh "$*"


################################################################################
build:
	@./bin/$@.sh $(VERSION)
build-%:
	@./bin/build.sh "$*"


################################################################################
push:
	@./bin/$@.sh $(VERSION)
push-%:
	@./bin/push.sh "$*"


################################################################################
github:
	@./bin/github.sh all


################################################################################
dockerhub:
	@./bin/$@.sh $(VERSION)
dockerhub-%:
	@./bin/dockerhub.sh "$*"


################################################################################
release:
	@./bin/$@.sh $(VERSION)
release-%:
	@./bin/release.sh "$*"


################################################################################
list:
	@./bin/$@.sh $(VERSION)
list-%:
	@./bin/list.sh "$*"


################################################################################
logs:
	@./bin/$@.sh $(VERSION)
logs-%:
	@./bin/logs.sh "$*"


################################################################################
inspect:
	@./bin/$@.sh $(VERSION)
inspect-%:
	@./bin/inspect.sh "$*"


################################################################################
test:
	@./bin/$@.sh $(VERSION)
test-%:
	@./bin/test.sh "$*"


################################################################################
rm:
	@./bin/$@.sh $(VERSION)
rm-%:
	@./bin/rm.sh "$*"


################################################################################
start:
	@./bin/$@.sh $(VERSION)
start-%:
	@./bin/start.sh "$*"


################################################################################
stop:
	@./bin/$@.sh $(VERSION)
stop-%:
	@./bin/stop.sh "$*"


################################################################################
ports:
	@./bin/$@.sh $(VERSION)
ports-%:
	@./bin/ports.sh "$*"


################################################################################
shell:
	@./bin/$@.sh $(VERSION)
shell-%:
	@./bin/shell.sh "$*"


################################################################################
ssh:
	@./bin/$@.sh $(VERSION)
ssh-%:
	@./bin/ssh.sh "$*"


