package main

import (
	"github.com/Icinga/terraform-provider-icinga2/icinga2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: icinga2.Provider})
}
