package node_acme_certificate

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

var _ resource.Resource = &NodeACMECertificateResource{}
var _ resource.ResourceWithConfigure = &NodeACMECertificateResource{}
var _ resource.ResourceWithImportState = &NodeACMECertificateResource{}

// NodeACMECertificateResource handles ordering and renewing the ACME TLS cert for a node.
// The node needs an ACME account and domain config in place before this resource can be applied.
type NodeACMECertificateResource struct {
	client *client.Client
}

type NodeACMECertificateResourceModel struct {
	ID       types.String `tfsdk:"id"`
	NodeName types.String `tfsdk:"node_name"`
	// force tells proxmox to order/renew even if the cert isnt due yet
	Force types.Bool `tfsdk:"force"`
}

func NewResource() resource.Resource {
	return &NodeACMECertificateResource{}
}

func (r *NodeACMECertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_environment_node_acme_certificate"
}

func (r *NodeACMECertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Orders or renews an ACME-issued TLS certificate for a Proxmox VE node. " +
			"The node must have an ACME account and domain configuration applied beforehand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The resource identifier (node name).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"node_name": schema.StringAttribute{
				Description: "The name of the Proxmox VE node for which the ACME certificate is ordered.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"force": schema.BoolAttribute{
				Description: "Whether to force a certificate order/renewal even if the certificate is not yet due for renewal.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *NodeACMECertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NodeACMECertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NodeACMECertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()
	plan.ID = types.StringValue(node)

	tflog.Debug(ctx, "Ordering ACME certificate", map[string]any{
		"node":  node,
		"force": plan.Force.ValueBool(),
	})

	upid, err := r.client.OrderNodeACMECertificate(ctx, node, plan.Force.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error ordering ACME certificate", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for ACME certificate order task", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeACMECertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NodeACMECertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()

	certs, err := r.client.GetNodeCertificates(ctx, node)
	if err != nil {
		resp.Diagnostics.AddError("Error reading node certificates", err.Error())
		return
	}

	// check that an ACME-issued cert exists — we look for the pveproxy cert filename
	// as a reliable indicator that ACME is active on this node
	acmeCertFound := false
	for _, cert := range certs {
		if cert.Filename == "pveproxy-ssl.pem" || cert.Filename == "/etc/pve/local/pveproxy-ssl.pem" {
			acmeCertFound = true
			break
		}
	}

	if !acmeCertFound && len(certs) == 0 {
		// nothing there at all, resource is gone
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NodeACMECertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NodeACMECertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := plan.NodeName.ValueString()

	tflog.Debug(ctx, "Renewing ACME certificate", map[string]any{
		"node":  node,
		"force": plan.Force.ValueBool(),
	})

	upid, err := r.client.RenewNodeACMECertificate(ctx, node, plan.Force.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error renewing ACME certificate", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for ACME certificate renewal task", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NodeACMECertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NodeACMECertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	node := state.NodeName.ValueString()

	tflog.Debug(ctx, "Revoking ACME certificate", map[string]any{"node": node})

	upid, err := r.client.RevokeNodeACMECertificate(ctx, node)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.IsNotFound() {
			return
		}
		resp.Diagnostics.AddError("Error revoking ACME certificate", err.Error())
		return
	}

	if upid != "" {
		if err := r.client.WaitForUPID(ctx, upid); err != nil {
			resp.Diagnostics.AddError("Error waiting for ACME certificate revocation task", err.Error())
		}
	}
}

func (r *NodeACMECertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// import uses just the node name
	node := req.ID

	state := NodeACMECertificateResourceModel{
		ID:       types.StringValue(node),
		NodeName: types.StringValue(node),
		Force:    types.BoolValue(false),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
