#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
export GO111MODULE=on

for SERVICE_DIR in  $(find $DIR/../microservices -type d -maxdepth 1 -mindepth 1)
do
    cd ${SERVICE_DIR} 
    go get -u
done
