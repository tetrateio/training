# Observability

To manage a microservice architecture, which can involve dozens of microservices and complex dependencies, you need to be able to see what’s going on.

Istio comes with out-of-the-box monitoring features - like generating consistent application metrics for every service in your mesh - and can be used with an array of backend systems to report telemetry about the mesh. It helps us solve some of the trickiest problems we face: identifying why and where a request is slow, distinguishing normal from deviant system performance, comparing apples-to-apples metrics across services regardless of programming language, and attaining a meaningful view of system performance.

Let's generate some traffic to observe. Get the application IP.

```shell
export INGRESS_IP=$(kubectl -n istio-system get svc istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $INGRESS_IP
```

You can navigate to that IP and sign up to receive some free ~~virtual~~ fake money. Then send some of that money between accounts. It doesn't really matter what you do here just generate some traffic!

Alternatively, if we have Go installed we can use a tool to automatically generate traffic. **This will run perpetually so remember to open a new shell**.

```shell
export GO111MODULE=on
cd training/samples/modernbank/tools/trafficGen
go run cmd/main.go --host $(kubectl -n istio-system get svc istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
```
<!-- DON'T change this to just go run training/samples/modernbank/tools/trafficGen/cmd/main.go as this doesn't work! Go gets confused about which modules to use! -->

## Exploring the UIs

### Grafana

Using these consistent metrics, we can build powerful dashboards and visualizations. Let's start by taking a look at our system with Grafana, which we installed alongside Istio.

This service is not exposed on our cluster, so we'll need to port-forward:

```shell
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod \
    -l app=grafana -o jsonpath='{.items[0].metadata.name}') 8083:3000 &
```

> We start the port-forwarding command in the background as we'll want to port-forward a few different services in the workshop.

We can go check out Grafana, and the default dashboards that Istio ships with, at http://localhost:3000/. Once there navigate to `Istio Service Dashboard` and select the Transaction service from the dropdown.

> If you're in Google Cloud Shell you access it via the web preview feature in the top right hand corner. You may need to change the port to 3000.

### Kiali

While metrics are awesome, for understanding a new system nothing beats seeing a graph of the services in the system communicating. We also installed [Kiali](https://www.kiali.io/) alongside Istio; it comes with some nice visualizations, including a graph. Like Grafana, it's not exposed outside of the cluster, so we'll need to port-forward it:

```shell
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod \
    -l app=kiali -o jsonpath='{.items[0].metadata.name}') 8081:20001 &
```

We can see the UI at http://localhost:8081/kiali with the username / password **admin / admin**. Once there navigate to `Graph` and select the `default` namespace from the dropdown.

> If you're in Google Cloud Shell you access it via the web preview feature in the top right hand corner. You may need to change the port to 8081.

### Tracing

Finally, we _also_ installed a tracing tool called Jaeger, which we can view in the same way:

```sh
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod \
    -l app=jaeger -o jsonpath='{.items[0].metadata.name}') 8082:16686 &
```

> If you're in Google Cloud Shell you access it via the web preview feature in the top right hand corner. You may need to change the port to 8082.

Which we can see at http://localhost:8082/. Set the service to `istio-ingressgateway` and run a search. We should see our account transactions in the traces.

## How It Works

<!-- TODO: This section needs updating -->

Mixer is called by every sidecar in the mesh for policy (the sidecar asks Mixer if each request is allowed). Prometheus scrapes telemetry from Envoy about each request. We'll cover the policy side in detail in the security section, but for now lets dig into telemetry. Mixer is Istio's intetgration point with external systems. A backend, for example Prometheus, can implement an integration with Mixer (called an "adapter"). Using Mixer's configuration, we can instantiate the adapter (called a "handler"), describe the data about each request we want to provide the adapter (an "instance") and when Mixer should call the handler with instances (a "rule").

The Prometheus adapter is compiled in to Mixer. The `prometheus` handler describes a specific instance of Mixer Prometheus adapter, which we can send `metric` data to.

```shell
kubectl -n istio-system get handler prometheus -o yaml
```

```yaml
apiVersion: config.istio.io/v1alpha2
kind: handler
metadata:
  name: prometheus
  namespace: istio-system
spec:
  compiledAdapter: prometheus
  params:
    metrics:
    - instance_name: requestcount.metric.istio-system
      kind: COUNTER
      label_names:
      - reporter
      - source_app
      - source_principal
      - source_workload
      - source_workload_namespace
      - source_version
      - destination_app
      - destination_principal
...
```

> A Handler is where we provide configuration specific to the adapter. For Prometheus, we need to configure the shape of the metrics we'll be emitting. The `prometheus` handler does exactly this.

Next we need to configure what data we'll send to the Prometheus adapter. This is called an `instance` in general, but there are a few special `instances` that have a proper name; `metric` is one of those. Typically an adapter is built to expect certain instances and the configuration for those instances will be provided alongside the adapter's other configuration. We can see the set of instances that the Prometheus adapter consumes:

```shell
$ kubectl -n istio-system get metrics
NAME                   AGE
requestcount           17m
requestduration        17m
requestsize            17m
responsesize           17m
tcpbytereceived        17m
tcpbytesent            17m
tcpconnectionsclosed   17m
tcpconnectionsopened   17m
```

And we can inspect one to see what it looks like:
```shell
kubectl -n istio-system get metrics requestcount -o yaml
```
```yaml
apiVersion: config.istio.io/v1alpha2
kind: metric
metadata:
  name: requestcount
  namespace: istio-system
spec:
  dimensions:
    connection_security_policy: conditional((context.reporter.kind | "inbound") ==
      "outbound", "unknown", conditional(connection.mtls | false, "mutual_tls", "none"))
    destination_app: destination.labels["app"] | "unknown"
    destination_principal: destination.principal | "unknown"
    destination_service: destination.service.host | "unknown"
    destination_service_name: destination.service.name | "unknown"
    destination_service_namespace: destination.service.namespace | "unknown"
    destination_version: destination.labels["version"] | "unknown"
    destination_workload: destination.workload.name | "unknown"
    destination_workload_namespace: destination.workload.namespace | "unknown"
    permissive_response_code: rbac.permissive.response_code | "none"
    permissive_response_policyid: rbac.permissive.effective_policy_id | "none"
    reporter: conditional((context.reporter.kind | "inbound") == "outbound", "source",
      "destination")
    request_protocol: api.protocol | context.protocol | "unknown"
    response_code: response.code | 200
    response_flags: context.proxy_error_code | "-"
    source_app: source.labels["app"] | "unknown"
    source_principal: source.principal | "unknown"
    source_version: source.labels["version"] | "unknown"
    source_workload: source.workload.name | "unknown"
    source_workload_namespace: source.workload.namespace | "unknown"
  monitored_resource_type: '"UNSPECIFIED"'
  value: "1"
```

Finally, Mixer needs to know when to generate this metric data and send it to Prometheus. This is defined as a `rule`. Every `rule` has a `match` condition that is evaluated; if the `match` is true, the `rule` is triggered. For example, we could use the `match` to receive only HTTP data, or only TCP data, etc. Prometheus does exactly this, and defines a rule for each set of protocols it has metric descriptions for:

```shell
$ kubectl -n istio-system get rules
NAME                      AGE
kubeattrgenrulerule       19m
promhttp                  19m
promtcp                   19m
promtcpconnectionclosed   19m
promtcpconnectionopen     19m
stdio                     19m
stdiotcp                  19m
tcpkubeattrgenrulerule    19m
```

And again we can inspect one to see what it looks like:
```shell
kubectl -n istio-system get rules promhttp -o yaml
```
```yaml
apiVersion: config.istio.io/v1alpha2
kind: rule
metadata:
  name: promhttp
  namespace: istio-system
spec:
  match: (context.protocol == "http" || context.protocol == "grpc") &&
            (match((request.useragent| "-"), "kube-probe*") == false)
  actions:
  - handler: prometheus
    instances:
    - requestcount.metric
    - requestduration.metric
    - requestsize.metric
    - responsesize.metric
```
