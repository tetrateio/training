# Workshop
variable "credz_file" {
  type = string
  description = "The fully qualified location of the terraform JSON GCP service account key for the training-infra-owner project."
}

variable "organization_id" {
  type = string
  description = "The GCP organization to deploy the infra into"
}

variable "billing_account" {
  type = string
  description = "The GCP billing account"
}

variable "region" {
  type = string
  default = "us-central1"
  description = "The region the workshop will be in."
}

variable "workshop_name" {
  type = string
  description = "The workshop name e.g. nist-2020. This will be used to prefix all created resources."
}

variable "participant_count" {
    type = number
    default = 0
    description = "The number of participants in the workshop"
}

# GKE
variable "project_id" {
    type = string
    default = "infra-provision"
    description = "The project ID to create the infra from (not to!). This will probably not change."
}

variable "zone" {
  type = string
  default = "us-central1-a"
  description = "The zone to deploy the GKE clusters."
}

variable "cluster_count" {
  type = number
  default = 0
  description = "Kube Clusters to deploy. Set this when you want to deploy the clusters and unset to spin down to save $$$ but keep the projects around."
}
