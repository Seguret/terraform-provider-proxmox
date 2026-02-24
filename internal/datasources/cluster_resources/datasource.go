package cluster_resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ datasource.DataSource = &ClusterResourcesDataSource{}
var _ datasource.DataSourceWithConfigure = &ClusterResourcesDataSource{}

type ClusterResourcesDataSource struct {
	client *client.Client
}

type ClusterResourceModel struct {
	ID      types.String  `tfsdk:"id"`
	Type    types.String  `tfsdk:"type"`
	Node    types.String  `tfsdk:"node"`
	Status  types.String  `tfsdk:"status"`
	Name    types.String  `tfsdk:"name"`
	VMID    types.Int64   `tfsdk:"vmid"`
	Pool    types.String  `tfsdk:"pool"`
	CPU     types.Float64 `tfsdk:"cpu"`
	MaxCPU  types.Int64   `tfsdk:"max_cpu"`
	Mem     types.Int64   `tfsdk:"mem"`
	MaxMem  types.Int64   `tfsdk:"max_mem"`
	Disk    types.Int64   `tfsdk:"disk"`
	MaxDisk types.Int64   `tfsdk:"max_disk"`
	Uptime  types.Int64   `tfsdk:"uptime"`
	Storage types.String  `tfsdk:"storage"`
}

type ClusterResourcesDataSourceModel struct {
	ID           types.String           `tfsdk:"id"`
	ResourceType types.String           `tfsdk:"resource_type"`
	Resources    []ClusterResourceModel `tfsdk:"resources"`
}

func NewDataSource() datasource.DataSource {
	return &ClusterResourcesDataSource{}
}

func (d *ClusterResourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_cluster_resources"
}

func (d *ClusterResourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves resources across the Proxmox VE cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder identifier.",
				Computed:    true,
			},
			"resource_type": schema.StringAttribute{
				Description: "Optional filter for resource type (e.g., 'vm', 'storage', 'node', 'sdn').",
				Optional:    true,
			},
			"resources": schema.ListNestedAttribute{
				Description: "The list of cluster resources.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The resource identifier.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The resource type.",
							Computed:    true,
						},
						"node": schema.StringAttribute{
							Description: "The node the resource resides on.",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The current status of the resource.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The resource name.",
							Computed:    true,
						},
						"vmid": schema.Int64Attribute{
							Description: "The VM or container ID.",
							Computed:    true,
						},
						"pool": schema.StringAttribute{
							Description: "The resource pool the resource belongs to.",
							Computed:    true,
						},
						"cpu": schema.Float64Attribute{
							Description: "The current CPU utilization (0.0-1.0).",
							Computed:    true,
						},
						"max_cpu": schema.Int64Attribute{
							Description: "The maximum number of CPUs allocated.",
							Computed:    true,
						},
						"mem": schema.Int64Attribute{
							Description: "The current memory usage in bytes.",
							Computed:    true,
						},
						"max_mem": schema.Int64Attribute{
							Description: "The maximum memory allocated in bytes.",
							Computed:    true,
						},
						"disk": schema.Int64Attribute{
							Description: "The current disk usage in bytes.",
							Computed:    true,
						},
						"max_disk": schema.Int64Attribute{
							Description: "The maximum disk space allocated in bytes.",
							Computed:    true,
						},
						"uptime": schema.Int64Attribute{
							Description: "The uptime in seconds.",
							Computed:    true,
						},
						"storage": schema.StringAttribute{
							Description: "The storage identifier (for storage resources).",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *ClusterResourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClusterResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ClusterResourcesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceType := ""
	if !config.ResourceType.IsNull() && !config.ResourceType.IsUnknown() {
		resourceType = config.ResourceType.ValueString()
	}

	tflog.Debug(ctx, "Reading Proxmox VE cluster resources", map[string]any{"resource_type": resourceType})

	resources, err := d.client.GetClusterResources(ctx, resourceType)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cluster resources", err.Error())
		return
	}

	idVal := "cluster_resources"
	if resourceType != "" {
		idVal = fmt.Sprintf("cluster_resources/%s", resourceType)
	}

	state := ClusterResourcesDataSourceModel{
		ID:           types.StringValue(idVal),
		ResourceType: config.ResourceType,
		Resources:    make([]ClusterResourceModel, len(resources)),
	}

	for i, r := range resources {
		state.Resources[i] = ClusterResourceModel{
			ID:      types.StringValue(r.ID),
			Type:    types.StringValue(r.Type),
			Node:    types.StringValue(r.Node),
			Status:  types.StringValue(r.Status),
			Name:    types.StringValue(r.Name),
			VMID:    types.Int64Value(int64(r.VMID)),
			Pool:    types.StringValue(r.Pool),
			CPU:     types.Float64Value(r.CPU),
			MaxCPU:  types.Int64Value(int64(r.MaxCPU)),
			Mem:     types.Int64Value(r.Mem),
			MaxMem:  types.Int64Value(r.MaxMem),
			Disk:    types.Int64Value(r.Disk),
			MaxDisk: types.Int64Value(r.MaxDisk),
			Uptime:  types.Int64Value(r.Uptime),
			Storage: types.StringValue(r.Storage),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
