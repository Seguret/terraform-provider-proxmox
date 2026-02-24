package pool

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &PoolResource{}
var _ resource.ResourceWithConfigure = &PoolResource{}
var _ resource.ResourceWithImportState = &PoolResource{}

type PoolResource struct {
	client *client.Client
}

type PoolResourceModel struct {
	ID      types.String `tfsdk:"id"`
	PoolID  types.String `tfsdk:"pool_id"`
	Comment types.String `tfsdk:"comment"`
}

func NewResource() resource.Resource {
	return &PoolResource{}
}

func (r *PoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_pool"
}

func (r *PoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE resource pool.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool_id": schema.StringAttribute{
				Description: "The pool identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "A comment for the pool.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *PoolResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &models.PoolCreateRequest{
		PoolID:  plan.PoolID.ValueString(),
		Comment: plan.Comment.ValueString(),
	}

	tflog.Debug(ctx, "Creating pool", map[string]any{"pool_id": createReq.PoolID})

	if err := r.client.CreatePool(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating pool", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.PoolID.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PoolResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	comment := plan.Comment.ValueString()
	updateReq := &models.PoolUpdateRequest{Comment: &comment}

	if err := r.client.UpdatePool(ctx, plan.PoolID.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating pool", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.PoolID.ValueString())
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PoolResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeletePool(ctx, state.PoolID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting pool", err.Error())
	}
}

func (r *PoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := PoolResourceModel{
		ID:     types.StringValue(req.ID),
		PoolID: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PoolResource) readIntoModel(ctx context.Context, model *PoolResourceModel, diagnostics *diag.Diagnostics) {
	p, err := r.client.GetPool(ctx, model.PoolID.ValueString())
	if err != nil {
		diagnostics.AddError("Error reading pool", err.Error())
		return
	}
	model.Comment = types.StringValue(p.Comment)
}
