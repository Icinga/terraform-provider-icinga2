# Configure a new user to be monitored by an Icinga2 Server
resource "icinga2_user" "tf-1" {
  name = "terraform-user-1"
}
