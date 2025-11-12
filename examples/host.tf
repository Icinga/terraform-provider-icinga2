# Configure a new host to be monitored by an Icinga2 Server
terraform {
  required_providers {
    icinga2 = {
      source = "Icinga/icinga2"
    }
  }
}

provider "icinga2" {}

resource "icinga2_host" "host" {
  hostname      = "terraform-host-1"
  address       = "10.10.10.1"
  groups        = ["linux-servers"]
  check_command = "hostalive"
  templates     = ["bp-host-web"]

  vars = {
    os        = "linux"
    osver     = "1"
    allowance = "none"
  }
}