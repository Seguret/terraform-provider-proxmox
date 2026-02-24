package sdn_dns

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

var _ resource.Resource = &SDNDnsResource{}
var _ resource.ResourceWithConfigure = &SDNDnsResource{}
var _ resource.ResourceWithImportState = &SDNDnsResource{}

type SDNDnsResource struct {
	client *client.Client
}

type SDNDnsResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Dns  types.String `tfsdk:"dns"`
	Type types.String `tfsdk:"type"`
	URL  types.String `tfsdk:"url"`
	Key  types.String `tfsdk:"key"`
	TTL  types.Int64  `tfsdk:"ttl"`
}

func NewResource() resource.Resource {
	return &SDNDnsResource{}
}

func (r *SDNDnsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_sdn_dns"
}

func (r *SDNDnsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Proxmox VE SDN DNS provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dns": schema.StringAttribute{
				Description: "The SDN DNS provider name (identifier).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The DNS provider plugin type (e.g. powerdns).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The DNS server API URL.",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "The DNS server API key. Sensitive.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
			},
			"ttl": schema.Int64Attribute{
				Description: "Default TTL value for DNS records.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *SDNDnsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SDNDnsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SDNDnsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating SDN DNS provider", map[string]any{
		"dns":  plan.Dns.ValueString(),
		"type": plan.Type.ValueString(),
	})

	if err := r.client.CreateSDNDns(ctx, &models.SDNDnsCreateRequest{
		Dns:  plan.Dns.ValueString(),
		Type: plan.Type.ValueString(),
		URL:  plan.URL.ValueString(),
		Key:  plan.Key.ValueString(),
		TTL:  int(plan.TTL.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error creating SDN DNS provider", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	plan.ID = plan.Dns
	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNDnsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SDNDnsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if isNotFound := r.readIntoModelWithNotFound(ctx, &state, &resp.Diagnostics); isNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SDNDnsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SDNDnsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSDNDns(ctx, plan.Dns.ValueString(), &models.SDNDnsUpdateRequest{
		URL: plan.URL.ValueString(),
		Key: plan.Key.ValueString(),
		TTL: int(plan.TTL.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Error updating SDN DNS provider", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)

	r.readIntoModel(ctx, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SDNDnsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SDNDnsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSDNDns(ctx, state.Dns.ValueString()); err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error deleting SDN DNS provider", err.Error())
		return
	}

	_ = r.client.ApplySDN(ctx)
}

func (r *SDNDnsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	state := SDNDnsResourceModel{
		ID:  types.StringValue(req.ID),
		Dns: types.StringValue(req.ID),
	}
	r.readIntoModel(ctx, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// readIntoModelWithNotFound fetches DNS provider state; returns true if it doesnt exist.
func (r *SDNDnsResource) readIntoModelWithNotFound(ctx context.Context, model *SDNDnsResourceModel, diags interface{ AddError(string, string) }) bool {
	dns, err := r.client.GetSDNDns(ctx, model.Dns.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return true
		}
		diags.AddError("Error reading SDN DNS provider", err.Error())
		return false
	}
	populateModel(model, dns)
	return false
}

// readIntoModel fetches DNS provider state and adds diagnostic errors on failure.
func (r *SDNDnsResource) readIntoModel(ctx context.Context, model *SDNDnsResourceModel, diags interface{ AddError(string, string) }) {
	dns, err := r.client.GetSDNDns(ctx, model.Dns.ValueString())
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			diags.AddError("SDN DNS provider not found", "The SDN DNS provider no longer exists.")
			return
		}
		diags.AddError("Error reading SDN DNS provider", err.Error())
		return
	}
	populateModel(model, dns)
}

func populateModel(model *SDNDnsResourceModel, dns *models.SDNDns) {
	model.Type = types.StringValue(dns.Type)
	model.URL = types.StringValue(dns.URL)
	// key is sensitive — only overwrite if the API actually returned something
	if dns.Key != "" {
		model.Key = types.StringValue(dns.Key)
	}
	model.TTL = types.Int64Value(int64(dns.TTL))
}
