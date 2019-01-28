#! /bin/bash -ex

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

for SERVICE_YAML in  $(find $DIR -name '*.yaml')
do 
    FILENAME=${SERVICE_YAML##*/}
    SERVICE_NAME=${FILENAME%.*}
    swagger generate server -f $SERVICE_YAML --target $DIR/../microservices/$SERVICE_NAME --model-package "pkg/model" --server-package "pkg/server" --api-package "restapi"
        swagger generate client -f $SERVICE_YAML --target $DIR/../microservices/$SERVICE_NAME --model-package "pkg/model" --client-package "pkg/client"

done
