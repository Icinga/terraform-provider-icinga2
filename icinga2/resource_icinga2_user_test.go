package icinga2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func TestAccCreateBasicUser(t *testing.T) {

	var testAccCreateBasicUser = fmt.Sprintf(`
		resource "icinga2_user" "tf-1" {
		name      = "terraform-user-1"
	}`)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateBasicUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("icinga2_user.tf-1"),
					testAccCheckResourceState("icinga2_user.tf-1", "name", "terraform-user-1"),
				),
			},
		},
	})
}

func TestAccCreateEmailUser(t *testing.T) {

	var testAccCreateBasicUser = fmt.Sprintf(`
		resource "icinga2_user" "tf-2" {
		name      = "terraform-user-2"
		email     = "email@example.com"
	}`)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCreateBasicUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("icinga2_user.tf-2"),
					testAccCheckResourceState("icinga2_user.tf-2", "name", "terraform-user-2"),
					testAccCheckResourceState("icinga2_user.tf-2", "email", "email@example.com"),
				),
			},
		},
	})
}

func testAccCheckUserExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("User resource not found: %s", rn)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*iapi.Server)
		_, err := client.GetUser(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting getting user: %s", err)
		}

		return nil
	}

}
