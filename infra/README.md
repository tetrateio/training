# Infra Setup

[Install Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html)

Download a JSON `terraform` GCP service account key for the `training-infra-owner` project.

Edit the [Terraform input variables](terraform.tfvars) file to meet your requirements.

Run the following commands:

```shell
terraform init
terraform apply
```

## Testing Istio installation

To ensure that there aren't other quota limits you may hit and to verify that Istio can be installed in the clusters you created you can use the go tests in this directory.

> These tests use the `gcloud`, `kubectl` and `istioctl` CLI tools so ensure you have them all and that `gcloud` is logged in with permission to retrieve the GKE cluster credentials.

```shell
go test . -v --timeout 60m
```
