package engineers

import (
	"context"

	"terraform-provider-devops/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &engineerDataSource{}
	_ datasource.DataSourceWithConfigure = &engineerDataSource{}
)

// NewEngineersDataSource is a helper function to simplify the provider implementation.
func NewEngineersDataSource() datasource.DataSource {
	return &engineerDataSource{}
}

// engineerDataSource is the data source implementation.
type engineerDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *engineerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

// Schema defines the schema for the data source.
func (d *engineerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"engineers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    schema.StringAttribute{Computed: true},
						"name":  schema.StringAttribute{Computed: true},
						"email": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

// Configure accepts provider configured data to set the API client.
func (d *engineerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			"Expected *client.Client for provider data, but received a different type.",
		)
		return
	}

	d.client = c
}

// Read refreshes the Terraform state with the latest data.
func (d *engineerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EngineerDataSourceModel
	// Ensure engineers is an empty list rather than null when there are no results
	state.Engineers = make([]engineerModel, 0)

	engineers, err := d.client.GetEngineers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Engineers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, engineer := range engineers {
		engineerState := engineerModel{
			ID:    types.StringValue(engineer.ID),
			Name:  types.StringValue(engineer.Name),
			Email: types.StringValue(engineer.Email),
		}

		state.Engineers = append(state.Engineers, engineerState)
	}

	// Set state
	// This writes to the Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// EngineerDataSourceModel maps the data source schema data
// Direction: Data Source -> Terraform State (output only)
type EngineerDataSourceModel struct {
	Engineers []engineerModel `tfsdk:"engineers"`
}

// engineerModel maps engineers schema data.
// Direction: Data Source -> Terraform State (output only)
type engineerModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}