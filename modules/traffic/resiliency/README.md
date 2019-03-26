# Improving Resiliency

Istio's traffic management features can be leveraged to both improve resiliency and prove out resiliency policies. In this section we will use fault injection to intentionally disrupt communication between services and then apply policies to remediate any user impact.

## Fault Injection - Aborts

We will inject an HTTP `fault` to simulate a transient failure in our `transaction-log` service. To do this we'll update `transaction-log`'s VirtualService to include a fault stanza to return an HTTP 500 response code for 100% of calls. This configures Envoy to return a 500 response code immediately when a client calls the `transaction-log` service instead of forwarding the request to the service itself.

<!-- TODO: @Liam add some words about what exactly this is doing! -->
<!-- TODO: @Liam come up with a sensible failure route match -->

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: transaction-log
spec:
  hosts:
  - transaction-log
  http:
  - match:
    - headers:
        end-user:
          exact: jason
    fault:
      abort:
        percent: 100
        httpStatus: 500
    route:
    - destination:
        host: transaction-log
        subset: v1
  - route:
    - destination:
        host: transaction-log
        subset: v1
```

<!-- TODO: @Liam Observe the effects of this in the UI -->

To ensure users aren't impacted by these request failures we have a couple tools in our toolbox: retries and outlier detection.

## Retries

Retries refer to the act of retrying a failed HTTP request to guard against transient failures. They can be used for any [safe](https://tools.ietf.org/html/rfc7231#section-4.2.1) (GET) or [idempotent](https://tools.ietf.org/html/rfc7231#section-4.2.2) (PUT/DELETE) requests.

<!-- TODO: @Liam Add Retry explanation of how this helps us -->
<!-- TODO: @Liam Add Retry Policy -->

## Outlier Detection

Outlier detection and ejection, a form of passive health checking, is the act of determining if some of the endpoints we're sending traffic to are performing differently than the others and avoiding sending traffic to them. We say that the endpoints we avoid have been "ejected from the active load balancing set". For any given service, Envoy always maintains a set of healthy endpoints for that service - the active load balancing set. Typically endpoints are ejected from the active load balancing set based on consecutive error responses. With HTTP services this is consecutive 5xx failures, in TCP services this is connect timeouts and connection errors/failures. Over time, Envoy will attempt to add ejected endpoints back into the active load balancing set by sending traffic to them. If the endpoint responds successfully it's reintroduced to the active load balancing set, otherwise it remains in the ejected set.

<!-- TODO: @Liam Add Outlier Detection explanation of how this helps us -->
<!-- TODO: @Liam Add Outlier Detection Policy -->

<!-- TODO: @Liam Summary of how these work together to help protect user from failure -->
<!-- TODO: @Liam Observe the effects of this in the UI -->

## Fault Injection - Delays

Next, we will inject an HTTP latency fault to simulate a transient latency into requests sent to one of our services. This will instruct the sidecar Envoy of the calling service to hold onto the request for the configured amount of time before proxying the request on to its intended target. This is also done using a VirtualService in a similar manner to the abort fault.

<!-- TODO: @Liam add some config for this (10s to service for 25% of requests) -->
<!-- TODO: @Liam come up with a sensible failure route match -->
<!-- TODO: @Liam add some words about the impact on the UI -->

## Timeouts & Retries

Networks are unreliable, consequently transient request latency is a common occurence in a distributed system. So how can we protect our users from this? Well, we can combine retries and timeouts. Rather than wait the 10s for the request to complete we will timeout after 1s and then retry, repeating until one of our requests succeed or we fail 4 consecutive times.

<!-- TODO: @Liam add some config for this (3 retries, 1s deadlines) -->
<!-- TODO: @Liam add some words about the impact on the UI -->
