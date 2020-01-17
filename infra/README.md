# Infra Setup

[Install Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html)

Download a JSON `terraform` GCP service account key for the `training-infra-owner` project.

Run the following commands:

```shell
terraform init
terraform apply -var="credz_file=<downloaded-key.json>" -var="cluster_prefix=<prefix>" -var="cluster_count=<number-of-clusters-needed>"
```
