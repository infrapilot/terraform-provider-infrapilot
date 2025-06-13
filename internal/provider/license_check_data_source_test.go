// Copyright (c) 2025 InfraPilot, LLC
// SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"testing"

	"terraform-provider-infrapilot/internal/model"
	"terraform-provider-infrapilot/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLicenseCheckDataSource(t *testing.T) {
	t.Parallel()

	testOrgID := "test-org"
	testTier := "pro"
	testModule := "sample_module"

	// Mock JWT validator that returns expected claims
	mockValidator := func(token, jwksURL string) (*model.LicenseClaims, error) {
		return &model.LicenseClaims{
			OrgID: testOrgID,
			Tier:  testTier,
		}, nil
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"infrapilot": providerserver.NewProtocol6WithError(
				provider.NewWithValidator("test", mockValidator)(),
			),
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "infrapilot" {
  token = "dummy-token"
}

data "infrapilot_license_check" "test" {
  module = "sample_module"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.infrapilot_license_check.test", "org_id", testOrgID),
					resource.TestCheckResourceAttr("data.infrapilot_license_check.test", "tier", testTier),
					resource.TestCheckResourceAttr("data.infrapilot_license_check.test", "module", testModule),
				),
			},
		},
	})
}
