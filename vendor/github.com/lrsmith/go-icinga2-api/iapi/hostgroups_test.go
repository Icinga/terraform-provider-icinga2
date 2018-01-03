package iapi

import "testing"

func TestGetValidHostgroup(t *testing.T) {

	name := "linux-servers"

	_, err := Icinga2_Server.GetHostgroup(name)

	if err != nil {
		t.Error(err)
	}

}

func TestGetInvalidHostgroup(t *testing.T) {

	name := "irix-servers"

	_, err := Icinga2_Server.GetHostgroup(name)
	if err != nil {
		t.Error(err)
	}

}

func TestCreateHostgroup(t *testing.T) {

	name := "docker-servers"
	displayName := "Docker Host Servers"
	_, err := Icinga2_Server.CreateHostgroup(name, displayName)

	if err != nil {
		t.Error(err)
	}

}

// func TestDeleteHostgroup
// Delete Hostgroup created via API. Should succeed
func TestDeleteHostgroup(t *testing.T) {

	name := "docker-servers"

	err := Icinga2_Server.DeleteHostgroup(name)
	if err != nil {
		t.Error(err)
	}
}

// func TestDeleteHostgroupNonAPI
func TestDeleteHostgroupNonAPI(t *testing.T) {

	name := "linux-servers"

	err := Icinga2_Server.DeleteHostgroup(name)
	if err.Error() != "Object cannot be deleted because it was not created using the API." {
		t.Error(err)
	}
}

func TestDeleteHostgroupDNE(t *testing.T) {

	name := "docker-servers"
	err := Icinga2_Server.DeleteHostgroup(name)

	if err.Error() != "No objects found." {
		t.Error(err)
	}
}
