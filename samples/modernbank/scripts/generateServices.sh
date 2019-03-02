#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
FLAT_DIR=${DIR}/../contracts


for SERVICE_YAML in  $(find $DIR/../contracts/src -name '*.yaml' -maxdepth 1)
do 
    FILENAME=${SERVICE_YAML##*/}
    SERVICE_NAME=${FILENAME%.*}
    SERVICE_DIR=${DIR}/../microservices/${SERVICE_NAME}

    # Swagger
    swagger validate ${SERVICE_YAML}
    swagger flatten ${SERVICE_YAML} --format yaml -o "${FLAT_DIR}/${SERVICE_NAME}.yaml"

    if [[ ${SERVICE_NAME} == *"-ingress"* ]]
    then
        # Only want a flat/readable contract for ingress. Not whole service.
        continue
    fi

    mkdir -p ${SERVICE_DIR}
    swagger generate server -f ${FLAT_DIR}/${SERVICE_NAME}.yaml --target ${SERVICE_DIR}  --model-package "pkg/model" --server-package "pkg/serve" --api-package "restapi" --flag-strategy pflag
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
