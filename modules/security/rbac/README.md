# Istio RBAC for Service to Service Communication

Istio provides every workload with a strong identity which is used to establish mTLS connections between services in the mesh. In Kubernetes, the workload's identity is the pod's `ServiceAccount`. While establishing mTLS connections, sidecars in the mesh will validate certificates according to [SPIFFE](https://github.com/spiffe/spiffe/blob/master/standards/X509-SVID.md), which means that after the connection is established we have the identity (_authenticated principal_) of the other party. Then, Istio allows you to write access control policies using those identities to describe which services can communicate.

Note that it means that you have to enforce mTLS between your apps for Authorization Policies to work.

## Setup a deny-all policy
First, let's start by locking down our application:
```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: deny-all
  namespace: hipstershopv1v2
spec:
  {}
EOF
```

At this moment, no traffic can go through the Mesh. If you try to browser the Hipstershop, you will see an error `RBAC: access denied` and some `error 403` in the logs.

## Allow the Ingress Gateway to talk to the Frontend

We set ourselfs out by applying a very restrictive policy. Now we need to clear communications we want to allow inside the mesh. As the policy is enforced at every Istio Proxy, we need two new policies: from outside to the Ingress Gateway and from the Ingress Gateway to the Frontend service. We are going to only allow `HTTP GET` requests first:

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: ingress-policy
spec:
  selector:
    matchLabels:
      app: istio-ingressgateway
  action: ALLOW
  rules:
  - to:
    - operation:
        methods:
        - GET
---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: ingress-policy-to-frontend
spec:
  selector:
    matchLabels:
      app: frontend
  action: ALLOW
  rules:
  - to:
    - operation:
        methods:
        - GET
EOF
```

If you browser the Hipstershop store now, you can access the frontend. Sadly, you are still getting an error:

```yaml
Uh, oh!
Something has failed. Below are some details for debugging.

HTTP Status: 500 Internal Server Error
```

If we look at the logs of the `frontend` service, we can see we still have a DENY: `could not retrieve currencies: rpc error: code = PermissionDenied desc = RBAC: access denied`

Let's add another policy for the `currency` service. This time it's an `HTTP POST` request (all GRPC requests are made using `POST`):

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: frontend-policy-to-currency
spec:
  selector:
    matchLabels:
      app: currencyservice
  action: ALLOW
  rules:
  - to:
    - operation:
        methods:
        - POST
EOF
```

If you browse again, you'll see the error has changed: `RBAC: access denied could not retrieve products`. It's now the `productservice` that can't be reached.

We could go on and add one policy for each services. Some companies will require this level of control.
Let's add a global rule for our demo. First, remove the 3 policies we just created:

```shell
kubectl -n hipstershopv1v2 delete AuthorizationPolicy frontend-policy-to-currency ingress-policy-to-frontend ingress-policy
```

And let's add a global one


```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: ingress-policy
spec:
  selector:
    matchLabels:
      app: istio-ingressgateway
  action: ALLOW
  rules:
  - to:
    - operation:
        methods:
        - GET
        - POST
---
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: policy-for-hipstershop
spec:
  selector:
    matchLabels:
      project: hipstershopv1v2
  action: ALLOW
  rules:
  - from:
    - source:
       namespaces: ["hipstershopv1v2"]
    to:
    - operation:
        methods:
        - POST
        - GET
EOF
```

If you browse again, you will see... another error ! What's wrong with this application ? 
Look at the logs of the `frontend` service: `POST /hipstershop.CartService/GetCart HTTP/2" 200 UH`. `UH` stands for Upstream Health. This log is telling us the `cartservice` is not healthy. Let's look at the pod's status:

```shell

kubectl -n hipstershopv1v2 get pods

NAME                                          READY   STATUS             RESTARTS   AGE
adservice-595b9b659b-qjm54                    2/2     Running            0          7d1h
adservice-v2-68d98869d9-jdb5d                 2/2     Running            0          7d
apiservice-fb574c858-pstgp                    2/2     Running            0          13d
cartservice-75467b59d-tft7t                   1/2     CrashLoopBackOff   16         2d20h
checkoutservice-67949869c-k6qj6               2/2     Running            0          13d
checkoutservice-v2-6bcd4c88b4-zf9mf           2/2     Running            0          10d
currencyservice-745fdb9f89-xwxbq              2/2     Running            0          13d
emailservice-5cb4654868-w9s4l                 2/2     Running            0          13d
frontend-8444d66974-pv284                     2/2     Running            0          13d
hipstershop-ingressgateway-7c95d4fbd6-xzs47   1/1     Running            0          11d
paymentservice-5c666d5fc7-kxf96               2/2     Running            0          13d
productcatalogservice-5df7987d54-wmdcx        2/2     Running            0          13d
productcatalogservice-slow-7497958d96-q2n9v   2/2     Running            0          13d
recommendationservice-648999bcb-8qr6d         2/2     Running            0          13d
redis-cart-5fddd6bd44-nrz5g                   2/2     Running            0          2d20h
shippingservice-76946cf5f6-s5g9j              2/2     Running            0          13d
```

Well, Istio is right, `cartservice`'s pod keeps crashing.
Let's look at this pod's logs:

```shell
kubectl -n hipstershopv1v2 logs cartservice-75467b59d-tft7t cartservice

Started as process with id 1
Reading host address from LISTEN_ADDR environment variable
Reading cart service port from PORT environment variable
Reading redis cache address from environment variable REDIS_ADDR
Connecting to Redis: redis-cart.hipstershopv1v2:6379,ssl=false,allowAdmin=true,connectRetry=5
StackExchange.Redis.RedisConnectionException: It was not possible to connect to the redis server(s). SocketClosed (ReadEndOfStream, 0-read, last-recv: 0) on redis-cart.hipstershopv1v2:6379/Interactive, Flushed/MarkProcessed, last: ECHO, origin: ReadFromPipe, outstanding: 6, last-read: 0s ago, last-write: 0s ago, keep-alive: 180s, state: ConnectedEstablishing, mgr: 8 of 10 available, last-heartbeat: never, global: 1s ago, v: 2.0.601.3402
   at StackExchange.Redis.ConnectionMultiplexer.ConnectImpl(Object configuration, TextWriter log) in C:\projects\stackexchange-redis\src\StackExchange.Redis\ConnectionMultiplexer.cs:line 955
   at cartservice.cartstore.RedisCartStore.EnsureRedisConnected() in /app/cartstore/RedisCartStore.cs:line 80
   at cartservice.cartstore.RedisCartStore.InitializeAsync() in /app/cartstore/RedisCartStore.cs:line 60
   at cartservice.Program.<>c__DisplayClass4_0.<<StartServer>b__0>d.MoveNext() in /app/Program.cs:line 54
```

This pod can't connect to the `redis-cart` service.

Oh ohhhhhh, Could it be that we broke the production while adding security rules ? The Redis protocol is not HTTP, and does not fall in the `GET`/`POST` category.

As the production is broken, let's quickly add a policy to allow `cartservice` to call the Redis port `6379`:

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
  name: redis-policy
spec:
  action: ALLOW
  rules:
  - from:
    - source:
        principals: ["cluster.local/ns/hipstershopv1v2/sa/cartservice"]
    to:
    - operation:
        ports: ["6379"]
EOF
```

Once applied the `cartservice` pod should be back up. If not, you can kill it to recycle it: `kubectl -n hipstershopv1v2 delete pod -l app=cartservice`. Few seconds later your should have a running pod, and the Hipstershop should be back online.

## Cleanup

    For now, lets clean up our Cluster Policy Config so we can carry on with the rest of the lab:

    ```shell
    kubectl -n hipstershopv1v2 delete AuthorizationPolicy deny-all ingress-policy policy-for-hipstershop redis-policy
    ```

## Takeaway

Once mTLS is enforced, we can use Istio create Authorization Policies to control the traffic in the mesh. While this is powerful, it may also lead to unintended traffic breaks. Plan, review and tests your rules before applying them in production!

---
This is the end of the Istio Workshop for now.
