BRANCH := $(shell git branch | grep \* | cut -d ' ' -f2)
HASH := $(shell git write-tree | cut -c1-6)

$(info $(BRANCH))
$(info $(HASH))

ifeq ($(BRANCH),master)
 VERSION=$(shell git describe --abbrev=0 --tags)
else
  VERSION=$(BRANCH)-$(HASH)
endif

DOCKER_REPO=registry.digitalocean.com/area51
DOCKER_IMAGE=compass

build:
	go build -v -ldflags="-extldflags=-static" -tags sqlite_omit_load_extension
	upx --no-color --no-progress --best -q compass

docker: build
	docker build -t ${DOCKER_REPO}/${DOCKER_IMAGE}:$(VERSION) .

push: docker
	docker push ${DOCKER_REPO}/${DOCKER_IMAGE}:$(VERSION)

