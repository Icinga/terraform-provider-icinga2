x# Configure the Icinga2 provider
provider "icinga2" {
  api_url                  = "https://192.168.33.5:5665/v1"
  api_user                 = "root"
  api_password             = "icinga"
  insecure_skip_tls_verify = true
}

# Configure a host
resource "icinga2_host" "web-server" {
  # ...
}