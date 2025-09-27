package devops

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-devops/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &devopsResource{}
	_ resource.ResourceWithConfigure = &devopsResource{}
)

// NewDevOpsResource is a helper function to simplify the provider implementation.
func NewDevOpsResource() resource.Resource { return &devopsResource{} }

// devopsResource is the resource implementation.
type devopsResource struct{
	client *client.Client
}

// Metadata returns the resource type name.
func (r *devopsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devops"
}

// Schema defines the schema for the resource.
func (r *devopsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			// Input: list of dev IDs
			"devs": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			// Input: list of ops IDs
			"ops": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *devopsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan devopsModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() { return }

	// Convert devs and ops lists (types.List of string IDs) to []string
	var devIDs []string
	diags = plan.Devs.ElementsAs(ctx, &devIDs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() { return }

	var opsIDs []string
	diags = plan.Ops.ElementsAs(ctx, &opsIDs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() { return }

	// Build slices of minimal objects with only ID set
	devObjs := make([]client.Dev, 0, len(devIDs))
	for _, id := range devIDs { devObjs = append(devObjs, client.Dev{ID: id}) }
	opsObjs := make([]client.Ops, 0, len(opsIDs))
	for _, id := range opsIDs { opsObjs = append(opsObjs, client.Ops{ID: id}) }

	reqDevOps := client.DevOps{ Dev: devObjs, Ops: opsObjs }

	created, err := r.client.CreateDevops(reqDevOps)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating DevOps",
			"Could not create DevOps, unexpected error: " + err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(created.ID)
	// keep devs/ops lists as provided
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *devopsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state devopsModel

    // Load current state to get the ID of this resource instance
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Fetch DevOps by ID
    found, err := r.client.GetDevOpsByID(state.ID.ValueString())
    if err != nil {
        // If the backend returns 404, treat as already deleted
        if strings.Contains(err.Error(), "status: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error Reading DevOps",
            "Could not read DevOps ID "+state.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // If the DevOps is not found, remove from state (resource drift)
    if found == nil {
        resp.State.RemoveResource(ctx)
        return
    }

    // Map found devops to state
    state.ID = types.StringValue(found.ID)
    // Extract IDs from []client.Dev and []client.Ops
    devIDs := make([]string, 0, len(found.Dev))
    for _, d := range found.Dev { devIDs = append(devIDs, d.ID) }
    devList, d1 := types.ListValueFrom(ctx, types.StringType, devIDs)
    resp.Diagnostics.Append(d1...)
    if resp.Diagnostics.HasError() { return }
    state.Devs = devList

    opsIDs := make([]string, 0, len(found.Ops))
    for _, o := range found.Ops { opsIDs = append(opsIDs, o.ID) }
    opsList, d2 := types.ListValueFrom(ctx, types.StringType, opsIDs)
    resp.Diagnostics.Append(d2...)
    if resp.Diagnostics.HasError() { return }
    state.Ops = opsList

    // Set state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *devopsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan devopsModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() { return }

    // Load current state to get the persisted ID (plan.ID may be unknown during update)
    var state devopsModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() { return }

    // Build request from plan
    var devIDs []string
    diags = plan.Devs.ElementsAs(ctx, &devIDs, false)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() { return }

    var opsIDs []string
    diags = plan.Ops.ElementsAs(ctx, &opsIDs, false)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() { return }

    devObjs := make([]client.Dev, 0, len(devIDs))
    for _, id := range devIDs { devObjs = append(devObjs, client.Dev{ID: id}) }
    opsObjs := make([]client.Ops, 0, len(opsIDs))
    for _, id := range opsIDs { opsObjs = append(opsObjs, client.Ops{ID: id}) }
    reqDevOps := client.DevOps{ Dev: devObjs, Ops: opsObjs }

    // Update existing devops by ID from state
    _, err := r.client.UpdateDevOps(state.ID.ValueString(), reqDevOps)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Updating DevOps",
            "Could not update DevOps ID "+state.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // Fetch updated DevOps by ID
    updated, err := r.client.GetDevOpsByID(state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading DevOps",
            "Could not read DevOps ID "+state.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // Update resource state
    plan.ID = types.StringValue(updated.ID)
    devIDs = make([]string, 0, len(updated.Dev))
    for _, d := range updated.Dev { devIDs = append(devIDs, d.ID) }
    devList, d1 := types.ListValueFrom(ctx, types.StringType, devIDs)
    resp.Diagnostics.Append(d1...)
    if resp.Diagnostics.HasError() { return }
    plan.Devs = devList

    opsIDs = make([]string, 0, len(updated.Ops))
    for _, o := range updated.Ops { opsIDs = append(opsIDs, o.ID) }
    opsList, d2 := types.ListValueFrom(ctx, types.StringType, opsIDs)
    resp.Diagnostics.Append(d2...)
    if resp.Diagnostics.HasError() { return }
    plan.Ops = opsList

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() { return }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *devopsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state devopsModel
    diags := req.State.Get(ctx, &state)

    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    err := r.client.DeleteDevOps(state.ID.ValueString())
    if err != nil {
        // If the backend returns 404, treat as already deleted
        if strings.Contains(err.Error(), "status: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error Deleting DevOps",
            fmt.Sprintf("Could not delete DevOps ID %s: %v", state.ID.ValueString(), err),
        )
        return
    }
    // Successful delete; nothing else to do. Terraform will drop state for this resource.
}

func (r *devopsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}