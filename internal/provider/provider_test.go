// Copyright (c) HashiCorp, Inc.

package provider_test

import (
	"errors"
	"regexp"
	"testing"

	"terraform-provider-infrapilot/internal/model"
	"terraform-provider-infrapilot/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func mockJWTValidator(token, _ string) (*model.LicenseClaims, error) {
	switch token {
	case "valid-token":
		return &model.LicenseClaims{
			OrgID: "test-org",
			Tier:  "pro",
		}, nil
	case "malformed-token":
		return nil, errors.New("invalid or malformed token")
	default:
		return nil, errors.New("token parse error")
	}
}

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"infrapilot": providerserver.NewProtocol6WithError(
		provider.NewWithValidator("test", mockJWTValidator)(),
	),
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

data "infrapilot_license_check" "check" {
  module = "example"
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
  # token is missing, no env var set
}

data "infrapilot_license_check" "test" {
  module = "example"
}
`,
				ExpectError: regexp.MustCompile(`no token was found`),
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
				ExpectError: regexp.MustCompile(`invalid or malformed token`),
			},
		},
	})
}

func TestProvider_UsesEnvVar(t *testing.T) {
	t.Setenv("INFRAPILOT_TOKEN", "valid-token")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "infrapilot" {
  # token deliberately omitted
}

data "infrapilot_license_check" "check" {
  module = "example"
}
`,
			},
		},
	})
}
