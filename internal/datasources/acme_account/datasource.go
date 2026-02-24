package acme_account

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ACMEAccountDataSource{}
var _ datasource.DataSourceWithConfigure = &ACMEAccountDataSource{}

type ACMEAccountDataSource struct {
	client *client.Client
}

type ACMEAccountDataSourceModel struct {
	ID        types.String   `tfsdk:"id"`
	Name      types.String   `tfsdk:"name"`
	Email     []types.String `tfsdk:"email"`
	Directory types.String   `tfsdk:"directory"`
	TOS       types.String   `tfsdk:"tos"`
	CreatedAt types.String   `tfsdk:"created_at"`
}

func NewDataSource() datasource.DataSource {
	return &ACMEAccountDataSource{}
}

func (d *ACMEAccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acme_account"
}

func (d *ACMEAccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a Proxmox VE ACME account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this data source.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The ACME account name.",
				Required:    true,
			},
			"email": schema.ListAttribute{
				Description: "The contact email addresses for this ACME account.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"directory": schema.StringAttribute{
				Description: "The ACME directory URL.",
				Computed:    true,
			},
			"tos": schema.StringAttribute{
				Description: "The Terms of Service URL.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The account creation timestamp.",
				Computed:    true,
			},
		},
	}
}

func (d *ACMEAccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *ACMEAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ACMEAccountDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := config.Name.ValueString()
	tflog.Debug(ctx, "Reading ACME account", map[string]any{"name": name})

	account, err := d.client.GetACMEAccount(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ACME account", err.Error())
		return
	}

	// Contact is a comma-separated list of mailto: URIs from proxmox.
	// parse them into plain email addresses by stripping the "mailto:" prefix.
	emails := parseContactEmails(account.Contact)

	state := ACMEAccountDataSourceModel{
		ID:        types.StringValue(name),
		Name:      types.StringValue(name),
		Email:     emails,
		Directory: types.StringValue(account.Directory),
		TOS:       types.StringValue(account.TosURL),
		CreatedAt: types.StringValue(account.AccountURL),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// parseContactEmails splits proxmox's comma-separated mailto: URI string
// into plain email addresses.
func parseContactEmails(contact string) []types.String {
	if contact == "" {
		return []types.String{}
	}
	parts := strings.Split(contact, ",")
	result := make([]types.String, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimPrefix(p, "mailto:")
		if p != "" {
			result = append(result, types.StringValue(p))
		}
	}
	return result
}
