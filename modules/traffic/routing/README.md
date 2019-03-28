# Routing

Traffic splitting allows you to distribute specific percentages of traffic across two or more versions of a service. This allows you to conduct A/B testing between versions and also allows for controlled, paced roll-out of untested features (i.e., canary deployment). Traffic splitting is achieved using routing rules that identify at least one weighted backend that corresponds to a specific version of the destination service that is expressed using labels.

If there are multiple registered instances associated with the label/s, routing will be based on the load balancing policy configured for the service, or round-robin by default. There are different load balancing methods you can configure via Destination Rule's Load Balancer Setting, e.g., Round Robin, Least Connection, and Random. It can also load balance with consistent hashing to achieve affinity by using Cookie, Header, or Source IP.

You may have noticed in previous sections two different emojis appearing in the banner of the account lists page: 1️⃣ and 2️⃣. If you haven’t, log in to your account and refresh the account list a couple of times. This is because we have two versions of the user microservice deployed and serving traffic to the UI. We can confirm this by checking the pods that are running.

```shell
$ kubectl get pods
NAME                                       READY   STATUS    RESTARTS   AGE
...
user-v1-86d76998b8-kvlrh                   2/2     Running   0          1d
user-v2-659f745bbb-8tg5d                   2/2     Running   0          1d
```

Both are serving traffic because we are using the default Kubernetes round robin load balancing. However, with Istio we can take advantage of more meaningful release patterns such as canary deploys, and we can do so decoupled from the application code, using configuration.

In order to control traffic we need to tell Istio how to distinguish between our two versions, we do this with a `DestinationRule`.

DestinationRules are in fact all about configuring clients. They allow a service operator to describe how a client in the mesh should call their service, including: subsets of the service (e.g. v1 and v2), the load balancing strategy the client should use, the conditions to use to mark endpoints of the service as unhealthy, L4 and L7 connection pool settings, and TLS settings for the server.

In this case, we have two versions of the user microservice, indicated by the `version` labels. Let’s add a `DestinationRule` to tell Istio how to identify the different versions we have deployed.

```shell
kubectl apply -f modules/traffic/routing/config/user-subset.yaml
```

The `DestinationRule` we just created describes two versions of the user service and how to identify which is which by creating a subset.

```shell
$ kubectl get destinationrule user -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: user
spec:
  host: user
  subsets:
  - name: v1
    labels:
      version: v1
  - name: v2
    labels:
      version: v2
```

Now that we have subsets defined we can pin all traffic to version 1 of the user service in its `VirtualService` (we will go over Virtual Services in a bit more detail in the resiliency section).

```shell
kubectl apply -f modules/traffic/routing/config/user-v1-100.yaml
```

If we take a look at the changes we just made to the `VirtualService`, we can see that there are now two destinations, one for each subset. However, we are weighting 100% of requests to version 1.

```shell
$ kubectl get virtualservice user-gateway -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user-gateway
spec:
  gateways:
  - ingress
  hosts:
  - '*'
  http:
  - match:
    - uri:
        prefix: /v1/users
    route:
    - destination:
        host: user
        port:
          number: 80
        subset: v1
      weight: 100
    - destination:
        host: user
        port:
          number: 80
        subset: v2
      weight: 0
```

Now if we refresh the accounts page several times we can see that we always get the same version with the 1️⃣ emoji.

We can also go in a change the weighting. Let’s edit the `VirtualService` to route 25% of traffic to version two. Note that the weight’s have to add up to 100 otherwise they will fail validation.

```shell
$ kubectl edit virtualservice user-gateway
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user-gateway
spec:
  gateways:
  - ingress
  hosts:
  - '*'
  http:
  - match:
    - uri:
        prefix: /v1/users
    route:
    - destination:
        host: user
        port:
          number: 80
        subset: v1
      weight: 75
    - destination:
        host: user
        port:
          number: 80
        subset: v2
      weight: 25
```

This enables you to release new versions to increments of a single percent of users. Istio also allows you to route traffic based on headers, URI or HTTP method. This enables you to do things like releasing to a single browser. Let’s release v2 to Chrome users only.

```shell
kubectl apply -f modules/traffic/routing/config/user-v2-chrome.yaml
```

If we look at the configuration we just created we can see that we added a stanza to the path match to search for the word `Chrome` anywhere in the user-agent. Note that ordering of the matches matters. Envoy will route to the first match it sees, not the most specific.

```shell
$ kubectl get virtualservice user-gateway -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user-gateway
spec:
  gateways:
  - ingress
  hosts:
  - '*'
  http:
  - match:
    - headers:
        user-agent:
          regex: .*Chrome.*
      uri:
        prefix: /v1/users
    route:
    - destination:
        host: user
        port:
          number: 80
        subset: v2
  - match:
    - uri:
        prefix: /v1/users
    route:
    - destination:
        host: user
        port:
          number: 80
        subset: v1
```

Now if you visit the application from Chrome you will always be routed to version 2 of the user service, but if you visit from any other browser you will be routed to version 1.

Let’s reset our VirtualService back to the default behavior.

```shell
kubectl apply -f modules/traffic/routing/config/virtualservice-reset.yaml
```
