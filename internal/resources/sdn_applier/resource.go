package sdn_applier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
)

var _ resource.Resource = &SDNApplierResource{}
var _ resource.ResourceWithConfigure = &SDNApplierResource{}

type SDNApplierResource struct {
	client *client.Client
}

type SDNApplierModel struct {
	ID            types.String `tfsdk:"id"`
	KeepUpToDate  types.Bool   `tfsdk:"keep_up_to_date"`
}

func NewResource() resource.Resource {
	return &SDNApplierResource{}
}

func (r *SDNApplierResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_applier"
}

func (r *SDNApplierResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Applies pending Proxmox VE SDN configuration changes. When created or updated, this resource calls the SDN apply endpoint to push pending changes into effect.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keep_up_to_date": schema.BoolAttribute{
				Description: "If true (default), call ApplySDN on every Create and Update to keep SDN applied. Set to false to manage when SDN is applied externally.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *SDNApplierResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNApplierResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNApplierModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.KeepUpToDate.ValueBool() {
		tflog.Debug(ctx, "Applying SDN configuration")
		if err := r.client.ApplySDN(ctx); err != nil {
			resp.Diagnostics.AddError("Error applying SDN", err.Error())
			return
		}
	}

	plan.ID = types.StringValue("sdn")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNApplierResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// SDN pending state cant be read meaningfully — keep state as-is
	var state SDNApplierModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNApplierResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNApplierModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.KeepUpToDate.ValueBool() {
		tflog.Debug(ctx, "Applying SDN configuration (update)")
		if err := r.client.ApplySDN(ctx); err != nil {
			resp.Diagnostics.AddError("Error applying SDN", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNApplierResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// SDN cant be "unapplied" — just remove from state
}
