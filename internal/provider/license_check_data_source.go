// Copyright (c) 2025 InfraPilot, LLC
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-infrapilot/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type licenseCheckDataSource struct {
	claims *model.LicenseMetadata
}

var (
	_ datasource.DataSource              = &licenseCheckDataSource{}
	_ datasource.DataSourceWithConfigure = &licenseCheckDataSource{}
)

func NewLicenseCheckDataSource() datasource.DataSource {
	return &licenseCheckDataSource{}
}

func (d *licenseCheckDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_license_check"
}

func (d *licenseCheckDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Validates the InfraPilot license token and exposes minimal metadata.",
		Attributes: map[string]schema.Attribute{
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "The organization ID tied to the license token.",
			},
			"tier": schema.StringAttribute{
				Computed:    true,
				Description: "The plan tier (e.g. pro, enterprise, etc.).",
			},
			"module": schema.StringAttribute{
				Required:    true,
				Description: "The name of the module using this data source.",
			},
		},
	}
}

func (d *licenseCheckDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "licenseCheckDataSource Configure() called")

	if req.ProviderData == nil {
		return
	}

	claims, ok := req.ProviderData.(*model.LicenseMetadata)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *model.LicenseMetadata, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.claims = claims
}

func (d *licenseCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.claims == nil {
		resp.Diagnostics.AddError("Missing License Claims", "The license token could not be validated or no claims were passed.")
		return
	}

	var config struct {
		Module types.String `tfsdk:"module"`
		OrgID  types.String `tfsdk:"org_id"`
		Tier   types.String `tfsdk:"tier"`
	}

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result := model.LicenseMetadata{
		OrgID:  d.claims.OrgID,
		Tier:   d.claims.Tier,
		Module: config.Module,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
