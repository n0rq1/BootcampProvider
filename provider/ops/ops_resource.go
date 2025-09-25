package ops

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
	_ resource.Resource = &opsResource{}
	_ resource.ResourceWithConfigure = &opsResource{}
)

// NewOpsResource is a helper function to simplify the provider implementation.
func NewOpsResource() resource.Resource { return &opsResource{} }

// opsResource is the resource implementation.
type opsResource struct{
	client *client.Client
}

// Metadata returns the resource type name.
func (r *opsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ops"
}

// Schema defines the schema for the resource.
func (r *opsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			// Input: name
			"name": schema.StringAttribute{
				Required: true,
			},
			// Input: list of engineer IDs
			"engineers": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *opsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan opsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() { return }

	// Convert engineers list (types.List of string IDs) to []client.Engineer
	var engineerIDs []string
	diags = plan.Engineers.ElementsAs(ctx, &engineerIDs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() { return }

	engs := make([]client.Engineer, 0, len(engineerIDs))
	for _, id := range engineerIDs {
		engs = append(engs, client.Engineer{ID: id})
	}

	reqOps := client.Ops{
		Name:      plan.Name.ValueString(),
		Engineers: engs,
	}

	created, err := r.client.CreateOps(reqOps)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Ops",
			"Could not create Ops group, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(created.ID)
	plan.Name = types.StringValue(created.Name)
	// keep engineers list as provided
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *opsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state opsResourceModel

    // Load current state to get the ID of this resource instance
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Fetch Ops by ID
    found, err := r.client.GetOpsByID(state.ID.ValueString())
    if err != nil {
        // If the backend returns 404, treat as already deleted
        if strings.Contains(err.Error(), "status: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error Reading Ops",
            "Could not read Ops ID "+state.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // If the Ops is not found, remove from state (resource drift)
    if found == nil {
        resp.State.RemoveResource(ctx)
        return
    }

    // Map found ops to state
    state.ID = types.StringValue(found.ID)
    state.Name = types.StringValue(found.Name)

    // Convert engineers to a list of engineer IDs as the schema expects list(string)
    engineerIDs := make([]string, 0, len(found.Engineers))
    for _, eng := range found.Engineers {
        engineerIDs = append(engineerIDs, eng.ID)
    }
    engList, diags2 := types.ListValueFrom(ctx, types.StringType, engineerIDs)
    resp.Diagnostics.Append(diags2...)
    if resp.Diagnostics.HasError() {
        return
    }
    state.Engineers = engList

    // Set state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}


// Update updates the resource and sets the updated Terraform state on success.
func (r *opsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan opsResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Load current state to get the persisted ID (plan.ID may be unknown during update)
    var state opsResourceModel
    diags = req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan using []client.Engineer
    var reqOps = client.Ops{
        Name: plan.Name.ValueString(),
    }
    var engineerIDs []string
    diags = plan.Engineers.ElementsAs(ctx, &engineerIDs, false)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    engs := make([]client.Engineer, 0, len(engineerIDs))
    for _, id := range engineerIDs {
        engs = append(engs, client.Engineer{ID: id})
    }
    reqOps.Engineers = engs

    // Update existing ops by ID from state
    _, err := r.client.UpdateOps(state.ID.ValueString(), reqOps)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Updating Ops",
            "Could not update Ops ID "+state.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // Fetch updated Ops by ID
    ops, err := r.client.GetOpsByID(state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading Ops",
            "Could not read Ops ID "+state.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // Update resource state with updated items and timestamp
    plan.ID = types.StringValue(ops.ID)
    plan.Name = types.StringValue(ops.Name)
    // map engineers back into list(string) of IDs
    updatedEngineerIDs := make([]string, 0, len(ops.Engineers))
    for _, eng := range ops.Engineers {
        updatedEngineerIDs = append(updatedEngineerIDs, eng.ID)
    }
    engList, diags2 := types.ListValueFrom(ctx, types.StringType, updatedEngineerIDs)
    resp.Diagnostics.Append(diags2...)
    if resp.Diagnostics.HasError() {
        return
    }
    plan.Engineers = engList

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *opsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state opsResourceModel
    diags := req.State.Get(ctx, &state)

    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    err := r.client.DeleteOps(state.ID.ValueString())
    if err != nil {
        // If the backend returns 404, treat as already deleted
        if strings.Contains(err.Error(), "status: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError(
            "Error Deleting Ops",
            fmt.Sprintf("Could not delete Ops ID %s: %v", state.ID.ValueString(), err),
        )
        return
    }
    // Successful delete; nothing else to do. Terraform will drop state for this resource.
}

func (r *opsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// opsResourceModel maps the resource schema data.
type opsResourceModel struct {
	ID          types.String    `tfsdk:"id"`
	Name        types.String    `tfsdk:"name"`
	Engineers   types.List      `tfsdk:"engineers"`
}
