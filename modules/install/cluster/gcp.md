Kubernetes Cluster Access
===

For this, we'll be using a Kubernetes cluster in GKE that we control in cloud shell. We'll provide credentials and a cluster name that you can use.
Cluster must be K8s 1.15 or newer (1.16+ recommended)

Log in to [Cloud Shell](https://ssh.cloud.google.com/cloudshell/editor) with the provided credentials.

> These credentials will conflict with your existing Google credentials, e.g. if you're using Chrome. For this reason, we recommend logging in to cloud shell in an incognito window.

Set the project for the Cloud Shell session to the provided project ID:

```shell
gcloud config set project <assigned-project-id>
```

Now you can access your Cluster:

```shell
gcloud container clusters get-credentials <assigned-project-id> --zone <assigned-zone>
```

Verify you have access via `kubectl`:

```shell
$ kubectl get pods --all-namespaces
NAMESPACE     NAME                                                       READY   STATUS    RESTARTS   AGE
kube-system   event-exporter-gke-6c9d8bd8d8-gs9cw                        2/2     Running   0          2m5s
kube-system   fluentd-gke-2hhlp                                          2/2     Running   0          80s
kube-system   fluentd-gke-f5vzm                                          2/2     Running   0          80s
kube-system   fluentd-gke-scaler-cd4d654d7-dr4ch                         1/1     Running   0          2m1s
kube-system   fluentd-gke-wlrbr                                          2/2     Running   0          40s
kube-system   gke-metrics-agent-bjzw7                                    1/1     Running   0          99s
kube-system   gke-metrics-agent-k9mwc                                    1/1     Running   0          99s
kube-system   gke-metrics-agent-zr6vb                                    1/1     Running   0          99s
kube-system   kube-dns-56d8cd994f-kl8vk                                  4/4     Running   0          92s
kube-system   kube-dns-56d8cd994f-qv4vs                                  4/4     Running   0          2m5s
kube-system   kube-dns-autoscaler-645f7d66cf-kdphv                       1/1     Running   0          2m
kube-system   kube-proxy-gke-demo-1-default-pool-a2959341-6v5d           1/1     Running   0          99s
kube-system   kube-proxy-gke-demo-1-default-pool-a2959341-mqzg           1/1     Running   0          99s
kube-system   kube-proxy-gke-demo-1-default-pool-a2959341-rjg6           1/1     Running   0          99s
kube-system   l7-default-backend-678889f899-98v8z                        1/1     Running   0          2m5s
kube-system   metrics-server-v0.3.6-7b7d6c7576-txf5s                     2/2     Running   0          93s
kube-system   prometheus-to-sd-gxlzz                                     1/1     Running   0          99s
kube-system   prometheus-to-sd-lct7t                                     1/1     Running   0          99s
kube-system   prometheus-to-sd-lnml8                                     1/1     Running   0          99s
kube-system   stackdriver-metadata-agent-cluster-level-bdc54b8bb-2bqqq   2/2     Running   0          74s
```

Finally, let's clone the training repository:

```shell
git clone https://github.com/tetrateio/training.git
```
