Istio RBAC for Service to Service Communication
===

Istio provides every workload with a strong identity - in Kubernetes, the pod's ServiceAccount - which is used to establish mTLS connections between services in the mesh. While establishing mTLS connections, sidecars in the mesh will validate certificates according to [SPIFFE's X.509 SVID spec](https://github.com/spiffe/spiffe/blob/master/standards/X509-SVID.md), which means that after the connection is established we have an authenticated identity of the other party. Istio allows you to write access control policies using those identities to describe which services can communicate.

This config comes in three parts: a cluster wide `ClusterRbacConfig` which sets the mesh's RBAC mode (on, off, inclusive list, exclusive list). We'll start by creating one which requires RBAC for the `default` namespace:

```yaml
$ cat modules/security/rbac/config/clusterrbacconfig.yaml
apiVersion: "rbac.istio.io/v1alpha1"
kind: ClusterRbacConfig
metadata:
  name: default
spec:
  mode: 'ON_WITH_INCLUSION'
  inclusion:
    namespaces:
    - "default"
```

Let's apply the rule and see what happens:

```shell
kubectl apply -f modules/security/rbac/config/clusterrbacconfig.yaml
```

> The other `mode`s are `ON`, `OFF`, and `ON_WITH_EXCLUSION`

If we go to the UI, we should see it fail to load with an RBAC related error. We can confirm this by looking at our graphs too.

```shell
$ curl $INGRESS_IP
RBAC: access denied
```

So, let's create a `Role` that allows read to use the UI:

```shell
$ cat modules/security/rbac/config/ui-viewer.yaml
apiVersion: "rbac.istio.io/v1alpha1"
kind: ServiceRole
metadata:
  name: ui-viewer
  namespace: default
spec:
  rules:
  - services: ["ui.default.svc.cluster.local"]
    methods: ["GET"]
```

And apply that as well:
```shell
kubectl apply -f modules/security/rbac/config/ui-viewer.yaml
```

Of course, a `Role` alone doesn't do us any good, we have to assign the role to ourselves so it applies. We'll do this by creating a `ServiceRoleBinding` which lets anyone see the UI:
```yaml
$ cat modules/security/rbac/config/allow-all-ui-access.yaml
apiVersion: "rbac.istio.io/v1alpha1"
kind: ServiceRoleBinding
metadata:
  name: bind-ui-viewer
  namespace: default
spec:
  subjects:
  - user: "*"
  roleRef:
    kind: ServiceRole
    name: "ui-viewer"
```

Which we'll apply like the others:
```shell
kubectl apply -f modules/security/rbac/config/allow-all-ui-access.yaml
```

We should see that we can now access the UI:
```shell
$ curl $INGRESS_IP
...<html>
```

Of course, the traffic in the rest of our app is still broken, so the UI page that gets rendered is incomplete. We'd need to describe a full set of RBAC policies for the application to allow all traffic. For now, lets clean up our Cluster RBAC Config so we can carry on with the rest of the lab:

```shell
kubectl delete -f modules/security/rbac/config
```