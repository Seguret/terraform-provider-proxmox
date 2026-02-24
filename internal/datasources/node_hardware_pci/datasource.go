package node_hardware_pci

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeHardwarePCIDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeHardwarePCIDataSource{}

type NodeHardwarePCIDataSource struct {
	client *client.Client
}

type NodeHardwarePCIModel struct {
	ID         types.String `tfsdk:"id"`
	Class      types.String `tfsdk:"class"`
	Device     types.String `tfsdk:"device"`
	Vendor     types.String `tfsdk:"vendor"`
	DeviceID   types.String `tfsdk:"device_id"`
	VendorID   types.String `tfsdk:"vendor_id"`
	IOMMUGroup types.Int64  `tfsdk:"iommu_group"`
}

type NodeHardwarePCIDataSourceModel struct {
	ID       types.String           `tfsdk:"id"`
	NodeName types.String           `tfsdk:"node_name"`
	Devices  []NodeHardwarePCIModel `tfsdk:"devices"`
}

func NewDataSource() datasource.DataSource {
	return &NodeHardwarePCIDataSource{}
}

func (d *NodeHardwarePCIDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_hardware_pci"
}

func (d *NodeHardwarePCIDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of PCI devices on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"devices": schema.ListNestedAttribute{
				Description: "The list of PCI devices on the node.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The PCI device identifier.",
							Computed:    true,
						},
						"class": schema.StringAttribute{
							Description: "The PCI device class.",
							Computed:    true,
						},
						"device": schema.StringAttribute{
							Description: "The PCI device name.",
							Computed:    true,
						},
						"vendor": schema.StringAttribute{
							Description: "The PCI device vendor name.",
							Computed:    true,
						},
						"device_id": schema.StringAttribute{
							Description: "The PCI device ID.",
							Computed:    true,
						},
						"vendor_id": schema.StringAttribute{
							Description: "The PCI vendor ID.",
							Computed:    true,
						},
						"iommu_group": schema.Int64Attribute{
							Description: "The IOMMU group the device belongs to.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *NodeHardwarePCIDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *NodeHardwarePCIDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeHardwarePCIDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE node PCI devices", map[string]any{"node": node})

	devices, err := d.client.ListNodeHardwarePCI(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node PCI devices", err.Error())
		return
	}

	state := NodeHardwarePCIDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/hardware/pci", node)),
		NodeName: config.NodeName,
		Devices:  make([]NodeHardwarePCIModel, len(devices)),
	}

	for i, d := range devices {
		state.Devices[i] = NodeHardwarePCIModel{
			ID:         types.StringValue(d.ID),
			Class:      types.StringValue(d.Class),
			Device:     types.StringValue(d.Device),
			Vendor:     types.StringValue(d.Vendor),
			DeviceID:   types.StringValue(d.DeviceID),
			VendorID:   types.StringValue(d.VendorID),
			IOMMUGroup: types.Int64Value(int64(d.IOMMUGroup)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
