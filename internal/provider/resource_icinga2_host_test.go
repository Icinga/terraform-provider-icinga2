package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCreateBasicHost(t *testing.T) {

	var testAccCreateBasicHost = `
		resource "icinga2_host" "tf-1" {
			hostname      = "terraform-host-1"
			address       = "10.10.10.1"
			check_command = "hostalive"
		}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateBasicHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHostExists("icinga2_host.tf-1"),
					resource.TestCheckResourceAttr("icinga2_host.tf-1", "hostname", "terraform-host-1"),
					resource.TestCheckResourceAttr("icinga2_host.tf-1", "address", "10.10.10.1"),
					resource.TestCheckResourceAttr("icinga2_host.tf-1", "check_command", "hostalive"),
				),
			},
			{
				ImportState:             true,
				ResourceName:            "icinga2_host.tf-1",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func TestAccCreateGroupHost(t *testing.T) {

	var testAccCreateBasicHost = `
		resource "icinga2_host" "tf-2" {
			hostname      = "terraform-host-2"
			address       = "10.10.10.2"
			check_command = "hostalive"
			groups = ["linux-servers"]
        }`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateBasicHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHostExists("icinga2_host.tf-2"),
					resource.TestCheckResourceAttr("icinga2_host.tf-2", "hostname", "terraform-host-2"),
					resource.TestCheckResourceAttr("icinga2_host.tf-2", "address", "10.10.10.2"),
					resource.TestCheckResourceAttr("icinga2_host.tf-2", "check_command", "hostalive"),
					resource.TestCheckResourceAttr("icinga2_host.tf-2", "groups.#", "1"),
					resource.TestCheckResourceAttr("icinga2_host.tf-2", "groups.0", "linux-servers"),
				),
			},
		},
	})
}

func TestAccCreateVariableHost(t *testing.T) {

	var testAccCreateVariableHost = `
		resource "icinga2_host" "tf-3" {
			hostname = "terraform-host-3"
			address = "10.10.10.3"
			check_command = "hostalive"
			vars = {
			  os = "linux"
			  osver = "1"
			  allowance = "none"
	        }
		}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateVariableHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHostExists("icinga2_host.tf-3"),
					resource.TestCheckResourceAttr("icinga2_host.tf-3", "hostname", "terraform-host-3"),
					resource.TestCheckResourceAttr("icinga2_host.tf-3", "address", "10.10.10.3"),
					resource.TestCheckResourceAttr("icinga2_host.tf-3", "check_command", "hostalive"),
					resource.TestCheckResourceAttr("icinga2_host.tf-3", "vars.%", "3"),
					resource.TestCheckResourceAttr("icinga2_host.tf-3", "vars.allowance", "none"),
					resource.TestCheckResourceAttr("icinga2_host.tf-3", "vars.os", "linux"),
					resource.TestCheckResourceAttr("icinga2_host.tf-3", "vars.osver", "1"),
				),
			},
		},
	})
}

func TestAccCreateTemplateHost(t *testing.T) {
	var testAccCreateTemplateHost = `resource "icinga2_host" "tf-4" {
		hostname = "terraform-host-4"
		address = "10.10.10.4"
		check_command = "hostalive"
		templates = ["generic", "az1"]
	}`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateTemplateHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHostExists("icinga2_host.tf-4"),
					resource.TestCheckResourceAttr("icinga2_host.tf-4", "hostname", "terraform-host-4"),
					resource.TestCheckResourceAttr("icinga2_host.tf-4", "address", "10.10.10.4"),
					resource.TestCheckResourceAttr("icinga2_host.tf-4", "check_command", "hostalive"),
					resource.TestCheckResourceAttr("icinga2_host.tf-4", "templates.#", "2"),
					resource.TestCheckResourceAttr("icinga2_host.tf-4", "templates.0", "generic"),
					resource.TestCheckResourceAttr("icinga2_host.tf-4", "templates.1", "az1"),
				),
			},
		},
	})
}

func TestAccCreateZoneHost(t *testing.T) {

	var testAccCreateZoneHost = `resource "icinga2_host" "tf-5" {
		hostname = "terraform-host-5"
		address = "10.10.10.5"
		check_command = "hostalive"
		zone = "master"
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateZoneHost,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHostExists("icinga2_host.tf-5"),
					resource.TestCheckResourceAttr("icinga2_host.tf-5", "hostname", "terraform-host-5"),
					resource.TestCheckResourceAttr("icinga2_host.tf-5", "address", "10.10.10.5"),
					resource.TestCheckResourceAttr("icinga2_host.tf-5", "check_command", "hostalive"),
					resource.TestCheckResourceAttr("icinga2_host.tf-5", "zone", "master"),
				),
			},
		},
	})
}

func testAccCheckHostExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Host resource not found: %s", rn)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client, err := testAccClient()
		if err != nil {
			return err
		}

		_, err = client.GetHost(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting getting host: %s", err)
		}

		return nil
	}

}
