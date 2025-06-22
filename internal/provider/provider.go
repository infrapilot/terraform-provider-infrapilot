// Copyright (c) 2025 InfraPilot, LLC
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-infrapilot/internal/jwt"
	"terraform-provider-infrapilot/internal/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const jwksURL = "https://license.infrapilot.ai/.well-known/jwks.json"

// Ensure provider satisfies the expected interfaces.
var (
	_ provider.Provider = &infraPilotProvider{}
)

type JWTValidator func(token, jwksURL string) (*model.LicenseClaims, error)

type infraPilotProvider struct {
	version      string
	jwtValidator JWTValidator
}

func New(version string) func() provider.Provider {
	return NewWithValidator(version, jwt.ValidateToken)
}

func NewWithValidator(version string, validator JWTValidator) func() provider.Provider {
	return func() provider.Provider {
		return &infraPilotProvider{
			version:      version,
			jwtValidator: validator,
		}
	}
}

func (p *infraPilotProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "infrapilot"
	resp.Version = p.version
}

func (p *infraPilotProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "InfraPilot provider used to validate access and license subscription to private Terraform modules.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "JWT token (valid for 7 days) used for license validation. May also be provided via INFRAPILOT_TOKEN env var.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

type infraPilotProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *infraPilotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring InfraPilot provider")

	var config infraPilotProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine token value (prefer explicit config, fallback to env var)
	var token string
	if !config.Token.IsNull() && !config.Token.IsUnknown() {
		token = config.Token.ValueString()
	} else {
		token = os.Getenv("INFRAPILOT_TOKEN")
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing InfraPilot Token",
			"The provider cannot be created as no token was found. Provide it either via the 'token' argument in the provider block or the INFRAPILOT_TOKEN environment variable.",
		)
		return
	}

	// Mask token for logs
	ctx = tflog.SetField(ctx, "infrapilot_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "infrapilot_token")

	tflog.Debug(ctx, "Using JWT token for InfraPilot authentication")

	// Validate the token using the configured validator
	claims, err := p.jwtValidator(token, jwksURL)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Validate InfraPilot Token",
			"An error occurred while validating the InfraPilot token:\n\n"+err.Error(),
		)
		return
	}

	// Pass claims to downstream resources/datasources
	resp.DataSourceData = &model.LicenseMetadata{
		OrgID: types.StringValue(claims.OrgID),
		Tier:  types.StringValue(claims.Tier),
	}
	resp.ResourceData = claims

	tflog.Info(ctx, "Configured InfraPilot client", map[string]any{"success": true})
}

func (p *infraPilotProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *infraPilotProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLicenseCheckDataSource,
	}
}
