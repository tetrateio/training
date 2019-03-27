# Improving Resiliency

Istio's traffic management features can be leveraged to both improve resiliency and prove out resiliency policies. In this section we will use fault injection to intentionally disrupt communication between services and then apply policies to remediate any user impact.

Istio makes fault injection simple because it’s decoupled from your application code. This means you can break things without making any changes to application code. The istio-proxy is intercepting all network traffic and can adjust responses and response speed. You can easily inject a variety of faults, including HTTP error codes (L7 faults) and network failures or delays (L4 faults).

## Fault Injection - Aborts

Now that you’re logged in to the banking application, you can see your accounts. You’re now going to inject an HTTP fault to simulate a transient failure in your accounts service, causing the accounts to fail to load. To do this, you’ll add a fault stanza to the account `VirtualService` causing an HTTP 500 response code for 50% of calls to the accounts service. This stanza configures Envoy to return a 500 response code immediately when a client calls the account service instead of forwarding the request to the service itself.

```shell
kubectl apply -f ./modules/traffic/resiliency/config/account-abort.yaml
```

Now if you refresh your account list a couple of times, sometimes you won’t see any accounts because we received a 500 from the backend. You can check the console in your browser’s developer tools in order to verify the response code. Let’s query the `VirtualService` to see the configuration we just added.

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
      abort:
        httpStatus: 500
        percent: 50
```

To ensure users aren't impacted by these request failures we have a couple tools in our toolbox: retries and outlier detection.

### Retries

Retries refer to the act of retrying a failed HTTP request to guard against transient failures. They can be used for any safe (GET) or idempotent (PUT/DELETE) requests.

### Outlier Detection

Outlier detection and ejection, a form of passive health checking, is the act of determining if some of the endpoints to which we're sending traffic are performing differently than the others and avoiding sending traffic to the outliers. We say that the endpoints we avoid have been "ejected from the active load balancing set." For any given service, Envoy always maintains a set of healthy endpoints for that service: the active load balancing set. Typically, endpoints are ejected from the active load balancing set based on consecutive error responses. With HTTP services this would be consecutive 5xx failures. With TCP services, connect timeouts and connection errors/failures would lead to ejection. Over time, Envoy will attempt to add ejected endpoints back into the active load balancing set by sending traffic to them. If the endpoint responds successfully, it's reintroduced to the active load balancing set. Otherwise, it remains in the ejected set. Typically, outlier detection and retries are used in conjunction to improve resiliency. Outlier detection will increase success rate by ejecting endpoints that are deemed unhealthy, and retries mask any failures to users calling our application. Let’s turn on retries.

```shell
kubectl apply -f ./modules/traffic/resiliency/config/account-abort-retry.yaml
```

If we inspect the account `VirtualService`, now we can see there is a retry stanza alongside the existing abort fault injection.

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
      abort:
        httpStatus: 500
        percent: 50
   retries:
      attempts: 5
      perTryTimeout: 1s
```

Now you should see your accounts almost every time you refresh. In fact, we reduced the failure rate our users experience from ~50% to ~3%.

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
kubectl apply -f ./modules/traffic/resiliency/config/account-delay.yaml
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
kubectl apply -f ./modules/traffic/resiliency/config/user-retry.yaml
```

```shell
$ kubectl get virtualservice user-gateway -o yaml
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: user-gateway
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

```shell
kubectl delete virtualservice account
```
