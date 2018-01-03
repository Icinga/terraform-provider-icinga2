package main

import (
	"github.com/caseyr232/terraform-provider-icinga2/icinga2"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: icinga2.Provider})
}
