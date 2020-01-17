# Istio

In this section we’ll get started with Istio on Kubernetes. Istio is infrastructure-agnostic and not tied to Kubernetes, but Kubernetes is the easiest place to run Istio because of its native support of sidecar deployments.

## Installing Istioctl

Download Istio CLI.

```shell
cd ~/
export ISTIO_VERSION=1.4.3
curl -L https://git.io/getLatestIstio | sh -
ln -sf istio-$ISTIO_VERSION istio
```

Add Istio binary path to `$PATH`.

```shell
export PATH=~/istio/bin:$PATH
echo 'export PATH=~/istio/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

## Installing Istio

Let’s install Istio.

``` bash
istioctl manifest apply --set profile=demo --set values.global.mtls.enabled=true --set values.global.controlPlaneSecurityEnabled=true
```

There were a lot of things deployed there so we will break down the important parts and verify that they installed successfully.

Istio extends Kubernetes using Custom Resource Definitions (CRDs). These enable Kubernetes to store configuration for Istio routing, security and telemetry. Let’s verify they were successfully added. Note, this is an abbreviated list covering the more frequently used CRDs, yours will contain a few more.

```shell
$ kubectl get crds | grep istio
adapters.config.istio.io                  2020-01-17T17:54:35Z
attributemanifests.config.istio.io        2020-01-17T17:54:35Z
authorizationpolicies.security.istio.io   2020-01-17T17:54:35Z
clusterrbacconfigs.rbac.istio.io          2020-01-17T17:54:35Z
destinationrules.networking.istio.io      2020-01-17T17:54:36Z
envoyfilters.networking.istio.io          2020-01-17T17:54:36Z
gateways.networking.istio.io              2020-01-17T17:54:36Z
handlers.config.istio.io                  2020-01-17T17:54:36Z
httpapispecbindings.config.istio.io       2020-01-17T17:54:36Z
httpapispecs.config.istio.io              2020-01-17T17:54:37Z
instances.config.istio.io                 2020-01-17T17:54:37Z
meshpolicies.authentication.istio.io      2020-01-17T17:54:37Z
policies.authentication.istio.io          2020-01-17T17:54:37Z
quotaspecbindings.config.istio.io         2020-01-17T17:54:37Z
quotaspecs.config.istio.io                2020-01-17T17:54:38Z
rbacconfigs.rbac.istio.io                 2020-01-17T17:54:38Z
rules.config.istio.io                     2020-01-17T17:54:38Z
serviceentries.networking.istio.io        2020-01-17T17:54:38Z
servicerolebindings.rbac.istio.io         2020-01-17T17:54:38Z
serviceroles.rbac.istio.io                2020-01-17T17:54:38Z
sidecars.networking.istio.io              2020-01-17T17:54:39Z
templates.config.istio.io                 2020-01-17T17:54:39Z
virtualservices.networking.istio.io       2020-01-17T17:54:39Z
```

Next, let’s verify the Istio control plane has installed successfully.

```shell
$ kubectl get pods -n istio-system
NAME                                      READY   STATUS    RESTARTS   AGE
grafana-f5585fb49-s2227                   1/1     Running   0          4m10s
istio-citadel-7db5544b59-s8q8z            1/1     Running   0          4m18s
istio-egressgateway-6f56f8d8c7-xgjqm      1/1     Running   0          4m18s
istio-galley-6b7b79cb59-tbl47             2/2     Running   0          4m16s
istio-ingressgateway-6d495cb95-jsk7f      1/1     Running   0          4m18s
istio-pilot-776b6cc999-p6srr              2/2     Running   0          4m17s
istio-policy-8c67c44bf-tswnk              2/2     Running   2          4m17s
istio-sidecar-injector-5677d969cf-zsswb   1/1     Running   0          4m16s
istio-telemetry-5c9fc5d887-nbwq4          2/2     Running   2          4m17s
istio-tracing-bc44d8d85-xrqwh             1/1     Running   0          4m18s
kiali-767f877b4-zbl54                     1/1     Running   0          4m17s
prometheus-5d56488ff6-5fkm6               1/1     Running   0          4m18s
```

This includes all the of Istio control-plane components (including Istio Sidecar Injector and Istio Gateways), as well as some Istio addons like Prometheus (metric collection), Grafana (metrics dashboard) and Kiali (to visualize how microservices on Istio service mesh are connected).

We can use the experimental `verify-install` command to fully validate that Istio successfully installed. This command may take up to a minute to complete.

```shell
$ istioctl manifest generate --set profile=demo --set values.global.mtls.enabled=true --set values.global.controlPlaneSecurityEnabled=true | istioctl verify-install -f -
Checked 23 crds
Checked 9 Istio Deployments
Istio is installed successfully
```
