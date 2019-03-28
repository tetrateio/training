Establish mTLS throughout the mesh
====

Enable mTLS for service-to-service
----

1. Apply updated Kubernetes configuration files

    ```
    kubectl apply -f modules/security/mtls/destinationrule.yaml
    ```

Validating mTLS
----

1. Create a new Namespace that doesn't have Istio automatic sidecar injection.

    ```
    $ kubectl create ns noistio
    ```

2. Run a Toolbox Pod in the noistio namespace. Download the sleep.yaml from [here](https://github.com/istio/istio/blob/master/samples/sleep/sleep.yaml).

    ```
    $ kubectl -n noistio apply -f sleep.yaml
    ```

3. Wait for this pod to start.

    ```
    $ kubectl get pods -n noistio -w
    ```

4. Connect to Account Service from the Toolbox that doesn't have Istio mTLS.

    ```
    $ kubectl -n noistio exec -it $(kubectl get pod -n noistio -l app=sleep -o jsonpath='{.items..metadata.name}') -- curl  http://account.default/v1/users/<USERNAME>/accounts
    ```

Apply strict policy
----

1. Apply strict mTLS mesh policy

    ```
    kubectl appy -f modules/security/mtls/strictpolicy.yaml
    ```

2. Run the last step in previous task to verify it no longer connects
