// Copyright (c) 2025 InfraPilot, LLC
// Originally derived from HashiCorp's terraform-provider-scaffolding-framework
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-infrapilot/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var version = "dev" // will be overridden at build time via ldflags

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/infra-pilot/infrapilot",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		log.Fatalf("error running provider: %s", err)
	}
}
