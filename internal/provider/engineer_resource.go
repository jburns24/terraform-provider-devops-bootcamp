// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	devops_resource "github.com/liatrio/devops-bootcamp/examples/ch6/devops-resources"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &EngineerResource{}
var _ resource.ResourceWithImportState = &EngineerResource{}

func NewEngineerResource() resource.Resource {
	return &EngineerResource{}
}

// EngineerResource defines the resource implementation.
type EngineerResource struct {
	client *Client
}

func (r *EngineerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_engineer"
}

func (r *EngineerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Engineer resource source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Engineer",
				Optional:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email of the Engineer",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the Engineer",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *EngineerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *EngineerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *EngineerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	engineer, err := r.client.CreateEngineer(plan.Name.ValueString(), plan.Email.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating engineer",
			"Could not create engineer,unexpected error:"+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(engineer.Id)
	plan.Email = types.StringValue(engineer.Email)
	plan.Name = types.StringValue(engineer.Name)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Explicitly setting the id

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EngineerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *EngineerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	engineer, err := r.client.GetEngineer(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	state.Email = types.StringValue(engineer.Email)
	state.Id = types.StringValue(engineer.Id)
	state.Name = types.StringValue(engineer.Name)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EngineerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *EngineerModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update via client api
	body, err := r.client.UpdateEngineer(plan.Id.ValueString(), plan.Name.ValueString(), plan.Email.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating engineer",
			"Could not create engineer,unexpected error:"+err.Error(),
		)
		return
	}

	// Update tf state
	plan.Email = types.StringValue(body.Email)
	plan.Name = types.StringValue(body.Name)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EngineerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *EngineerModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	engineer := devops_resource.Engineer{
		Id:    data.Id.ValueString(),
		Name:  data.Name.ValueString(),
		Email: data.Email.ValueString(),
	}

	err := r.client.DeleteEngineer(&engineer)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Engineer",
			"Could not delete engineer, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *EngineerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
