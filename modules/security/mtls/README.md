Establish mTLS throughout the mesh
====

Enable mTLS for service-to-service
----

1. Apply updated Kubernetes configuration files

    ```
    $ kubectl apply -f modules/security/mtls/config/destinationrule.yaml
    ```

    This creates two DestinationRules, the first of which forces Istio to use mTLS everywhere:
    ```
    $ kubectl describe -n istio-system destinationrule.networking.istio.io/default

    ...
    Spec:
      Host:  *.local
      Traffic Policy:
        Tls:
           Mode:  ISTIO_MUTUAL
    ```

    And the second which excludes the Kubernetes API server from mTLS:
    ```
    $ kubectl describe -n istio-system destinationrule.networking.istio.io/api-server

    ...
    Spec:
      Host:  kubernetes.default.svc.cluster.local
      Traffic Policy:
        Tls:
          Mode:  DISABLE
    ```

Validating mTLS
----

1. Create a new Namespace that doesn't have Istio automatic sidecar injection.

    ```
    $ kubectl create ns noistio
    ```

2. Run the `sleep` Toolbox Pod in the noistio namespace.:.

    ```
    $ kubectl -n noistio apply -f modules/security/mtls/config/sleep.yaml
    ```

3. Wait for this pod to start.

    ```
    $ kubectl get pods -n noistio -w
    ```

4. Connect to Account Service from the Toolbox that doesn't have Istio mTLS.

    ```
    $ export USERNAME=<login name from demo site>
    $ kubectl -n noistio exec -it $(kubectl get pod -n noistio -l app=sleep -o jsonpath='{.items..metadata.name}') -- curl  http://user.default/v1/users/${USERNAME}/accounts
    ```

Apply strict policy
----

1. Apply strict mTLS mesh policy

    ```
    $ kubectl apply -f modules/security/mtls/config/strictpolicy.yaml
    ```

    We can verify the policy matches what we expect:
    ```
    $ kubectl describe meshpolicy.authentication.istio.io/default

    ...
    Spec:
      Peers:
        Mtls:
          Mode:  STRICT
    ```

1. Run the last step in previous task to verify it no longer connects:

    ```
    $ kubectl -n noistio exec -it $(kubectl get pod -n noistio -l app=sleep -o jsonpath='{.items..metadata.name}') -- curl  http://user.default/v1/users/${USERNAME}/accounts

    curl: (56) Recv failure: Connection reset by peer
    command terminated with exit code 56
    ```

    If everything worked correctly, we should see curl return an error code 56, indicating if failed to establish a TLS connection.
