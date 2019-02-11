
CONTAINER_REGISTRY ?= gcr.io/tetratelabs/modernbank
TAG ?= v1.0.0

generate-services:
	./scripts/generateServices.sh

update-dependencies:
	./scripts/updateDependencies.sh

build:
	./scripts/buildBinaries.sh

docker-build: build
	./scripts/buildImages.sh ${CONTAINER_REGISTRY} ${TAG}

docker-push: build docker-build
	./scripts/pushImages.sh ${CONTAINER_REGISTRY} ${TAG}

deploy-kube:
	./scripts/deployKube.sh bankingservices default
	kubectl apply -f networking/bankingservices/gateway.yaml
