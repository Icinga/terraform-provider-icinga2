# Configure a new host to be monitored by an Icinga2 Server
terraform {
  required_providers {
    icinga2 = {
      source = "Icinga/icinga2"
    }
  }
}

resource "icinga2_host" "tf-1" {
  hostname      = "terraform-host-1"
  address       = "10.10.10.1"
  check_command = "hostalive"
}
