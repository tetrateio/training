#! /bin/bash -ex

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

for SERVICE_YAML in  $(find $DIR -name '*.yaml')
do 
    FILENAME=${SERVICE_YAML##*/}
    SERVICE_NAME=${FILENAME%.*}
    swagger generate server -f $SERVICE_YAML --target $DIR/../$SERVICE_NAME --model-package "pkg/model" --server-package "pkg/rest" --api-package "api"
done
