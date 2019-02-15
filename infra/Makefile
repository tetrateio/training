ISTIO_VERSION ?= 1.1.0-snapshot.5

provision-infra-gke:
	./gke/provision

destroy-infra-gke:
	./gke/destroy

install-istio-gke:
	kubectl apply --as=admin --as-group=system:masters -f ./istio/${ISTIO_VERSION}/kubernetes/crds.yaml
	kubectl apply --as=admin --as-group=system:masters -f ./istio/${ISTIO_VERSION}/kubernetes/istio-demo-auth.yaml
	kubectl label namespace default istio-injection=enabled --overwrite
	kubectl delete pods -n default --all

gke: provision-infra-gke install-istio-gke
