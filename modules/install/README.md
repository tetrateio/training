Installation
===

This section covers initial cluster creation and the deployment of Istio and our victim (err, demo) App. You can dive into the directory for each section to see detailed explanations for each installation, including how to check that the installation was successful. You should approach them in the order:
1. [Cluster](cluster/)
2. [Istio](istio/)
4. [App](app/)

If you want to skip all of that, you can use the abbreviated guide below to create and initialize a cluster that's ready to start the workshop.

1. Configure `kubectl` to use our credentials:
    
    ```shell
    # point kubectl at this config file
    export KUBECONFIG=${PWD}/k8s-01/config

    # ensure it works
    kubectl get ns
    NAME          STATUS   AGE
    default       Active   31m
    kube-public   Active   31m
    kube-system   Active   31m
    ```

    > Yours won't be literally `k8s-01`, but some number 1-25, e.g. `k8s-25`.

    If you get an error like:
    ```
    Error in configuration:
    * unable to read client-cert /home/ubuntu/cloud-native-summit/k8s-01/cert.pem for admin due to open /home/ubuntu/cloud-native-summit/k8s-01/cert.pem: no such file or directory
    * unable to read client-key /home/ubuntu/cloud-native-summit/k8s-01/key.pem for admin due to open /home/ubuntu/cloud-native-summit/k8s-01/key.pem: no such file or directory
    * unable to read certificate-authority /home/ubuntu/cloud-native-summit/k8s-01/ca.pem for k8s-01 due to open /home/ubuntu/cloud-native-summit/k8s-01/ca.pem: no such file or directory
    ```
    Then you'll need to update the configuration file's certificate paths. You can set the paths in the config file to be relative, so your final config should look like:
    ```
    apiVersion: v1
    clusters:
    - cluster:
        server: ""
    name: cloud-native-summit
    - cluster:
        certificate-authority: ca.pem
        server: -
    name: k8s-01
    contexts:
    - context:
        cluster: k8s-01
        user: admin
    name: default
    current-context: default
    kind: Config
    preferences: {}
    users:
    - name: admin
    user:
        client-certificate: cert.pem
        client-key: key.pem
    ```

1. Deploy Istio
See [the detailed guide](istio/README.md) or skip it and just:
    ```shell
    kubectl apply -f modules/install/istio/config --as=admin --as-group=system:masters
    ```

1. Deploy our demo app
See [the detailed guide](app/README.md) or skip it and just:
    ```shell
    kubectl label namespace default istio-injection=enabled
    kubectl apply -f modules/install/app/config
    ```

1. Export the addresses for the cluster's ingress:
    
    ```shell
    export INGRESS_IP=$(kubectl -n istio-system get svc istio-ingressgateway \
          -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    echo Ingress: $INGRESS_IP
    ```
