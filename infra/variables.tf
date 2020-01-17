# Provider
variable "credz_file" {
    type = string
}

variable "project_id" {
    type = string
    default = "training-infra-owner"
}

variable "region" {
  type = string
  default = "us-central1"
}

variable "zone" {
  type = string
  default = "us-central1-c"
}

# GKE
variable "cluster_prefix" {
  type = string
}

variable "cluster_count" {
    type = number
    default = 1
}
