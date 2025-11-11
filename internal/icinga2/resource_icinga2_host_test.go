package icinga2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCreateBasicHost(t *testing.T) {
	var testAccCreateBasicHost = fmt.Sprintf(`
resource "icinga2_host" "tf-1" {
  hostname      = "terraform-host-1"
  address       = "10.10.10.1"
  check_command = "hostalive"
}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCreateBasicHost,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-1",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("terraform-host-1"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-1",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.10.10.1"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-1",
						tfjsonpath.New("check_command"),
						knownvalue.StringExact("hostalive"),
					),
				},
			},
		},
	})
}

func TestAccCreateHost(t *testing.T) {

	var testAccCreateBasicHost = fmt.Sprintf(`
resource "icinga2_host" "tf-2" {
  hostname      = "terraform-host-2"
  address       = "10.10.10.2"
  check_command = "hostalive"
  groups        = ["linux-servers"]
}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCreateBasicHost,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-2",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("terraform-host-2"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-2",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.10.10.2"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-2",
						tfjsonpath.New("check_command"),
						knownvalue.StringExact("hostalive"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-2",
						tfjsonpath.New("groups.#"),
						knownvalue.StringExact("1"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-2",
						tfjsonpath.New("groups.0"),
						knownvalue.StringExact("linux-servers"),
					),
				},
			},
		},
	})
}

func TestAccCreateVariableHost(t *testing.T) {

	var testAccCreateVariableHost = fmt.Sprintf(`
resource "icinga2_host" "tf-3" {
  hostname      = "terraform-host-3"
  address       = "10.10.10.3"
  check_command = "hostalive"
  vars          = {
    os        = "linux"
    osver     = "1"
    allowance = "none"
  }
}`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCreateVariableHost,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-3",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("terraform-host-3"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-3",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.10.10.3"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-3",
						tfjsonpath.New("check_command"),
						knownvalue.StringExact("hostalive"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-3",
						tfjsonpath.New("vars.%"),
						knownvalue.StringExact("3"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-3",
						tfjsonpath.New("vars.allowance"),
						knownvalue.StringExact("none"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-3",
						tfjsonpath.New("vars.os"),
						knownvalue.StringExact("linux"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-3",
						tfjsonpath.New("vars.osver"),
						knownvalue.StringExact("1"),
					),
				},
			},
		},
	})
}

func TestAccCreateTemplateHost(t *testing.T) {
	var testAccCreateTemplateHost = `
resource "icinga2_host" "tf-4" {
  hostname      = "terraform-host-4"
  address       = "10.10.10.4"
  check_command = "hostalive"
  templates     = ["generic", "az1"]
}`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCreateTemplateHost,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-4",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("terraform-host-4"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-4",
						tfjsonpath.New("address"),
						knownvalue.StringExact("10.10.10.4"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-4",
						tfjsonpath.New("check_command"),
						knownvalue.StringExact("hostalive"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-4",
						tfjsonpath.New("templates.#"),
						knownvalue.StringExact("2"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-4",
						tfjsonpath.New("templates.0"),
						knownvalue.StringExact("generic"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_host.tf-4",
						tfjsonpath.New("templates.1"),
						knownvalue.StringExact("az1"),
					),
				},
			},
		},
	})
}
