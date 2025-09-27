package engineers

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
	_ resource.Resource              = &EngineerResource{}
	_ resource.ResourceWithConfigure = &EngineerResource{}
)

// NewEngineerResource is a helper function to simplify the provider implementation.
func NewEngineerResource() resource.Resource {
	return &EngineerResource{}
}

// EngineerResource is the resource implementation.
type EngineerResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *EngineerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the resource.
func (r *EngineerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *EngineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan engineerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var engineer = client.Engineer{
		Name:  plan.Name.ValueString(),
		Email: plan.Email.ValueString(),
	}

	createdEngineer, err := r.client.CreateEngineer(engineer)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Engineer",
			"Could not create order, unexpected error: "+err.Error(),
		)

		return
	}

	plan.ID = types.StringValue(createdEngineer.ID)
	plan.Name = types.StringValue(createdEngineer.Name)
	plan.Email = types.StringValue(createdEngineer.Email)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *EngineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state engineerModel

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	engineer, err := r.client.GetEngineer(state.ID.ValueString())

	if err != nil {
		// Treat backend not-found as drift and remove from state
		if strings.Contains(err.Error(), "status: 404") || strings.Contains(err.Error(), "No engineer with that ID") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Engineer",
			"Could not read Engineer: "+err.Error(),
		)

		return
	}

	state.ID = types.StringValue(engineer.ID)
	state.Name = types.StringValue(engineer.Name)
	state.Email = types.StringValue(engineer.Email)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *EngineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan engineerModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var reqEngineer = client.Engineer{
		Name:  plan.Name.ValueString(),
		Email: plan.Email.ValueString(),
	}
	reqEngineer.ID = plan.ID.ValueString()

	// Update existing order
	_, err := r.client.UpdateEngineer(plan.ID.ValueString(), reqEngineer)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Engineer",
			"Could not update Engineer, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	engineer, err := r.client.GetEngineer(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Engineer",
			"Could not read Engineer ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state with updated items and timestamp
	plan.ID = types.StringValue(engineer.ID)
	plan.Name = types.StringValue(engineer.Name)
	plan.Email = types.StringValue(engineer.Email)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *EngineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state engineerModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing engineer
	err := r.client.DeleteEngineer(state.ID.ValueString())
	if err != nil {
		// If backend returns 404, treat as already deleted
		if strings.Contains(err.Error(), "status: 404") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Engineer: "+state.ID.ValueString(),
			"Could not delete Engineer, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *EngineerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = client
}
