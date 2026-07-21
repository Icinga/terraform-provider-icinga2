# Configure a new user to be monitored by an Icinga2 Server
resource "icinga2_user" "tf-3" {
  name      = "terraform-user-3"
  vars      = {
    field = "data"
  }
}
