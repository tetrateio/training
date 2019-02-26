#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

TENANT=$1
NAMESPACE=$2

for VALUES_YAML in  $(find $DIR/../kubernetes/helm/values/$TENANT -name '*.yaml')
do 
    FILENAME=${VALUES_YAML##*/}
    SERVICE_NAME=${FILENAME%-*}
    helm dependencies update $DIR/../kubernetes/helm/microservice
    helm template --name ${SERVICE_NAME} --namespace $NAMESPACE $DIR/../kubernetes/helm/microservice -f $VALUES_YAML | kubectl apply --as=admin --as-group=system:masters -f -
done
