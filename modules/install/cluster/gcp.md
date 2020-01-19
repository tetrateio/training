Kubernetes Cluster Access
===

For this training we'll be using a Kubernetes cluster in GKE that we control in cloud shell. We will provide credentials and a cluster name that you can use.

Log in to [Cloud Shell](https://ssh.cloud.google.com/cloudshell/editor) with the provided credentials.

> These credentials will conflict with your existing Google credentials, e.g. if you're using Chrome. For this reason, we recommend logging in to cloud shell in an incognito window.

Set the project for the Cloud Shell session to the provided project ID:

```shell
gcloud config set project <assigned-project-id>
```

Now we can get access to our Cluster:

```shell
gcloud container clusters get-credentials <assigned-project-id> --zone <assigned-zone>
```

Verify we have access via `kubectl`:

```shell
$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                                      READY   STATUS    RESTARTS   AGE
kube-system   kube-dns-6cd7bbdf65-hk7v2                                 4/4     Running   0          56m
kube-system   kube-dns-6cd7bbdf65-vqggd                                 4/4     Running   0          55m
kube-system   kube-dns-autoscaler-bb58c6784-dkfnp                       1/1     Running   0          55m
kube-system   kube-proxy-gke-nist-2020-000-default-pool-1131490b-c2r3   1/1     Running   0          56m
kube-system   kube-proxy-gke-nist-2020-000-default-pool-1131490b-q6l8   1/1     Running   0          56m
kube-system   kube-proxy-gke-nist-2020-000-default-pool-1131490b-xt0z   1/1     Running   0          56m
kube-system   l7-default-backend-fd59995cd-8t5qm                        1/1     Running   0          56m
kube-system   metrics-server-v0.3.1-57c75779f-9xj8r                     2/2     Running   0          55m
```

Finally, let's clone the training repository:

```shell
git clone https://github.com/tetrateio/training.git
```
