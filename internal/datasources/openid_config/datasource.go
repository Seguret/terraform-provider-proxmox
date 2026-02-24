package openid_config

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &OpenIDConfigDataSource{}

type OpenIDConfigDataSource struct {
	client *client.Client
}

type OpenIDConfigDataSourceModel struct {
	ID                     types.String `tfsdk:"id"`
	Realm                  types.String `tfsdk:"realm"`
	Issuer                 types.String `tfsdk:"issuer"`
	AuthorizationEndpoint  types.String `tfsdk:"authorization_endpoint"`
	TokenEndpoint          types.String `tfsdk:"token_endpoint"`
	UserinfoEndpoint       types.String `tfsdk:"userinfo_endpoint"`
	JwksURI                types.String `tfsdk:"jwks_uri"`
	ScopesSupported        types.List   `tfsdk:"scopes_supported"`
	ResponseTypesSupported types.List   `tfsdk:"response_types_supported"`
}

func NewDataSource() datasource.DataSource {
	return &OpenIDConfigDataSource{}
}

func (d *OpenIDConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_openid_config"
}

func (d *OpenIDConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves OpenID Connect configuration for a Proxmox VE authentication realm.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"realm": schema.StringAttribute{
				Description: "The authentication realm name (e.g., 'pve', 'openid-realm').",
				Required:    true,
			},
			"issuer": schema.StringAttribute{
				Description: "The OpenID issuer URL.",
				Computed:    true,
			},
			"authorization_endpoint": schema.StringAttribute{
				Description: "The authorization endpoint URL.",
				Computed:    true,
			},
			"token_endpoint": schema.StringAttribute{
				Description: "The token endpoint URL.",
				Computed:    true,
			},
			"userinfo_endpoint": schema.StringAttribute{
				Description: "The userinfo endpoint URL.",
				Computed:    true,
			},
			"jwks_uri": schema.StringAttribute{
				Description: "The JWKS (JSON Web Key Set) URI.",
				Computed:    true,
			},
			"scopes_supported": schema.ListAttribute{
				Description: "List of supported OAuth2 scopes.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"response_types_supported": schema.ListAttribute{
				Description: "List of supported OAuth2 response types.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *OpenIDConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cl, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = cl
}

func (d *OpenIDConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OpenIDConfigDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	realm := config.Realm.ValueString()

	oidcConfig, err := d.client.GetOpenIDConfig(ctx, realm)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OpenID configuration", err.Error())
		return
	}

	config.ID = types.StringValue(realm)
	config.Issuer = types.StringValue(oidcConfig.Issuer)
	config.AuthorizationEndpoint = types.StringValue(oidcConfig.AuthorizationEndpoint)
	config.TokenEndpoint = types.StringValue(oidcConfig.TokenEndpoint)
	config.UserinfoEndpoint = types.StringValue(oidcConfig.UserinfoEndpoint)
	config.JwksURI = types.StringValue(oidcConfig.JwksURI)

	scopes, diags := types.ListValueFrom(ctx, types.StringType, oidcConfig.ScopesSupported)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.ScopesSupported = scopes

	responseTypes, diags := types.ListValueFrom(ctx, types.StringType, oidcConfig.ResponseTypesSupported)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.ResponseTypesSupported = responseTypes

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
