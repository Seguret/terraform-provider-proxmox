package hardware_mapping_pci

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &PCIHardwareMappingDataSource{}
var _ datasource.DataSourceWithConfigure = &PCIHardwareMappingDataSource{}

type PCIHardwareMappingDataSource struct {
	client *client.Client
}

type PCIHardwareMappingDataSourceModel struct {
	ID   types.String   `tfsdk:"id"`
	IDs  []types.String `tfsdk:"ids"`
}

func NewDataSource() datasource.DataSource {
	return &PCIHardwareMappingDataSource{}
}

func (d *PCIHardwareMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_hardware_mapping_pci_list"
}

func (d *PCIHardwareMappingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of Proxmox VE cluster PCI hardware mapping names.",
		Attributes: map[string]schema.Attribute{
			"id":  schema.StringAttribute{Computed: true},
			"ids": schema.ListAttribute{Description: "PCI hardware mapping names.", Computed: true, ElementType: types.StringType},
		},
	}
}

func (d *PCIHardwareMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PCIHardwareMappingDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading PCI hardware mappings list")

	mappings, err := d.client.GetPCIHardwareMappings(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading PCI hardware mappings", err.Error())
		return
	}

	state := PCIHardwareMappingDataSourceModel{
		ID:  types.StringValue("pci_hardware_mappings"),
		IDs: make([]types.String, len(mappings)),
	}
	for i, m := range mappings {
		state.IDs[i] = types.StringValue(m.ID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
