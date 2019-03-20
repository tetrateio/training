# Installing Demo Application

Let’s take a look at the user microservice deployment.

```yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: user
  labels:
    app: user
    version: v1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: user
  template:
    metadata:
      labels:
        app: user
        version: v1
    spec:
      serviceAccountName: user
      containers:
        - name: user
          image: "gcr.io/tetratelabs/modernbank/user:v1.0.0"
          imagePullPolicy: Always
          args: ["--port", "8080"]
          ports:
          - name: http
            containerPort: 8080
```

This is a typical Kubernetes deployment, there is nothing here specific to Istio.

However, In order for Istio to intercept and proxy the requests, an Istio sidecar must be installed alongside of the application container. There are 2 ways to do this:

- Manual sidecar injection
- Automatic sidecar injection via a Mutating Admission Webhook.

## Manual Sidecar Injection

Use istioctl to see what using manual sidecar injection will add to the deployment.

```bash
$ istioctl kube-inject -f config/user.yaml | less
...
        image: gcr.io/tetratelabs/modernbank/user:v1.0.0
        name: user
...
        image: docker.io/istio/proxyv2:1.1.0
        name: istio-proxy
...
        image: docker.io/istio/proxy_init:1.1.0
        name: istio-init
```

In addition to an application container, this output now has an `istio-init` container and an `istio-proxy` container.

The `istio-init` container will set up the IP table rules in the pod’s network namespace to intercept incoming and outgoing connections and direct them to the `istio-proxy` container. The `istio-proxy` container is an Envoy binary wrapped with Pilot-Agent to manage its lifecycle.

If you want to use manual Istio sidecar injection, then you would always filter your Kubernetes deployment file through the istioctl utility, and deploy the resulting deployment specification. However, most users tend to take advantage of Istio’s automatic sidecar injection.

## Automatic Sidecar Injection

You can turn on Automatic Sidecar Injection on a per Kubernetes namespace level by setting the label istio-injection to enabled.

```bash
kubectl label namespace default istio-injection=enabled
```

Deploy the entire application.

```bash
kubectl apply -f config/
```

This will:

- Deploy our microservices:
  - User (manages customer information)
  - Account (manages customer accounts)
  - Transaction (orchestrates customer transactions)
  - Transaction-Log (append only log of transactions)
  - UI (serves user interface)
- Deploy MongoDB instances for microservices that store state (user, account, transaction-log).

Check that all components have the Running status, and that Ready column shows 2/2. This signifies that there are 2 containers running in each of the pods (the application container and the Istio Proxy container) and that both of these are running.

```bash
$ kubectl get pods
NAME                                       READY   STATUS        RESTARTS   AGE
account-mongodb-6757fb69c5-cg7hd           2/2     Running       0          1m
account-v1-85598cdcbb-2n9xf                2/2     Running       0          1m
account-v2-76dbb58d49-ks6n7                2/2     Running       0          1m
details-v1-8c85d99c7-ct7pw                 2/2     Running       0          1m
transaction-log-mongodb-556fbc7499-8l4s4   2/2     Running       0          1m
transaction-log-v1-6fd7878454-sn8pr        2/2     Running       0          1m
transaction-v1-84747c99b5-rhrpd            2/2     Running       0          1m
ui-v1-5c498fd8b5-ql69l                     2/2     Running       0          1m
user-mongodb-7867988d7c-nv6xr              2/2     Running       0          1m
user-v1-86d76998b8-hwh7b                   2/2     Running       0          1m
```

Istio has automatically injected the sidecar proxy into the pod. You can see this here:

```bash
$ kubectl get pods -l app=user -o yaml
containerStatuses:
image: istio/proxyv2:1.1.0
      name: istio-proxy
...
   -  image: gcr.io/tetratelabs/modernbank/user:v1.0.0
      name: user
...    
initContainerStatuses:
    - image: istio/proxy_init:1.1.0
      name: istio-init
...
```

You should see the `istio-init` container, and well as a container named `istio-proxy` automatically injected into the pod. Awesome. Now that we have installed Istio, we can start to check out what it can do.
