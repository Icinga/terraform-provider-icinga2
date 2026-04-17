package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

var (
	_ resource.Resource              = &downtimeResource{}
	_ resource.ResourceWithConfigure = &downtimeResource{}
)

func Downtime() resource.Resource {
	return &downtimeResource{}
}

type downtimeResourceModel struct {
	Names        types.List   `tfsdk:"names"`
	LastUpdated  types.String `tfsdk:"last_updated"`
	Type         types.String `tfsdk:"type"`
	Filter       types.String `tfsdk:"filter"`
	Author       types.String `tfsdk:"author"`
	Comment      types.String `tfsdk:"comment"`
	StartTime    types.Int64  `tfsdk:"start_time"`
	EndTime      types.Int64  `tfsdk:"end_time"`
	Fixed        types.Bool   `tfsdk:"fixed"`
	Duration     types.Int64  `tfsdk:"duration"`
	AllServices  types.Bool   `tfsdk:"all_services"`
	TriggerName  types.String `tfsdk:"trigger_name"`
	ChildOptions types.String `tfsdk:"child_options"`
}

// hostResource defines the resource implementation.
type downtimeResource struct {
	client *iapi.Server
}

func (r *downtimeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_downtime"
}

// https://icinga.com/docs/icinga-2/latest/doc/12-icinga2-api/#icinga2-api-actions-schedule-downtime
func (r *downtimeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"names": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of downtime (Host or Service).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"filter": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "[Filter](https://icinga.com/docs/icinga-2/latest/doc/12-icinga2-api/#icinga2-api-filters) to apply the downtime.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"author": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the author.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Comment text.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"start_time": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Unix timestamp marking the beginning of the downtime.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"end_time": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Unix timestamp marking the end of the downtime.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"fixed": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Defaults to true. If true, the downtime is fixed otherwise flexible",
				Default:             booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"duration": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Required for flexible downtimes. Duration of the downtime in seconds if fixed is set to false.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"all_services": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional for host downtimes. Sets downtime for all services for the matched host objects. If child_options are set, all child hosts and their services will schedule a downtime too. Defaults to false.",
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"trigger_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Sets the trigger for a triggered downtime.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"child_options": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Schedule child downtimes. DowntimeNoChildren does not do anything, DowntimeTriggeredChildren schedules child downtimes triggered by this downtime, DowntimeNonTriggeredChildren schedules non-triggered downtimes. Defaults to DowntimeNoChildren.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *downtimeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *downtimeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan downtimeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	names, err := r.client.ScheduleDowntime(plan.Type.ValueString(), plan.Filter.ValueString(), plan.Author.ValueString(), plan.Comment.ValueString(), plan.StartTime.ValueInt64(), plan.EndTime.ValueInt64(), plan.Fixed.ValueBool(), plan.Duration.ValueInt64(), plan.AllServices.ValueBool(), plan.TriggerName.ValueString(), plan.ChildOptions.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating downtime",
			fmt.Sprintf("%v", err),
		)
	}

	namesAttr := make([]attr.Value, 0, len(names))
	for _, name := range names {
		namesAttr = append(namesAttr, types.StringValue(name))
	}
	plan.Names, _ = types.ListValue(types.StringType, namesAttr)

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *downtimeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state downtimeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *downtimeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan downtimeResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO

	// Set refreshed plan
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *downtimeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state downtimeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	downtimes := make([]types.String, 0, len(state.Names.Elements()))
	state.Names.ElementsAs(ctx, &downtimes, false)

	for _, downtime := range downtimes {
		err := r.client.RemoveDowntime(downtime.ValueString(), state.Author.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting downtime",
				fmt.Sprintf("%v", err),
			)
			return
		}
	}
}
