# Istio Workshop

This workshop will introduce you to Istio and Envoy, service mesh tools that have changed the landscape of cloud-native applications. You will learn how and why to use these tools to solve the challenges of observability, security, networking, and multi-cloud deployments.

Throughout this workshop, we will be deploying onto, and interacting with Kubernetes. In this lab, we’ll provide the configuration inline to help you understand what’s happening, but the config is also provided in YAML files that you can apply directly.

A basic understanding of Kubernetes is a prerequisite to these labs. For background on Kubernetes or additional sources on service mesh, we recommend labs like [Kubernetes 101](http://saturnism.me/talk/kubernetes-101) and [Kubernetes Code Lab](http://bit.ly/k8s-lab) by [Ray Tsang](https://saturnism.me/about/), the [Kubernetes website for background](https://kubernetes.io/docs/tutorials/kubernetes-basics/), and environments like the [Play with Kubernetes Classroom](https://training.play-with-kubernetes.com/) to play with a live system. We also recommend and gratefully acknowledge [Ryan Knight's](https://twitter.com/knight_cloud) [Istio Workshop](https://github.com/retroryan/istio-workshop), as well as Ray Tsang's.

We'll go over the content in the following order:
1. [Installation](modules/install/), where we'll:
    1. [Setup your cluster](modules/install/cluster/)
    2. [Deploy Istio](modules/install/istio/)
    3. [Deploy our demo App](modules/install/app/)
2. [Look at Istio's observability features](modules/observability/)
3. [Learn about Istio's traffic management capabilities](modules/traffic)
    1. [Setting up ingress with Gateways](modules/traffic/ingress)
    2. [See fine-grained traffic splitting by canarying a release](modules/traffic/routing)
    3. [Explore resiliency features](modules/traffic/resiliency)
4. [Use Istio's security features](modules/security) to:
    1. [Terminate TLS at ingress](modules/security/ingress)
    2. [Establish mTLS throughout the mesh](modules/security/mtls)
    3. [Use Istio Authz](modules/security/authz)
