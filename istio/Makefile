ISTIO_VERSION ?= 1.0.5

install-istio-gke:
	kubectl apply --as=admin --as-group=system:masters -f ./${ISTIO_VERSION}/kubernetes/crds.yaml
	kubectl apply --as=admin --as-group=system:masters -f ./${ISTIO_VERSION}/kubernetes/istio-demo-auth.yaml
	kubectl label namespace default istio-injection=enabled --overwrite
	kubectl delete pods -n default --all
