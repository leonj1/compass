ifndef VERSION
 VERSION=$(shell git rev-parse --short HEAD)
endif

DOCKER_REPO=registry.digitalocean.com/area51
DOCKER_IMAGE=compass

build:
	go build -v
	upx --no-color --no-progress --best -q compass

docker: build
	docker build -t ${DOCKER_REPO}/${DOCKER_IMAGE}:$(VERSION) .

push: docker
	docker push ${DOCKER_REPO}/${DOCKER_IMAGE}:$(VERSION)

