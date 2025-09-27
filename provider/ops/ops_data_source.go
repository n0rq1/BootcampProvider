package ops

import (
	"context"

	"terraform-provider-devops/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &opsDataSource{}
	_ datasource.DataSourceWithConfigure = &opsDataSource{}
)

// NewOpsDataSource is a helper function to simplify the provider implementation.
func NewOpsDataSource() datasource.DataSource {
	return &opsDataSource{}
}

// opsDataSource is the data source implementation.
type opsDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *opsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ops"
}

// Schema defines the schema for the data source.
func (d *opsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ops": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"engineers": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Configure accepts provider configured data to set the API client.
func (d *opsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *opsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state opsDataSourceModel

	state.Ops = make([]opsModel, 0)
	devs, err := d.client.GetOps()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Ops groups",
			err.Error(),
		)
		return
	}

	for _, dv := range devs {
		dvm := opsModel{
			ID:   types.StringValue(dv.ID),
			Name: types.StringValue(dv.Name),
		}

		// Map engineers (objects) to a list of engineer ID strings
		engineerIDs := make([]string, 0, len(dv.Engineers))
		for _, eng := range dv.Engineers {
			engineerIDs = append(engineerIDs, eng.ID)
		}
		engList, diags := types.ListValueFrom(ctx, types.StringType, engineerIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		dvm.Engineers = engList

		state.Ops = append(state.Ops, dvm)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// DataSourceModel maps the data source schema data.
type opsDataSourceModel struct {
	Ops []opsModel `tfsdk:"ops"`
}

// opsModel maps Ops schema data.
type opsModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Engineers types.List   `tfsdk:"engineers"`
}
