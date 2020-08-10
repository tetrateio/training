# Ingress Gateway

By default, the components deployed on the service mesh are not exposed outside of the cluster. In this section, you’ll configure Istio to expose our example UI service outside of the service mesh via the Istio Ingress.

You won’t configure Istio Ingress with the Kubernetes Ingress definitions you would be using in the absence of Istio. The Istio Ingress uses a different configuration model, the `Istio Gateway`, to allow features such as monitoring and route rules to be applied to traffic entering the cluster.

Note that Istio `Gateway` resource only configures the L4-L6 functions (e.g., ports to expose, TLS configuration). Users can then use standard Istio rules to control HTTP requests as well as TCP traffic entering a Gateway by binding it to a `VirtualService`.

We say a `VirtualService` binds to a `Gateway` if the `VirtualService` lists the Gateway’s name in its gateways field and at least one host claimed by the `VirtualService` is exposed by the `Gateway`.


In this module we are going to configure the `Istio Ingress Gateway` to allow HTTP and HTTPS traffic to our Hipstershop application.

[![Ingress Gateway](/assets/hipstershop-istio-ingress-IngressGW-1.svg)](/assets/hipstershop-istio-ingress-IngressGW-1.svg)

## Preliminary checks
When using the `demo` install Profile, a default (global) Gateway is installed in the `istio-system` Namespace. Let's first use it to get access to our application.
Check that the gateway is here and working:

```shell
kubectl -n istio-system get pods -l app=istio-ingressgateway

NAME                                    READY   STATUS    RESTARTS   AGE
istio-ingressgateway-69bb66cb4b-xvctp   1/1     Running   0          8h

kubectl -n istio-system get svc -l app=istio-ingressgateway

NAME                   TYPE           CLUSTER-IP     EXTERNAL-IP     PORT(S)                                                                      AGE
istio-ingressgateway   LoadBalancer   10.122.11.61   35.221.16.159   15020:32232/TCP,80:30545/TCP,443:32139/TCP,31400:32604/TCP,15443:30777/TCP   8h
```

We will also use the `sslip.io` service, which is a dynamic DNS that always reply with the IP provided in the name. For example, a DNS lookup for `hipstershop.35.221.16.159.sslip.io` will answer `35.221.16.159`.

First we need to grab the IP of the Ingress Gateway.
Run this on EKS cluster (as a DNS name is returned instead of an IP):

```shell
INGRESSHOST=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath="{.status.loadBalancer.ingress[0].hostname}")
INGRESSIP=$(dig +noall +answer +short ${INGRESSHOST} | head -1)
export INGRESSIP
echo $INGRESSIP
```
Run this on any other cluster (GKE):

```shell
INGRESSIP=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
export INGRESSIP
echo $INGRESSIP
```

## Using the default Ingress Gateway

Try to connect to the Istio Ingress IP.

```shell
$ curl hipstershop.${INGRESSIP}.sslip.io
curl: (7) Failed to connect to ... port 80: Connection refused
```

The connection is refused because nothing is binding to this port at this address.

Let’s tell Istio to bind a Gateway to the Istio Ingress.

We can apply the Gateway manifest:

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: hipstershop
spec:
  selector:
    istio: ingressgateway 
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "hipstershop.${INGRESSIP}.sslip.io"
EOF
```

> Note: This action will make the application accessible on the public internet. In some situations, you might consider alternative configurations for privacy:
>
> - Accessing services within the cluster (which does not require the gateway).
> - Accessing internal services across the cluster: consider using GKE private cluster, VPC, IP aliasing, and internal load balancer.
> - Accessing services on the public internet: require credentials, consider using Identity Access Proxy.
>

As you can see in the spec, this gateway is listening for request that match `hipstershop.${INGRESSIP}.sslip.io` hosts on port 80.

Curl the Istio Ingress IP again. Blast! You get a 404 error...

```shell
$ curl -v hipstershop.${INGRESSIP}.sslip.io
...
< HTTP/1.1 404 Not Found
< location: ...
< date: Thu, 21 Mar 2019 14:40:33 GMT
< server: istio-envoy
< content-length: 0
```

It's returning 404 error because once the gateway receives the request we haven’t told it where it needs to send it! To do this we need to attach a `VirtualService` (or two) to the gateway.

A `VirtualService` defines the rules that control how requests for a service are routed within an Istio service mesh. For example, a `VirtualService` can route requests to different versions of a service or to a completely different service than was requested. Requests can be routed based on the request source and destination, HTTP paths and header fields, and weights associated with individual service versions.

Note that within a `VirtualService`, the match conditions are checked at runtime in the order that they appear. This means that the most specific match clauses should appear first, and less specific clauses later. For safety, a “default” route, with no match conditions, should be provided; a request that does not match any condition of a `VirtualService` will result in a 404 for the sender (or some connection refused error for non-HTTP protocols).

Let's create a `VirtualService` that will link the Ingress Gateway to the Hipstershop `frontend` and `apiservice` services

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hipstershop
spec:
  hosts:
  - "hipstershop.${INGRESSIP}.sslip.io"
  gateways:
  - hipstershop
  http:
  - match:
    - uri:
        prefix: /api
    route:
    - destination:
        host: apiservice
        port:
          number: 8080
  - route:
    - destination:
        host: frontend
        port:
          number: 8080
EOF
echo "you can now connect to http://hipstershop.${INGRESSIP}.sslip.io/"
```

This `VirtualService` is bound to the ingress gateway we just created, named `hipstershop`. It handles API requests for the `/api` and `/` paths on our Hipstershop application, routing all requests to either the `apiservice` or `frontend` microservice depending on their path.

You can open your browser at [http://hipstershop.${INGRESSIP}.sslip.io/](http://hipstershop.${INGRESSIP}.sslip.io/) and enjoy buying stuff for your  shelves.

![Hipstershop App Home Screen](/assets/Hipster_Shop.png)


## Configuring HTTPS

Now that we have a working HTTP gateway, let's add SSL.
To do that we are going to use Cert-Manager, an application that can generate SSL certificates into Kubernetes secrets. Istio-proxies will then requests the SSL certificates to Istiod using the SDS protocol. Istiod will read the secrets and send the certificate to the proxies.

Let's install Cert-Manager:

```shell
kubectl create namespace cert-manager
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.1/cert-manager.yaml
```

We will use Self-Signed certificates for this training, but Cert-Manager is able to handle a lot of other Issuers, like ACME/Let's Encrypt:

```yaml
kubectl apply -n cert-manager -f - <<EOF
apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
  namespace: cert-manager
spec:
  selfSigned: {}
EOF
```

Now, request an SSL certificate for the Hipstershop deployment:

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: hipstershop
spec:
  secretName: hipstershop-cert
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  organization:
  - hipstershop
  commonName: hipstershop.${INGRESSIP}.sslip.io
  isCA: false
  keySize: 2048
  keyAlgorithm: rsa
  keyEncoding: pkcs1
  usages:
    - server auth
    - client auth
  dnsNames:
  - hipstershop.${INGRESSIP}.sslip.io
  ipAddresses:
  - ${INGRESSIP}
  issuerRef:
    name: selfsigned-issuer
    kind: ClusterIssuer
EOF
```

If everything is working, you should have a new secret:

```shell
kubectl get secret hipstershop-cert -o yaml

apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  annotations:
    cert-manager.io/alt-names: hipstershop.35.221.16.159.sslip.io
    cert-manager.io/certificate-name: hipstershop
    cert-manager.io/common-name: hipstershop.35.221.16.159.sslip.io
    cert-manager.io/ip-sans: 35.221.16.159
    cert-manager.io/issuer-kind: ClusterIssuer
    cert-manager.io/issuer-name: selfsigned-issuer
    cert-manager.io/uri-sans: ""
  name: hipstershop-cert
  namespace: hipstershopv1v2
data:
  ca.crt: LS0tLS1CRUdJTi...LS0tCg==
  tls.crt: LS0tLS1CRUdJT...LS0tCg==
  tls.key: LS0tLS1CRFFDGH...S0tLQo=
```

We can now update the Gateway resource to add support for HTTPS protocol on port 443:

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: hipstershop
spec:
  selector:
    istio: ingressgateway 
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "hipstershop.${INGRESSIP}.sslip.io"
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      mode: SIMPLE
      credentialName: "hipstershop-cert" 
    hosts:
    - "hipstershop.${INGRESSIP}.sslip.io"
EOF
```

Now, point your browser to your Hipstershop uprl using the HTTPS protocol. 
You should see an error !

Yes, this setup is not working !

Let's understand why.
By using the `istioctl log` command, we can check the logs of the Ingress Gateway:

```shell
kubectl -n istio-system logs deployment/istio-ingressgateway

...
        info	sds	resource:hipstershop-cert new connection
        warn	secretfetcher	Cannot find secret hipstershop-cert, searching for fallback secret gateway-fallback
        error	secretfetcher	cannot find secret hipstershop-cert and cannot find fallback secret gateway-fallback
        warn	cache	resource:hipstershop-cert SecretFetcher cannot find secret hipstershop-cert from cache
```

While the error message does not offer the solution to our problem, it points out that the secret `hipstershop-cert` is not found. The problem here is that we created the secret in the `hipstershopv1v2` namespace, while the Ingress Gateway is running in the `istio-syste` namespace. The secret must reside in the same namespace as the Ingress Gateway.

Let's move it:

```yaml
kubectl -n hipstershopv1v2 delete  secret hipstershop-cert

kubectl apply -n istio-system -f - <<EOF
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: hipstershop
spec:
  secretName: hipstershop-cert
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  organization:
  - hipstershop
  commonName: hipstershop.${INGRESSIP}.sslip.io
  isCA: false
  keySize: 2048
  keyAlgorithm: rsa
  keyEncoding: pkcs1
  usages:
    - server auth
    - client auth
  dnsNames:
  - hipstershop.${INGRESSIP}.sslip.io
  ipAddresses:
  - ${INGRESSIP}
  issuerRef:
    name: selfsigned-issuer
    kind: ClusterIssuer
EOF
```

If we look at the logs, we now see:

```yaml
        info	sds	resource:hipstershop-cert pushed key/cert pair to proxy
        info	sds	Dynamic push for secret hipstershop-cert
```

You have just opened your Hipstershop to the world with HTTPS support.

## Creating a dedicated Ingress Gateway

To ensure separation of concerns and isolate teams in their `namespace`, it is possible to create one Ingress Gateway per application.
By restricting each team into it's own namespace and Gateway you prevent one team to steal or block other's team ingress traffic.

We are going to create a new Ingress Gateway in the `hipstershopv1v2` namespace.

As we are using the IstioOperator, we need to update it and add our Gateway definition:

```shell
kubectl edit -n istio-system istiooperator
```

Then, add the following values. This will change depending if you deployed the Operator from ground or if you used `istioctl manifest apply` then `istioctl operator init`:

```yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  namespace: istio-system
  name: installed-state
spec:
  profile: demo

  ...

  components:
    ingressGateways:
    - enabled: true
      name: hipstershop-ingressgateway
      namespace: hipstershopv1v2
      k8s:
        resources:
          requests:
            cpu: 10m
            memory: 40Mi
        service:
          ports:
          - name: http2
            port: 80
            targetPort: 8080
          - name: https
            port: 443
            targetPort: 8443
```

After a few seconds, the IstioOperator will create the new gateway and an associated `LoadBalancer` service:

```shell
kubectl -n hipstershopv1v2 get pods -l istio=ingressgateway

NAME                                          READY   STATUS    RESTARTS   AGE
hipstershop-ingressgateway-57dbd74c87-fg88m   1/1     Running   0          1m17s


kubectl -n hipstershopv1v2 get svc -l istio=ingressgateway

NAME                         TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                                                                      AGE
hipstershop-ingressgateway   LoadBalancer   10.122.13.198   35.111.222.333   15020:30802/TCP,80:30025/TCP,443:31450/TCP,31400:30489/TCP,15443:30344/TCP   11m
```

Let's delete the `Certificate`, `Gateway` and `VirtualService` we created before:


```shell
kubectl -n hipstershopv1v2 delete gw hipstershop
kubectl -n hipstershopv1v2 delete vs hipstershop
kubectl -n istio-system    delete certificate hipstershop
```

As we created a new `LoadBalancer` service, we need to get it's public IP:

```shell
INGRESSIP=$(kubectl -n hipstershopv1v2 get service hipstershop-ingressgateway -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
export INGRESSIP
echo $INGRESSIP
```

Now re-create the `Certificate` in our shop Namespace:

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: hipstershop
spec:
  secretName: hipstershop-cert
  duration: 2160h # 90d
  renewBefore: 360h # 15d
  organization:
  - hipstershop
  commonName: hipstershop.${INGRESSIP}.sslip.io
  isCA: false
  keySize: 2048
  keyAlgorithm: rsa
  keyEncoding: pkcs1
  usages:
    - server auth
    - client auth
  dnsNames:
  - hipstershop.${INGRESSIP}.sslip.io
  ipAddresses:
  - ${INGRESSIP}
  issuerRef:
    name: selfsigned-issuer
    kind: ClusterIssuer
EOF
```

Then, re-create the `Gateway` and change the `selector` to target the local `Gateway`. To do so, we use one of the specific `Labels` of the service (which we can list using `kubectl get pods --show-labels`):

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: hipstershop
spec:
  selector:
    istio: ingressgateway 
    service.istio.io/canonical-name: hipstershop-ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "hipstershop.${INGRESSIP}.sslip.io"
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      mode: SIMPLE
      credentialName: "hipstershop-cert" 
    hosts:
    - "hipstershop.${INGRESSIP}.sslip.io"
EOF
```

Finaly we re-create the `VirtualService`:

```yaml
kubectl apply -n hipstershopv1v2 -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: hipstershop
spec:
  hosts:
  - "hipstershop.${INGRESSIP}.sslip.io"
  gateways:
  - hipstershop
  http:
  - match:
    - uri:
        prefix: /api
    route:
    - destination:
        host: apiservice
        port:
          number: 8080
  - route:
    - destination:
        host: frontend
        port:
          number: 8080
EOF

echo "you can now connect to http://hipstershop.${INGRESSIP}.sslip.io/"
```

## Takeaway

In this module you leveraged the Cert-Manager application to create SSL certificates for public URLs and you configured a gateway to allow traffic to your applications.

---
Next step: [Traffic Routing](/modules/traffic/routing)