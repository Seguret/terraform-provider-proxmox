package ceph_pool

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &CephPoolResource{}
var _ resource.ResourceWithConfigure = &CephPoolResource{}
var _ resource.ResourceWithImportState = &CephPoolResource{}

type CephPoolResource struct {
	client *client.Client
}

type CephPoolResourceModel struct {
	ID              types.String `tfsdk:"id"`
	NodeName        types.String `tfsdk:"node_name"`
	Name            types.String `tfsdk:"name"`
	Size            types.Int64  `tfsdk:"size"`
	MinSize         types.Int64  `tfsdk:"min_size"`
	PGNum           types.Int64  `tfsdk:"pg_num"`
	PGAutoscaleMode types.String `tfsdk:"pg_autoscale_mode"`
	Application     types.String `tfsdk:"application"`
	CrushRule       types.String `tfsdk:"crush_rule"`
	AddStorages     types.Bool   `tfsdk:"add_storages"`
}

func NewResource() resource.Resource {
	return &CephPoolResource{}
}

func (r *CephPoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ceph_pool"
}

func (r *CephPoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Ceph pool on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The Proxmox node on which to manage the Ceph pool.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Ceph pool.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "Number of replicas.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(3),
			},
			"min_size": schema.Int64Attribute{
				Description: "Minimum number of replicas for I/O.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(2),
			},
			"pg_num": schema.Int64Attribute{
				Description: "Number of placement groups.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(128),
			},
			"pg_autoscale_mode": schema.StringAttribute{
				Description: "PG autoscale mode (on, off, warn).",
				Optional:    true,
				Computed:    true,
			},
			"application": schema.StringAttribute{
				Description: "Pool application (rbd, cephfs, rgw).",
				Optional:    true,
				Computed:    true,
			},
			"crush_rule": schema.StringAttribute{
				Description: "CRUSH rule name.",
				Optional:    true,
				Computed:    true,
			},
			"add_storages": schema.BoolAttribute{
				Description: "Whether to add a storage entry in the Proxmox config (creation only).",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CephPoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CephPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CephPoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	size := int(plan.Size.ValueInt64())
	minSize := int(plan.MinSize.ValueInt64())
	pgNum := int(plan.PGNum.ValueInt64())

	createReq := &models.CephPoolCreateRequest{
		Name:            plan.Name.ValueString(),
		Size:            &size,
		MinSize:         &minSize,
		PGNum:           &pgNum,
		PGAutoscaleMode: plan.PGAutoscaleMode.ValueString(),
		Application:     plan.Application.ValueString(),
		CrushRule:       plan.CrushRule.ValueString(),
	}

	if !plan.AddStorages.IsNull() && !plan.AddStorages.IsUnknown() && plan.AddStorages.ValueBool() {
		v := 1
		createReq.AddStorages = &v
	}

	tflog.Debug(ctx, "Creating Ceph pool", map[string]any{
		"node": plan.NodeName.ValueString(),
		"name": plan.Name.ValueString(),
	})

	if err := r.client.CreateCephPool(ctx, plan.NodeName.ValueString(), createReq); err != nil {
		resp.Diagnostics.AddError("Error creating Ceph pool", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + "/" + plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading Ceph pool after create", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CephPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CephPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Ceph pool", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CephPoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	size := int(plan.Size.ValueInt64())
	minSize := int(plan.MinSize.ValueInt64())
	pgNum := int(plan.PGNum.ValueInt64())

	updateReq := &models.CephPoolUpdateRequest{
		Size:            &size,
		MinSize:         &minSize,
		PGNum:           &pgNum,
		PGAutoscaleMode: plan.PGAutoscaleMode.ValueString(),
		Application:     plan.Application.ValueString(),
		CrushRule:       plan.CrushRule.ValueString(),
	}

	if err := r.client.UpdateCephPool(ctx, plan.NodeName.ValueString(), plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating Ceph pool", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading Ceph pool after update", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CephPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CephPoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCephPool(ctx, state.NodeName.ValueString(), state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Ceph pool", err.Error())
	}
}

func (r *CephPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: {node_name}/{pool_name}")
		return
	}

	state := CephPoolResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		Name:     types.StringValue(parts[1]),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing Ceph pool", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephPoolResource) readIntoModel(ctx context.Context, model *CephPoolResourceModel) error {
	pool, err := r.client.GetCephPool(ctx, model.NodeName.ValueString(), model.Name.ValueString())
	if err != nil {
		return err
	}
	model.Name = types.StringValue(pool.Name)
	model.Size = types.Int64Value(int64(pool.Size))
	model.MinSize = types.Int64Value(int64(pool.MinSize))
	model.PGNum = types.Int64Value(int64(pool.PGNum))
	model.PGAutoscaleMode = types.StringValue(pool.PGAutoscaleMode)
	model.CrushRule = types.StringValue(pool.CrushRule)

	// pull the application name out of the ApplicationMeta map keys
	if len(pool.ApplicationMeta) > 0 {
		for k := range pool.ApplicationMeta {
			model.Application = types.StringValue(k)
			break
		}
	}
	return nil
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
