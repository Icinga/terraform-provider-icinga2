package icinga2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func TestAccCreateBasicHostGroup(t *testing.T) {
	var (
		hostgroup iapi.HostgroupStruct

		hostgroupName     = "terraform-hostgroup-1"
		firstDisplayName  = "Terraform Test HostGroup"
		secondDisplayName = "Some New HostGroup DisplayName"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a basic HostGroup
			{
				Config: testAccCreateHostGroupBasic(hostgroupName, firstDisplayName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the HostGroup object
					testAccCheckHostgroupExists("icinga2_hostgroup.tf-hg-1", &hostgroup),
					// verify remote values
					testAccCheckHostgroupValues(&hostgroup, firstDisplayName),
					// verify local values
					resource.TestCheckResourceAttr("icinga2_hostgroup.tf-hg-1", "name", hostgroupName),
					resource.TestCheckResourceAttr("icinga2_hostgroup.tf-hg-1", "display_name", firstDisplayName),
				),
			},
			// Update recently created HostGroup
			{
				Config: testAccCreateHostGroupBasic(hostgroupName, secondDisplayName),
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the HostGroup object
					testAccCheckHostgroupExists("icinga2_hostgroup.tf-hg-1", &hostgroup),
					// verify remote values
					testAccCheckHostgroupValues(&hostgroup, secondDisplayName),
					// verify local values
					resource.TestCheckResourceAttr("icinga2_hostgroup.tf-hg-1", "name", hostgroupName),
					resource.TestCheckResourceAttr("icinga2_hostgroup.tf-hg-1", "display_name", secondDisplayName),
				),
			},
		},
	})
}

func testAccCheckHostgroupExists(resourceName string, hg *iapi.HostgroupStruct) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Hostgroup resource not found: %s", resourceName)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("Hostgroup resource id not set")
		}

		client := testAccProvider.Meta().(*iapi.Server)
		storedHostgroup, err := client.GetHostgroup(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting getting hostgroup: %s", err)
		}

		*hg = storedHostgroup[0]
		return nil
	}
}

func testAccCheckHostgroupValues(hg *iapi.HostgroupStruct, displayName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if hg.Attrs.DisplayName != displayName {
			return fmt.Errorf("expected displayName to be set, got not set")
		}
		return nil
	}
}

func testAccCreateHostGroupBasic(name, displayName string) string {
	return fmt.Sprintf(`
resource "icinga2_hostgroup" "tf-hg-1" {
	name = "%s"
	display_name = "%s"
}
`, name, displayName)
}
