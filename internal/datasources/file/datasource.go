package file

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

var _ datasource.DataSource = &FileDataSource{}
var _ datasource.DataSourceWithConfigure = &FileDataSource{}

type FileDataSource struct {
	client *client.Client
}

type FileDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	DatastoreID types.String `tfsdk:"datastore_id"`
	VolumeID    types.String `tfsdk:"volume_id"`
	ContentType types.String `tfsdk:"content_type"`
	Size        types.Int64  `tfsdk:"size"`
	FileName    types.String `tfsdk:"file_name"`
}

func NewDataSource() datasource.DataSource {
	return &FileDataSource{}
}

func (d *FileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_file"
}

func (d *FileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a specific file stored in Proxmox VE storage.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The volume ID of the file (same as volume_id).",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"datastore_id": schema.StringAttribute{
				Description: "The storage ID (e.g. 'local').",
				Required:    true,
			},
			"volume_id": schema.StringAttribute{
				Description: "The volume ID to look up (e.g. 'local:iso/ubuntu.iso').",
				Required:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "The content type of the file (e.g. 'iso', 'vztmpl').",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "The file size in bytes as reported by Proxmox.",
				Computed:    true,
			},
			"file_name": schema.StringAttribute{
				Description: "The filename portion of the volume ID.",
				Computed:    true,
			},
		},
	}
}

func (d *FileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config FileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := config.NodeName.ValueString()
	storage := config.DatastoreID.ValueString()
	volumeID := config.VolumeID.ValueString()

	tflog.Debug(ctx, "Reading Proxmox VE file", map[string]any{
		"node":      node,
		"storage":   storage,
		"volume_id": volumeID,
	})

	contents, err := d.client.GetStorageContent(ctx, node, storage)
	if err != nil {
		resp.Diagnostics.AddError("Error reading storage content",
			fmt.Sprintf("Unable to list storage content on node '%s', storage '%s': %s", node, storage, err.Error()))
		return
	}

	for _, item := range contents {
		if item.VolID == volumeID {
			// pull the filename from the last path component of the volid
			fileName := volumeID
			if idx := strings.LastIndex(volumeID, "/"); idx >= 0 {
				fileName = volumeID[idx+1:]
			}

			state := FileDataSourceModel{
				ID:          types.StringValue(item.VolID),
				NodeName:    types.StringValue(node),
				DatastoreID: types.StringValue(storage),
				VolumeID:    types.StringValue(item.VolID),
				ContentType: types.StringValue(item.Content),
				Size:        types.Int64Value(item.Size),
				FileName:    types.StringValue(fileName),
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	resp.Diagnostics.AddError("File not found",
		fmt.Sprintf("No file with volume ID '%s' was found in storage '%s' on node '%s'.",
			volumeID, storage, node))
}
