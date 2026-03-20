# Configure a new checkcommand to be monitored by an Icinga2 Server
terraform {
  required_providers {
    icinga2 = {
      source = "Icinga/icinga2"
    }
  }
}

resource "icinga2_checkcommand" "checkcommand" {
  name      = "terraform-test-checkcommand-1"
  templates = []
  command   = "/usr/local/bin/check_command"
  arguments = {
    "-I" = "$IARG$",
    "-J" = "$JARG$"
  }
}