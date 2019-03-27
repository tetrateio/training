Cluster Creation
===

For this training we'll be using a Kubernetes cluster in GKE. We will provide credentials that you can use. These credentials will conflict with your existing Google credentials, e.g. if you're using Chrome. For this reason, we recommend logging in to console.cloud.google.com in an incognito window.

First, log in to the new account via `gcloud`, which will launch a new browser window you can log in with:
```shell
gcloud auth login
```
> This will open a new tab in your browser; this is the page you should open in an incognito tab if you're using Chrome.

We'll be using the GKE API, which itself depends on the Google Compute API; we'll use `gcloud` to enable both:
```shell
gcloud services enable compute.googleapis.com container.googleapis.com
```

The last `gcloud` setup step we need to set the default zone and region; since we're in San Francisco we'll use `us-west1`:
```shell
gcloud config set compute/region us-west1
gcloud config set compute/zone us-west1-b
```

Finally, we can create our Cluster:
```shell
gcloud container clusters create --machine-type n1-standard-2 --num-nodes 4 workshop

gcloud container clusters get-credentials workshop
```

> We create a larger cluster than most demos, 4 standard-2 nodes, because some of the SkyWalking deployments require a decent bit of RAM (as they're JVM based).
