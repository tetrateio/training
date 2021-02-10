Istio Authz
===

Istio provides every workload with a strong identity - in Kubernetes, the pod's ServiceAccount - which is used to establish mTLS connections between services in the mesh. While establishing mTLS connections, sidecars in the mesh will validate certificates according to [SPIFFE's X.509 SVID spec](https://github.com/spiffe/spiffe/blob/master/standards/X509-SVID.md), which means that after the connection is established we have an authenticated identity of the other party. Istio allows you to write access control policies using those identities to describe which services can communicate.

Lets use the `AuthorizationPolicy` CRD to test end-user-to-workload authorization policies

1. Create the `deny-all` Authz Policy

    Now, let's create a `AuthorizationPolicy` that doesnt allow any access to the workloads in the `default` namespace:
    ```shell
    $ kubectl apply -f modules/security/authz/config/deny-all.yaml
    ```

    We can describe the object to see what the policy looks like in Istio:
    ```shell
    $ kubectl describe authorizationpolicy.security.istio.io/deny-all 

    ...
    Spec:
    ```

    Try accessing the UI, it should fail.
    ```shell
    $ curl $INGRESS_IP
    RBAC: access denied
    ```

2. Create a AuthorizationPolicy to allow access to the UI

    Lets create a `AuthorizationPolicy` to allow you to access the UI from the ingress gateway. The `principals` field is set to allow all
    traffic originating from the `istio-ingressgateway-service-account` service account in the `istio-system` namespace, to access any `GET`
    API exposed by Istio

    ```shell
    kubectl apply -f modules/security/authz/config/allow-ui.yaml
    ```

    And we can look at the `AuthorizationPolicy` to see how it has allowed access to users:

    ```yaml
    $ kubectl describe authorizationpolicy.security.istio.io/allow-ui 

    ...
    Spec:
      Action:  ALLOW
      Rules:
        From:
          Source:
            Principals:
              cluster.local/ns/istio-system/sa/istio-ingressgateway-service-account
        To:
          Operation:
            Methods:
              GET
            Paths:
              *

    ```

    We should see that we can now access the UI:
    ```shell
    $ curl $INGRESS_IP
    ...<html>
    ```

4. Cleanup

    Unfortunately, traffic in the rest of our app is still broken by our `deny-all` policy so the UI page that gets rendered is incomplete. We'd need to describe a full set of Authz policies for the application to allow all traffic. For now, lets clean up our Authz config so we can carry on with the rest of the lab:

    ```shell
    kubectl delete -f modules/security/authz/config 
    ```
