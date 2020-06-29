# Istio

In this section we’ll get started with Istio on Kubernetes. Istio is infrastructure-agnostic and not tied to Kubernetes, but Kubernetes is the easiest place to run Istio because of its native support of sidecar deployments.

## Installing Istioctl

Download Istio CLI.

```shell
cd ~/
export ISTIO_VERSION=1.6.3
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

You have two ways to install Istio: Using the Istio Operator or using `istioctl`.
We will start by using the `istioctl` command line.

Istio Install is based on `profiles` that defines a set of options to apply during the deployment. In the training we will deploy using the `demo` profile. Note that it is strongly advised NOT to use the `demo` profile for anything else than a demo (even on a dev or pre-prod env, you're better using a `default` profile and add the needed resources if needed).


Let’s install Istio.

``` bash
istioctl install --set profile=demo
```

There were a lot of things deployed there so we will break down the important parts and verify that they installed successfully.

Istio extends Kubernetes using Custom Resource Definitions (CRDs). These enable Kubernetes to store configuration for Istio routing, security and telemetry. Let’s verify they were successfully added.

```shell
$ kubectl get crds | grep istio
adapters.config.istio.io                    2020-06-29T19:22:15Z
attributemanifests.config.istio.io          2020-06-29T19:22:15Z
authorizationpolicies.security.istio.io     2020-06-29T19:22:15Z
clusterrbacconfigs.rbac.istio.io            2020-06-29T19:22:15Z
destinationrules.networking.istio.io        2020-06-29T19:22:15Z
envoyfilters.networking.istio.io            2020-06-29T19:22:15Z
gateways.networking.istio.io                2020-06-29T19:22:15Z
handlers.config.istio.io                    2020-06-29T19:22:16Z
httpapispecbindings.config.istio.io         2020-06-29T19:22:16Z
httpapispecs.config.istio.io                2020-06-29T19:22:16Z
instances.config.istio.io                   2020-06-29T19:22:16Z
istiooperators.install.istio.io             2020-06-29T19:22:16Z
peerauthentications.security.istio.io       2020-06-29T19:22:16Z
quotaspecbindings.config.istio.io           2020-06-29T19:22:17Z
quotaspecs.config.istio.io                  2020-06-29T19:22:17Z
rbacconfigs.rbac.istio.io                   2020-06-29T19:22:17Z
requestauthentications.security.istio.io    2020-06-29T19:22:17Z
rules.config.istio.io                       2020-06-29T19:22:18Z
serviceentries.networking.istio.io          2020-06-29T19:22:18Z
servicerolebindings.rbac.istio.io           2020-06-29T19:22:18Z
serviceroles.rbac.istio.io                  2020-06-29T19:22:18Z
sidecars.networking.istio.io                2020-06-29T19:22:18Z
templates.config.istio.io                   2020-06-29T19:22:18Z
virtualservices.networking.istio.io         2020-06-29T19:22:18Z
workloadentries.networking.istio.io         2020-06-29T19:22:19Z
```

Next, let’s verify the Istio control plane has installed successfully.

```shell
$ kubectl get pods -n istio-system
NAME                                    READY   STATUS    RESTARTS   AGE
grafana-5dc4b4676c-f47gp                1/1     Running   0          63s
istio-egressgateway-5db676495d-ffrss    1/1     Running   0          65s
istio-ingressgateway-69bb66cb4b-2qkr4   1/1     Running   0          65s
istio-tracing-8584b4d7f9-bctz2          1/1     Running   0          63s
istiod-5cdccfd474-qlzvm                 1/1     Running   0          78s
kiali-6f457f5964-zgp9s                  1/1     Running   0          63s
prometheus-d8b7c5949-fqdvj              2/2     Running   0          63s
```

> Some components may take a couple of minutes to become fully ready. Re-run the command until everything is ready.

This includes all the of Istio control-plane components (including Istio Sidecar Injector and Istio Gateways), as well as some Istio addons like Prometheus (metric collection), Grafana (metrics dashboard) and Kiali (to visualize how microservices on Istio service mesh are connected).

We can use the experimental `verify-install` command to fully validate that Istio successfully installed. This command may take up to a minute to complete.

```shell
$ istioctl manifest generate --set profile=demo | istioctl verify-install -f -
...
Checked 23 crds
Checked 3 Istio Deployments
Istio is installed successfully
```
