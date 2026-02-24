package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &DNSDataSource{}
var _ datasource.DataSourceWithConfigure = &DNSDataSource{}

type DNSDataSource struct {
	client *client.Client
}

type DNSDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	NodeName types.String   `tfsdk:"node_name"`
	Domain   types.String   `tfsdk:"domain"`
	Servers  []types.String `tfsdk:"servers"`
}

func NewDataSource() datasource.DataSource {
	return &DNSDataSource{}
}

func (d *DNSDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_dns"
}

func (d *DNSDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the DNS configuration of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The DNS datasource identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The node to retrieve DNS configuration from.",
				Required:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The search domain configured on the node.",
				Computed:    true,
			},
			"servers": schema.ListAttribute{
				Description: "The list of DNS server addresses configured on the node.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *DNSDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DNSDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DNSDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading node DNS", map[string]any{"node": node})

	dnsConfig, err := d.client.GetNodeDNS(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node DNS", err.Error())
		return
	}

	// build the list from DNS1/2/3, skipping empty slots
	servers := []types.String{}
	for _, srv := range []string{dnsConfig.DNS1, dnsConfig.DNS2, dnsConfig.DNS3} {
		if srv != "" {
			servers = append(servers, types.StringValue(srv))
		}
	}

	state := DNSDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/dns", node)),
		NodeName: types.StringValue(node),
		Domain:   types.StringValue(dnsConfig.Search),
		Servers:  servers,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
