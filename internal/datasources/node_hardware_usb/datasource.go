package node_hardware_usb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &NodeHardwareUSBDataSource{}
var _ datasource.DataSourceWithConfigure = &NodeHardwareUSBDataSource{}

type NodeHardwareUSBDataSource struct {
	client *client.Client
}

type NodeHardwareUSBModel struct {
	BusNum       types.Int64  `tfsdk:"bus_num"`
	DevNum       types.Int64  `tfsdk:"dev_num"`
	Manufacturer types.String `tfsdk:"manufacturer"`
	Product      types.String `tfsdk:"product"`
	ProdID       types.String `tfsdk:"prod_id"`
	VendorID     types.String `tfsdk:"vendor_id"`
	Speed        types.String `tfsdk:"speed"`
	Serialnumber types.String `tfsdk:"serialnumber"`
}

type NodeHardwareUSBDataSourceModel struct {
	ID       types.String           `tfsdk:"id"`
	NodeName types.String           `tfsdk:"node_name"`
	Devices  []NodeHardwareUSBModel `tfsdk:"devices"`
}

func NewDataSource() datasource.DataSource {
	return &NodeHardwareUSBDataSource{}
}

func (d *NodeHardwareUSBDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_hardware_usb"
}

func (d *NodeHardwareUSBDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of USB devices on a Proxmox VE node.",
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
				Description: "The list of USB devices on the node.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"bus_num": schema.Int64Attribute{
							Description: "The USB bus number.",
							Computed:    true,
						},
						"dev_num": schema.Int64Attribute{
							Description: "The USB device number.",
							Computed:    true,
						},
						"manufacturer": schema.StringAttribute{
							Description: "The device manufacturer name.",
							Computed:    true,
						},
						"product": schema.StringAttribute{
							Description: "The product name.",
							Computed:    true,
						},
						"prod_id": schema.StringAttribute{
							Description: "The USB product ID.",
							Computed:    true,
						},
						"vendor_id": schema.StringAttribute{
							Description: "The USB vendor ID.",
							Computed:    true,
						},
						"speed": schema.StringAttribute{
							Description: "The USB device speed (e.g., '480', '5000').",
							Computed:    true,
						},
						"serialnumber": schema.StringAttribute{
							Description: "The device serial number.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *NodeHardwareUSBDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NodeHardwareUSBDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeHardwareUSBDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	tflog.Debug(ctx, "Reading Proxmox VE node USB devices", map[string]any{"node": node})

	devices, err := d.client.ListNodeHardwareUSB(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node USB devices", err.Error())
		return
	}

	state := NodeHardwareUSBDataSourceModel{
		ID:       types.StringValue(fmt.Sprintf("%s/hardware/usb", node)),
		NodeName: config.NodeName,
		Devices:  make([]NodeHardwareUSBModel, len(devices)),
	}

	for i, dev := range devices {
		state.Devices[i] = NodeHardwareUSBModel{
			BusNum:       types.Int64Value(int64(dev.BusNum)),
			DevNum:       types.Int64Value(int64(dev.DevNum)),
			Manufacturer: types.StringValue(dev.Manufacturer),
			Product:      types.StringValue(dev.Product),
			ProdID:       types.StringValue(dev.ProdID),
			VendorID:     types.StringValue(dev.VendorID),
			Speed:        types.StringValue(dev.Speed),
			Serialnumber: types.StringValue(dev.Serialnumber),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
