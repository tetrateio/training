
CONTAINER_REGISTRY ?= gcr.io/tetratelabs/modernbank
TAG ?= v1.0.0

MONGO_REPOSITORY ?= bitnami/mongodb
MONGO_TAG ?= 4.0.6

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

docker-push-third-party:
	docker pull ${MONGO_REPOSITORY}:${MONGO_TAG}
	docker tag ${MONGO_REPOSITORY}:${MONGO_TAG} gcr.io/tetratelabs/modernbank/mongodb:${MONGO_TAG}
	docker push gcr.io/tetratelabs/modernbank/mongodb:${MONGO_TAG}

deploy-kube:
	./scripts/deployKube.sh bankingservices default
	kubectl apply -f networking/bankingservices/gateway.yaml
