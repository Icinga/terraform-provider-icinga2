package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCreateBasicDowntime(t *testing.T) {
	hostname := "docker-icinga2"
	Groups := []string{"linux-servers"}
	createHost := func() {
		client, _ := testAccClient()
		_, err := client.CreateHost(hostname, "10.0.0.2", "", "hostalive", nil, nil, Groups, "")
		if err != nil {
			t.Errorf("Error creating host before test: %s", err)
		}
	}

	testAccCreateDowntimeBasic := fmt.Sprintf(`
resource "icinga2_downtime" "tf-downtime-1" {
  type         = "Host"
  filter       = "host.name==\"docker-icinga2\""
  author       = "terraform"
  comment      = "Initial downtime"
  start_time   = %d
  end_time     = %d
  all_services = false
}`, time.Now().Unix(), time.Now().Unix()+3600)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: createHost,
				Config:    providerConfig + testAccCreateDowntimeBasic,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"icinga2_downtime.tf-downtime-1",
						tfjsonpath.New("type"),
						knownvalue.StringExact("Host"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_downtime.tf-downtime-1",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("host.name==\"docker-icinga2\""),
					),
					statecheck.ExpectKnownValue(
						"icinga2_downtime.tf-downtime-1",
						tfjsonpath.New("author"),
						knownvalue.StringExact("terraform"),
					),
					statecheck.ExpectKnownValue(
						"icinga2_downtime.tf-downtime-1",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("Initial downtime"),
					),
				},
			},
		},
	})

	client, _ := testAccClient()
	err := client.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error deleting host object after test completed: %s", err)
	}
}
