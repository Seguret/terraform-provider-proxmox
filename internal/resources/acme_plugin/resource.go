package acme_plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &ACMEPluginResource{}
var _ resource.ResourceWithConfigure = &ACMEPluginResource{}
var _ resource.ResourceWithImportState = &ACMEPluginResource{}

type ACMEPluginResource struct {
	client *client.Client
}

type ACMEPluginResourceModel struct {
	ID              types.String `tfsdk:"id"`
	PluginID        types.String `tfsdk:"plugin_id"`
	Type            types.String `tfsdk:"type"`
	API             types.String `tfsdk:"api"`
	Data            types.String `tfsdk:"data"`
	Nodes           types.String `tfsdk:"nodes"`
	ValidationDelay types.Int64  `tfsdk:"validation_delay"`
}

func NewResource() resource.Resource {
	return &ACMEPluginResource{}
}

func (r *ACMEPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_acme_plugin"
}

func (r *ACMEPluginResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE ACME DNS challenge plugin.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plugin_id": schema.StringAttribute{
				Description: "The plugin identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Plugin type ('dns' or 'standalone').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"api": schema.StringAttribute{
				Description: "DNS API name (e.g. 'cf' for Cloudflare, 'aws' for Route53).",
				Optional:    true,
				Computed:    true,
			},
			"data": schema.StringAttribute{
				Description: "DNS API credentials in key=value format (one per line). Sensitive.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"nodes": schema.StringAttribute{
				Description: "Comma-separated list of nodes to restrict the plugin to.",
				Optional:    true,
				Computed:    true,
			},
			"validation_delay": schema.Int64Attribute{
				Description: "Delay in seconds to wait for DNS propagation.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *ACMEPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ACMEPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ACMEPluginResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating ACME plugin", map[string]any{"id": plan.PluginID.ValueString()})

	if err := r.client.CreateACMEPlugin(ctx, &models.ACMEPluginCreateRequest{
		ID:              plan.PluginID.ValueString(),
		Type:            plan.Type.ValueString(),
		API:             plan.API.ValueString(),
		Data:            plan.Data.ValueString(),
		Nodes:           plan.Nodes.ValueString(),
		ValidationDelay: int(plan.ValidationDelay.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating ACME plugin", err.Error())
		return
	}

	plan.ID = plan.PluginID
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACMEPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ACMEPluginResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACMEPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ACMEPluginResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateACMEPlugin(ctx, plan.PluginID.ValueString(), &models.ACMEPluginUpdateRequest{
		API:             plan.API.ValueString(),
		Data:            plan.Data.ValueString(),
		Nodes:           plan.Nodes.ValueString(),
		ValidationDelay: int(plan.ValidationDelay.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating ACME plugin", err.Error())
		return
	}

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ACMEPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ACMEPluginResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteACMEPlugin(ctx, state.PluginID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting ACME plugin", err.Error())
	}
}

func (r *ACMEPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := ACMEPluginResourceModel{
		ID:       types.StringValue(req.ID),
		PluginID: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ACMEPluginResource) readIntoModel(ctx context.Context, model *ACMEPluginResourceModel, diagnostics interface{ AddError(string, string) }) {
	plugin, err := r.client.GetACMEPlugin(ctx, model.PluginID.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diagnostics.AddError("ACME plugin not found", "The ACME plugin no longer exists.")
			return
		}
		diagnostics.AddError("Error reading ACME plugin", err.Error())
		return
	}

	model.Type = types.StringValue(plugin.Type)
	model.API = types.StringValue(plugin.API)
	// data is sensitive — only overwrite if the API actually returned something
	if plugin.Data != "" {
		model.Data = types.StringValue(plugin.Data)
	}
	model.Nodes = types.StringValue(plugin.Nodes)
	model.ValidationDelay = types.Int64Value(int64(plugin.ValidationDelay))
}
