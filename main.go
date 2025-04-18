package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/Icinga/terraform-provider-icinga2/internal/icinga2"
)

var (
	version = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(
		context.Background(),
		icinga2.New(version),
		providerserver.ServeOpts{
			Address: "registry.terraform.io/Icinga/icinga2",
			Debug:   debug,
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
