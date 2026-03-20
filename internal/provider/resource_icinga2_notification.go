package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/lrsmith/go-icinga2-api/iapi"
)

var (
	_ resource.Resource              = &notificationResource{}
	_ resource.ResourceWithConfigure = &notificationResource{}
)

func Notification() resource.Resource {
	return &notificationResource{}
}

type notificationResourceModel struct {
	ID          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Hostname    types.String `tfsdk:"hostname"`
	Servicename types.String `tfsdk:"servicename"`
	Command     types.String `tfsdk:"command"`
	Users       types.List   `tfsdk:"users"`
	Vars        types.Map    `tfsdk:"vars"`
	Templates   types.List   `tfsdk:"templates"`
	Interval    types.Int64  `tfsdk:"interval"`
}

type notificationResource struct {
	client *iapi.Server
}

func (r *notificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notification"
}

func (r *notificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"hostname": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"servicename": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"command": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"users": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"vars": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"templates": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"interval": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1800),
			},
		},
	}
}

func (r *notificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *notificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostname := plan.Hostname.ValueString()
	servicename := plan.Servicename.ValueString()
	var name string
	if servicename != "" {
		name = hostname + "!" + servicename + "!" + hostname + "-" + servicename
	} else {
		name = hostname + "!" + hostname
	}

	var users []string
	if !plan.Users.IsNull() && !plan.Users.IsUnknown() {
		for _, user := range plan.Users.Elements() {
			if strVal, ok := user.(types.String); ok {
				users = append(users, strVal.ValueString())
			} else {
				resp.Diagnostics.AddError(
					"Error creating Notification",
					"User not a string",
				)
			}
		}
	}

	vars := make(map[string]string)
	if !plan.Vars.IsNull() && !plan.Vars.IsUnknown() {
		for key, value := range plan.Vars.Elements() {
			if strVal, ok := value.(types.String); ok {
				vars[key] = strVal.ValueString()
			} else {
				resp.Diagnostics.AddError(
					"Error creating Notification",
					"Variable not a string",
				)
			}
		}
	}

	var templates []string
	if !plan.Templates.IsNull() && !plan.Templates.IsUnknown() {
		for _, template := range plan.Templates.Elements() {
			if strVal, ok := template.(types.String); ok {
				templates = append(templates, strVal.ValueString())
			} else {
				resp.Diagnostics.AddError(
					"Error creating Notification",
					"Template not a string",
				)
			}
		}
	}

	interval := int(plan.Interval.ValueInt64())

	notifications, err := r.client.CreateNotification(name, hostname, plan.Command.ValueString(), servicename, interval, users, vars, templates)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Notification",
			"Could not create notification unexpected error: "+err.Error(),
		)
		return
	}

	for _, notification := range notifications {
		if notification.Name == name {
			plan.ID = types.StringValue(notification.Name)
			plan.Hostname = types.StringValue(hostname)
			plan.Command = types.StringValue(notification.Attrs.Command)
			plan.Servicename = types.StringValue(notification.Attrs.Servicename)
			plan.Interval = types.Int64Value(int64(notification.Attrs.Interval))
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *notificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostname := state.Hostname.ValueString()
	servicename := state.Servicename.ValueString()
	var name string
	if servicename != "" {
		name = hostname + "!" + servicename + "!" + hostname + "-" + servicename
	} else {
		name = hostname + "!" + hostname
	}

	notifications, err := r.client.GetNotification(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Notification",
			"Could not read notification "+name+": "+err.Error(),
		)
		return
	}

	for _, notification := range notifications {
		if notification.Name == name {
			state.ID = types.StringValue(notification.Name)
			state.Hostname = types.StringValue(hostname)
			state.Command = types.StringValue(notification.Attrs.Command)
			state.Servicename = types.StringValue(notification.Attrs.Servicename)
			state.Interval = types.Int64Value(int64(notification.Attrs.Interval))
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *notificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Updates are currently not supported for notification resources",
	)
}

func (r *notificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostname := state.Hostname.ValueString()
	servicename := state.Servicename.ValueString()
	var name string
	if servicename != "" {
		name = hostname + "!" + servicename + "!" + hostname + "-" + servicename
	} else {
		name = hostname + "!" + hostname
	}

	err := r.client.DeleteNotification(name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Notification",
			"Could not delete notification, unexpected error: "+err.Error(),
		)
		return
	}
}
