# Istio

In this section we’ll get started with Istio on Kubernetes. Istio is infrastructure-agnostic and not tied to Kubernetes, but Kubernetes is the easiest place to run Istio because of its native support of sidecar deployments.

## Installing Istio

Let’s install Istio.

``` bash
kubectl apply -f ./modules/install/istio/config/istio-demo.yaml --as=admin --as-group=system:masters
```

> We use `--as=admin --as-group=system:masters` to escalate our privilege while installing Istio, which is required to configure the webhook Istio uses to automate sidecar injection. We'll talk about sidecar injection in depth later.

There were a lot of things deployed there so we will break down the important parts and verify that they installed successfully.

Istio extends Kubernetes using Custom Resource Definitions (CRDs). These enable Kubernetes to store configuration for Istio routing, security and telemetry. Let’s verify they were successfully added. Note, this is an abbreviated list covering the more frequently used CRDs, yours will contain a few more.

```bash
$ kubectl get crds | grep istio | head
authorizations.config.istio.io          2019-03-19T20:41:53Z
clusterrbacconfigs.rbac.istio.io        2019-03-19T20:41:44Z
envoyfilters.networking.istio.io        2019-03-19T20:41:44Z
gateways.networking.istio.io            2019-03-19T20:41:43Z
meshpolicies.authentication.istio.io    2019-03-19T20:41:45Z
metrics.config.istio.io                 2019-03-19T20:41:55Z
policies.authentication.istio.io        2019-03-19T20:41:44Z
rbacconfigs.rbac.istio.io               2019-03-19T20:41:56Z
rbacs.config.istio.io                   2019-03-19T20:41:50Z
rules.config.istio.io                   2019-03-19T20:41:46Z
serviceentries.networking.istio.io      2019-03-19T20:41:43Z
sidecars.networking.istio.io            2019-03-19T20:41:59Z
virtualservices.networking.istio.io     2019-03-19T20:41:42Z
```

Next, let’s verify the Istio control plane has installed successfully.

```bash
$ kubectl get pods -n istio-system
NAME                                      READY   STATUS      RESTARTS   AGE
grafana-7b46bf6b7c-hbr74                  1/1     Running     0          11m
istio-citadel-75fdb679db-wjzbg            1/1     Running     0          11m
istio-cleanup-secrets-1.1.0-kjr92         0/1     Completed   0          11m
istio-egressgateway-78759f4cd8-g8vqr      1/1     Running     0          11m
istio-galley-59b7b685f-vvkqx              1/1     Running     0          11m
istio-grafana-post-install-1.1.0-gjf7f    0/1     Completed   0          11m
istio-ingressgateway-55b87dc69c-dzwbt     1/1     Running     0          11m
istio-pilot-5694cb4ff-qfmf2               2/2     Running     0          11m
istio-policy-8df48ccbd-ph7fg              2/2     Running     0          11m
istio-security-post-install-1.1.0-r2hz4   0/1     Completed   0          11m
istio-sidecar-injector-7b47cb4689-zkkhb   1/1     Running     0          11m
istio-telemetry-64f5b9d4d8-zbv6j          2/2     Running     0          11m
istio-tracing-75dd89b8b4-fcdxj            1/1     Running     0          11m
kiali-5d68f4c676-n4bdw                    1/1     Running     0          11m
prometheus-89bc5668c-crdzs                1/1     Running     0          11m
```

This includes all the of Istio control-plane components (including Istio Sidecar Injector and Istio Gateways), as well as some Istio addons like Prometheus (metric collection), Grafana (metrics dashboard) and Kiali (to visualize how microservices on Istio service mesh are connected).

## Installing Istioctl

Istio also has its own command line tool for debugging, verifying configuration, manually injecting sidecars and various other things. We can also use istioctl to validate everything has installed correctly.
 
Download Istio CLI.

```bash
cd ~/
export ISTIO_VERSION=1.1.0
curl -L https://git.io/getLatestIstio | sh -
ln -sf istio-$ISTIO_VERSION istio
```

Add Istio binary path to `$PATH`.

```bash
export PATH=~/istio/bin:$PATH
echo 'export PATH=~/istio/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

We can use the experimental `verify-install` command to fully validate that Istio successfully installed. This command may take up to a minute to complete.

```bash
$ istioctl experimental verify-install -f config/istio-demo-auth.yaml
Istio is installed successfully
```
