#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DEMO_APP_DIR=${DIR}/../../../modules/installation/app

TENANT=$1

# Deployment YAMLs
for VALUES_YAML in  $(find $DIR/../kubernetes/helm/values/$TENANT -name '*.yaml')
do 
    FILENAME=${VALUES_YAML##*/}
    SERVICE_NAME=${FILENAME%-*}
    helm dependencies update $DIR/../kubernetes/helm/microservice

    if [[ ${SERVICE_NAME} != *"-v2"* ]]
    then
        # Only want normal app for base install, v1 only.
        helm template --name ${SERVICE_NAME} ${DIR}/../kubernetes/helm/microservice -f $VALUES_YAML > ${DEMO_APP_DIR}/short/config/${SERVICE_NAME}.yaml
    fi
done

# Default Istio Gateway
cp ${DIR}/../networking/${TENANT}/gateway.yaml ${DEMO_APP_DIR}/short/config/gateway.yaml

# Ingress Swagger/API Contract
cat ${DIR}/../contracts/banking-ingress.yaml | yq . > ${DEMO_APP_DIR}/short/swagger/ingress.json
