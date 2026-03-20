# Configure a new hostgroup to be monitored by an Icinga2 Server
terraform {
  required_providers {
    icinga2 = {
      source = "Icinga/icinga2"
    }
  }
}

resource "icinga2_hostgroup" "my-hostgroup" {
  name         = "terraform-hostgroup-1"
  display_name = "Terraform Test HostGroup"
}
