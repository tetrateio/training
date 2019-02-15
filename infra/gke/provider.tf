provider "google" {
  version = "v1.15.0"

  project = "${var.project_id}"
  zone    = "${var.zone}"
}
