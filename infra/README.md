# Infra Setup

## Pre-existing Organization

[Install Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html)

Download a JSON `terraform` GCP service account key for the `infra-provision` project.

Edit the [Terraform input variables](terraform.tfvars) file to meet your requirements.

Run the following commands:

```shell
terraform init
terraform apply
```

> Note: You may have to run apply a couple of times if you get errors like `Permission iam.serviceAccounts.create is required to perform this operation on project`.

## New Google Organization

There are several gotchas when setting this up in a new organization. Including, that you can only set org level permissions on service accounts via the command line. When setting up a new organization in GCP, [this](https://cloud.google.com/community/tutorials/managing-gcp-projects-with-terraform) guide is a good starting point.

## Testing Istio installation

To ensure that there aren't other quota limits you may hit and to verify that Istio can be installed in the clusters you created you can use the go tests in this directory.

> These tests use the `gcloud`, `kubectl` and `istioctl` CLI tools so ensure you have them all and that `gcloud` is logged in with permission to retrieve the GKE cluster credentials.

Remember to set `clusterCount`, `zone` and `workshopName` as required.

```shell
go test . -v --timeout 60m
```

## Teardown

In order to teardown all infrastructure run:

```shell
terraform destroy
```
