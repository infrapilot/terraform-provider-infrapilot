package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/infra-pilot/terraform-provider-infrapilot/internal/license"
)

// Ensure interface compliance
var _ datasource.DataSource = &checkAccessDataSource{}

func NewCheckAccessDataSource() datasource.DataSource { return &checkAccessDataSource{} }

type checkAccessDataSource struct {
	token string
}

type checkAccessModel struct {
	ModuleName        types.String `tfsdk:"module_name"`
	Token             types.String `tfsdk:"token"`
	AccessGranted     types.Bool   `tfsdk:"access_granted"`
	SubscriptionLevel types.String `tfsdk:"subscription_level"`
	ErrorMessage      types.String `tfsdk:"error_message"`
}

func (d *checkAccessDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_check_access"
}

func (d *checkAccessDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"module_name": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Optional: true,
			},
			"access_granted": schema.BoolAttribute{
				Computed: true,
			},
			"subscription_level": schema.StringAttribute{
				Computed: true,
			},
			"error_message": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *checkAccessDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if t, ok := req.ProviderData.(string); ok {
		d.token = t
	}
}

func (d *checkAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data checkAccessModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := d.token
	if !data.Token.IsNull() && !data.Token.IsUnknown() {
		token = data.Token.ValueString()
	}

	result, err := license.Validate(token)
	if err != nil {
		data.AccessGranted = types.BoolValue(false)
		data.ErrorMessage = types.StringValue(err.Error())
		if result != nil {
			data.SubscriptionLevel = types.StringValue(result.SubscriptionLevel)
		}
	} else {
		data.AccessGranted = types.BoolValue(true)
		data.SubscriptionLevel = types.StringValue(result.SubscriptionLevel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
