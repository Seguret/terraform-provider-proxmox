package hosts

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &HostsDataSource{}
var _ datasource.DataSourceWithConfigure = &HostsDataSource{}

type HostsDataSource struct {
	client *client.Client
}

type HostEntryModel struct {
	Address   types.String   `tfsdk:"address"`
	Hostnames []types.String `tfsdk:"hostnames"`
}

type HostsDataSourceModel struct {
	ID       types.String     `tfsdk:"id"`
	NodeName types.String     `tfsdk:"node_name"`
	Digest   types.String     `tfsdk:"digest"`
	Entries  []HostEntryModel `tfsdk:"entries"`
}

func NewDataSource() datasource.DataSource {
	return &HostsDataSource{}
}

func (d *HostsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_hosts"
}

func (d *HostsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the /etc/hosts configuration of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The hosts datasource identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The node to retrieve hosts configuration from.",
				Required:    true,
			},
			"digest": schema.StringAttribute{
				Description: "The SHA1 digest of the current hosts content.",
				Computed:    true,
			},
			"entries": schema.ListNestedAttribute{
				Description: "The list of host entries parsed from /etc/hosts.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Description: "The IP address of the host entry.",
							Computed:    true,
						},
						"hostnames": schema.ListAttribute{
							Description: "The hostnames associated with this address.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *HostsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// parseHostsData turns raw /etc/hosts content into structured entries.
// Comments and blank lines are ignored.
func parseHostsData(data string) []HostEntryModel {
	var entries []HostEntryModel

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip blank lines and comment lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// strip anything after a # on the same line
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		address := fields[0]
		hostnames := make([]types.String, len(fields)-1)
		for i, h := range fields[1:] {
			hostnames[i] = types.StringValue(h)
		}

		entries = append(entries, HostEntryModel{
			Address:   types.StringValue(address),
			Hostnames: hostnames,
		})
	}

	if entries == nil {
		entries = []HostEntryModel{}
	}

	return entries
}

func (d *HostsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config HostsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading node hosts", map[string]any{"node": node})

	hostsConfig, err := d.client.GetNodeHosts(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node hosts", err.Error())
		return
	}

	state := HostsDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/hosts", node)),
		NodeName: types.StringValue(node),
		Digest:   types.StringValue(hostsConfig.Digest),
		Entries:  parseHostsData(hostsConfig.Data),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
