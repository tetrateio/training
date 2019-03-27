Installation
===

This section covers initial cluster creation and the deployment of Istio, SkyWalking, and our victim (err, demo) App. You can dive into the directory for each section to see detailed explanations for each installation, including how to check that the installation was successful. You should approach them in the order:
1. [Cluster](cluster/)
2. [Istio](istio/)
3. [SkyWalking](skywalking/)
4. [App](app/)

If you want to skip all of that, you can use the abbreviated guide below to create and initialize a cluster that's ready to start the workshop.

1. Create your cluster:

```shell
# Create the cluster
gcloud container clusters create --machine-type n1-standard-2 --num-nodes 4 workshop

# Get the credentials for kubectl
gcloud container clusters get-credentials workshop
```

2. Deploy Istio
See [the detailed guide](istio/README.md) or skip it and just:
```shell
kubectl apply -f modules/install/istio/config --as=admin --as-group=system:masters
```

3. Deploy Skywalking
See [the detailed guide](skywalking/README.md) or skip it and just:
```shell
kubectl apply -f modules/install/skywalking/config --as=admin --as-group=system:masters
```

4. Deploy our demo app
See [the detailed guide](app/README.md) or skip it and just:
```shell
kubectl label namespace default istio-injection=enabled
kubectl apply -f modules/install/app/config
```

5. Export the addresses for the cluster's ingress and SkyWalking's UI:

```shell
export INGRESS_IP=$(kubectl -n istio-system get svc istio-ingressgateway \
      -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo Ingress: $INGRESS_IP

export SKYWALKING_IP=$(kubectl -n skywalking get svc ui \
      -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo SkyWalking: $SKYWALKING_IP
```
