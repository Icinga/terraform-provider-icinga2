package icinga2

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/lrsmith/go-icinga2-api/iapi"
)

var (
	_ resource.Resource              = &hostResource{}
	_ resource.ResourceWithConfigure = &hostResource{}
)

func Host() resource.Resource {
	return &hostResource{}
}

type hostResourceModel struct {
	ID            types.String `tfsdk:"id"`
	LastUpdated   types.String `tfsdk:"last_updated"`
	Hostname      types.String `tfsdk:"hostname"`
	Address       types.String `tfsdk:"address"`
	Check_command types.String `tfsdk:"check_command"`
	Groups        []string     `tfsdk:"groups"`
	Vars          types.Map    `tfsdk:"vars"`
	Templates     []string     `tfsdk:"templates"`
}

// hostResource defines the resource implementation.
type hostResource struct {
	client *iapi.Server
}

func (r *hostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *hostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "Name of the Host",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "Address of Host",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"check_command": schema.StringAttribute{
				Required:    true,
				Description: "Address of Host",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"groups": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"vars": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"templates": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *hostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*iapi.Server)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *iapi.Server, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan hostResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vars := make(map[string]string)
	iterator := plan.Vars.Elements()
	for key, value := range iterator {
		vars[key] = value.String()
	}

	hosts, err := r.client.CreateHost(
		plan.Hostname.ValueString(),
		plan.Address.ValueString(),
		plan.Check_command.ValueString(),
		vars,
		plan.Templates,
		plan.Groups)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Host",
			"Could not create host unexpected error: "+err.Error(),
		)
		return
	}

	for _, host := range hosts {
		if host.Name == plan.Hostname.ValueString() {
			plan.ID = types.StringValue(host.Name)
			plan.Hostname = types.StringValue(host.Name)
			plan.Address = types.StringValue(host.Attrs.Address)
			plan.Check_command = types.StringValue(host.Attrs.CheckCommand)
			plan.Groups = host.Attrs.Groups
			plan.Templates = host.Attrs.Templates
			//plan.Vars = types.MapValue(types.StringType, host.Attrs.Vars)
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state hostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hosts, err := r.client.GetHost(state.Hostname.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Host",
			"Could not read host "+state.Hostname.ValueString()+": "+err.Error(),
		)
		return
	}

	for _, host := range hosts {
		if host.Name == state.Hostname.ValueString() {
			state.ID = types.StringValue(host.Name)
			state.Hostname = types.StringValue(host.Name)
			state.Address = types.StringValue(host.Attrs.Address)
			state.Check_command = types.StringValue(host.Attrs.CheckCommand)
			state.Groups = host.Attrs.Groups
			state.Templates = host.Attrs.Templates
			//state.Vars = types.MapValue(types.String{}, host.Attrs.Vars)
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state hostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHost(state.Hostname.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Host",
			"Could not delete host, unexpected error: "+err.Error(),
		)
		return
	}
}
