package dns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Seguret/terraform-provider-proxmox/internal/client"
	"github.com/Seguret/terraform-provider-proxmox/internal/client/models"
)

var _ resource.Resource = &DNSResource{}
var _ resource.ResourceWithConfigure = &DNSResource{}
var _ resource.ResourceWithImportState = &DNSResource{}

type DNSResource struct {
	client *client.Client
}

type DNSResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	Search   types.String `tfsdk:"search"`
	DNS1     types.String `tfsdk:"dns1"`
	DNS2     types.String `tfsdk:"dns2"`
	DNS3     types.String `tfsdk:"dns3"`
}

func NewResource() resource.Resource {
	return &DNSResource{}
}

func (r *DNSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_dns"
}

func (r *DNSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the DNS configuration of a Proxmox VE node.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The node name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"search": schema.StringAttribute{
				Description: "The DNS search domain.",
				Required:    true,
			},
			"dns1": schema.StringAttribute{
				Description: "The first DNS server.",
				Optional:    true,
				Computed:    true,
			},
			"dns2": schema.StringAttribute{
				Description: "The second DNS server.",
				Optional:    true,
				Computed:    true,
			},
			"dns3": schema.StringAttribute{
				Description: "The third DNS server.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *DNSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSResource) applyAndRead(ctx context.Context, model *DNSResourceModel) error {
	node := model.NodeName.ValueString()

	updateReq := &models.NodeDNSUpdateRequest{
		Search: model.Search.ValueString(),
		DNS1:   model.DNS1.ValueString(),
		DNS2:   model.DNS2.ValueString(),
		DNS3:   model.DNS3.ValueString(),
	}

	if err := r.client.UpdateNodeDNS(ctx, node, updateReq); err != nil {
		return fmt.Errorf("error setting node DNS: %w", err)
	}

	return r.readState(ctx, model)
}

func (r *DNSResource) readState(ctx context.Context, model *DNSResourceModel) error {
	dns, err := r.client.GetNodeDNS(ctx, model.NodeName.ValueString())
	if err != nil {
		return fmt.Errorf("error reading node DNS: %w", err)
	}
	model.Search = types.StringValue(dns.Search)
	model.DNS1 = types.StringValue(dns.DNS1)
	model.DNS2 = types.StringValue(dns.DNS2)
	model.DNS3 = types.StringValue(dns.DNS3)
	return nil
}

func (r *DNSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.NodeName.ValueString())

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error creating DNS config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error reading DNS config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.applyAndRead(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating DNS config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// DNS config cant be deleted from proxmox, just drop it from state
}

func (r *DNSResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	node := req.ID
	state := DNSResourceModel{
		ID:       types.StringValue(node),
		NodeName: types.StringValue(node),
	}

	if err := r.readState(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Error importing DNS config", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
