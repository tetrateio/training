#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source $DIR/prettyPrint.sh
CONTAINER_REGISTRY=$1
TAG=$2

for SERVICE_DIR in  $(find $DIR/../microservices -type d -maxdepth 1 -mindepth 1)
do
    SERVICE_NAME=${SERVICE_DIR##*/}
    IMAGE_NAME=${CONTAINER_REGISTRY}/${SERVICE_NAME}:${TAG}
    prettyPrint "Pushing ${IMAGE_NAME}"
    docker push ${IMAGE_NAME}
done
