package iapi

import (
	"strings"
	"testing"
)

func TestGetValidService(t *testing.T) {

	hostname := "c1-mysql-1"
	servicename := "ssh"

	_, err := Icinga2_Server.GetService(servicename, hostname)

	if err != nil {
		t.Error(err)
	}
}

func TestGetInvalidService(t *testing.T) {

	hostname := "c1-mysql-1"
	servicename := "foo"

	_, err := Icinga2_Server.GetService(servicename, hostname)

	if err != nil {
		t.Error(err)
	}
}

// func TestCreateServiceHostDNE
// Try and create a service for a host that does not exist.
// Should fail with an error about the host not existing.
func TestCreateServiceHostDNE(t *testing.T) {

	hostname := "c1-host-dne-1"
	servicename := "ssh"
	check_command := "ssh"

	_, err := Icinga2_Server.CreateService(servicename, hostname, check_command)

	if !strings.Contains(err.Error(), "type 'Host' does not exist.") {
		t.Error(err)
	}

}

// func TestCreateHostAndService
// Create a host and service via the API
func TestCreateHostAndService(t *testing.T) {

	hostname := "c1-test-1"
	servicename := "ssh"
	check_command := "ssh"
	Group := []string{"linux-servers"}

	_, _ = Icinga2_Server.CreateHost(hostname, "127.0.0.1", "hostalive", nil, nil, Group)

	_, err := Icinga2_Server.CreateService(servicename, hostname, check_command)

	if err != nil {
		t.Errorf("Error : Failed to create service %s!%s : %s", hostname, servicename, err)
	}

}

// func TestCreateServiceAlreadyExists
// Test creating a host/service pair that already exists. Should get error about it already existing.
func TestCreateServiceAlreadyExists(t *testing.T) {

	hostname := "c1-test-1"
	servicename := "ssh"
	check_command := "ssh"

	_, err := Icinga2_Server.CreateService(servicename, hostname, check_command)

	if !strings.HasSuffix(err.Error(), " already exists.") {
		t.Error(err)
	}

}

// func TestDeleteHostAndService
// Delete a service which was create via the API. NOTE : Host also is created via the API in previous test.
// Should not get an error
func TestDeleteHostAndService(t *testing.T) {

	hostname := "c1-test-1"
	servicename := "ssh"

	err := Icinga2_Server.DeleteService(servicename, hostname)
	if err != nil {
		_ = Icinga2_Server.DeleteHost(hostname)
		t.Error(err)
	}

	_ = Icinga2_Server.DeleteHost(hostname)
}

// func TestDeleteServiceHostDNE
// Try and delet a service, where the host does not exists.
// Should get an error abot no object found
func TestDeleteServiceHostDNE(t *testing.T) {

	hostname := "c1-test-1"
	servicename := "ssh"

	err := Icinga2_Server.DeleteService(servicename, hostname)
	if err.Error() != "No objects found." {
		t.Error(err)
	}
}

// func TestDeleteServiceDNS
// Try and delete a service, where the host exists but the service does not.
// Should get an error abot no object found
func TestDeleteServiceDNE(t *testing.T) {

	hostname := "c1-mysql-1"
	servicename := "foo"

	err := Icinga2_Server.DeleteService(servicename, hostname)
	if err.Error() != "No objects found." {
		t.Error(err)
	}
}

// func TestDeleteServiceNonAPI
// Services that were not created via the API, cannot be deleted via the API
// Should get an error about not being created via the API
func TestDeleteServiceNonAPI(t *testing.T) {

	hostname := "docker-icinga2"
	servicename := "random-001"

	err := Icinga2_Server.DeleteService(servicename, hostname)
	if err.Error() != "No objects found." {
		t.Error(err)
	}
}
