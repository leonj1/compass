#!/bin/bash

set -e

export CORE_NETWORK=core_net
export PROJECT=compass
export container=${PROJECT}; docker stop $container || true; docker rm $container || true

docker network create core_net || true

export HTTP_INTERNAL=80
export HTTP_EXTERNAL=3244
export DOCKER_IMAGE_TAG=$(python get_docker_build_version.py)

docker run -d --name ${PROJECT} \
-p ${HTTP_EXTERNAL}:${HTTP_INTERNAL} \
-e HTTPPORT=${HTTP_INTERNAL} \
--net ${CORE_NETWORK} \
www.dockerhub.us/${PROJECT}:${DOCKER_IMAGE_TAG}

curl http://localhost:${HTTP_EXTERNAL}/public/health

