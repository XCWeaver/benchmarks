provider "google" {
  credentials = file(var.gcp_credentials_path)
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

module "gcp_create_storage_instance_manager" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-dsb-db-manager"
  hostname      = "xcweaver-dsb-db-manager.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_storage_instance_eu" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-dsb-db-eu"
  hostname      = "xcweaver-dsb-db-eu.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_storage_instance-us" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-dsb-db-us"
  hostname      = "xcweaver-dsb-db-us.prod"
  zone          = "us-central1-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_app_instance-eu" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-dsb-app-eu"
  hostname      = "xcweaver-dsb-app-eu.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_app_instance-us" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-dsb-app-us"
  hostname      = "xcweaver-dsb-app-us.prod"
  zone          = "us-central1-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_app_instance-wrk2" {
  source        = "./modules/gcp_create_wrk2_instance"
  instance_name = "xcweaver-dsb-app-wrk2"
  hostname      = "xcweaver-dsb-app-wrk2.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}
