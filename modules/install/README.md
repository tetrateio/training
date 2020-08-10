# Installation

This section covers initial cluster creation and the deployment of Istio and our victim (err, demo) App. You can dive into the directory for each section to see detailed explanations for each installation, including how to check that the installation was successful. You should approach them in the order:

1. [Cluster](cluster/)
2. [Istio](istio/)
3. [App](app/)

If you want to skip all of that, you can use the abbreviated guide below to create and initialize a cluster that's ready to start the workshop.

## Needed tools

stern
kubectx/kubens

## Abbreviated Guide

Configure `kubectl` to use our credentials:

- [Google Cloud Platform](cluster/gcp.md)

Deploy Istio. See [the detailed guide](istio/README.md) or skip it and just:

```shell
istioctl install --set profile=demo
```

Deploy our demo app. See [the detailed guide](app/README.md) or skip it and just:

```shell
kubectl create namespace hipstershopv1v2
kubectl label namespace hipstershopv1v2 istio-injection=enabled
kubectl apply -f https://raw.githubusercontent.com/tetratelabs/microservices-demo/master/release/hipstershop-v1-v2.yaml
```

---
Next step: [Traffic Management](/modules/traffic)