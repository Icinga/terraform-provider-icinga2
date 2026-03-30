package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCreateService(t *testing.T) {

	var testAccCreateService = `
		resource "icinga2_service" "tf-service-1" {
			hostname      = "docker-icinga2"
			name          = "ssh3"
			check_command = "ssh"
		}`
	hostname := "docker-icinga2"
	Groups := []string{"linux-servers"}
	createHost := func() {
		client, _ := testAccClient()
		_, err := client.CreateHost(hostname, "10.0.0.1", "", "hostalive", nil, nil, Groups, "")
		if err != nil {
			t.Errorf("Error creating host object before test started: %s", err)
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: createHost,
				Config:    testAccCreateService,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceExists("icinga2_service.tf-service-1"),
					resource.TestCheckResourceAttr("icinga2_service.tf-service-1", "hostname", hostname),
					resource.TestCheckResourceAttr("icinga2_service.tf-service-1", "name", "ssh3"),
					resource.TestCheckResourceAttr("icinga2_service.tf-service-1", "check_command", "ssh"),
				),
			},
		},
	})

	client, _ := testAccClient()
	err := client.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error deleting host object after test completed: %s", err)
	}
}

func testAccCheckServiceExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Service resource not found: %s", rn)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client, err := testAccClient()
		if err != nil {
			return err
		}
		tokens := strings.Split(resource.Primary.ID, "!")

		_, err = client.GetService(tokens[1], tokens[0])
		if err != nil {
			return fmt.Errorf("error getting getting Service: %s", err)
		}

		return nil
	}

}
