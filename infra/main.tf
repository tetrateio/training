provider "google" {
  credentials = file(var.credz_file)

  project = var.project_id
  region  = var.region
  zone    = var.zone
}

resource "google_container_cluster" "k8s" {
  name     = format("%s-%03d", var.cluster_prefix, count.index)
  location = var.zone
  count = var.cluster_count

  # Manage node pools in separate resource
  remove_default_node_pool = true
  initial_node_count       = 1

  logging_service = "none"
  monitoring_service = "none"

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }
}

resource "google_container_node_pool" "k8s_nodes" {
  name       = "nodes"
  location   = var.zone
  count = var.cluster_count

  cluster    = google_container_cluster.k8s[count.index].name
  node_count = 4

  node_config {
    machine_type = "n1-standard-2"

    metadata = {
      disable-legacy-endpoints = "true"
    }
  }
}
