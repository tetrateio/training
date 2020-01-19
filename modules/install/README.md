Installation
===

This section covers initial cluster creation and the deployment of Istio and our victim (err, demo) App. You can dive into the directory for each section to see detailed explanations for each installation, including how to check that the installation was successful. You should approach them in the order:

1. [Cluster](cluster/)
1. [Istio](istio/)
1. [App](app/)

If you want to skip all of that, you can use the abbreviated guide below to create and initialize a cluster that's ready to start the workshop.

Configure `kubectl` to use our credentials:

- [Google Cloud Platform](cluster/gcp.md)

Deploy Istio. See [the detailed guide](istio/README.md) or skip it and just:

```shell
istioctl manifest apply --set values.global.mtls.enabled=true --set values.global.controlPlaneSecurityEnabled=true
```

Deploy our demo app. See [the detailed guide](app/README.md) or skip it and just:

```shell
cd training
kubectl label namespace default istio-injection=enabled
kubectl apply -f modules/install/app/config
```
