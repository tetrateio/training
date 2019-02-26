# Istio

## Prerequisites

<!-- TODO: @Liam update once this section is complete
 - kubectl
 - Kubernetes Cluster Context
  
GKE perms one liner: kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user liam@tetrate.io
 -->

## Install

Istio extends Kubernetes using Custom Resource Definitions (CRDs). These enable Kubernetes to store configuration for Istio routing, security and telemetry. To install Istio's CRDs run the following command.

```bash
kubectl apply -f config/crds.yaml
```

Next, install Istio.

```bash
kubectl apply -f config/istio-demo-auth.yaml -f config/kiali-secret.yaml --as=admin --as-group=system:masters
```

Verify Istio components have been installed successfully.

<!-- TODO: @Liam update with 1.1 Istio release -->

```bash
$ kubectl get pods -n istio-system
NAME                                                 READY   STATUS      RESTARTS   AGE
grafana-7b46bf6b7c-6nmxp                             1/1     Running     0          1m
istio-citadel-644bfb5c47-5bhd2                       1/1     Running     0          1m
istio-cleanup-secrets-1.1.0-snapshot.5-2mwf9         0/1     Completed   0          1m
istio-egressgateway-5c5f6dc485-rvdmg                 1/1     Running     0          1m
istio-galley-598f4886cb-bjfzt                        1/1     Running     0          1m
istio-grafana-post-install-1.1.0-snapshot.5-nn785    0/1     Completed   0          1m
istio-ingressgateway-686c8695dd-sqfmb                1/1     Running     0          1m
istio-pilot-5cbc4cfcc8-jcl6v                         1/2     Running     0          1m
istio-policy-844f5656bc-ld9n7                        2/2     Running     0          1m
istio-security-post-install-1.1.0-snapshot.5-qjffc   0/1     Completed   0          1m
istio-sidecar-injector-ff98d6cff-7wbkl               1/1     Running     0          1m
istio-telemetry-6c7f67ffc7-8wn9d                     2/2     Running     0          1m
istio-tracing-759fbf95b7-gf29q                       1/1     Running     0          1m
kiali-647c9649d8-ld7xk                               1/1     Running     0          1m
prometheus-c4b6997b-8prk9                            1/1     Running     0          1m
servicegraph-5dd66d7c7b-266g6                        1/1     Running     0          1m
```
