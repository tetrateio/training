# Terminate TLS at ingress


We already leveraged the default Ingress Gateway from the `istio-system` Namespace at the setup time. Then we installed a specific gateway in the Hipstershop's namespace.
For the sake of the demo, we are now going to create an SSL certificate and use it to configure another Gateway to the same app.

Let's configure our Mesh se we can reach the application through the name `storefront.<your-ingress-ip>.sslip.io`

## Find the Gateway's IP

As we did before, we need the IP address of the Ingress Gateway. If you don't already have it in a shell variable, run this command:

```shell
INGRESSIP=$(kubectl -n hipstershopv1v2 get service hipstershop-ingressgateway -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
export INGRESSIP
echo "our URL is storefront.$INGRESSIP.sslip.io"
```

You can check connection to the URL:

```shell
curl -vs http://storefront.$INGRESSIP.sslip.io

HTTP/1.1 404 Not Found
```

You should get an HTTP 404 error. This is because the DNS is pointing to the same IP as the Gateway we created earlier. So, as of now, Istio (Envoy) already bind the port 80 and is receiving our request. As it does not know about our `storefront.$INGRESSIP.sslip.io` name, it's answering with a 404.

Let's try with HTTPS:

```shell
curl -vs https://storefront.$INGRESSIP.sslip.io

* TLSv1.2 (OUT), TLS handshake, Client hello (1):
* LibreSSL SSL_connect: SSL_ERROR_SYSCALL in connection to storefront.$INGRESSIP.sslip.io:443
```

Here we have an issue negociating the SSL handshake. This is because we're asking for `storefront` while the server only knows about `hipstershop`.

## Generate client and server certificates and keys

1.  Change directory to the `mtls-go-exmaple` directory:

    ```
    pushd modules/security/ingress/mtls-go-example
    ```

    > It's important to use `pushd` here so our command to gather up all the certificates we're producing works correctly - if you `cd` then in step 3 you'll have to find where the `storefront.$INGRESSIP.sslip.io` directory is created.

1.  Generate the certificates for `storefront.$INGRESSIP.sslip.io`. Change `password` to any value you like in the following command:

    ```
    ./generate.sh storefront.$INGRESSIP.sslip.io password
    ```

    When prompted, select `y` for all the questions. The command will generate four directories: `1_root`,
   `2_intermediate`, `3_application`, and `4_client` containing the client and server certificates you use in the
    procedures below.

1.  Move the certificates into a directory named `storefront`:

    ```
    mkdir ~-/storefront && mv 1_root 2_intermediate 3_application 4_client ~-/storefront
    ```

1.  Go back to your previous directory:

    ```
    popd
    ```

## Configure a TLS ingress gateway

1.  Create a Kubernetes secret to hold the server's certificate and private key.
    Use `kubectl` to create the secret `istio-ingressgateway-certs` in namespace
    `hipstershopv1v2` .
    This secret is always used (mounted) by Istio Ingress Gateways. If you don't want to use this secret you can use the SDS mecanisme that we used at the installation time. 
    Remember the secret have to be in the same Namespace as the Gateway.

    ```
    kubectl create -n hipstershopv1v2 secret tls istio-ingressgateway-certs --key storefront/3_application/private/storefront.$INGRESSIP.sslip.io.key.pem --cert storefront/3_application/certs/storefront.$INGRESSIP.sslip.io.cert.pem

    secret/istio-ingressgateway-certs created
    ```

1.  Define a `Gateway` with an HTTPS port and redirect traffic

    ```yaml
    kubectl apply -n hipstershopv1v2 -f - <<EOF
    apiVersion: networking.istio.io/v1alpha3
    kind: Gateway
    metadata:
      name: storefront
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
          - "storefront.${INGRESSIP}.sslip.io"
      - port:
          number: 443
          name: https
          protocol: HTTPS
        tls:
          mode: SIMPLE
          serverCertificate: /etc/istio/ingressgateway-certs/tls.crt
          privateKey: /etc/istio/ingressgateway-certs/tls.key
        hosts:
        - "storefront.${INGRESSIP}.sslip.io"
    EOF
    ```
1. Check again using the curl command

    As we did before, let's start with an HTTP curl command:

    ```shell
    curl -vs http://storefront.$INGRESSIP.sslip.io

    HTTP/1.1 404 Not Found
    ```

    Nothing changed... Well, the Ingress Gateway was already listening on the HTTP port 80, so, effectively, nothing changed.
    Let's try using HTTPS:

    ```shell
    curl -vs https://storefront.$INGRESSIP.sslip.io

    * TLSv1.2 (OUT), TLS handshake, Client hello (1):
    * TLSv1.2 (IN), TLS handshake, Server hello (2):
    * TLSv1.2 (IN), TLS handshake, Certificate (11):
    * TLSv1.2 (OUT), TLS alert, unknown CA (560):
    * SSL certificate problem: unable to get local issuer certificate
    ```

    Now the error changed. We have an `unknown CA`. 

    This is perfectly true as we are using a self-signed certificate and `curl` is trying to validate it by default. We can test again using the `-k` option to skip the verification:

    ```shell
    curl -kvs https://storefront.$INGRESSIP.sslip.io

    * Server certificate:
    *  subject: C=US; ST=Denial; L=Springfield; O=Dis; CN=storefront.$INGRESSIP.sslip.io
    *  start date: Jul 13 12:10:27 2020 GMT
    *  expire date: Jul 23 12:10:27 2021 GMT
    *  issuer: C=US; ST=Denial; O=Dis; CN=storefront.$INGRESSIP.sslip.io
    *  SSL certificate verify result: unable to get local issuer certificate (20), continuing anyway.

    < HTTP/2 404
    ```

    We now have a 404 error, as with HTTP.

1. Adding a VirtualService

    Thanks to the Gateway, our request is accepted by Istio (Envoy). But it is still not linked to any Service known by the Mesh.

    As we did before, we need a `VirtualService` to be added to make the link between the traffic reveived by the Gateway and the Service:

    ```yaml
    kubectl apply -n hipstershopv1v2 -f - <<EOF
    apiVersion: networking.istio.io/v1alpha3
    kind: VirtualService
    metadata:
      name: storefront
    spec:
      hosts:
      - "storefront.${INGRESSIP}.sslip.io"
      gateways:
      - storefront
      http:
      - route:
        - destination:
            host: frontend.hipstershopv1v2.svc.cluster.local
            port:
              number: 8080
    EOF

    echo "you can now connect to http://storefront.${INGRESSIP}.sslip.io/"
    ```


1. Validate you can still reach the UI in your browser: you should now get a certificate warning (because we used a self-signed certificate), but the page _is_ served over HTTPS. We can verify the same with `curl`, setting the CA cert so `curl` accept's Envoy's certificate: 

    ```
    curl -vs --cacert storefront/2_intermediate/certs/ca-chain.cert.pem https://storefront.${INGRESSIP}.sslip.io/
    ```

1. Discussion

    In this example, we re-created a full set of `Gateway` and `VirtualService` for the new URL. Do we have any other solution ?

    Answer is: yes.

    We could have edited the `hipstershop` resources and just added the new name in them:

    ```yaml
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
        - "storefront.${INGRESSIP}.sslip.io"
      - port:
          number: 443
          name: https
          protocol: HTTPS
        tls:
          mode: SIMPLE
          credentialName: "hipstershop-cert" 
        hosts:
        - "hipstershop.${INGRESSIP}.sslip.io"
        - "storefront.${INGRESSIP}.sslip.io"
    ---
    apiVersion: networking.istio.io/v1alpha3
    kind: VirtualService
    metadata:
      name: hipstershop
    spec:
      hosts:
      - "hipstershop.${INGRESSIP}.sslip.io"
      - "storefront.${INGRESSIP}.sslip.io"
      gateways:
      - hipstershop
      http:
      - match:
        - uri:
            prefix: /api
        route:
        - destination:
            host: apiservice.hipstershopv1v2.svc.cluster.local
            port:
              number: 8080
      - route:
        - destination:
            host: frontend.hipstershopv1v2.svc.cluster.local
            port:
              number: 8080
    ```

    By creating another set of resources, we were able to only forward traffig to the `frontend` service, and nothing to the `/api` endpoint. 
    Do you see any other way to achieve that ? 

1. Finally, let's clean up so we can continue to `curl` in future examples without all the crazy flags!

    ```sh
    kubectl delete gw storefront
    kubectl delete vs storefront
    ```

    And a `curl` should verify we're back to our original state:

    ```sh
    curl -kvs https://storefront.$INGRESSIP.sslip.io
    ...
    ```

## Takeaway

Once again, we used the Istio Ingress Gateway to demonstrate how to accept and route traffic for two different applications.

---
Next step: [Mutual TLS](/modules/security/mtls)
