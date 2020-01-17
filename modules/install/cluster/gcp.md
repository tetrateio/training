Kubernetes Cluster Access
===

For this training we'll be using a Kubernetes cluster in GKE that we control in cloud shell. We will provide credentials and a cluster name that you can use.

Log in to [Cloud Shell](https://ssh.cloud.google.com/cloudshell/editor) with the provided credentials.

> These credentials will conflict with your existing Google credentials, e.g. if you're using Chrome. For this reason, we recommend logging in to cloud shell in an incognito window.

Now we can get access to our Cluster:

```shell
gcloud container clusters get-credentials <assigned-cluster-name> --zone <assigned-zone>
```

We should also clone the training repository:

```shell
git clone https://github.com/tetrateio/training.git
```
