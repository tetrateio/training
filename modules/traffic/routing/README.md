# Routing

Traffic splitting allows you to distribute specific percentages of traffic across two or more versions of a service. This allows you to conduct A/B testing between versions and also allows for controlled, paced roll-out of untested features (i.e., canary deployment). Traffic splitting is achieved using routing rules that identify at least one weighted backend that corresponds to a specific version of the destination service that is expressed using labels.

If there are multiple registered instances associated with the label/s, routing will be based on the load balancing policy configured for the service, or round-robin by default. There are different load balancing methods you can configure via Destination Rule's Load Balancer Setting, e.g., Round Robin, Least Connection, and Random. It can also load balance with consistent hashing to achieve affinity by using Cookie, Header, or Source IP.

You may have noticed in previous sections two different banner color, one grey and one red. If you did not noticed, go back to the Hipstershop and reload the page multiple time.
We actually deployed two versions of the `frontend` microservice, a `v1` and a `v2`:

```shell
 kubectl -n hipstershopv1v2 get pods -l app=frontend

NAME                          READY   STATUS    RESTARTS   AGE
frontend-8444d66974-pv284     2/2     Running   0          2d20h
frontend-v2-76c878ff5-f5zpg   2/2     Running   0          2d20h
```

Both are serving traffic because we are using the default Kubernetes round robin load balancing. However, with Istio we can take advantage of more meaningful release patterns such as canary deploys, and we can do so decoupled from the application code, using configuration.

In order to control traffic we need to tell Istio how to distinguish between our two versions, we do this with a `DestinationRule`.

DestinationRules are in fact all about configuring clients. They allow a service operator to describe how a client in the mesh should call their service, including: subsets of the service (e.g. v1 and v2), the load balancing strategy the client should use, the conditions to use to mark endpoints of the service as unhealthy, L4 and L7 connection pool settings, and TLS settings for the server.

In this case, we have two versions of the `frontend` microservice, indicated by the `version` labels. Let’s add a `DestinationRule` to tell Istio how to identify the different versions we have deployed.


```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: frontend-subset
spec:
  host: frontend.hipstershopv1v2.svc.cluster.local
  subsets:
  - labels:
      version: v1
    name: v1
  - labels:
      version: v2
    name: v2
EOF
```

The `DestinationRule` we just created describes two versions of the frontend service and how to identify which is which by creating a subset based on the `version` Label.


Now that we have subsets defined we can pin all traffic to version 1 of the user service in its `VirtualService` (we will go over Virtual Services in a bit more detail in the resiliency section).

```shell
kubectl apply -f modules/traffic/routing/config/user-v1-100.yaml
```

If we take a look at the changes we just made to the `VirtualService`, we can see that there are now two destinations, one for each subset. However, we are weighting 100% of requests to version 1.

```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hipstershop
spec:
  hosts:
  - "hipstershop.${INGRESSIP}.sslip.io"
  gateways:
  - hipstershop
  http:
  - match:
    - uri:
        prefix: /api
    route:
    - destination:
        host: apiservice
        port:
          number: 8080
  - route:
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v1
      weight: 100
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v2
      weight: 0
EOF
```

Now if we refresh the Hipstershop page several times we can see that we always get the same version with the grey banner.

We can also go in a change the weighting. Let’s edit the `VirtualService` again to route 25% of traffic to version two. Note that the weight’s have to add up to 100 otherwise they will fail validation.

```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hipstershop
spec:
  hosts:
  - "hipstershop.${INGRESSIP}.sslip.io"
  gateways:
  - hipstershop
  http:
  - match:
    - uri:
        prefix: /api
    route:
    - destination:
        host: apiservice
        port:
          number: 8080
  - route:
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v1
      weight: 75
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v2
      weight: 25
EOF
```

This enables you to release new versions to increments of a single percent of users. Istio also allows you to route traffic based on headers, URI or HTTP method. This enables you to do things like releasing to a single browser. Let’s release v2 to Chrome users only.

```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hipstershop
spec:
  hosts:
  - "hipstershop.${INGRESSIP}.sslip.io"
  gateways:
  - hipstershop
  http:
  - match:
    - uri:
        prefix: /api
    route:
    - destination:
        host: apiservice
        port:
          number: 8080
  - match:
    - headers:
        user-agent:
          regex: ".*Chrome.*"
    route:
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v2
  - route:
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v1
EOF
```

If we look at the configuration we just created we can see that we added a stanza to the headers match to search for the word `Chrome` anywhere in the user-agent. Note that ordering of the matches matters. Envoy will route to the first match it sees, not the most specific. We keep the default route at the end, which sends 100% of the traffic to the `v1` version. 


Now if you visit the application from Chrome you will always be routed to version 2 of the frontend service, with the red banner, but if you visit from any other browser you will be routed to version 1 with the grey banner.

Let’s reset our VirtualService back to the default behavior, with a 50/50 between v1 and v2.

```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hipstershop
spec:
  hosts:
  - "hipstershop.${INGRESSIP}.sslip.io"
  gateways:
  - hipstershop
  http:
  - match:
    - uri:
        prefix: /api
    route:
    - destination:
        host: apiservice
        port:
          number: 8080
  - route:
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v1
      weight: 50
    - destination:
        host: frontend
        port:
          number: 8080
        subset: v2
      weight: 50
EOF
```

What will happen  when we stop the `v2` pods ? Let's try it out:

```shell
kubectl -n hipstershopv1v2  scale --replicas=0 deployment/frontend-v2
```

Wait few seconds for the pod to be terminated, then reload the front page in your browser. 

You should see `no healthy upstream` messages every other request.

What's happening here ?

Envoy manages a list of endpoints for every destination it knows. In our case, we defined two destinations, v1 and v2, then we terminated the v2 endpoints. 
In this scenario, Envoy still send 50% requests to the `v2` destination, as requested, and then return an error.

To further diagnose, we can check the `hipstershop-ingressgateway logs`:

```yaml
hipstershop-ingressgateway-7c95d4fbd6-xzs47 istio-proxy [2020-07-02T17:07:07.619Z] "GET /product/2ZYFJ3GM2N HTTP/1.1" 503 UH "-" "-" 0 19 0 - "10.56.0.1" "python-requests/2.21.0" "1fe3d871-4b0d-4e6d-b568-23195b67c7c6" "hipstershop.35.245.245.42.sslip.io" "-" - - 10.56.0.46:8443 10.56.0.1:39760 hipstershop.35.245.245.42.sslip.io -
```
Envoy is reporting a `503` error code with a `UF` flag. Envoy always add a respose flag in case of error. In this case:
- UH: No healthy upstream hosts in upstream cluster in addition to 503 response code.

You can learn more about Envoy Response Flags further down the page [here](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#configuration).

Once you're done playing, don't forget to restart the `v2` frontend:

```shell
kubectl -n hipstershopv1v2  scale --replicas=1 deployment/frontend-v2
```