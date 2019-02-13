resource "google_container_cluster" "cluster" {
  name = "${var.name}"

  # NB: if TF complains that these "mandatory" fields aren't present,
  # they're not mandatory, they just need setting in the Provider
  # region = from provider

  project = "${var.project}"
  zone    = "${var.zone}"

  # Versions

  min_master_version = "${var.k8s_version}"
  node_version       = "${var.k8s_version}"

  # Disable both basic auth and x509 auth (GCP service account external
  # auth provider only)

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }

  # The top-level resource support *initial*_node_count, but doesn't offer
  # all of the config options for full node-pools, e.g. the ability to
  # enable autoscaling. google_container_node_pool (external, or in-line
  # like here) only supports node_count, so we have to ignore changes
  # otherwise we try to undo the work of the autoscaler.

  lifecycle {
    ignore_changes = ["node_pool.0.node_count"]
  }
  
  node_pool = [
    {
      name       = "ondemand"
      node_count = 2

      # Enables cluster autoscaling by implication
      autoscaling {
        min_node_count = 1
        max_node_count = 5
      }

      node_config {
        disk_size_gb = 100
        machine_type = "n1-highmem-2"
        preemptible  = false

        oauth_scopes = [
          "https://www.googleapis.com/auth/compute",
          "https://www.googleapis.com/auth/devstorage.read_only",
          "https://www.googleapis.com/auth/logging.write",
          "https://www.googleapis.com/auth/monitoring",
        ]
      }
    },
  ]
  addons_config {
    horizontal_pod_autoscaling {
      disabled = false
    }

    # Disable GCE ingress controller
    http_load_balancing {
      disabled = true
    }

    kubernetes_dashboard {
      disabled = false
    }
  }
}
