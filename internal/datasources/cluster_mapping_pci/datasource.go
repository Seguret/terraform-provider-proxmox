package cluster_mapping_pci

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *client.Client
}

type DataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Mappings types.List   `tfsdk:"mappings"`
}

type MappingModel struct {
	ID          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Map         types.List   `tfsdk:"map"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_mapping_pci"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves PCI hardware mappings from the cluster.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Data source identifier.",
				Computed:            true,
			},
			"mappings": schema.ListNestedAttribute{
				MarkdownDescription: "List of PCI hardware mappings.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Mapping ID.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Mapping description.",
							Computed:            true,
						},
						"map": schema.ListAttribute{
							MarkdownDescription: "List of PCI device mappings.",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *client.Client, got something else.",
		)
		return
	}

	d.client = c
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	mappings, err := d.client.GetHardwareMappingPCI(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", "Failed to retrieve PCI mappings: "+err.Error())
		return
	}

	data.ID = types.StringValue("cluster_mapping_pci")

	var mappingModels []MappingModel
	for _, m := range mappings {
		mapList, diags := types.ListValueFrom(ctx, types.StringType, m.Map)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		mappingModels = append(mappingModels, MappingModel{
			ID:          types.StringValue(m.ID),
			Description: types.StringValue(m.Description),
			Map:         mapList,
		})
	}

	mappingsList, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"description": types.StringType,
			"map":         types.ListType{ElemType: types.StringType},
		},
	}, mappingModels)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Mappings = mappingsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
