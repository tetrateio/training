Observability
=====

To manage microservice architecture, which can involve dozens of microservices and complex dependencies, you need to be able to see what’s going on.

Istio comes with out-of-the-box monitoring features - like generating consistent application metrics for every service in your mesh - and can be used with an array of backend systems to report telemetry about the mesh. It helps us solve some of the trickiest problems we face: identifying why and where a request is slow, distinguishing normal from deviant system performance, comparing apples-to-apples metrics across services regardless of programming language, and attaining a meaningful view of system performance.

Exploring the UI
----

Using these consistent metrics, we can build powerful dashboards and visualizations. Let's start by taking a look at our system with RocketBot, the SkyWalking UI we installed at the beginning.

This service is exposed on our cluster, and we can get the address from Kubernetes:
```sh
export SKYWALKING_UI=$(kubectl -n skywalking get svc ui -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo $SKYWALKING_UI
```

Let's take that address and paste it into our browsers to see the RocketBot UI. We'll be promted to login, and we can use the default credentials **admin / admin**.

> Of course, we could just use `kubectl get svc ui -n skywalking` and copy and paste out the address manually.

![RocketBot UI Login screen, with username admin and password admin](/assets/rocketbot-login.png)

Once we're logged in, we're greeted with an empty screen. We can enter the admin view (bottom left) and select the elements to show in the UI.

![RocketBot UI has an edit button in the bottom left corner that we use to set up the metric display](/assets/rocketbot-editmode.png)

And configure the metrics we want to see:

![Use the RocketBot edit mode to show per-service metrics, then leave edit mode to view the stats](/assets/rocketbot-editmode-selectmetrics.png)

> Unfortunately, RocketBot doesn't support importing/exporting profiles. However, these views are stored client side in your browser as a Cookie. If you create a cookie named "dashboard" with the contents:
> 
> ```[{"name":"Service Dashboard","type":"service","query":{"service":{},"endpoint":{},"instance":{}},"children":[{"name":"service","children":[{"o":"Service","name":"Service Avg Response","comp":"ChartAvgResponse","title":"Service Avg Response","type":"serviceInfo","width":3},{"o":"Service","name":"Service Avg Throughput","comp":"ChartAvgThroughput","title":"Service Avg Throughput","type":"serviceInfo","width":3},{"o":"Service","name":"Service Avg SLA","comp":"ChartAvgSLA","title":"Service Avg SLA","type":"serviceInfo","width":3},{"o":"Service","name":"Service Percent Response","comp":"ChartResponse","title":"Service Percent Response","type":"serviceInfo","width":3},{"o":"Service","name":"Service Top Slow Endpoint","comp":"ChartSlow","title":"Service Top Slow Endpoint","type":"serviceInfo.getSlowEndpoint","width":3},{"o":"Service","name":"Running ServiceInstance","comp":"ChartTroughput","title":"Running ServiceInstance","type":"serviceInfo.getInstanceThroughput","width":3}]}]},{"name":"Database Dashboard","type":"database","query":{"service":{}},"children":[{"name":"Database","children":[]}]}]```
> 
> you'll have a view that matches mine.
> ![Using the Chrome developer console to set the cookie value for the RocketBot UI](/assets/rocketbot-chromeconsole-cookie.png)


When we're done, exit edit mode by clicking on it and view the stats:
> You may need to adjust the time settings in the bottom right to see data.

We can also view our service graph, via the `Topology` tab at the top.

![Graph of our deployment via RocketBot's Topology view](/assets/rocketbot-graph.png)

Finally, we can view traces of individual requests. (These traces are used to color the slow nodes red in the Topology view.)

![Trace screen, which shows a trace through the service graph as a tree](TODO: get graph)

How it works
---

Mixer is called by every sidecar in the mesh for policy (the sidecar asks Mixer if each request is allowed) and to report telemetry about each request. We'll cover the policy side in detail in the security section, but for now lets dig into telemetry. Mixer is Istio's intetgration point with external systems. A backend, for example SkyWalking, can implement an integration with Mixer (called an "adapter"). Using Mixer's configuration, we can instantiate the adapter (called a "handler"), describe the data about each request we want to provide the adapter (an "instance") and when Mixer should call the handler with instances (a "rule").


We configured all of this when we installed SkyWalking. To see the config, we can query Kubernetes about each of the types we list above. First, we can view the `adapter` config, which states that `skywalking-adapter` implements the `metric` template and is called per-request (`session_based: false`):
```shell
kubectl get adapter skywalking-adapter -n istio-system -o yaml
```
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
```

Our handler describes a specific instance of `skywalking-adapter` which we can send `metric` data to. In this case, we'll forward metric data to the SkyWalking collector we deployed as the `oap.skywalking` service.
```shell
kubectl get handlerskywalking-handler -n istio-system -o yaml
```
```yaml
apiVersion: "config.istio.io/v1alpha2"
kind: handler
metadata:
 name: skywalking-handler
 namespace: istio-system
spec:
 adapter: skywalking-adapter
 connection:
   address: "oap.skywalking.svc.cluster.local:11800"
```

Next we need to configure what data we'll send to SkyWalking; this is called an `instance`. Typically an adapter is built to expect certain instances and they'll be provided alongside the adapter's other configuration. We can the single metric that SkyWalking consumes, which is a metric instance with a bunch of dimensions: 
```shell
kubectl get instanceskywalking-metric -n istio-system -o yaml
```
```yaml
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
```

Finally, Mixer needs to know when to generate this metric data and send it to SkyWalking. This is defined as a `rule`. Every `rule` has a `match` condition that is evaluated; it the `match` is true, the `rule` is triggered. For example, we could use the `match` to receive only HTTP data, or only TCP data, etc. In our rule we omit the `match` because it defaults to `true` - we'll get data for every request. So this rule fires for every request, and constructs a `skywalking-metric` value it passes to the `skywalking-handler` (our `oap` service).

```shell
kubectl get rule skywalking-rule -n istio-system -o yaml
```
```yaml
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