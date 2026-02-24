package notification_endpoint_gotify

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

var _ resource.Resource = &GotifyEndpointResource{}
var _ resource.ResourceWithConfigure = &GotifyEndpointResource{}
var _ resource.ResourceWithImportState = &GotifyEndpointResource{}

type GotifyEndpointResource struct {
	client *client.Client
}

type GotifyEndpointResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Server  types.String `tfsdk:"server"`
	Token   types.String `tfsdk:"token"`
	Comment types.String `tfsdk:"comment"`
	Disable types.Bool   `tfsdk:"disable"`
}

func NewResource() resource.Resource {
	return &GotifyEndpointResource{}
}

func (r *GotifyEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_notification_endpoint_gotify"
}

func (r *GotifyEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE Gotify notification endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the Gotify endpoint.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server": schema.StringAttribute{
				Description: "The Gotify server URL.",
				Required:    true,
			},
			"token": schema.StringAttribute{
				Description: "The Gotify application token.",
				Optional:    true,
				Sensitive:   true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the Gotify endpoint.",
				Optional:    true,
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the Gotify endpoint is disabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *GotifyEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *GotifyEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GotifyEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating Gotify notification endpoint", map[string]any{"name": plan.Name.ValueString()})

	createReq := &models.NotificationEndpointGotifyCreateRequest{
		Name:    plan.Name.ValueString(),
		Server:  plan.Server.ValueString(),
		Token:   plan.Token.ValueString(),
		Comment: plan.Comment.ValueString(),
		Disable: plan.Disable.ValueBool(),
	}

	if err := r.client.CreateNotificationEndpointGotify(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating Gotify notification endpoint", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading Gotify notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GotifyEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GotifyEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading Gotify notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GotifyEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GotifyEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disable := plan.Disable.ValueBool()
	updateReq := &models.NotificationEndpointGotifyUpdateRequest{
		Server:  plan.Server.ValueString(),
		Token:   plan.Token.ValueString(),
		Comment: plan.Comment.ValueString(),
		Disable: &disable,
	}

	if err := r.client.UpdateNotificationEndpointGotify(ctx, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating Gotify notification endpoint", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading Gotify notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GotifyEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GotifyEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNotificationEndpointGotify(ctx, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Gotify notification endpoint", err.Error())
	}
}

func (r *GotifyEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := GotifyEndpointResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing Gotify notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModel fetches the Gotify endpoint and fills in the model.
// Token is sensitive and not returned by the API — we keep whatever is in state.
func (r *GotifyEndpointResource) readIntoModel(ctx context.Context, model *GotifyEndpointResourceModel) error {
	ep, err := r.client.GetNotificationEndpointGotify(ctx, model.Name.ValueString())
	if err != nil {
		return err
	}
	model.Name = types.StringValue(ep.Name)
	model.Server = types.StringValue(ep.Server)
	// token is sensitive — the API wont return it so we keep what was in state
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
