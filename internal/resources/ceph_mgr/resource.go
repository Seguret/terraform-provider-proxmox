package ceph_mgr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ resource.Resource = &CephMGRResource{}
var _ resource.ResourceWithConfigure = &CephMGRResource{}
var _ resource.ResourceWithImportState = &CephMGRResource{}

type CephMGRResource struct {
	client *client.Client
}

type CephMGRResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	MgrID    types.String `tfsdk:"mgr_id"`
	State    types.String `tfsdk:"state"`
	Addr     types.String `tfsdk:"addr"`
}

func NewResource() resource.Resource {
	return &CephMGRResource{}
}

func (r *CephMGRResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_ceph_mgr"
}

func (r *CephMGRResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Ceph MGR daemon on a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The Proxmox node on which to create the MGR.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mgr_id": schema.StringAttribute{
				Description: "The MGR daemon ID (typically the node name).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The MGR daemon state (active or standby).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"addr": schema.StringAttribute{
				Description: "The MGR daemon address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CephMGRResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CephMGRResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CephMGRResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Ceph MGR", map[string]any{
		"node":   plan.NodeName.ValueString(),
		"mgr_id": plan.MgrID.ValueString(),
	})

	if err := r.client.CreateCephMGR(ctx, plan.NodeName.ValueString(), plan.MgrID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error creating Ceph MGR", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString() + "/" + plan.MgrID.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		// MGR might not show up right away — just zero these out
		plan.State = types.StringValue("")
		plan.Addr = types.StringValue("")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CephMGRResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CephMGRResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Ceph MGR", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephMGRResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// all fields are ForceNew so Update is never actually called
}

func (r *CephMGRResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CephMGRResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteCephMGR(ctx, state.NodeName.ValueString(), state.MgrID.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Ceph MGR", err.Error())
	}
}

func (r *CephMGRResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			"Expected format: {node_name}/{mgr_id}")
		return
	}

	state := CephMGRResourceModel{
		ID:       types.StringValue(req.ID),
		NodeName: types.StringValue(parts[0]),
		MgrID:    types.StringValue(parts[1]),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing Ceph MGR", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CephMGRResource) readIntoModel(ctx context.Context, model *CephMGRResourceModel) error {
	mgrs, err := r.client.GetCephMGRList(ctx, model.NodeName.ValueString())
	if err != nil {
		return err
	}

	mgrID := model.MgrID.ValueString()

	// check if this is the active MGR
	if mgrs.Active != nil && mgrs.Active.ID == mgrID {
		model.State = types.StringValue("active")
		model.Addr = types.StringValue(mgrs.Active.Addr)
		return nil
	}

	// check standbys
	for _, m := range mgrs.Standbys {
		if m.ID == mgrID {
			model.State = types.StringValue("standby")
			model.Addr = types.StringValue(m.Addr)
			return nil
		}
	}

	return &client.APIError{StatusCode: 404, Status: "404 Not Found"}
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
