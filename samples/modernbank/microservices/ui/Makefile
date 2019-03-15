# Copyright (c) Tetrate, Inc 2019 All Rights Reserved.

CONTAINER_REGISTRY ?= gcr.io/tetratelabs/modernbank
NAME ?= ui
TAG ?= v1.0.0

# The following are for docker run.
API_HOST ?= training-stag.tetr8labs.net
API_PORT ?= 80
HOST_PORT ?= 8000

docker-build: build
	docker build -t $(CONTAINER_REGISTRY)/$(NAME):$(TAG) -f Dockerfile .

docker-push:
	docker push $(CONTAINER_REGISTRY)/$(NAME):$(TAG)

docker-run:
	docker run --rm -it -p $(HOST_PORT):80 \
		-e API_HOST=$(API_HOST) \
		-e API_PORT=$(API_PORT) \
		$(CONTAINER_REGISTRY)/$(NAME):$(TAG)

build:
	npm run build

dev:
	npm start

clean:
	rm -fr build

.PHONY: docker-build docker-push docker-run clean
