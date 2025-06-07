// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/infra-pilot/terraform-provider-infrapilot/internal/license"
	"github.com/infra-pilot/terraform-provider-infrapilot/internal/telemetry"
)

type infrapilotProvider struct {
	version string
}

type infrapilotProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &infrapilotProvider{
			version: version,
		}
	}
}

func (p *infrapilotProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "infrapilot"
	resp.Version = p.version
}

func (p *infrapilotProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "License token for accessing InfraPilot modules.",
			},
		},
	}
}

func (p *infrapilotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config infrapilotProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var token string
	if !config.Token.IsUnknown() && !config.Token.IsNull() {
		token = config.Token.ValueString()
	}
	if token == "" {
		token = os.Getenv("INFRAPILOT_LICENSE_TOKEN")
	}
	if token == "" {
		resp.Diagnostics.AddError("Missing Token", "A token must be provided via the provider configuration or INFRAPILOT_LICENSE_TOKEN environment variable.")
		return
	}

	result, err := license.Validate(token)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Token", fmt.Sprintf("Token validation failed: %s", err))
		return
	}

	telemetry.Log(req.TerraformVersion, p.version, token)

	resp.DataSourceData = token
	resp.ResourceData = token

	_ = result // result reserved for future use
}

func (p *infrapilotProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *infrapilotProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewExampleEphemeralResource,
	}
}

func (p *infrapilotProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		NewCheckAccessDataSource,
	}
}

func (p *infrapilotProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}
