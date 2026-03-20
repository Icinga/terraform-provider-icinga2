# Configure a new service to be monitored by an Icinga2 Server
terraform {
  required_providers {
    icinga2 = {
      source = "Icinga/icinga2"
    }
  }
}

resource "icinga2_service" "tf-service-1" {
  hostname      = "terraform-host-1"
  name          = "ssh3"
  check_command = "ssh"
}
