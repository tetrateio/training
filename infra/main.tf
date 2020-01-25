provider "google" {
  credentials = file(var.credz_file)

  project = var.project_id
  region  = var.region
}

resource "google_folder" "training" {
  display_name = var.workshop_name
  parent       = format("organizations/%s", var.organization_id)
}

resource "google_project" "training" {
  name       = format("%s-%03d", var.workshop_name, count.index)
  count = var.participant_count

  project_id = format("%s-%03d", var.workshop_name, count.index)
  folder_id  = google_folder.training.name
  billing_account = var.billing_account
}

resource "google_project_service" "container" {
  project = google_project.training[count.index].project_id
  count = var.participant_count

  service = "container.googleapis.com"
  disable_dependent_services = true
}

resource "google_service_account" "compute" {
  account_id = "compute"
  count = var.participant_count

  project = format("%s-%03d", var.workshop_name, count.index)
}

resource "google_project_iam_member" "project" {
  project = format("%s-%03d", var.workshop_name, count.index)
  count = var.participant_count

  role    = "roles/editor"
  member  = format("serviceAccount:%s", google_service_account.compute[count.index].email)
}

resource "google_container_cluster" "k8s" {
  name     = format("%s-%03d", var.workshop_name, count.index)
  count    = min(var.cluster_count, var.participant_count)
  
  project = google_project.training[count.index].project_id
  location = var.zone

  logging_service = "none"
  monitoring_service = "none"

  initial_node_count       = 3
  node_config {
    machine_type = "n1-standard-2"
    service_account = google_service_account.compute[count.index].email

    metadata = {
      disable-legacy-endpoints = "true"
    }
  }

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }

  depends_on = [
    google_project_service.container,
  ]
}
