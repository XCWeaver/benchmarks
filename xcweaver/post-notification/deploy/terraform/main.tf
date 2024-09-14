provider "google" {
  credentials = file(var.gcp_credentials_path)
  project = var.gcp_project
  region  = var.gcp_region
  zone    = var.gcp_zone
}

module "gcp_create_storage_instance_manager" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-pn-db-manager"
  hostname      = "xcweaver-pn-db-manager.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_storage_instance_eu" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-pn-db-eu"
  hostname      = "xcweaver-pn-db-eu.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_storage_instance-us" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-pn-db-us"
  hostname      = "xcweaver-pn-db-us.prod"
  zone          = "us-central1-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_app_instance-eu" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-pn-app-eu"
  hostname      = "xcweaver-pn-app-eu.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_app_instance-us" {
  source        = "./modules/gcp_create_instance"
  instance_name = "xcweaver-pn-app-us"
  hostname      = "xcweaver-pn-app-us.prod"
  zone          = "us-central1-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}

module "gcp_create_app_instance-wrk2" {
  source        = "./modules/gcp_create_wrk2_instance"
  instance_name = "xcweaver-pn-app-wrk2"
  hostname      = "xcweaver-pn-app-wrk2.prod"
  zone          = "europe-west3-a"
  image         = "debian-cloud/debian-11"
  providers = {
    google = google
  }
}
