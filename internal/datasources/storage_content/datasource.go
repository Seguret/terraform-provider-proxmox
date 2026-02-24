package storage_content

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

var _ datasource.DataSource = &StorageContentDataSource{}
var _ datasource.DataSourceWithConfigure = &StorageContentDataSource{}

type StorageContentDataSource struct {
	client *client.Client
}

type StorageContentDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	Storage     types.String `tfsdk:"storage"`
	ContentType types.String `tfsdk:"content_type"`

	VolIDs  []types.String `tfsdk:"volids"`
	Content []types.String `tfsdk:"content"`
	Format  []types.String `tfsdk:"format"`
	Size    []types.Int64  `tfsdk:"size"`
	Used    []types.Int64  `tfsdk:"used"`
	CTime   []types.Int64  `tfsdk:"ctime"`
	Notes   []types.String `tfsdk:"notes"`
}

func NewDataSource() datasource.DataSource {
	return &StorageContentDataSource{}
}

func (d *StorageContentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_storage_content"
}

func (d *StorageContentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the list of files stored in a Proxmox VE node storage.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node.",
				Required:    true,
			},
			"storage": schema.StringAttribute{
				Description: "The storage name (e.g., 'local', 'local-lvm').",
				Required:    true,
			},
			"content_type": schema.StringAttribute{
				Description: "Optional filter for content type (e.g., 'iso', 'vztmpl', 'backup').",
				Optional:    true,
			},
			"volids": schema.ListAttribute{
				Description: "The volume IDs for each file.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"content": schema.ListAttribute{
				Description: "The content type for each file.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"format": schema.ListAttribute{
				Description: "The format for each file (if provided by Proxmox).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"size": schema.ListAttribute{
				Description: "The size in bytes for each file.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"used": schema.ListAttribute{
				Description: "The used size in bytes for each file (if provided).",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"ctime": schema.ListAttribute{
				Description: "The creation time (unix timestamp) for each file (if provided).",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"notes": schema.ListAttribute{
				Description: "Notes for each file (if provided).",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *StorageContentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *StorageContentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config StorageContentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nodeName := config.NodeName.ValueString()
	storage := config.Storage.ValueString()
	contentFilter := ""
	if !config.ContentType.IsNull() && !config.ContentType.IsUnknown() {
		contentFilter = strings.TrimSpace(config.ContentType.ValueString())
	}

	tflog.Debug(ctx, "Reading Proxmox VE storage content", map[string]any{
		"node":    nodeName,
		"storage": storage,
		"content": contentFilter,
	})

	items, err := d.client.GetStorageContent(ctx, nodeName, storage)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Proxmox VE Storage Content",
			fmt.Sprintf("An error occurred while reading storage content for node '%s' and storage '%s': %s", nodeName, storage, err.Error()),
		)
		return
	}

	filtered := make([]int, 0, len(items))
	for i, item := range items {
		if contentFilter == "" || strings.EqualFold(item.Content, contentFilter) {
			filtered = append(filtered, i)
		}
	}

	state := StorageContentDataSourceModel{
		ID:          types.StringValue(fmt.Sprintf("storage-content/%s/%s", nodeName, storage)),
		NodeName:    types.StringValue(nodeName),
		Storage:     types.StringValue(storage),
		ContentType: config.ContentType,
		VolIDs:      make([]types.String, len(filtered)),
		Content:     make([]types.String, len(filtered)),
		Format:      make([]types.String, len(filtered)),
		Size:        make([]types.Int64, len(filtered)),
		Used:        make([]types.Int64, len(filtered)),
		CTime:       make([]types.Int64, len(filtered)),
		Notes:       make([]types.String, len(filtered)),
	}

	for i, idx := range filtered {
		item := items[idx]
		state.VolIDs[i] = types.StringValue(item.VolID)
		state.Content[i] = types.StringValue(item.Content)
		state.Format[i] = types.StringValue(item.Format)
		state.Size[i] = types.Int64Value(item.Size)
		state.Used[i] = types.Int64Value(item.Used)
		state.CTime[i] = types.Int64Value(item.CTime)
		state.Notes[i] = types.StringValue(item.Notes)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}