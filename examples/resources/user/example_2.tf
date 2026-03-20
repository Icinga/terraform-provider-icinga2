# Configure a new user to be monitored by an Icinga2 Server
resource "icinga2_user" "tf-2" {
  name      = "terraform-user-2"
  email     = "email@example.com"
}
