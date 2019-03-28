Terminate TLS at ingress
====

Generate client and server certificates and keys
----

1.  Clone the <https://github.com/nicholasjackson/mtls-go-example> repository:

    ```
    $ git clone https://github.com/nicholasjackson/mtls-go-example
    ```

1.  Change directory to the cloned repository:

    ```
    $ pushd mtls-go-example
    ```

1.  Generate the certificates for `httpbin.example.com`. Change `password` to any value you like in the following command:

    ```
    $ ./generate.sh httpbin.example.com password
    ```

    When prompted, select `y` for all the questions. The command will generate four directories: `1_root`,
   `2_intermediate`, `3_application`, and `4_client` containing the client and server certificates you use in the
    procedures below.

1.  Move the certificates into a directory named `httpbin.example.com`:

    ```
    $ mkdir ~+1/httpbin.example.com && mv 1_root 2_intermediate 3_application 4_client ~+1/httpbin.example.com
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

1. Validate that you can still visit the UI from the browser and now it is redirect to HTTPS.

    ```
    $ export INGRESS_IP=$(kubectl -n istio-system get svc istio-ingressgateway \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    $ curl -v $INGRESS_IP
    ```