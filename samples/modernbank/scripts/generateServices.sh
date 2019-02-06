#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
FLAT_DIR=${DIR}/flat
mkdir -p ${FLAT_DIR}

for SERVICE_YAML in  $(find $DIR/../contracts -name '*.yaml' -maxdepth 1 | grep -v 'ingress.yaml')
do 
    FILENAME=${SERVICE_YAML##*/}
    SERVICE_NAME=${FILENAME%.*}
    SERVICE_DIR=${DIR}/../microservices/${SERVICE_NAME}

    mkdir -p ${SERVICE_DIR}

    # Swagger
    swagger validate ${SERVICE_YAML}
    swagger flatten ${SERVICE_YAML} --format yaml -o "${FLAT_DIR}/${SERVICE_NAME}.yaml"
    swagger generate server -f ${FLAT_DIR}/${SERVICE_NAME}.yaml --target ${SERVICE_DIR}  --model-package "pkg/model" --server-package "pkg/server" --api-package "restapi"
    swagger generate client -f ${FLAT_DIR}/${SERVICE_NAME}.yaml --target ${SERVICE_DIR} --model-package "pkg/model" --client-package "pkg/client"

    # Dockerfile
    cat > ${SERVICE_DIR}/Dockerfile << EOF
FROM gcr.io/tetratelabs/tetrate-base:v0.1
ADD build/bin/${SERVICE_NAME}-server-static /usr/local/bin/${SERVICE_NAME}-server
ENTRYPOINT [ "/usr/local/bin/${SERVICE_NAME}-server" ]
EOF

    # Go modules
    if [ ! -f ${SERVICE_DIR}/go.mod ]
    then
        export GO111MODULE=on
        cd ${SERVICE_DIR} && go mod init
    fi


done
