# Ingress Gateway

Istio provide an Ingress Gateway, used to replace the default Kubernetes `Ingress` resources. The Istio Ingress Gateway allows us to use all the Mesh features at the level 7 OSI layer.
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

We can now apply the Gateway manifest:

```shell
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

Then we create a `VirtualService` that will link the Ingress Gateway to the Hipstershop `frontend` service:

```shell
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

You can open your browser at [http://hipstershop.${INGRESSIP}.sslip.io/](http://hipstershop.${INGRESSIP}.sslip.io/)

Now that we have a working HTTP gateway, let's add SSL.
To do that we are going to use Cert-Manager, an application that can generate SSL certificates into Kubernetes secrets. Istio-proxies will then requests the SSL certificates to Istiod using the SDS protocol. Istiod will read the secrets and send the certificate to the proxies.

Let's install Cert-Manager:

```shell
kubectl create namespace cert-manager
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v0.14.1/cert-manager.yaml
```

We will use Self-Signed certificates for this training, but Cert-Manager is able to handle a lot of other Issuers, like ACME/Let's Encrypt:

```shell
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

```shell
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
k get secret hipstershop-cert -o yaml

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

```shell
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

While the error message does not offer the solution to our problem, it points out that the secret `hipstershop-cert` is not found. The problem here is that we created the secret in the `hipstershopv1v2` namespace, while the Ingress Gateway is running in the `istio-syste` namespace.
Let's move it:

```shell
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
