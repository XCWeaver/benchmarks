provider "google" {
  credentials = file(var.gcp_credentials_path)
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

module "gcp_create_storage_instance_manager" {
  source        = "./modules/gcp_create_instance"
  instance_name = "weaver-pn-db-manager"
  hostname      = "weaver-pn-db-manager.prod"
  zone          = "europe-west6-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_storage_instance_eu" {
  source        = "./modules/gcp_create_instance"
  instance_name = "weaver-pn-db-eu"
  hostname      = "weaver-pn-db-eu.prod"
  zone          = "europe-west6-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_storage_instance-us" {
  source        = "./modules/gcp_create_instance"
  instance_name = "weaver-pn-db-us"
  hostname      = "weaver-pn-db-us.prod"
  zone          = "us-central1-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_app_instance-wrk2" {
  source        = "./modules/gcp_create_wrk2_instance"
  instance_name = "weaver-pn-app-wrk2"
  hostname      = "weaver-pn-app-wrk2.prod"
  zone          = "europe-west6-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}