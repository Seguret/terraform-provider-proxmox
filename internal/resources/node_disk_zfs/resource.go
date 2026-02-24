package node_disk_zfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &NodeDiskZFSResource{}
var _ resource.ResourceWithConfigure = &NodeDiskZFSResource{}
var _ resource.ResourceWithImportState = &NodeDiskZFSResource{}

type NodeDiskZFSResource struct {
	client *client.Client
}

type NodeDiskZFSResourceModel struct {
	ID          types.String `tfsdk:"id"`
	NodeName    types.String `tfsdk:"node_name"`
	Devices     types.String `tfsdk:"devices"`
	Name        types.String `tfsdk:"name"`
	RAIDLevel   types.String `tfsdk:"raid_level"`
	Ashift      types.Int64  `tfsdk:"ashift"`
	Compression types.String `tfsdk:"compression"`
	AddStorage  types.Bool   `tfsdk:"add_storage"`
}

func NewResource() resource.Resource {
	return &NodeDiskZFSResource{}
}

func (r *NodeDiskZFSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_disk_zfs"
}

func (r *NodeDiskZFSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ZFS pool on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"devices": schema.StringAttribute{
				Description: "Comma-separated list of device paths to use for the ZFS pool (e.g. /dev/sdb,/dev/sdc).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The ZFS pool name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"raid_level": schema.StringAttribute{
				Description: "The RAID level (mirror, raidz, raidz2, raidz3, single, draid, draid2, draid3).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ashift": schema.Int64Attribute{
				Description: "Pool sector size exponent (ashift value, e.g. 12 for 4K sectors).",
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"compression": schema.StringAttribute{
				Description: "ZFS compression algorithm (e.g. lz4, zstd, on, off).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"add_storage": schema.BoolAttribute{
				Description: "Whether to automatically add the created ZFS pool as a Proxmox storage.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NodeDiskZFSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *NodeDiskZFSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NodeDiskZFSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	createReq := &models.NodeDiskZFSCreateRequest{
		Devices:     plan.Devices.ValueString(),
		Name:        plan.Name.ValueString(),
		RAIDLevel:   plan.RAIDLevel.ValueString(),
		Compression: plan.Compression.ValueString(),
		AddStorage:  plan.AddStorage.ValueBool(),
	}
	if !plan.Ashift.IsNull() && !plan.Ashift.IsUnknown() {
		createReq.Ashift = int(plan.Ashift.ValueInt64())
	}

	tflog.Debug(ctx, "Creating node disk ZFS pool", map[string]any{
		"node":    node,
		"name":    createReq.Name,
		"devices": createReq.Devices,
		"raid":    createReq.RAIDLevel,
	})

	upid, err := r.client.CreateNodeDiskZFS(ctx, node, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating node disk ZFS pool", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for node disk ZFS pool creation", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", node, plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeDiskZFSResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// no GET for individual ZFS pools via the disks endpoint — keep state as-is
}

func (r *NodeDiskZFSResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all attributes are ForceNew so Update is never called
}

func (r *NodeDiskZFSResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NodeDiskZFSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()
	name := state.Name.ValueString()

	tflog.Debug(ctx, "Deleting node disk ZFS pool", map[string]any{"node": node, "name": name})

	upid, err := r.client.DeleteNodeDiskZFS(ctx, node, name)
	if err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting node disk ZFS pool", err.Error())
		return
	}

	if err := r.client.WaitForUPID(ctx, upid); err != nil {
		resp.Diagnostics.AddError("Error waiting for node disk ZFS pool deletion", err.Error())
	}
}

func (r *NodeDiskZFSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// import format: node_name:name
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID",
			"Format: <node_name>:<name> (e.g. 'pve:rpool')")
		return
	}

	state := NodeDiskZFSResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Name:     types.StringValue(parts[1]),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
