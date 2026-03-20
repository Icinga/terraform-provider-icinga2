package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCreateCheckcommand(t *testing.T) {
	var testAccCreateCheckcommand = `
		resource "icinga2_checkcommand" "checkcommand" {
			name      = "terraform-test-checkcommand-1"
			templates = []
			command   = "/usr/local/bin/check_command"
			arguments = {
				"-I" = "$IARG$"
				"-J" = "$JARG$" }
		}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateCheckcommand,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCheckcommandExists("icinga2_checkcommand.checkcommand"),
					resource.TestCheckResourceAttr("icinga2_checkcommand.checkcommand", "name", "terraform-test-checkcommand-1"),
					resource.TestCheckResourceAttr("icinga2_checkcommand.checkcommand", "command", "/usr/local/bin/check_command"),
					resource.TestCheckResourceAttr("icinga2_checkcommand.checkcommand", "arguments.%", "2"),
					resource.TestCheckResourceAttr("icinga2_checkcommand.checkcommand", "arguments.-I", "$IARG$"),
					resource.TestCheckResourceAttr("icinga2_checkcommand.checkcommand", "arguments.-J", "$JARG$"),
				),
			},
		},
	})
}

func testAccCheckCheckcommandExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Checkcommand resource not found: %s", rn)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("Checkcommand resource id not set")
		}

		client, err := testAccClient()
		if err != nil {
			return err
		}

		_, err = client.GetCheckcommand(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting getting Checkcommand: %s", err)
		}

		return nil
	}
}
