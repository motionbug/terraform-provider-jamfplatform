// Copyright 2025 Jamf Software LLC.

package main

import (
	"context"
	"log"

	jamfcompliancebenchmarkengine "github.com/Jamf-Concepts/terraform-provider-jamfcompliancebenchmarkengine/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), jamfcompliancebenchmarkengine.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/Jamf-Concepts/jamfcompliancebenchmarkengine",
	})
	if err != nil {
		log.Fatalf("Error serving provider: %v", err)
	}
}
