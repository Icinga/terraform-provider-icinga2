terraform {
  required_providers {
    icinga2 = {
      source = "Icinga/icinga2"
    }
  }
}

provider "icinga2" {
  api_url = "https://127.0.0.1:5665/v1"
  api_user = "icingaweb"
  api_password = "icingaweb"
}