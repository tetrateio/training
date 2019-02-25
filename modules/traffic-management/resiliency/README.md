# Improving Resiliency

Istio's traffic management features can be leveraged to both improve resiliency and prove out resiliency policies. In this section we will use fault injection to intentionally disrupt communication between services and then apply policies to remediate any user impact.

## Request Failures

First, we will inject an HTTP abort fault to simulate a transient failure of one of our services. This will instruct the sidecar proxy of the calling service to return 500 response codes instead of proxying on the request. To do this we need to add a VirtualService. A VirtualService defines a set of traffic routing rules to apply when a host is addressed and is responsible for configuring which service a request should be routed to.

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

Retries refer to the act of retrying a failed HTTP request to guard against transient failures. They can be used for any [safe](https://tools.ietf.org/html/rfc7231#section-4.2.1) (GET) and/or [idempotent](https://tools.ietf.org/html/rfc7231#section-4.2.2) (PUT/DELETE) requests.

<!-- TODO: @Liam Add Retry explanation of how this helps us -->
<!-- TODO: @Liam Add Retry Policy -->

## Outlier Detection

Outlier detection, or passive health checking, is the act of ejecting endpoints from the load balancing pool. Endpoints are ejected based on consecutive error responses. With HTTP services this is consecutive 5xx failures, in TCP services this is connect timeouts or connection errors/failures.

<!-- TODO: @Liam Add Outlier Detection explanation of how this helps us -->
<!-- TODO: @Liam Add Outlier Detection Policy -->

<!-- TODO: @Liam Summary of how these work together to help protect user from failure -->
<!-- TODO: @Liam Observe the effects of this in the UI -->


## Request Latency

Next, we will inject an HTTP latency fault to simulate a transient latency into requests sent to one of our services. This will instruct the sidecar proxy of the calling service to hold onto the request for the configured amount of time before proxying the request on to its intended target. This is also done using a VirtualService in a similar manner to the abort fault.

<!-- TODO: @Liam add some config for this (10s to service for 25% of requests) -->
<!-- TODO: @Liam come up with a sensible failure route match -->
<!-- TODO: @Liam add some words about the impact on the UI -->

## Timeouts & Retries

In order reduce the impact to our users whilst 25% of requests are suffering from severe network latency we can use a combination of retries and timeouts. Rather than wait the 10s for the request to complete we will timeout after 1s and then retry, repeating until one of our requests succeed or we fail 4 consecutive times.

<!-- TODO: @Liam add some config for this (3 retries, 1s deadlines) -->
<!-- TODO: @Liam add some words about the impact on the UI -->
