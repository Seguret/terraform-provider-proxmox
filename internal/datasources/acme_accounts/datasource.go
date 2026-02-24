package acme_accounts

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ACMEAccountsDataSource{}
var _ datasource.DataSourceWithConfigure = &ACMEAccountsDataSource{}

type ACMEAccountsDataSource struct {
	client *client.Client
}

type ACMEAccountEntryModel struct {
	Name types.String `tfsdk:"name"`
}

type ACMEAccountsDataSourceModel struct {
	ID       types.String             `tfsdk:"id"`
	Accounts []ACMEAccountEntryModel  `tfsdk:"accounts"`
}

func NewDataSource() datasource.DataSource {
	return &ACMEAccountsDataSource{}
}

func (d *ACMEAccountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acme_accounts"
}

func (d *ACMEAccountsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE ACME accounts.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for this data source.",
				Computed:    true,
			},
			"accounts": schema.ListNestedAttribute{
				Description: "The list of ACME accounts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The ACME account name.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ACMEAccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ACMEAccountsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading ACME accounts list")

	accounts, err := d.client.GetACMEAccounts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading ACME accounts", err.Error())
		return
	}

	state := ACMEAccountsDataSourceModel{
		ID:       types.StringValue("acme_accounts"),
		Accounts: make([]ACMEAccountEntryModel, len(accounts)),
	}

	for i, a := range accounts {
		state.Accounts[i] = ACMEAccountEntryModel{
			Name: types.StringValue(a.Name),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
