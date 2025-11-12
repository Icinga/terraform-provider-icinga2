# Configure a new hostgroup to be monitored by an Icinga2 Server
provider "icinga2" {
  api_url = "https://192.168.33.5:5665/v1"
}

resource "icinga2_hostgroup" "my-hostgroup" {
  name         = "terraform-hostgroup-1"
  display_name = "Terraform Test HostGroup"
}