package ceph_osd

import (
	"context"
	"fmt"
	"strconv"
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

var _ resource.Resource = &CephOSDResource{}
var _ resource.ResourceWithConfigure = &CephOSDResource{}
var _ resource.ResourceWithImportState = &CephOSDResource{}

type CephOSDResource struct {
	client *client.Client
}

type CephOSDResourceModel struct {
	ID        types.String `tfsdk:"id"`
	NodeName  types.String `tfsdk:"node_name"`
	Dev       types.String `tfsdk:"dev"`
	Encrypted types.Bool   `tfsdk:"encrypted"`
	DBDev     types.String `tfsdk:"db_dev"`
	WALDev    types.String `tfsdk:"wal_dev"`
	OSDID     types.Int64  `tfsdk:"osd_id"`
}

func NewResource() resource.Resource {
	return &CephOSDResource{}
}

func (r *CephOSDResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ceph_osd"
}

func (r *CephOSDResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Ceph OSD on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The Proxmox node on which to manage the OSD.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dev": schema.StringAttribute{
				Description: "The block device path (e.g., /dev/sdb).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"encrypted": schema.BoolAttribute{
				Description: "Whether to encrypt the OSD.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"db_dev": schema.StringAttribute{
				Description: "Block device for the OSD DB.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"wal_dev": schema.StringAttribute{
				Description: "Block device for the OSD WAL.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"osd_id": schema.Int64Attribute{
				Description: "The OSD ID assigned by Ceph.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CephOSDResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CephOSDResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CephOSDResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()

	// snapshot existing OSD IDs so we can find the new one afterward
	before, err := r.client.GetCephOSDs(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error listing OSDs before creation", err.Error())
		return
	}
	existingIDs := map[int]struct{}{}
	for _, n := range before.Nodes {
		if n.Type == "osd" {
			existingIDs[n.ID] = struct{}{}
		}
	}

	createReq := &models.CephOSDCreateRequest{
		Dev:    plan.Dev.ValueString(),
		DBDev:  plan.DBDev.ValueString(),
		WALDev: plan.WALDev.ValueString(),
	}
	if !plan.Encrypted.IsNull() && !plan.Encrypted.IsUnknown() && plan.Encrypted.ValueBool() {
		v := 1
		createReq.Encrypted = &v
	}

	tflog.Debug(ctx, "Creating Ceph OSD", map[string]any{
		"node": node,
		"dev":  plan.Dev.ValueString(),
	})

	if _, err := r.client.CreateCephOSD(ctx, node, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating Ceph OSD", err.Error())
		return
	}

	// find which OSD ID was just created by diffing before/after
	after, err := r.client.GetCephOSDs(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error listing OSDs after creation", err.Error())
		return
	}

	newOSDID := -1
	for _, n := range after.Nodes {
		if n.Type == "osd" {
			if _, existed := existingIDs[n.ID]; !existed {
				newOSDID = n.ID
				break
			}
		}
	}

	if newOSDID < 0 {
		resp.Diagnostics.AddError("Error identifying new OSD", "Could not find new OSD ID after creation")
		return
	}

	plan.OSDID = types.Int64Value(int64(newOSDID))
	plan.ID = types.StringValue(fmt.Sprintf("%s/%d", node, newOSDID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CephOSDResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CephOSDResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	osdID := int(state.OSDID.ValueInt64())
	osds, err := r.client.GetCephOSDs(ctx, state.NodeName.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Ceph OSDs", err.Error())
		return
	}

	found := false
	for _, n := range osds.Nodes {
		if n.Type == "osd" && n.ID == osdID {
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephOSDResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all fields are ForceNew so Update is never actually called
}

func (r *CephOSDResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CephOSDResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	osdID := int(state.OSDID.ValueInt64())
	if err := r.client.DeleteCephOSD(ctx, state.NodeName.ValueString(), osdID); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Ceph OSD", err.Error())
	}
}

func (r *CephOSDResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: {node_name}/{osd_id}")
		return
	}

	osdID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid OSD ID", fmt.Sprintf("Could not parse OSD ID: %s", err))
		return
	}

	state := CephOSDResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		OSDID:    types.Int64Value(osdID),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
