# Configure a new host to be monitored by an Icinga2 Server
provider "icinga2" {
  api_url = "https://192.168.33.5:5665/v1"
}

resource "icinga2_host" "host" {
  hostname      = "terraform-host-1"
  address       = "10.10.10.1"
  groups        = ["example-hostgroup"]
  check_command = "hostalive"
  templates     = ["bp-host-web"]

  vars = {
    os        = "linux"
    osver     = "1"
    allowance = "none"
  }
}