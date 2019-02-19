#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

TENANT=$1

# Deployment YAMLs
for VALUES_YAML in  $(find $DIR/../kubernetes/helm/values/$TENANT -name '*.yaml')
do 
    FILENAME=${VALUES_YAML##*/}
    SERVICE_NAME=${FILENAME%.*}
    helm dependencies update $DIR/../kubernetes/helm/microservice

    helm template --name ${SERVICE_NAME} ${DIR}/../kubernetes/helm/microservice -f $VALUES_YAML > ${DIR}/../../../modules/installation/demo-app-short/config/${SERVICE_NAME}.yaml
done

# Default Istio Gateway
cp ${DIR}/../networking/${TENANT}/gateway.yaml ${DIR}/../../../modules/installation/demo-app-short/config/gateway.yaml

# Ingress Swagger/API Contract
cat ${DIR}/../contracts/banking-ingress.yaml | yq . > ${DIR}/../../../modules/installation/demo-app-short/swagger/ingress.json
