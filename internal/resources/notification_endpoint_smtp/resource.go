package notification_endpoint_smtp

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

var _ resource.Resource = &SmtpEndpointResource{}
var _ resource.ResourceWithConfigure = &SmtpEndpointResource{}
var _ resource.ResourceWithImportState = &SmtpEndpointResource{}

type SmtpEndpointResource struct {
	client *client.Client
}

type SmtpEndpointResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Server     types.String `tfsdk:"server"`
	Port       types.Int64  `tfsdk:"port"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	Mode       types.String `tfsdk:"mode"`
	Mailto     types.List   `tfsdk:"mailto"`
	MailtoUser types.List   `tfsdk:"mailto_user"`
	From       types.String `tfsdk:"from"`
	Comment    types.String `tfsdk:"comment"`
	Disable    types.Bool   `tfsdk:"disable"`
}

func NewResource() resource.Resource {
	return &SmtpEndpointResource{}
}

func (r *SmtpEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_notification_endpoint_smtp"
}

func (r *SmtpEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SMTP notification endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the SMTP endpoint.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server": schema.StringAttribute{
				Description: "The SMTP server hostname or IP address.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "The SMTP server port.",
				Optional:    true,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "The SMTP username for authentication.",
				Optional:    true,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "The SMTP password for authentication.",
				Optional:    true,
				Sensitive:   true,
			},
			"mode": schema.StringAttribute{
				Description: "The SMTP encryption mode: insecure, starttls, or tls.",
				Optional:    true,
				Computed:    true,
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
			"from": schema.StringAttribute{
				Description: "The sender email address.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Comment for the SMTP endpoint.",
				Optional:    true,
				Computed:    true,
			},
			"disable": schema.BoolAttribute{
				Description: "Whether the SMTP endpoint is disabled.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SmtpEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SmtpEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SmtpEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SMTP notification endpoint", map[string]any{"name": plan.Name.ValueString()})

	createReq := &models.NotificationEndpointSmtpCreateRequest{
		Name:       plan.Name.ValueString(),
		Server:     plan.Server.ValueString(),
		Port:       int(plan.Port.ValueInt64()),
		Username:   plan.Username.ValueString(),
		Password:   plan.Password.ValueString(),
		Mode:       plan.Mode.ValueString(),
		Mailto:     listToStrings(plan.Mailto),
		MailtoUser: listToStrings(plan.MailtoUser),
		From:       plan.From.ValueString(),
		Comment:    plan.Comment.ValueString(),
		Disable:    plan.Disable.ValueBool(),
	}

	if err := r.client.CreateNotificationEndpointSmtp(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error creating SMTP notification endpoint", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SMTP notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SmtpEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SmtpEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading SMTP notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SmtpEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SmtpEndpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disable := plan.Disable.ValueBool()
	updateReq := &models.NotificationEndpointSmtpUpdateRequest{
		Server:     plan.Server.ValueString(),
		Port:       int(plan.Port.ValueInt64()),
		Username:   plan.Username.ValueString(),
		Password:   plan.Password.ValueString(),
		Mode:       plan.Mode.ValueString(),
		Mailto:     listToStrings(plan.Mailto),
		MailtoUser: listToStrings(plan.MailtoUser),
		From:       plan.From.ValueString(),
		Comment:    plan.Comment.ValueString(),
		Disable:    &disable,
	}

	if err := r.client.UpdateNotificationEndpointSmtp(ctx, plan.Name.ValueString(), updateReq); err != nil {
		resp.Diagnostics.AddError("Error updating SMTP notification endpoint", err.Error())
		return
	}

	if err := r.readIntoModel(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error reading SMTP notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SmtpEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SmtpEndpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteNotificationEndpointSmtp(ctx, state.Name.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting SMTP notification endpoint", err.Error())
	}
}

func (r *SmtpEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SmtpEndpointResourceModel{
		ID:   types.StringValue(req.ID),
		Name: types.StringValue(req.ID),
	}

	if err := r.readIntoModel(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing SMTP notification endpoint", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModel fetches the SMTP endpoint and fills in the model.
// Password is sensitive and not returned by the API — we keep whatever is in state.
func (r *SmtpEndpointResource) readIntoModel(ctx context.Context, model *SmtpEndpointResourceModel) error {
	ep, err := r.client.GetNotificationEndpointSmtp(ctx, model.Name.ValueString())
	if err != nil {
		return err
	}
	model.Name = types.StringValue(ep.Name)
	model.Server = types.StringValue(ep.Server)
	model.Port = types.Int64Value(int64(ep.Port))
	model.Username = types.StringValue(ep.Username)
	// password is sensitive — the API wont return it so we keep what was in state
	model.Mode = types.StringValue(ep.Mode)
	model.Mailto = stringsToList(ep.Mailto)
	model.MailtoUser = stringsToList(ep.MailtoUser)
	model.From = types.StringValue(ep.From)
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
