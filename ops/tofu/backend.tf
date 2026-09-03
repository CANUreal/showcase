terraform {
  backend "s3" {
    bucket  = "tofu-state"
    key     = "showcase/terraform.tfstate"
  

    endpoints = {
      s3 = "http://s3-api.garage.local"
    }

    region = "garage"

    skip_credentials_validation = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_s3_checksum            = true
    use_path_style              = true
  }
}
