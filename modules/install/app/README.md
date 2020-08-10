# Installing Demo Application and Istio Proxy Injection

## Overview

The Demo application is using the [microservice-demo](https://github.com/tetratelabs/microservices-demo) app, also known as the Hipstershop app.
It's an application composed of many micro-services, developped using many different languages and protocols.

[![Architecture of microservices](/assets/hipstershop-arch.svg)](/assets/hipstershop-arch.svg)

## Deployment

Let's deploy our entire microservices application.

```shell
kubectl create namespace hipstershopv1v2
kubectl apply -f https://raw.githubusercontent.com/tetratelabs/microservices-demo/master/release/hipstershop-v1-v2.yaml
```

This will:

- Deploy the Hipstershop Demo Store microservices:
- Deploy Redis instances for microservices that store state (cart)
- Create a Kubernetes Service for each microservice
- Create a Kubernetes Service Account for each microservice to be used to prove it's identity later in the security section.

Let’s take a look at the microservice deployments.

```shell
kubectl -n hipstershopv1v2 get deployment

NAME                         READY   UP-TO-DATE   AVAILABLE   AGE
adservice                    1/1     1            1           2m3s
adservice-v2                 1/1     1            1           2m3s
apiservice                   1/1     1            1           2m3s
cartservice                  1/1     1            1           2m3s
checkoutservice              1/1     1            1           2m3s
checkoutservice-v2           1/1     1            1           2m2s
currencyservice              1/1     1            1           2m2s
emailservice                 1/1     1            1           2m2s
frontend                     1/1     1            1           2m2s
...
```
Here we have a bunch of micro-services, all having 1 pod ready out of 1. (Load-generator pod may be crashing at this point as it is not configured)


Now let’s take a look at the frontend microservice deployment.

```yaml
$ kubectl -n hipstershopv1v2 get deployment frontend -o yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: frontend
    project: hipstershopv1v2
    version: v1
  name: frontend
  namespace: hipstershopv1v2
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: frontend
        name: frontend
        project: hipstershopv1v2
        version: v1
    spec:
      containers:
      - env:
        - name: AD_SERVICE_ADDR
          value: adservice.hipstershopv1v2:9555
        - name: CART_SERVICE_ADDR
          value: cartservice.hipstershopv1v2:7070
        - name: CHECKOUT_SERVICE_ADDR
          value: checkoutservice.hipstershopv1v2:5050
        - name: CURRENCY_SERVICE_ADDR
          value: currencyservice.hipstershopv1v2:7000
        - name: PRODUCT_CATALOG_SERVICE_ADDR
          value: productcatalogservice.hipstershopv1v2:3550
        - name: RECOMMENDATION_SERVICE_ADDR
          value: recommendationservice.hipstershopv1v2:8080
        - name: SHIPPING_SERVICE_ADDR
          value: shippingservice.hipstershopv1v2:50051
        - name: SRVURL
          value: :8080
        image: microservicesdemomesh/frontend:v0.1.8
        imagePullPolicy: Always
        livenessProbe:
          failureThreshold: 3
          httpGet:
            httpHeaders:
            - name: Cookie
              value: shop_session-id=x-readiness-probe
            path: /healthz
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: frontend
        ports:
        - containerPort: 8080
          name: http
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            httpHeaders:
            - name: Cookie
              value: shop_session-id=x-readiness-probe
            path: /healthz
            port: 8080
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
      serviceAccountName: frontend

...
```

This is a typical Kubernetes deployment, there is nothing here specific to Istio.

For Istio to intercept and proxy the requests, an Istio sidecar must be installed alongside of the application container. There are 2 ways to do this:

- Manual sidecar injection
- Automatic sidecar injection via a Mutating Admission Webhook.

## Manual Sidecar Injection

Use `istioctl` to see what using manual sidecar injection will add to the deployment.

```shell
$ istioctl kube-inject -f modules/install/app/config/user-v1.yaml | grep "image:" -A2 | head -11
        image: microservicesdemomesh/adservice2:v0.1.8
        imagePullPolicy: Always
        livenessProbe:
--
        image: docker.io/istio/proxyv2:1.6.3
        imagePullPolicy: Always
        name: istio-proxy
--
        image: docker.io/istio/proxyv2:1.6.3
        imagePullPolicy: Always
        name: istio-init
```

In addition to an application container, this output now has an `istio-init` container and an `istio-proxy` container.

The `istio-init` container will set up the IP table rules in the pod’s network namespace to intercept incoming and outgoing connections and direct them to the `istio-proxy` container. The `istio-proxy` container is an Envoy binary wrapped with `pilot-agent` to manage its lifecycle.

If you want to use manual Istio sidecar injection, then you would always filter your Kubernetes deployment file through the `istioctl` utility, and deploy the resulting deployment specification. However, most users tend to take advantage of Istio’s automatic sidecar injection.

## Automatic Sidecar Injection

Auto-Injection will depend on the global `Injection Policy` defined when installing Istio. You can configure it by changing the value `values.global.proxy.autoInject` from the profile:
- `enabled`: all deployments inside a namespace with the label `istio-injection=enabled` will have a sidecar. Setting `sidecar.istio.io/inject = "false"` will prevent the injection. This is the default option in all profiles
- `disabled`: NO sidcar will be added even if a namespace have the label `istio-injection=enabled` set. To enable sidecar auto-injection you need to set *both* `istio-injection=enabled` label on the namespace *and* `sidecar.istio.io/inject = "true"` annotation on the Pod/Deployment


For now we enable Automatic Sidecar Injection on a per Kubernetes namespace level by setting the label `istio-injection` to enabled:

```shell
kubectl label namespace hipstershopv1v2 istio-injection=enabled
```

This instructs Istio to mutate all new pods with an injected sidecar via Kubernetes' Mutating Admission Webhooks. In order to trigger this mutation and inject our Envoy sidecars, delete all the current application's pods (This can also be done using a rolling-restart: `kubectl rollout restart`):

```shell
kubectl -n hipstershopv1v2 delete pods --all
```

Check that all components have the Running status, and that Ready column shows 2/2. This signifies that there are now 2 containers running in each of the pods (the application container and the Istio Proxy container) and that both of these are ready.

```shell
$ kubectl -n hipstershopv1v2 get pods
NAME                                          READY   STATUS             RESTARTS   AGE
adservice-5f87986f57-n6kfm                    2/2     Running            0          43s
adservice-v2-65c57c74db-8j7rp                 2/2     Running            0          43s
apiservice-fb574c858-pstgp                    2/2     Running            0          43s
cartservice-75467b59d-czfzl                   1/2     Running            2          43s
checkoutservice-67949869c-k6qj6               2/2     Running            0          42s
checkoutservice-v2-6bcd4c88b4-c5sk2           2/2     Running            0          42s
currencyservice-745fdb9f89-xwxbq              2/2     Running            0          42s
emailservice-5cb4654868-w9s4l                 2/2     Running            0          42s
frontend-8444d66974-pv284                     2/2     Running            0          42s
frontend-v2-76c878ff5-f5zpg                   2/2     Running            0          42s
...
```

Istio has now automatically injected the sidecar proxy into the pod. You can see this here:

```shell
kubectl -n hipstershopv1v2 get pods -l app=apiservice -o yaml | grep "image:" -A2 | head -11
      image: microservicesdemomesh/apiservice:v0.1.8
      imagePullPolicy: Always
      livenessProbe:
--
      image: docker.io/istio/proxyv2:1.6.3
      imagePullPolicy: Always
      name: istio-proxy
--
      image: docker.io/istio/proxyv2:1.6.3
      imagePullPolicy: Always
      name: istio-init
```

You should see the `istio-init` container as well as a container named `istio-proxy` automatically injected into the pod. Awesome. Now that we have installed Istio and integrated it with our application, we can start to check out what it can do.

## Testing the application

Now that the Hipstershop app is deployed, we can connect to it using our browser. For the moment the shop is not opened to the world and can only be reached internally to the Kubernetes cluster. We can still access it by forwarding the traffic to our local computer. 
We can do this by using the `port-forward` command of `kubectl`:

```shell
kubectl -n hipstershopv1v2 port-forward deployment/frontend 8080
```
You can now open your browser at [http://localhost:8080](http://localhost:8080). Play around to see how the Hipstershop application works. You're done installing the application.

---
Next step: [Observability](modules/observability/)