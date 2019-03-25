#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DEMO_APP_DIR=${DIR}/../../../modules/installation/app
TRAFFIC_MAN_DIR=${DIR}/../../../modules/traffic-management

# Deployment YAMLs
helm dependencies update $DIR/../kubernetes/helm/microservice
for VALUES_YAML in  $(find $DIR/../kubernetes/helm/values/$TENANT -name '*.yaml')
do 
    FILENAME=${VALUES_YAML##*/}
    SERVICE_NAME=${FILENAME%-*}

    if [[ ${FILENAME} != *"-v2"* ]]
    then
        helm template --name ${SERVICE_NAME} ${DIR}/../kubernetes/helm/microservice -f $VALUES_YAML > ${DEMO_APP_DIR}/config/${SERVICE_NAME}-v1.yaml
    else
        helm template --name ${SERVICE_NAME} ${DIR}/../kubernetes/helm/microservice -f $VALUES_YAML > ${DEMO_APP_DIR}/config/${SERVICE_NAME}-v2.yaml
    fi
done

# Default Istio Gateway
cp -R ${DIR}/../networking/ingress/ ${TRAFFIC_MAN_DIR}/ingress/config/

