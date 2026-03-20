package provider

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

const (
	providerConfig = `
provider "icinga2" {}
`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"icinga2": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	v := os.Getenv("ICINGA2_API_URL")
	if v == "" {
		t.Fatal("ICINGA2_API_URL must be set for acceptance tests")
	}

	v = os.Getenv("ICINGA2_API_USER")
	if v == "" {
		t.Fatal("ICINGA2_API_USER must be set for acceptance tests")
	}

	v = os.Getenv("ICINGA2_API_PASSWORD")
	if v == "" {
		t.Fatal("ICINGA2_API_PASSWORD must be set for acceptance tests")
	}
}

func testAccClient() (*iapi.Server, error) {
	api_url := os.Getenv("ICINGA2_API_URL")
	api_user := os.Getenv("ICINGA2_API_USER")
	api_password := os.Getenv("ICINGA2_API_PASSWORD")
	tlsVerify, _ := strconv.ParseBool(os.Getenv("ICINGA2_INSECURE_SKIP_TLS_VERIFY"))

	return iapi.New(
		api_user,
		api_password,
		api_url,
		tlsVerify,
	)
}
