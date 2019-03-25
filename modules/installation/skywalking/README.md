# Installing Apache Skywalking

Skywalking is an open source Application Performance Management (APM) system owned by the Apache Foundation. It supports monitoring and tracing applications as an all-in-one package, but focuses on higher level primitives (like Services) than traditional systems used for visibility (like Prometheus and Zipkin). We'll use SkyWalking's open source UI, RocketBot, to help understand our deployment.

First we'll install SkyWalking in three deployments: its collector, storage (Elastic Search), and the RocketBot UI.

```shell
$ kubectl apply -f config/skywalking.yaml --as=admin --as-group=system:masters
```
> We need to execute this command as an admin to allow the creation of SkyWalking's RBAC role, which allows SkyWalking to monitor pods in the cluster.

This creates a new namespace, `skywalking`:

```shell
$ kubectl get pods -n skywalking
NAME                              READY   STATUS    RESTARTS   AGE
elasticsearch-0                   1/1     Running   0          1h
elasticsearch-1                   1/1     Running   0          1h
elasticsearch-2                   1/1     Running   0          1h
oap-deployment-7d86485855-kb4kx   1/1     Running   0          1h
ui-deployment-78666987b6-s5bp2    1/1     Running   0          1h
```

> Note: Elastic Search and SkyWalking itself require a decent bit of RAM to schedule; sometimes Kubernetes will fail to schedule the pods even if there's enough RAM globally in the cluster. If this happens in your cluster, delete a pod or two to shake the scheduler out of its rut and get it to reschedule.

The UI itself is exposed as a `Service`; we use this rather than accessing the UI via Istio's Gateway to keep things simple for us.

```shell
$ kubectl get services -n skywalking
NAME            TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)               AGE
elasticsearch   ClusterIP      None            <none>         9200/TCP,9300/TCP     1h
oap             ClusterIP      10.47.251.115   <none>         12800/TCP,11800/TCP   1h
ui              LoadBalancer   10.47.243.220   <YOUR IP>      80:30845/TCP          1h
```

Keep note of this IP address, we'll want to open it in a browser later.

Finally, we also installed a bit a Istio configuration to point Mixer at SkyWalking:

```yaml
apiVersion: "config.istio.io/v1alpha2"
kind: adapter
metadata:
  name: skywalking-adapter
  namespace: istio-system
spec:
  description:
  session_based: false
  templates:
  - metric
---
apiVersion: "config.istio.io/v1alpha2"
kind: handler
metadata:
 name: skywalking-handler
 namespace: istio-system
spec:
 adapter: skywalking-adapter
 connection:
   address: "oap.skywalking.svc.cluster.local:11800"
---
apiVersion: "config.istio.io/v1alpha2"
kind: instance
metadata:
 name: skywalking-metric
 namespace: istio-system
spec:
 template: metric
 params:
   value: request.size | 0
   dimensions:
     sourceService: source.workload.name | ""
     sourceNamespace: source.workload.namespace | ""
     sourceUID: source.uid | ""
     destinationService: destination.workload.name | ""
     destinationNamespace: destination.workload.namespace | ""
     destinationUID: destination.uid | ""
     requestMethod: request.method | ""
     requestPath: request.path | ""
     requestScheme: request.scheme | ""
     requestTime: request.time
     responseTime: response.time
     responseCode: response.code | 200
     reporter: conditional((context.reporter.kind | "inbound") == "outbound", "source", "destination")
     apiProtocol: api.protocol | ""
---
apiVersion: "config.istio.io/v1alpha2"
kind: rule
metadata:
 name: skywalking-rule
 namespace: istio-system
spec:
 actions:
 - handler: skywalking-handler
   instances:
   - skywalking-metric
```