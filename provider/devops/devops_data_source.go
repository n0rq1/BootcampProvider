package devops

import (
	"context"

	"terraform-provider-devops/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &devopsDataSource{}
	_ datasource.DataSourceWithConfigure = &devopsDataSource{}
)

// NewDevDataSource is a helper function to simplify the provider implementation.
func NewDevopsDataSource() datasource.DataSource {
	return &devopsDataSource{}
}

// devDataSource is the data source implementation.
type devopsDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *devopsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devops"
}

// Schema defines the schema for the data source.
func (d *devopsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"devops": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Computed: true},
						"devs": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"ops": schema.ListAttribute{
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
func (d *devopsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *devopsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state DevopsDataSourceModel

    state.Devops = make([]devopsModel, 0)
    items, err := d.client.GetDevOps()
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to read DevOps groups",
            err.Error(),
        )
        return
    }

    for _, it := range items {
        row := devopsModel{
            ID: types.StringValue(it.ID),
        }

        // Map dev IDs from []client.Dev
        devIDs := make([]string, 0, len(it.Dev))
        for _, d := range it.Dev { devIDs = append(devIDs, d.ID) }
        devList, diags := types.ListValueFrom(ctx, types.StringType, devIDs)
        resp.Diagnostics.Append(diags...)
        if resp.Diagnostics.HasError() {
            return
        }
        row.Devs = devList

        // Map ops IDs from []client.Ops
        opsIDs := make([]string, 0, len(it.Ops))
        for _, o := range it.Ops { opsIDs = append(opsIDs, o.ID) }
        opsList, diags2 := types.ListValueFrom(ctx, types.StringType, opsIDs)
        resp.Diagnostics.Append(diags2...)
        if resp.Diagnostics.HasError() {
            return
        }
        row.Ops = opsList

        state.Devops = append(state.Devops, row)
    }

    diags := resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// DataSourceModel maps the data source schema data.
type DevopsDataSourceModel struct {
	Devops []devopsModel `tfsdk:"devops"`
}

// devModel maps Dev schema data.
type devopsModel struct {
	ID   types.String `tfsdk:"id"`
	Devs types.List   `tfsdk:"devs"`
	Ops  types.List   `tfsdk:"ops"`
}
