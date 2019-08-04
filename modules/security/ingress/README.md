Terminate TLS at ingress
====

Generate client and server certificates and keys
----

1.  Change directory to the `mtls-go-exmaple` directory:

    ```
    $ pushd modules/security/ingress/mtls-go-example
    ```

    > It's important to use `pushd` here so our command to gather up all the certificates we're producing works correctly - if you `cd` then in step 3 you'll have to find where the `httpbin.example.com` directory is created.

1.  Generate the certificates for `httpbin.example.com`. Change `password` to any value you like in the following command:

    ```
    $ ./generate.sh httpbin.example.com password
    ```

    When prompted, select `y` for all the questions. The command will generate four directories: `1_root`,
   `2_intermediate`, `3_application`, and `4_client` containing the client and server certificates you use in the
    procedures below.

1.  Move the certificates into a directory named `httpbin.example.com`:

    ```
    $ mkdir ~-/httpbin.example.com && mv 1_root 2_intermediate 3_application 4_client ~-/httpbin.example.com
    ```

1.  Go back to your previous directory:

    ```
    $ popd
    ```

Configure a TLS ingress gateway
----

1.  Create a Kubernetes secret to hold the server's certificate and private key.
   Use `kubectl` to create the secret `istio-ingressgateway-certs` in namespace
   `istio-system` . The Istio gateway will load the secret automatically.

    ```
    $ kubectl create -n istio-system secret tls istio-ingressgateway-certs --key httpbin.example.com/3_application/private/httpbin.example.com.key.pem --cert httpbin.example.com/3_application/certs/httpbin.example.com.cert.pem
    secret "istio-ingressgateway-certs" created
    ```

1.  Define a `Gateway` with an HTTPS port and redirect traffic

    ```
    $ kubectl apply -f modules/security/ingress/config/gateway.yaml
    ```

1. Validate you can still reach the UI in your browser: you should now get a certificate warning (because we used a self-signed certificate), but the page _is_ served over HTTPS. We can verify the same with `curl`, setting the CA cert so `curl` accept's Envoy's certificate: 

    ```
    $ export INGRESS_IP=$(kubectl -n istio-system get svc istio-ingressgateway \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    $ curl -v -HHost:httpbin.example.com --resolve httpbin.example.com:443:$INGRESS_IP --cacert httpbin.example.com/2_intermediate/certs/ca-chain.cert.pem https://httpbin.example.com
    ```
    > All of the fancy flags on the command configure the way `curl` sets the Server Name Indication (SNI) on the request; our Ingress Envoy uses that SNI to serve its certificate, but also for routing.

1. Finally, let's clean up so we can continue to `curl` in future examples without all the crazy flags!

    ```sh
    $ kubectl apply -f modules/traffic/ingress/config/gateway.yaml
    ```

    And a `curl` should verify we're back to our original state:

    ```sh
    $ curl $INGRESS_IP/
    ...
    ```