package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-icinga2/icinga2"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: icinga2.Provider})
}
