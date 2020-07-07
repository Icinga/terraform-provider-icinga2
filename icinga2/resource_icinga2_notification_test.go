package icinga2

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func TestAccCreateHostNotification(t *testing.T) {

	var testAccCreateNotification = fmt.Sprintf(`
		resource "icinga2_notification" "tf-notification-1" {
		hostname      = "docker-icinga2"
		command       = "mail-host-notification"
		users         = ["user"]
	}`)
	hostname := "docker-icinga2"
	username := "user"
	createResources := func() {
		icinga2Server := testAccProvider.Meta().(*iapi.Server)
		icinga2Server.CreateHost(hostname, "10.0.0.1", "hostalive", nil, nil, nil)
		icinga2Server.CreateUser(username, "email@example.com")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: createResources,
				Config:    testAccCreateNotification,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists("icinga2_notification.tf-notification-1"),
					testAccCheckResourceState("icinga2_notification.tf-notification-1", "hostname", hostname),
					testAccCheckResourceState("icinga2_notification.tf-notification-1", "command", "mail-host-notification"),
					testAccCheckResourceState("icinga2_notification.tf-notification-1", "users.#", "1"),
					testAccCheckResourceState("icinga2_notification.tf-notification-1", "users.0", "user"),
				),
			},
		},
	})

	icinga2Server := testAccProvider.Meta().(*iapi.Server)
	err := icinga2Server.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error deleting host object after test completed: %s", err)
	}

	err = icinga2Server.DeleteUser(username)
	if err != nil {
		t.Errorf("Error deleting user object after test completed: %s", err)
	}
}

func TestAccCreateServiceNotification(t *testing.T) {

	var testAccCreateNotification = fmt.Sprintf(`
		resource "icinga2_notification" "tf-notification-2" {
		hostname      = "docker-icinga2"
		command       = "mail-service-notification"
		users         = ["user"]
		servicename   = "ping"
	}`)
	hostname := "docker-icinga2"
	username := "user"
	servicename := "ping"
	createResources := func() {
		icinga2Server := testAccProvider.Meta().(*iapi.Server)
		icinga2Server.CreateHost(hostname, "10.0.0.1", "hostalive", nil, nil, nil)
		icinga2Server.CreateUser(username, "email@example.com")
		icinga2Server.CreateService(servicename, hostname, "ping", nil, nil)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: createResources,
				Config:    testAccCreateNotification,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists("icinga2_notification.tf-notification-2"),
					testAccCheckResourceState("icinga2_notification.tf-notification-2", "hostname", hostname),
					testAccCheckResourceState("icinga2_notification.tf-notification-2", "command", "mail-service-notification"),
					testAccCheckResourceState("icinga2_notification.tf-notification-2", "users.#", "1"),
					testAccCheckResourceState("icinga2_notification.tf-notification-2", "users.0", username),
					testAccCheckResourceState("icinga2_notification.tf-notification-2", "servicename", servicename),
				),
			},
		},
	})

	icinga2Server := testAccProvider.Meta().(*iapi.Server)

	err := icinga2Server.DeleteService(servicename, hostname)
	if err != nil {
		t.Errorf("Error deleting service object after test completed: %s", err)
	}

	err = icinga2Server.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error deleting host object after test completed: %s", err)
	}

	err = icinga2Server.DeleteUser(username)
	if err != nil {
		t.Errorf("Error deleting user object after test completed: %s", err)
	}
}

func testAccCheckNotificationExists(rn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("Notification resource not found: %s", rn)
		}

		if resource.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*iapi.Server)

		_, err := client.GetNotification(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting getting Notification: %s", err)
		}

		return nil
	}

}
