# Modern Bank Demo Application

## Prerequisites

<!-- TODO: @Liam update once this section is complete
 - kubectl
 - Istio installed
 - Kubernetes Cluster Context
 -->

## Modern Bank

Modern Bank are a bank that allows customers to make payments to other customers on their platform.

They have five microservices:

- User (manages customer information)
- Account (manages customer accounts)
- Transaction (orchestrates customer transactions)
- Transaction-Log (append only log of transactions)
- UI (serves user interface)

<!-- TODO: @Gus Insert architecture diagram here with various interactions between the microservices -->

<!-- TODO: @Liam Discuss API contract here, was unable to get Swagger-UI docker image to do what I wanted it to do! -->

## Installation

<!-- TODO: @Liam @Zack Work out the kube semantics -->

Label the `default` namespace to tell Istio to inject sidecar proxies to the application.

```bash
kubectl label namespace default istio-injection=enabled
```

Install the demo application.

```bash
kubectl apply -f config/
```

Verify the demo application has installed correctly.

<!-- TODO: update response with UI microservice -->

```bash
$ kubectl get pods
NAME                                       READY   STATUS    RESTARTS   AGE
account-ccc775744-nnbwt                    2/2     Running   0          4h
account-mongodb-6757fb69c5-nx4v6           2/2     Running   0          4h
transaction-b6b9c96f5-qjc7v                2/2     Running   0          4h
transaction-log-d5976d644-xrnbs            2/2     Running   0          4h
transaction-log-mongodb-556fbc7499-jkslx   2/2     Running   0          4h
user-7ccc796d46-l9x8n                      2/2     Running   0          4h
user-mongodb-7867988d7c-ptld5              2/2     Running   0          4h
```

We have installed all the microservices detailed above and databases for thoses that store state. We have also installed a default Istio Gateway to expose our application to the internet, we will go into more detail on this in the traffic routing section.

To make a request to our application we need to know the external IP of the Istio Ingress Gateway.

```bash
$ kubectl get service -n istio-system istio-ingressgateway
NAME                   TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)                         AGE
istio-ingressgateway   LoadBalancer   10.51.250.86   35.246.69.176   80:31380/TCP,443:31390/TCP...   3d
```

To set a variable to this IP address run the following command

```bash
export CLUSTER_IP=$(kubectl get service -n istio-system istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
```

Now we know the external IP, we are able to make a request to the application.

<!-- TODO: @Liam move this to view the UI when it exists -->

```bash
$ curl -v http://$CLUSTER_IP/v1/users/not-yet-created
...
* Connected to 35.197.239.230 (35.197.239.230) port 80 (#0)
> GET /v1/users/not-yet-created HTTP/1.1
> Host: 35.197.239.230
...
< HTTP/1.1 404 Not Found
```

We received a 404, this is expected as we tried to get information about a user that doesn't yet exist.
