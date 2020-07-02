# Observability

To manage a microservice architecture, which can involve dozens of microservices and complex dependencies, you need to be able to see what’s going on.

Istio comes with out-of-the-box monitoring features - like generating consistent application metrics for every service in your mesh - and can be used with an array of backend systems to report telemetry about the mesh. It helps us solve some of the trickiest problems we face: identifying why and where a request is slow, distinguishing normal from deviant system performance, comparing apples-to-apples metrics across services regardless of programming language, and attaining a meaningful view of system performance.

Let's generate some traffic to observe. Get the application IP.
Hipstershop comes with a `LoadGenerator`, which is an HTTP load-tester, pre-configured to simulate a real user behaviour. For the moment it is in a failed state as we need to configure the external IP of our service.
If you have followed the [Setup Ingress Gateway](/modules/ingressGateway/) module, you should have a `$INGRESSIP` variable with the public IP of the Ingress Gateway. You will need to copy this IP and replace the `add-your-ip` text from the `LoadGenerator` deployment:

```shell
kubectl -n hipstershopv1v2 get deployment loadgenerator -o yaml | sed "s/add-your-ip/$INGRESSIP/g" | kubectl apply -f -
```

You should soon see the pod as `Running`:

```shell
kubectl -n hipstershopv1v2 get pods -l app=loadgenerator

NAME                             READY   STATUS    RESTARTS   AGE
loadgenerator-7d45d6cbfc-pfx6h   1/1     Running   0          2m47s
```

## Logs

Now that we have a "client" running, we can check our application's logs.
Using the `kubectl logs` command, or better, the [Stern](https://github.com/wercker/stern) command, let see what is going on on the `frontend` pods:

```shell
kubectl -n hipstershopv1v2 logs deployment/frontend istio-proxy

[2020-07-02T14:22:48.554Z] "GET /product/2ZYFJ3GM2N HTTP/1.1" 200 - "-" "-" 0 7957 6037 6036 "10.56.0.1" "python-requests/2.21.0" "fd03dbd0-5643-4e2d-8058-e414840d0bc4" "hipstershop.35.245.245.42.sslip.io" "127.0.0.1:8080" inbound|8080|http-frontend|frontend.hipstershopv1v2.svc.cluster.local 127.0.0.1:48300 10.56.2.23:8080 10.56.0.1:0 outbound_.8080_._.frontend.hipstershopv1v2.svc.cluster.local default
[2020-07-02T14:22:54.610Z] "POST /hipstershop.CurrencyService/GetSupportedCurrencies HTTP/2" 200 - "-" "-" 5 170 2 1 "-" "grpc-go/1.26.0" "13453d3f-6a36-4e3a-a434-e52e1de8cac0" "currencyservice.hipstershopv1v2:7000" "10.56.2.22:7000" outbound|7000||currencyservice.hipstershopv1v2.svc.cluster.local 10.56.2.23:34828 10.122.7.139:7000 10.56.2.23:47962 - default
[2020-07-02T14:22:54.613Z] "POST /hipstershop.CartService/GetCart HTTP/2" 200 - "-" "-" 43 21 2 2 "-" "grpc-go/1.26.0" "06d010d4-ab97-4b72-83a8-2d6b30336bc3" "cartservice.hipstershopv1v2:7070" "10.56.2.20:7070" outbound|7070||cartservice.hipstershopv1v2.svc.cluster.local 10.56.2.23:36912 10.122.14.171:7070 10.56.2.23:57254 - default
[2020-07-02T14:22:54.616Z] "POST /hipstershop.RecommendationService/ListRecommendations HTTP/2" 200 - "-" "-" 55 65 4 4 "-" "grpc-go/1.26.0" "c8b23546-bf52-44ad-98c7-cb8682d8e07a" "recommendationservice.hipstershopv1v2:8080" "10.56.2.25:8080" outbound|8080||recommendationservice.hipstershopv1v2.svc.cluster.local 10.56.2.23:51064 10.122.12.93:8080 10.56.2.23:40336 - default
[2020-07-02T14:22:55.561Z] "POST /hipstershop.CurrencyService/GetSupportedCurrencies HTTP/2" 200 - "-" "-" 5 170 1 1 "-" "grpc-go/1.26.0" "4d3fed9c-ff73-472e-bd54-dabd96cd429f" "currencyservice.hipstershopv1v2:7000" "10.56.2.22:7000" outbound|7000||currencyservice.hipstershopv1v2.svc.cluster.local 10.56.2.23:34828 10.122.7.139:7000 10.56.2.23:47962 - default
[2020-07-02T14:22:55.563Z] "POST /hipstershop.ProductCatalogService/ListProducts HTTP/2" 200 - "-" "-" 5 1439 2 1 "-" "grpc-go/1.26.0" "7d780e90-7ddf-9df5-ba77-df5034b97871" "productcatalogservice.hipstershopv1v2:3550" "10.56.1.33:3550" outbound|3550||productcatalogservice.hipstershopv1v2.svc.cluster.local 10.56.2.23:34542 10.122.14.228:3550 10.56.2.23:50940 - default
[2020-07-02T14:22:55.567Z] "POST /hipstershop.CartService/GetCart HTTP/2" 200 - "-" "-" 43 21 2 2 "-" "grpc-go/1.26.0" "a09771bf-9622-497a-87fe-e2819b4a4c9d" "cartservice.hipstershopv1v2:7070" "10.56.2.20:7070" outbound|7070||cartservice.hipstershopv1v2.svc.cluster.local 10.56.2.23:36912 10.122.14.171:7070 10.56.2.23:57254 - default
[2020-07-02T14:22:55.570Z] "POST /hipstershop.CurrencyService/Convert HTTP/2" 200 - "-" "-" 25 17 1 1 "-" "grpc-go/1.26.0" "1dd990eb-741f-4a21-8840-e870ea7dc2b1" "currencyservice.hipstershopv1v2:7000" "10.56.2.22:7000" outbound|7000||currencyservice.hipstershopv1v2.svc.cluster.local 10.56.2.23:34828 10.122.7.139:7000 10.56.2.23:47962 - default
```

For each request comming in, Istio-proxy will log an entry using the `accesslog` format. You can find more informations on [Envoy's documentation page](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#configuration).
We can also check the logs of the `frontend-v2` deployment:

```shell
kubectl -n hipstershopv1v2 logs deployment/frontend-v2 istio-proxy

[2020-07-02T14:25:23.557Z] "GET /product/2ZYFJ3GM2N HTTP/1.1" 200 - "-" "-" 0 7958 6039 6038 "10.56.0.1" "python-requests/2.21.0" "34dbbf8e-e1a6-43db-bd5f-b80b1d9ba19f" "hipstershop.35.245.245.42.sslip.io" "127.0.0.1:8080" inbound|8080|http-frontend|frontend.hipstershopv1v2.svc.cluster.local 127.0.0.1:46144 10.56.1.29:8080 10.56.0.1:0 outbound_.8080_._.frontend.hipstershopv1v2.svc.cluster.local default
[2020-07-02T14:25:29.815Z] "POST /hipstershop.ProductCatalogService/GetProduct HTTP/2" 200 - "-" "-" 17 185 1 1 "-" "grpc-go/1.26.0" "1d45e05a-1736-4016-b827-467e211b73b5" "productcatalogservice.hipstershopv1v2:3550" "10.56.1.33:3550" outbound|3550||productcatalogservice.hipstershopv1v2.svc.cluster.local 10.56.1.29:47842 10.122.14.228:3550 10.56.1.29:37812 - default
[2020-07-02T14:25:29.816Z] "POST /hipstershop.CurrencyService/GetSupportedCurrencies HTTP/2" 200 - "-" "-" 5 170 2 2 "-" "grpc-go/1.26.0" "31949ad6-4214-498e-8778-938759178c1a" "currencyservice.hipstershopv1v2:7000" "10.56.2.22:7000" outbound|7000||currencyservice.hipstershopv1v2.svc.cluster.local 10.56.1.29:45092 10.122.7.139:7000 10.56.1.29:33418 - default
[2020-07-02T14:25:29.820Z] "POST /hipstershop.CartService/GetCart HTTP/2" 200 - "-" "-" 43 69 6 5 "-" "grpc-go/1.26.0" "3031cfdf-c80d-40ca-a1dc-b621551ba1e2" "cartservice.hipstershopv1v2:7070" "10.56.2.20:7070" outbound|7070||cartservice.hipstershopv1v2.svc.cluster.local 10.56.1.29:59836 10.122.14.171:7070 10.56.1.29:38750 - default
...
```

They are quite similar, which is the expected behaviour as both `frontend` and `frontend-v2` are target destination for the `frontend` service that we declared in the `VirtualService`.

Let's move to the graphical tools

## Exploring the UIs

### Grafana

Using these consistent metrics, we can build powerful dashboards and visualizations. Let's start by taking a look at our system with Grafana, which we installed alongside Istio.

This service is not exposed on our cluster. You can use the `istioctl dashboard grafana` command to do the port-forward and open your browser all at once or use a port-forward:

```shell
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod \
    -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000 &
```

> We start the port-forwarding command in the background as we'll want to port-forward a few different services in the workshop.

We can go check out Grafana, and the default dashboards that Istio ships with, at http://localhost:3000/. Once there navigate to `Istio Service Dashboard` and select the Transaction service from the dropdown.

> If you're in Google Cloud Shell you access it via the web preview feature in the top right hand corner. You may need to change the port to 3000.

### Kiali

While metrics are awesome, for understanding a new system nothing beats seeing a graph of the services in the system communicating. We also installed [Kiali](https://www.kiali.io/) alongside Istio; it comes with some nice visualizations, including a graph. 
Like Grafana, it's not exposed outside of the cluster, so we'll use the `istioctl dashboard kiali` or the port-forward command:

```shell
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod \
    -l app=kiali -o jsonpath='{.items[0].metadata.name}') 8081:20001 &
```

We can see the UI at http://localhost:8081/kiali with the username / password **admin / admin**. Once there navigate to `Graph` and select the `hipstershopv1v2` namespace from the dropdown.

> If you're in Google Cloud Shell you access it via the web preview feature in the top right hand corner. You may need to change the port to 8081.

[![Kiali](/assets/Kiali_Console.png)](/assets/Kiali_Console.png)

### Tracing

Finally, we _also_ installed a tracing tool called Jaeger, which we can view in the same way: `istioctl dashboard jaeger` or:

```sh
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod \
    -l app=jaeger -o jsonpath='{.items[0].metadata.name}') 8082:16686 &
```

> If you're in Google Cloud Shell you access it via the web preview feature in the top right hand corner. You may need to change the port to 8082.

Then open http://localhost:8082/ in your browser. Set the service to `hipstershop-ingressgateway` and run a search. We should see our store transactions in the traces.

## How It Works

With Istio 1.5+ (also in 1.4.x but opt-in) and Istiod, all the metrics of the Istio-proxy are generated/served by the Envoy process using an embeded WASM extension. This is called Telemetry v2. Prometheus scrapes the metrics from each Envoy in the cluster. 
You can extend the metrics exported by Istio by using an `EnvoyFiler` resource. This process is described in [Istio's documentation](https://istio.io/latest/docs/tasks/observability/metrics/customize-metrics/).

Istio Sidecar also generate Opentracing traces. Many protocols are supported, including Jaeger and Zipkin. When Envoy process inside the sidecar receive a request, a trace is sent to the tracing agent with the required informations. 
If the trace already have a `traceID` header, it is used by the next Envoy Proxy to add a new span. You can configure Istio's sampling rate during installation and it defaults to 1% of the requests. 
Demo profile is set to 100%, which can have a huge performance impact on a loaded system.
