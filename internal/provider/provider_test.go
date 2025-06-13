// Copyright (c) HashiCorp, Inc.

package provider_test

import (
	"regexp"
	"testing"

	"terraform-provider-infrapilot/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"infrapilot": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestProvider_ValidConfiguration(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "infrapilot" {
  token = "valid-token"
}
`,
			},
		},
	})
}

func TestProvider_MissingToken(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "infrapilot" {
  # token is missing
}

data "infrapilot_license_check" "test" {
  module = "example"
}
`,
				ExpectError: regexp.MustCompile(`The argument "token" is required`),
			},
		},
	})
}

func TestProvider_MalformedToken(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "infrapilot" {
  token = "malformed-token"
}

data "infrapilot_license_check" "fail" {
  module = "fail_module"
}
`,
				ExpectError: regexp.MustCompile(`token parse error|invalid or malformed token`),
			},
		},
	})
}
