package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCreateBasicUser(t *testing.T) {

	var testAccCreateBasicUser = `
	resource "icinga2_user" "tf-1" {
		name      = "terraform-user-1"
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateBasicUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("icinga2_user.tf-1"),
					resource.TestCheckResourceAttr("icinga2_user.tf-1", "name", "terraform-user-1"),
				),
			},
			{
				ImportState:             true,
				ResourceName:            "icinga2_user.tf-1",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"last_updated"},
			},
		},
	})
}

func TestAccCreateEmailUser(t *testing.T) {

	var testAccCreateBasicUser = `
	resource "icinga2_user" "tf-2" {
		name      = "terraform-user-2"
		email     = "email@example.com"
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateBasicUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("icinga2_user.tf-2"),
					resource.TestCheckResourceAttr("icinga2_user.tf-2", "name", "terraform-user-2"),
					resource.TestCheckResourceAttr("icinga2_user.tf-2", "email", "email@example.com"),
				),
			},
		},
	})
}

func TestAccCreateVarsUser(t *testing.T) {

	var testAccCreateBasicUser = `
	resource "icinga2_user" "tf-3" {
		name      = "terraform-user-3"
        vars      = {
          field = "data"
        }
	}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateBasicUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("icinga2_user.tf-3"),
					resource.TestCheckResourceAttr("icinga2_user.tf-3", "name", "terraform-user-3"),
					resource.TestCheckResourceAttr("icinga2_user.tf-3", "vars.field", "data"),
				),
			},
		},
	})
}

func testAccCheckUserExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		userResource, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("User userResource not found: %s", rn)
		}

		if userResource.Primary.ID == "" {
			return fmt.Errorf("userResource id not set")
		}

		client, err := testAccClient()
		if err != nil {
			return err
		}

		_, err = client.GetUser(context.Background(), userResource.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting getting user: %s", err)
		}

		return nil
	}

}
