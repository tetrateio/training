# Modern Bank

Modern Bank is the sample application used in training and in the future demonstrations of Tetrate software. Currently it consists of a single tenant, the banking services team. This team has an all Golang microservice architecture.

## Developing New Microservices

Modern Bank microservices are mostly automatically generated from an initial contract. These contracts can be found in the [contracts](./contracts).

To create a new microservice define a new API contract with a `<service-name>.yaml` file in the `./contracts/src` directory. Contracts should reference as much as possible from the `./contracts/src/paths` and `./contracts/src/definitions` directories in order to allow re-use across other microservices and the ingress gateway.

Once the contract has been defined run `make generate-services` and it will produced a flattened contract and the microservice code in the `./microservices` directory.

> Note: You need to have [go-swagger installed](https://goswagger.io/install.html) for code generation.

To add the handler for each of the paths defined edit the `./microservices/<service-name>/pkg/server/configure_<service-name>.go` file.

Finally add the corresponding Helm values file to `./kubernetes/helm/values/bankingservices/<service-name>.yaml`.

Optional: If this microservice's API is exposed via ingress you will need to add an Istio Virtual Service into the [gateway.yaml](./networking/bankingservices/gateway.yaml) and add the path to the [banking-ingress.yaml](./contracts/src/banking-ingress.yaml) contract.

## Publishing Docker images

Optional: Update the Go dependencies to the latest minor version.

```bash
make update-dependencies
```

Build and push the Docker images.

```bash
make docker-build
```

## Running in Kubernetes

Currently, there is only one tenancy configuration (as there is only one tenant!). Deploy the banking services to the default namespace, you must be pointed at a Kubernetes cluster with Istio installed. See the [Istio directory](../../istio) in this repo for installation configuration.

```bash
make deploy-kube
```

## Traffic Generation

Once you have deployed the application you may want to hit it with some artificial traffic, to do this use the traffic generation tools. It creates users and accounts, and sends random amounts of money between two random accounts.

```bash
go run tools/trafficGen/cmd/main.go --host <KUBE_LOADBALANCER_IP>
```
