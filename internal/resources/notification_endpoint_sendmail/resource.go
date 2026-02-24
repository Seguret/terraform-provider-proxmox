package notification_endpoint_sendmail

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &SendmailEndpointResource{}
var _ resource.ResourceWithConfigure = &SendmailEndpointResource{}
var _ resource.ResourceWithImportState = &SendmailEndpointResource{}

type SendmailEndpointResource struct {
	client *client.Client
}

type SendmailEndpointResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Mailto      types.List   `tfsdk:"mailto"`
	MailtoUser  types.List   `tfsdk:"mailto_user"`
	FromAddress types.String `tfsdk:"from_address"`
	Author      types.String `tfsdk:"author"`
	Comment     types.String `tfsdk:"comment"`
	Disable     types.Bool   `tfsdk:"disable"`
}

func NewResource() resource.Resource {
	return &SendmailEndpointResource{}
}

func (r *SendmailEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_notification_endpoint_sendmail"
}

func (r *SendmailEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE sendmail notification endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the sendmail endpoint.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mailto": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of email addresses to send notifications to.",
				Optional:    true,
			},
			"mailto_user": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "List of users to send notifications to (by Proxmox user ID).",
				Optional:    true,
			},
			"from_address": schema.StringAttribute{
				Description: "The sender email address.",
				Optional:    true,
				Computed:    true,
			},
			"author": schema.StringAttribute{
				Description: "The author name used in the notification.",
				Optional:    true,
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the sendmail endpoint.",
				Optional:    true,
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the sendmail endpoint is disabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SendmailEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SendmailEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SendmailEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating sendmail notification endpoint", map[string]any{"name": plan.Name.ValueString()})

	createReq := &models.NotificationEndpointSendmailCreateRequest{
		Name:        plan.Name.ValueString(),
		Mailto:      listToStrings(plan.Mailto),
		MailtoUser:  listToStrings(plan.MailtoUser),
		FromAddress: plan.FromAddress.ValueString(),
		Author:      plan.Author.ValueString(),
		Comment:     plan.Comment.ValueString(),
		Disable:     plan.Disable.ValueBool(),
	}

	if err := r.client.CreateNotificationEndpointSendmail(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating sendmail notification endpoint", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading sendmail notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SendmailEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SendmailEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading sendmail notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SendmailEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SendmailEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disable := plan.Disable.ValueBool()
	updateReq := &models.NotificationEndpointSendmailUpdateRequest{
		Mailto:      listToStrings(plan.Mailto),
		MailtoUser:  listToStrings(plan.MailtoUser),
		FromAddress: plan.FromAddress.ValueString(),
		Author:      plan.Author.ValueString(),
		Comment:     plan.Comment.ValueString(),
		Disable:     &disable,
	}

	if err := r.client.UpdateNotificationEndpointSendmail(ctx, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating sendmail notification endpoint", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading sendmail notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SendmailEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SendmailEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNotificationEndpointSendmail(ctx, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting sendmail notification endpoint", err.Error())
	}
}

func (r *SendmailEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SendmailEndpointResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing sendmail notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SendmailEndpointResource) readIntoModel(ctx context.Context, model *SendmailEndpointResourceModel) error {
	ep, err := r.client.GetNotificationEndpointSendmail(ctx, model.Name.ValueString())
	if err != nil {
		return err
	}
	model.Name = types.StringValue(ep.Name)
	model.Mailto = stringsToList(ep.Mailto)
	model.MailtoUser = stringsToList(ep.MailtoUser)
	model.FromAddress = types.StringValue(ep.FromAddress)
	model.Author = types.StringValue(ep.Author)
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

func stringsToList(strs []string) types.List {
	if len(strs) == 0 {
		return types.ListValueMust(types.StringType, []attr.Value{})
	}
	vals := make([]attr.Value, len(strs))
	for i, s := range strs {
		vals[i] = types.StringValue(s)
	}
	return types.ListValueMust(types.StringType, vals)
}

func listToStrings(list types.List) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var strs []string
	for _, v := range list.Elements() {
		strs = append(strs, v.(types.String).ValueString())
	}
	return strs
}
