provider "google" {
  credentials = file(var.gcp_credentials_path)
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

module "gcp_create_app_instance-wrk2" {
  source        = "./modules/gcp_create_wrk2_instance"
  instance_name = "trainticket-wrk2"
  hostname      = "trainticket-wrk2.prod"
  zone          = "europe-west6-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}