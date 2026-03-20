# Configure a new notification to be monitored by an Icinga2 Server
terraform {
  required_providers {
    icinga2 = {
      source = "Icinga/icinga2"
    }
  }
}

resource "icinga2_user" "tf-1" {
  name = "terraform-user-1"
}

resource "icinga2_notification" "tf-notification-1" {
  hostname   = "terraform-host-1"
  command    = "mail-host-notification"
  users      = ["terraform-user-1"]
  depends_on = [icinga2_user.tf-1]
}