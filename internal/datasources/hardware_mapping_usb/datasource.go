package hardware_mapping_usb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &USBHardwareMappingDataSource{}
var _ datasource.DataSourceWithConfigure = &USBHardwareMappingDataSource{}

type USBHardwareMappingDataSource struct {
	client *client.Client
}

type USBHardwareMappingDataSourceModel struct {
	ID  types.String   `tfsdk:"id"`
	IDs []types.String `tfsdk:"ids"`
}

func NewDataSource() datasource.DataSource {
	return &USBHardwareMappingDataSource{}
}

func (d *USBHardwareMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_hardware_mapping_usb_list"
}

func (d *USBHardwareMappingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE cluster USB hardware mapping names.",
		Attributes: map[string]schema.Attribute{
			"id":  schema.StringAttribute{Computed: true},
			"ids": schema.ListAttribute{Description: "USB hardware mapping names.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *USBHardwareMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *USBHardwareMappingDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading USB hardware mappings list")

	mappings, err := d.client.GetUSBHardwareMappings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading USB hardware mappings", err.Error())
		return
	}

	state := USBHardwareMappingDataSourceModel{
		ID:  types.StringValue("usb_hardware_mappings"),
		IDs: make([]types.String, len(mappings)),
	}
	for i, m := range mappings {
		state.IDs[i] = types.StringValue(m.ID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
