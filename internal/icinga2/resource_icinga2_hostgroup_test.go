package icinga2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCreateBasicHostGroup(t *testing.T) {
	var (
		hostgroupName     = "terraform-hostgroup-1"
		firstDisplayName  = "Terraform Test HostGroup"
		secondDisplayName = "Some New HostGroup DisplayName"
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCreateHostGroupBasic(hostgroupName, firstDisplayName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"icinga2_hostgroup.tf-hg-1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(hostgroupName),
					),
					statecheck.ExpectKnownValue(
						"icinga2_hostgroup.tf-hg-1",
						tfjsonpath.New("display_name"),
						knownvalue.StringExact(firstDisplayName),
					),
				},
			},
			// Update recently created HostGroup
			{
				Config: testAccCreateHostGroupBasic(hostgroupName, secondDisplayName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"icinga2_hostgroup.tf-hg-1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(hostgroupName),
					),
					statecheck.ExpectKnownValue(
						"icinga2_hostgroup.tf-hg-1",
						tfjsonpath.New("display_name"),
						knownvalue.StringExact(secondDisplayName),
					),
				},
			},
		},
	})
}

func testAccCreateHostGroupBasic(name, displayName string) string {
	return fmt.Sprintf(`
resource "icinga2_hostgroup" "tf-hg-1" {
	name = "%s"
	display_name = "%s"
}
`, name, displayName)
}
