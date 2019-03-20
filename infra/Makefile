ISTIO_VERSION ?= 1.1.0

provision-infra-gke:
	./gke/provision

destroy-infra-gke:
	./gke/destroy

install-istio-gke:
	kubectl apply --as=admin --as-group=system:masters -f ./istio/${ISTIO_VERSION}/kubernetes/crds.yaml
	kubectl apply --as=admin --as-group=system:masters -f ./istio/${ISTIO_VERSION}/kubernetes/istio-demo-auth.yaml
	kubectl apply --as=admin --as-group=system:masters -f ./istio/${ISTIO_VERSION}/kubernetes/kiali-secret.yaml
	kubectl apply --as=admin --as-group=system:masters -f ./istio/${ISTIO_VERSION}/kubernetes/kiali-gateway.yaml
	kubectl label namespace default istio-injection=enabled --overwrite
	kubectl delete pods -n default --all

gke: provision-infra-gke install-istio-gke

delete.istio:
	-kubectl delete $$(kubectl get crds -o name | grep istio)
	-kubectl delete $$(kubectl get clusterrole -o name | grep istio)
	-kubectl delete $$(kubectl get clusterrolebinding -o name | grep istio)
	-kubectl delete $$(kubectl get mutatingwebhookconfiguration -o name | grep istio)
	-kubectl delete namespace istio-system

build-module-config:
	cp ./istio/${ISTIO_VERSION}/kubernetes/* ../modules/installation/istio/config/
