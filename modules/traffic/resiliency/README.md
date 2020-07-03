# Improving Resiliency

Istio's traffic management features can be leveraged to both improve resiliency and prove out resiliency policies. In this section we will use fault injection to intentionally disrupt communication between services and then apply policies to remediate any user impact.

Istio makes fault injection simple because it’s decoupled from your application code. This means you can break things without making any changes to application code. The istio-proxy is intercepting all network traffic and can adjust responses and response speed. You can easily inject a variety of faults, including HTTP error codes (L7 faults) and network failures or delays (L4 faults).

## Fault Injection - Aborts

Connect to the Hipstershop application, and select one item in the shop, like the Terrarium.

![Product view](/assets/Hipster_Shop-ads.png)

This page contains a detailed view of the product you selected, a button to buy the product then 4 other items that the shop propose, and a text advertisment.

We're going to inject an HTTP fault to simulate a transient failure in your `adservice` service, causing the ads to fail to load. To do this, you’ll add a fault stanza to the account `VirtualService` causing an HTTP 500 response code for 50% of calls to the accounts service. This stanza configures Envoy to return a 500 response code immediately when a client calls the account service instead of forwarding the request to the service itself.

```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: adservice
spec:
  hosts:
  - adservice.hipstershopv1v2.svc.cluster.local
  http:
  - route:
    - destination:
        host: adservice
    fault:
      abort:
        httpStatus: 500
        percentage: 
          value: 50
EOF
```

Now if you refresh your account list a couple of times, sometimes you won’t see any adds because we received a 500 from the backend. You can check the console in your browser’s developer tools in order to verify the response code.

![Product view noadds](/assets/Hipster_Shop-noads.png)


To ensure users aren't impacted by these request failures we have a couple tools in our toolbox: retries and outlier detection.

### Retries

Retries refer to the act of retrying a failed HTTP request to guard against transient failures. They can be used for any safe (GET) or idempotent (PUT/DELETE) requests.

### Outlier Detection

Outlier detection and ejection, a form of passive health checking, is the act of determining if some of the endpoints to which we're sending traffic are performing differently than the others and avoiding sending traffic to the outliers. We say that the endpoints we avoid have been "ejected from the active load balancing set." For any given service, Envoy always maintains a set of healthy endpoints for that service: the active load balancing set.

Typically, endpoints are ejected from the active load balancing set based on consecutive error responses. With HTTP services this would be consecutive 5xx failures. With TCP services, connect timeouts and connection errors/failures would lead to ejection. Over time, Envoy will attempt to add ejected endpoints back into the active load balancing set by sending traffic to them. If the endpoint responds successfully, it's reintroduced to the active load balancing set. Otherwise, it remains in the ejected set.

Outlier detection and retries are used in conjunction to improve resiliency. Outlier detection will increase success rate by ejecting endpoints that are deemed unhealthy, and retries mask any failures to users calling our application. 

> There is currently a [bug](https://github.com/istio/istio/issues/13705) in Istio where if you set fault injection AND retries then the retries do not take effect. If you set only one or the other then they work.

To be able to demonstrate this scenario, we are going to use a feature from the `adservice` microservice. 
For that, remove the `VirtualService` we just created. We also scale the `adservice-v2` microservice to 0:

```shell
kubectl -n hipstershopv1v2 delete vs adservice
kubectl -n hipstershopv1v2  scale --replicas=0 deployment adservice-v2
kubectl -n hipstershopv1v2  scale --replicas=0 deployment frontend-v2
```

Then edit the `adservice` deployment and add the `CONSECUTIVEERROR=2` environment variable so 2 requests out of 3 will return an error. 
We are also adding a small lattency of 2 seconds, to better demonstrate what's going on during a retry.

```shell
kubectl -n hipstershopv1v2 patch deployment adservice --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/env", value: [{"name":"CONSECUTIVEERROR","value":"2"},{"name":"EXTRA_LATENCY","value":"2s"},{"name":"LOGLEVEL","value":"debug"},{"name":"SRVURL","value":":9555"}]}]'
```

The default behaviour of Istio is to retry two times, so we are going to turn it down to 0 so we can experience a failure:

Let’s turn off retries.

```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: adservice
spec:
  hosts:
  - adservice.hipstershopv1v2.svc.cluster.local
  http:
  - route:
    - destination:
        host: adservice
    retries:
      attempts: 0
EOF
```

As before, 50% of the requests have the `advertisements` block at the bottom of the page.
If we look at the access logs of the frontend, we clearly see that one out of three requests to the `adservice` service is having a 500 error code:

```yaml
  "GET /ad HTTP/1.1" 200 - "-" "-" 0 123 2031 2030 "-" "Go-http-client/1.1" "09f236b7-dd59-4ae2-8a67-11a49417da31" "adservice.hipstershopv1v2:9555" "10.56.2.27:9555" outbound|9555||adservice.hipstershopv1v2.svc.cluster.local 10.56.0.48:40966 10.122.15.107:9555 10.56.0.48:51888 - -
  "GET /ad HTTP/1.1" 500 - "-" "-" 0 52 2003 2003 "-" "Go-http-client/1.1" "fddbbb61-5289-4d6f-8d18-a885245f1992" "adservice.hipstershopv1v2:9555" "10.56.2.27:9555" outbound|9555||adservice.hipstershopv1v2.svc.cluster.local 10.56.0.48:40966 10.122.15.107:9555 10.56.0.48:51888 - -
```

To be precise, let's look at the start of the access logs. Remember the formating: `<METHOD> <PATH> <PROTOCOL> <RESPONSE_CODE> <RESPONSE_FLAG> "x" "x" <BYTES_RECEIVED> <BYTES_SENT>` then the more interresting part here: `<BYTES_RECEIVED> <BYTES_SENT> <REQ_DURATION> <UPSTREAM_DURATION>`

We see both requests beeing arount 2000ms of response-time. This is on par with the configuration we made, the 2s delay.

Let's add 2 retries on errors. We also add `retryOn` configuration to add, beside default cases, the `5xx`, which will retry on any `500, 501, 502...` errors. We need this as a `500` error is not retriable by default:

```shell
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: adservice
spec:
  hosts:
  - adservice.hipstershopv1v2.svc.cluster.local
  http:
  - route:
    - destination:
        host: adservice
    retries:
      attempts: 2
      perTryTimeout: 5s
      retryOn: connect-failure,refused-stream,unavailable,cancelled,retriable-status-codes,5xx
EOF
```

After few reloads on the Hipstershop frontend, we always have the `Advertisement` block. Our problem is solved !

If we look again at the logs, truncated to the interresting part:

```yaml
 "GET /ad HTTP/1.1" 200 - "-" "-" 0 113 2030 2029 "
 "GET /ad HTTP/1.1" 200 - "-" "-" 0 137 6074 6074 "
```

Here, the first request took `2030ms`. We were lucky and the first call to the `adservice` was successful.

The second request took `6074ms` (6 seconds) to succeed... Let's look at the `adservice` logs:

```yaml
 "GET /ad HTTP/1.1" 500 - "-" "-" 0 52 2001 2000 
 "GET /ad HTTP/1.1" 500 - "-" "-" 0 52 2001 2001 
 "GET /ad HTTP/1.1" 200 - "-" "-" 0 123 2001 2000 
```
Here, from the `adservice` perspective, 3 requests came in and the first two failed. Each request tool `2000ms`, which correspond to the overal 6s as seen by the `frontend` service.


# re-install adservice-v2
# see how adservice is ejected

We can also add outlier detection by adding a `DestinationRule` (more on those in later sections). However, in this case it is the entire service rather than individual service instances seeing the errors so it won’t help us! An example of how we would configure outlier detection is below. In it, a service instance will be ejected if it returns a 5xx on 5 consecutive attempts. It will only be allowed back in after 5 minutes multiplied by the number of times it has been ejected.

```yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: account
spec:
  host: account
  trafficPolicy:
   outlierDetection:
      consecutiveErrors: 5
      baseEjectionTime: 5m
```

## Fault Injection - Delays

Next, we will inject an HTTP latency fault to simulate a transient latency into requests sent to the accounts service. This will instruct the sidecar Envoy of the calling service to hold onto the request for the configured amount of time before proxying the request on to its intended target. This is also done using a VirtualService in a similar manner to the abort fault.

```shell
kubectl apply -f modules/traffic/resiliency/config/account-delay.yaml
```

Now our `VirtualService` has a delay stanza that will cause 50% of requests to have 10 seconds of latency.

```shell
$ kubectl get virtualservice account -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: account
spec:
  hosts:
  - account
  http:
  - route:
    - destination:
        host: account
    fault:
      delay:
        percentage:
          value: 50
        fixedDelay: 10s
```

If we refresh our accounts page you will notice that now sometimes accounts appear after 10 seconds.

### Timeouts & Retries

Networks are unreliable. Consequently, transient request latency is a common occurance in a distributed system. So how can we protect our users from this? Well, we can combine retries and timeouts. Rather than wait the 10s for the request to complete, we will timeout after 0.5s and then retry, repeating until one of our requests succeed or we fail 3 consecutive times. Let’s apply that configuration to the microservice that calls the accounts service, the user service.

```shell
kubectl apply -f modules/traffic/resiliency/config/user-retry.yaml
```

```shell
$ kubectl get virtualservice user -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user
spec:
  hosts:
  - "*"
  gateways:
  - ingress
  http:
...
    timeout: 3s
    retries:
      attempts: 3
      perTryTimeout: 0.5s
      retryOn: gateway-error,connect-failure,unavailable
```

Now when you refresh the page you will notice that account load time will fluctuate but it has reduced from the 10s on half of requests previously. It’s safe to leave the user retry change we made, but we should clean up the latency injection.

> There is currently a [bug](https://github.com/istio/istio/issues/13705) in Istio where if you set fault injection AND retries then the retries do not take effect. If you set only one or the other then they work.

```shell
kubectl delete virtualservice account
```
