#! /bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DEMO_APP_DIR=${DIR}/../../../modules/installation/app
TRAFFIC_MAN_DIR=${DIR}/../../../modules/traffic-management

# Deployment YAMLs
for VALUES_YAML in  $(find $DIR/../kubernetes/helm/values/$TENANT -name '*.yaml')
do 
    FILENAME=${VALUES_YAML##*/}
    SERVICE_NAME=${FILENAME%-*}
    helm dependencies update $DIR/../kubernetes/helm/microservice

    if [[ ${SERVICE_NAME} != *"-v2"* ]]
    then
        # Only want normal app for base install, v1 only.
        helm template --name ${SERVICE_NAME} ${DIR}/../kubernetes/helm/microservice -f $VALUES_YAML > ${DEMO_APP_DIR}/config/${SERVICE_NAME}.yaml
    fi
done

# Default Istio Gateway
cp -R ${DIR}/../networking/ingress/ ${TRAFFIC_MAN_DIR}/ingress/config/

# Default Istio VirtualServices
cp ${DIR}/../networking/default/vs.yaml ${TRAFFIC_MAN_DIR}/resiliency/config/default-vs.yaml
