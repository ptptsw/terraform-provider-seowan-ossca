package provider

import (
    "context"
    "fmt"
    "strconv"
    "time"

	"github.com/ptptsw/hashicups-client-go"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"    
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &foodResource{}
    _ resource.ResourceWithConfigure = &foodResource{}
)
// NewfoodResource is a helper function to simplify the provider implementation.
func NewFoodResource() resource.Resource {
    return &foodResource{}
}

// foodResource is the resource implementation.
type foodResource struct{
	client *hashicups.Client
}

type foodResourceModel struct {
    ID          types.String     `tfsdk:"id"`
    Items       []foodItemModel `tfsdk:"items"`
    LastUpdated types.String     `tfsdk:"last_updated"`
}


type foodItemModel struct {
	Name types.String `tfsdk:"name"`
	Price types.Float64 `tfsdk:"price"`
}

// Configure adds the provider configured client to the resource.
func (r *foodResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

	client, ok := req.ProviderData.(*hashicups.Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}


// Metadata returns the resource type name.
func (r *foodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_food"
}

// Schema defines the schema for the resource.
func (r *foodResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
                Computed: true,
				PlanModifiers: []planmodifier.String{                    
					stringplanmodifier.UseStateForUnknown(),                
				},
            },
            "last_updated": schema.StringAttribute{
                Computed: true,
			},
			"items": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"price": schema.Float64Attribute{
							Required: true,
						},
					},
				},
			},
			
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *foodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan foodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

	var items []hashicups.FoodItem
	for _, item := range plan.Items{
		items = append(items, hashicups.FoodItem{
			Name: item.Name.ValueString(),
			Price: item.Price.ValueFloat64(),
		})

	}

	food, err := r.client.CreateFood(items,nil)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating food",
            "Could not create , unexpected error: "+err.Error(),
        )
        return
	}
	
	plan.ID = types.StringValue(strconv.Itoa(food.ID))
	for foodItemIndex, foodItem := range food.Items {
		plan.Items[foodItemIndex] = foodItemModel{
			Name: types.StringValue(foodItem.Name),
			Price : types.Float64Value(foodItem.Price),
		}
	}
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

}

// Read refreshes the Terraform state with the latest data.
func (r *foodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state foodResourceModel
	diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
	}
	food, err := r.client.GetFood(state.ID.ValueString(), nil)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading hashicups Food",
            "Could not read hashicups food ID "+state.ID.ValueString()+": "+err.Error(),
        )
        return
	}
	
	state.Items = []foodItemModel{}
	for _, item := range food.Items{
		state.Items = append(state.Items, foodItemModel{
			Name: types.StringValue(item.Name),
			Price: types.Float64Value(item.Price),
		})
	}

	diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }


}

// Update updates the resource and sets the updated Terraform state on success.
func (r *foodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Retrieve values from plan
    var plan foodResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Generate API request body from plan
    var hashicupsItems []hashicups.FoodItem
    for _, item := range plan.Items {
        hashicupsItems = append(hashicupsItems, hashicups.FoodItem{
            Name: item.Name.ValueString(),
            Price: item.Price.ValueFloat64(),
        })
    }


    //Update existing food
    _, err := r.client.UpdateFood(plan.ID.ValueString(), hashicupsItems, nil)
    if err != nil {
        resp.Diagnostics.AddError(
			"Error Updating HashiCups Food",
            "Could not update food, unexpected error: "+err.Error(),
        )
        return
    }

    // Fetch updated items from GetFood as UpdateFood items are not
    // populated.
    food, err := r.client.GetFood(plan.ID.ValueString(), nil)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error Reading HashiCups Food",
            "Could not read HashiCups food ID "+plan.ID.ValueString()+": "+err.Error(),
        )
        return
    }

    // Update resource state with updated items and timestamp
    plan.Items = []foodItemModel{}
    for _, item := range food.Items {
        plan.Items = append(plan.Items, foodItemModel{
			Name : types.StringValue(item.Name),
            Price: types.Float64Value(item.Price),
        })
    }
    plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

    diags = resp.State.Set(ctx, plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *foodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
		var state foodResourceModel
		diags := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	
		// Delete existing food
		err := r.client.DeleteFood(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Deleting HashiCups Food",
				"Could not delete food, unexpected error: "+err.Error(),
			)
			return
		}
	}