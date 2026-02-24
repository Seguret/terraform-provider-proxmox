package notification_endpoint_webhook

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

var _ resource.Resource = &WebhookEndpointResource{}
var _ resource.ResourceWithConfigure = &WebhookEndpointResource{}
var _ resource.ResourceWithImportState = &WebhookEndpointResource{}

type WebhookEndpointResource struct {
	client *client.Client
}

type WebhookEndpointResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	URL     types.String `tfsdk:"url"`
	Method  types.String `tfsdk:"method"`
	Comment types.String `tfsdk:"comment"`
	Disable types.Bool   `tfsdk:"disable"`
}

func NewResource() resource.Resource {
	return &WebhookEndpointResource{}
}

func (r *WebhookEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_notification_endpoint_webhook"
}

func (r *WebhookEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE webhook notification endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the webhook endpoint.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The webhook URL to send notifications to.",
				Required:    true,
			},
			"method": schema.StringAttribute{
				Description: "The HTTP method used for the webhook request: GET, POST, or PUT.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the webhook endpoint.",
				Optional:    true,
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the webhook endpoint is disabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *WebhookEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebhookEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebhookEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating webhook notification endpoint", map[string]any{"name": plan.Name.ValueString()})

	createReq := &models.NotificationEndpointWebhookCreateRequest{
		Name:    plan.Name.ValueString(),
		URL:     plan.URL.ValueString(),
		Method:  plan.Method.ValueString(),
		Comment: plan.Comment.ValueString(),
		Disable: plan.Disable.ValueBool(),
	}

	if err := r.client.CreateNotificationEndpointWebhook(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating webhook notification endpoint", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading webhook notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebhookEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebhookEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading webhook notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebhookEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WebhookEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disable := plan.Disable.ValueBool()
	updateReq := &models.NotificationEndpointWebhookUpdateRequest{
		URL:     plan.URL.ValueString(),
		Method:  plan.Method.ValueString(),
		Comment: plan.Comment.ValueString(),
		Disable: &disable,
	}

	if err := r.client.UpdateNotificationEndpointWebhook(ctx, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating webhook notification endpoint", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading webhook notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebhookEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebhookEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNotificationEndpointWebhook(ctx, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting webhook notification endpoint", err.Error())
	}
}

func (r *WebhookEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := WebhookEndpointResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing webhook notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebhookEndpointResource) readIntoModel(ctx context.Context, model *WebhookEndpointResourceModel) error {
	ep, err := r.client.GetNotificationEndpointWebhook(ctx, model.Name.ValueString())
	if err != nil {
		return err
	}
	model.Name = types.StringValue(ep.Name)
	model.URL = types.StringValue(ep.URL)
	model.Method = types.StringValue(ep.Method)
	model.Comment = types.StringValue(ep.Comment)
	model.Disable = types.BoolValue(ep.Disable)
	return nil
}

func isNotFound(err error) bool {
	if apiErr, ok := err.(*client.APIError); ok {
		return apiErr.IsNotFound()
	}
	return false
}
