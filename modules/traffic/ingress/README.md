# Ingress Traffic

By default, the components deployed on the service mesh are not exposed outside the cluster. In this section, you’ll configure Istio to expose our example UI service outside of the service mesh via the Istio Ingress.

You won’t configure Istio Ingress with the Kubernetes Ingress definitions you would be using in the absence of Istio. The Istio Ingress uses a different configuration model, the Istio Gateway, to allow features such as monitoring and route rules to be applied to traffic entering the cluster.

Note that Istio Gateway only configures the L4-L6 functions (e.g., ports to expose, TLS configuration). Users can then use standard Istio rules to control HTTP requests as well as TCP traffic entering a Gateway by binding it to a `VirtualService`.

We say a `VirtualService` binds to a `Gateway` if the `VirtualService` lists the Gateway’s name in its gateways field and at least one host claimed by the `VirtualService` is exposed by the `Gateway`.

First, lets clean up the Istio config we have already created in previous sections so we can learn by building it back up again!

```shell
kubectl delete virtualservices --all
kubectl delete gateway --all
```

Now, find the Istio Ingress IP address.

```shell
export INGRESS_IP=$(kubectl -n istio-system get svc istio-ingressgateway \
      -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $INGRESS_IP
```

Try to connect to the Istio Ingress IP.

```shell
$ curl $INGRESS_IP
curl: (7) Failed to connect to ... port 80: Connection refused
```

The connection is refused because nothing is binding to this port at this address.

Let’s tell Istio to bind a Gateway to the Istio Ingress.

```shell
kubectl apply -f modules/traffic/ingress/config/gateway.yaml
```

> Note: This action will make the application accessible on the public internet. In some situations, you might consider alternative configurations for privacy:
>
> - Accessing services within the cluster (which does not require the gateway).
> - Accessing internal services across the cluster: consider using GKE private cluster, VPC, IP aliasing, and internal load balancer.
> - Accessing services on the public internet: require credentials, consider using Identity Access Proxy.
>

An Istio Gateway describes a load balancer operating at the edge of the mesh receiving incoming or outgoing HTTP/TCP connections. The specification describes a set of ports that should be exposed, the type of protocol to use, virtual host name to listen to, etc. Here’s the gateway we just created.

```shell
$ kubectl get gateway ingress -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: ingress
  namespace: default
spec:
  selector:
    istio: ingressgateway
  servers:
  - hosts:
    - '*'
    port:
      name: http
      number: 80
      protocol: HTTP
```

As you can see in the spec, this gateway is listening for request that match all hosts on port 80.

Curl the Istio Ingress IP again. Blast! You get a 404 error...

```shell
$ curl -v $INGRESS_IP
...
< HTTP/1.1 404 Not Found
< location: ...
< date: Thu, 21 Mar 2019 14:40:33 GMT
< server: istio-envoy
< content-length: 0
```

It's returning 404 error because once the gateway receives the request we haven’t told it where it needs to send it! To do this we need to attach a `VirtualService` (or two) to the gateway.

```shell
kubectl apply -f modules/traffic/ingress/config/virtualservice.yaml
```

A `VirtualService` defines the rules that control how requests for a service are routed within an Istio service mesh. For example, a `VirtualService` can route requests to different versions of a service or to a completely different service than was requested. Requests can be routed based on the request source and destination, HTTP paths and header fields, and weights associated with individual service versions.

Note that within a `VirtualService`, the match conditions are checked at runtime in the order that they appear. This means that the most specific match clauses should appear first, and less specific clauses later. For safety, a “default” route, with no match conditions, should be provided; a request that does not match any condition of a `VirtualService` will result in a 404 for the sender (or some connection refused error for non-HTTP protocols).

Let’s take a look at our user `VirtualService`.

```shell
$ kubectl get virtualservice user -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: transaction-gateway
spec:
  hosts:
  - "*"
  gateways:
  - ingress
  http:
  - match:
    - uri:
        prefix: /v1/transaction
    route:
    - destination:
        host: transaction
        port:
          number: 80
  - match:
    - uri:
        prefix: /v1/account
    route:
    - destination:
        host: transaction-log
        port:
          number: 80
```

This `VirtualService` is bound to the ingress gateway we just created. It handles API requests for the `/v1/transaction` and `/v1/account` paths on our API, routing all requests to either the transaction or transaction-log microservice depending on their path.

Find the Ingress IP again, and open it up in the browser.

```shell
echo http://$INGRESS_IP
```

While we're here sign up for another account and more free fake money.

![Banking App Home Screen](/assets/banking-app-home.png)
