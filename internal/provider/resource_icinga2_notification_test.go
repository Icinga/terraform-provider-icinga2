package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCreateHostNotification(t *testing.T) {

	var testAccCreateNotification = `
		resource "icinga2_notification" "tf-notification-1" {
			hostname      = "docker-icinga2"
			command       = "mail-host-notification"
			users         = ["user"]
		}`
	hostname := "docker-icinga2"
	username := "user"
	createResources := func() {
		client, _ := testAccClient()
		_, errH := client.CreateHost(hostname, "10.0.0.1", "hostalive", nil, nil, nil)
		if errH != nil {
			t.Errorf("Error creating host object before test start: %s", errH)
		}
		_, errU := client.CreateUser(username, "email@example.com")
		if errU != nil {
			t.Errorf("Error creating user object before test start: %s", errU)
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: createResources,
				Config:    testAccCreateNotification,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists("icinga2_notification.tf-notification-1"),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-1", "hostname", hostname),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-1", "command", "mail-host-notification"),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-1", "users.#", "1"),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-1", "users.0", "user"),
				),
			},
		},
	})

	client, _ := testAccClient()
	err := client.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error deleting host object after test completed: %s", err)
	}

	err = client.DeleteUser(username)
	if err != nil {
		t.Errorf("Error deleting user object after test completed: %s", err)
	}
}

func TestAccCreateServiceNotification(t *testing.T) {

	var testAccCreateNotification = `
		resource "icinga2_notification" "tf-notification-2" {
			hostname      = "docker-icinga2"
			command       = "mail-service-notification"
			users         = ["user"]
			servicename   = "ping"
		}`
	hostname := "docker-icinga2"
	username := "user"
	servicename := "ping"
	createResources := func() {
		client, _ := testAccClient()
		_, errH := client.CreateHost(hostname, "10.0.0.1", "hostalive", nil, nil, nil)
		if errH != nil {
			t.Errorf("Error creating host object before test start: %s", errH)
		}
		_, errU := client.CreateUser(username, "email@example.com")
		if errU != nil {
			t.Errorf("Error creating user object before test start: %s", errU)
		}
		_, errS := client.CreateService(servicename, hostname, "ping", nil, nil)
		if errS != nil {
			t.Errorf("Error creating user object before test start: %s", errS)
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: createResources,
				Config:    testAccCreateNotification,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists("icinga2_notification.tf-notification-2"),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-2", "hostname", hostname),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-2", "command", "mail-service-notification"),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-2", "users.#", "1"),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-2", "users.0", username),
					resource.TestCheckResourceAttr("icinga2_notification.tf-notification-2", "servicename", servicename),
				),
			},
		},
	})

	client, _ := testAccClient()

	err := client.DeleteService(servicename, hostname)
	if err != nil {
		t.Errorf("Error deleting service object after test completed: %s", err)
	}

	err = client.DeleteHost(hostname)
	if err != nil {
		t.Errorf("Error deleting host object after test completed: %s", err)
	}

	err = client.DeleteUser(username)
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

		client, err := testAccClient()
		if err != nil {
			return err
		}

		_, err = client.GetNotification(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting getting Notification: %s", err)
		}

		return nil
	}

}
