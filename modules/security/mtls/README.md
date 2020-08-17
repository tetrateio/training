Establish mTLS throughout the mesh
====

By default, Istio's _authentication policy_, which controls mTLS, is automatically configure as `permissive`. This means Istio will use mTLS between workloads that both have a Sidecar and will revert to plaintext otherwise (HTTP instead of HTTPS, for example). This is a great way to onboard new clusters and workloads.
Permissive mode is not really secure because any client can request to downgrade to clear text - it exists only to help transition your mesh into mTLS without breaking existing applications. At some point, you're going to need to limit non-mTLS traffic.


## Installing client

To demontrate how traffic is managed we are going to install a client in a namespace without Istio. 
We are using the `sleep` pod, which is a simple pod with few tools, like curl.

1. Create a new Namespace that doesn't have Istio automatic sidecar injection.

    ```sh
    kubectl create namespace noistio
    ```

2. Run the `sleep` Toolbox Pod in the noistio namespace.

    ```sh
    kubectl -n noistio apply -f modules/security/mtls/config/sleep.yaml
    ```

3. Wait for this pod to start.

    ```sh
    kubectl get pods -n noistio -w
    ```

4. Connect to Account Service from the Toolbox that doesn't have Istio mTLS.

    ```sh
    kubectl -n noistio exec -it $(kubectl get pod -n noistio -l app=sleep -o jsonpath='{.items..metadata.name}') -- curl  -s http://frontend.hipstershopv1v2:8080 -o /dev/null -w '%{http_code}'
    ```

    You should get the full HTML page of the shop.


## Applying strict policy

You can disable permissive mode and enforce mTLS across the entire mesh by setting a global `PeerAuthentication` in the Istio ROOT namespace (`istio-system` by default):

1. Apply strict mTLS mesh policy

    ```yaml
    kubectl apply -f - <<EOF
    apiVersion: "security.istio.io/v1beta1"
    kind: "PeerAuthentication"
    metadata:
      name: "default"
      namespace: "istio-system"
    spec:
      mtls:
        mode: STRICT
    EOF
    ```

2. Re-run the curl command from previous step in to verify you can no longer connect:

    ```sh
    kubectl -n noistio exec -it $(kubectl get pod -n noistio -l app=sleep -o jsonpath='{.items..metadata.name}') -- curl  -s http://frontend.hipstershopv1v2:8080 -o /dev/null -w '%{http_code}'

    000
    curl: (56) Recv failure: Connection reset by peer
    command terminated with exit code 56
    ```

    If the mTLS configuration applied correctly, we should see curl return an error code 56, indicating it failed to establish a TLS connection.


3. Allow the sleep pod to reach some workloads

    While you usually don't want any pod from external namespaces to reach your application's workloads (your Hipstershop micro-services), you may still want some of them to be reachable. 
    Let's allow our sleep container to reach the frontend service:

    ```yaml
    kubectl apply -f - <<EOF
    apiVersion: "security.istio.io/v1beta1"
    kind: "PeerAuthentication"
    metadata:
      name: "frontend"
      namespace: "hipstershopv1v2"
    spec:
      selector:
        matchLabels:
          app: frontend.hipstershopv1v2.svc.cluster.local
      mtls:
        mode: PERMISSIVE
    EOF
    ```

    Let's play the same curl command:

    ```sh
    kubectl -n noistio exec -it $(kubectl get pod -n noistio -l app=sleep -o jsonpath='{.items..metadata.name}') -- curl  -s http://frontend.hipstershopv1v2:8080 -o /dev/null -w '%{http_code}'

    200
    ```

    Now, check the `adservice`:

    ```sh
    kubectl -n noistio exec -it $(kubectl get pod -n noistio -l app=sleep -o jsonpath='{.items..metadata.name}') -- curl  -s http://adservice.hipstershopv1v2:9555 -o /dev/null -w '%{http_code}'

    000
    curl: (56) Recv failure: Connection reset by peer
    command terminated with exit code 56
    ```

    As you can see, we are now allowed to reach the `frontend` micro-service but we are still denied to reach others as our client can't enforce the required mTLS. 
    Note that this can't be considered a security mesure. We never denied traffic to go to the `adservice`. It's just that the `adservice` sidecar does not know us, and don't let us in.

    This is called Authentication (AuthN) policy and ensures we all know each other before allowing communication. In other words, we use the certificates from Istio mutual TLS to identify services to each other in a verifiable way (i.e. we _authenticate_ them).

    In the next chapter we are going to dive into Authorization (AuthZ), where we set the rules deciding `who` can talk `to` whom, and `when` they're allowed to.

---
Next step: [Authorization](/modules/security/authorization)
