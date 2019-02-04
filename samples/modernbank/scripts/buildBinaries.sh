#! /bin/bash


DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
source $DIR/prettyPrint.sh

export GO111MODULE=on

DOCKER=$1

for SERVICE_DIR in  $(find $DIR/../microservices -type d -maxdepth 1 -mindepth 1)
do  
    SERVICE_NAME=${SERVICE_DIR##*/}
    BUILD_DIR=${SERVICE_DIR}/build/bin
    cd ${SERVICE_DIR}
    prettyPrint "Building ${SERVICE_NAME} binary!"
    GOOS=linux go build \
		-a --ldflags '-extldflags "-static"' \
		-o ${BUILD_DIR}/${SERVICE_NAME}-server-static ./cmd/${SERVICE_NAME}-server
	    chmod +x  ${BUILD_DIR}/${SERVICE_NAME}-server-static
done
