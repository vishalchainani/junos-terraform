package main

import (
	"context"
	"encoding/xml"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

// Junos XML Hierarchy

type xml_Configuration struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName        xml.Name             `xml:"groups"`
		Name           *string              `xml:"name"`
		Policy_options []xml_Policy_options `xml:"policy-options,omitempty"`
	}
}
type xml_Policy_options struct {
	XMLName          xml.Name                              `xml:"policy-options"`
	Policy_statement []xml_Policy_options_Policy_statement `xml:"policy-statement,omitempty"`
	Community        []xml_Policy_options_Community        `xml:"community,omitempty"`
}

type xml_Policy_options_Policy_statement struct {
	XMLName xml.Name                                   `xml:"policy-statement"`
	Name    *string                                    `xml:"name,omitempty"`
	Term    []xml_Policy_options_Policy_statement_Term `xml:"term,omitempty"`
}
type xml_Policy_options_Community struct {
	XMLName xml.Name `xml:"community"`
	Name    *string  `xml:"name,omitempty"`
	Members *string  `xml:"members,omitempty"`
}

type xml_Policy_options_Policy_statement_Term struct {
	XMLName xml.Name                                        `xml:"term"`
	Name    *string                                         `xml:"name,omitempty"`
	From    []xml_Policy_options_Policy_statement_Term_From `xml:"from,omitempty"`
	Then    []xml_Policy_options_Policy_statement_Term_Then `xml:"then,omitempty"`
}

type xml_Policy_options_Policy_statement_Term_From struct {
	XMLName  xml.Name `xml:"from"`
	Protocol *string  `xml:"protocol,omitempty"`
}
type xml_Policy_options_Policy_statement_Term_Then struct {
	XMLName   xml.Name                                                  `xml:"then"`
	Community []xml_Policy_options_Policy_statement_Term_Then_Community `xml:"community,omitempty"`
	Accept    *string                                                   `xml:"accept,omitempty"`
	Reject    *string                                                   `xml:"reject,omitempty"`
}

type xml_Policy_options_Policy_statement_Term_Then_Community struct {
	XMLName        xml.Name `xml:"community"`
	Add            *string  `xml:"add,omitempty"`
	Community_name *string  `xml:"community-name,omitempty"`
}

// Collecting objects from the .tf file
type Groups_Model struct {
	ResourceName   types.String `tfsdk:"resource_name"`
	Policy_options types.List   `tfsdk:"policy_options"`
}

func (o Groups_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"policy_options": types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()}},
	}
}
func (o Groups_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"resource_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "xpath is `config.Groups.resource_name`",
		},
		"policy_options": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Model{}.Attributes(),
			},
		},
	}
}

type Policy_options_Model struct {
	Policy_statement types.List `tfsdk:"policy_statement"`
	Community        types.List `tfsdk:"community"`
}

func (o Policy_options_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"policy_statement": types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Model{}.AttrTypes()}},
		"community":        types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Community_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"policy_statement": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Policy_statement_Model{}.Attributes(),
			},
		},
		"community": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Community_Model{}.Attributes(),
			},
		},
	}
}

type Policy_options_Policy_statement_Model struct {
	Name types.String `tfsdk:"name"`
	Term types.List   `tfsdk:"term"`
}

func (o Policy_options_Policy_statement_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"term": types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Policy_statement_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Name.Policy_statement`",
		},
		"term": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Policy_statement_Term_Model{}.Attributes(),
			},
		},
	}
}

type Policy_options_Community_Model struct {
	Name    types.String `tfsdk:"name"`
	Members types.String `tfsdk:"members"`
}

func (o Policy_options_Community_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":    types.StringType,
		"members": types.StringType,
	}
}
func (o Policy_options_Community_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Name.Community`",
		},
		"members": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Members.Community`",
		},
	}
}

type Policy_options_Policy_statement_Term_Model struct {
	Name types.String `tfsdk:"name"`
	From types.List   `tfsdk:"from"`
	Then types.List   `tfsdk:"then"`
}

func (o Policy_options_Policy_statement_Term_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"from": types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Model{}.AttrTypes()}},
		"then": types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Model{}.AttrTypes()}},
	}
}
func (o Policy_options_Policy_statement_Term_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Name.Term`",
		},
		"from": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Policy_statement_Term_From_Model{}.Attributes(),
			},
		},
		"then": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Policy_statement_Term_Then_Model{}.Attributes(),
			},
		},
	}
}

type Policy_options_Policy_statement_Term_From_Model struct {
	Protocol types.String `tfsdk:"protocol"`
}

func (o Policy_options_Policy_statement_Term_From_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"protocol": types.StringType,
	}
}
func (o Policy_options_Policy_statement_Term_From_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"protocol": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Protocol.From`",
		},
	}
}

type Policy_options_Policy_statement_Term_Then_Model struct {
	Community types.List `tfsdk:"community"`
	Accept    types.Bool `tfsdk:"accept"`
	Reject    types.Bool `tfsdk:"reject"`
}

func (o Policy_options_Policy_statement_Term_Then_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"community": types.ListType{ElemType: types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Community_Model{}.AttrTypes()}},
		"accept":    types.BoolType,
		"reject":    types.BoolType,
	}
}
func (o Policy_options_Policy_statement_Term_Then_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"community": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Policy_options_Policy_statement_Term_Then_Community_Model{}.Attributes(),
			},
		},
		"accept": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Accept.Then`",
		},
		"reject": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Reject.Then`",
		},
	}
}

type Policy_options_Policy_statement_Term_Then_Community_Model struct {
	Add            types.Bool   `tfsdk:"add"`
	Community_name types.String `tfsdk:"community_name"`
}

func (o Policy_options_Policy_statement_Term_Then_Community_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"add":            types.BoolType,
		"community_name": types.StringType,
	}
}
func (o Policy_options_Policy_statement_Term_Then_Community_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"add": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Add.Community`",
		},
		"community_name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Community-name.Community`",
		},
	}
}

// Collects the data for the crud work
type resource_Apply_Groups struct {
	client ProviderConfig
}

func (r *resource_Apply_Groups) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(ProviderConfig)
}

// Metadata implements resource.Resource.
func (r *resource_Apply_Groups) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_Apply_Groups"
}

// Schema implements resource.Resource.
func (r *resource_Apply_Groups) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"policy_options": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Policy_options_Model{}.Attributes(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *resource_Apply_Groups) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var plan Groups_Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
	var config xml_Configuration
	config.Groups.Name = plan.ResourceName.ValueStringPointer()

	var var_policy_options []Policy_options_Model
	if plan.Policy_options.IsNull() {
		var_policy_options = []Policy_options_Model{}
	} else {
		resp.Diagnostics.Append(plan.Policy_options.ElementsAs(ctx, &var_policy_options, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	config.Groups.Policy_options = make([]xml_Policy_options, len(var_policy_options))

	for i_policy_options, v_policy_options := range var_policy_options {
		var var_policy_options_policy_statement []Policy_options_Policy_statement_Model
		resp.Diagnostics.Append(v_policy_options.Policy_statement.ElementsAs(ctx, &var_policy_options_policy_statement, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Policy_options[i_policy_options].Policy_statement = make([]xml_Policy_options_Policy_statement, len(var_policy_options_policy_statement))
		for i_policy_options_policy_statement, v_policy_options_policy_statement := range var_policy_options_policy_statement {
			config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Name = v_policy_options_policy_statement.Name.ValueStringPointer()
			var var_policy_options_policy_statement_term []Policy_options_Policy_statement_Term_Model
			resp.Diagnostics.Append(v_policy_options_policy_statement.Term.ElementsAs(ctx, &var_policy_options_policy_statement_term, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term = make([]xml_Policy_options_Policy_statement_Term, len(var_policy_options_policy_statement_term))
			for i_policy_options_policy_statement_term, v_policy_options_policy_statement_term := range var_policy_options_policy_statement_term {
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Name = v_policy_options_policy_statement_term.Name.ValueStringPointer()
				var var_policy_options_policy_statement_term_from []Policy_options_Policy_statement_Term_From_Model
				resp.Diagnostics.Append(v_policy_options_policy_statement_term.From.ElementsAs(ctx, &var_policy_options_policy_statement_term_from, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From = make([]xml_Policy_options_Policy_statement_Term_From, len(var_policy_options_policy_statement_term_from))
				for i_policy_options_policy_statement_term_from, v_policy_options_policy_statement_term_from := range var_policy_options_policy_statement_term_from {
					config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Protocol = v_policy_options_policy_statement_term_from.Protocol.ValueStringPointer()
				}
				var var_policy_options_policy_statement_term_then []Policy_options_Policy_statement_Term_Then_Model
				resp.Diagnostics.Append(v_policy_options_policy_statement_term.Then.ElementsAs(ctx, &var_policy_options_policy_statement_term_then, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then = make([]xml_Policy_options_Policy_statement_Term_Then, len(var_policy_options_policy_statement_term_then))
				for i_policy_options_policy_statement_term_then, v_policy_options_policy_statement_term_then := range var_policy_options_policy_statement_term_then {
					var var_policy_options_policy_statement_term_then_community []Policy_options_Policy_statement_Term_Then_Community_Model
					resp.Diagnostics.Append(v_policy_options_policy_statement_term_then.Community.ElementsAs(ctx, &var_policy_options_policy_statement_term_then_community, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community = make([]xml_Policy_options_Policy_statement_Term_Then_Community, len(var_policy_options_policy_statement_term_then_community))
					for i_policy_options_policy_statement_term_then_community, v_policy_options_policy_statement_term_then_community := range var_policy_options_policy_statement_term_then_community {
						if v_policy_options_policy_statement_term_then_community.Add.ValueBool() {
							empty := ""
							config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Add = &empty
						}
						config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Community_name = v_policy_options_policy_statement_term_then_community.Community_name.ValueStringPointer()
					}
					if v_policy_options_policy_statement_term_then.Accept.ValueBool() {
						empty := ""
						config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Accept = &empty
					}
					if v_policy_options_policy_statement_term_then.Reject.ValueBool() {
						empty := ""
						config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Reject = &empty
					}
				}
			}
		}
		var var_policy_options_community []Policy_options_Community_Model
		resp.Diagnostics.Append(v_policy_options.Community.ElementsAs(ctx, &var_policy_options_community, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Policy_options[i_policy_options].Community = make([]xml_Policy_options_Community, len(var_policy_options_community))
		for i_policy_options_community, v_policy_options_community := range var_policy_options_community {
			config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Name = v_policy_options_community.Name.ValueStringPointer()
			config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Members = v_policy_options_community.Members.ValueStringPointer()
		}
	}

	err := r.client.SendTransaction(plan.ResourceName.ValueString(), config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while adding group", err.Error())
		return
	}
	commit_err := r.client.SendCommit()
	if commit_err != nil {
		resp.Diagnostics.AddError("Failed while committing apply-group", commit_err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *resource_Apply_Groups) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state Groups_Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config xml_Configuration
	err := r.client.MarshalGroup(state.ResourceName.ValueString(), &config)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read group", err.Error())
		return
	}
	state.Policy_options = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
	policy_options_List := make([]Policy_options_Model, len(config.Groups.Policy_options))
	for i_policy_options, v_policy_options := range config.Groups.Policy_options {
		var policy_options_model Policy_options_Model
		policy_options_policy_statement_List := make([]Policy_options_Policy_statement_Model, len(v_policy_options.Policy_statement))
		for i_policy_options_policy_statement, v_policy_options_policy_statement := range v_policy_options.Policy_statement {
			var policy_options_policy_statement_model Policy_options_Policy_statement_Model
			policy_options_policy_statement_model.Name = types.StringPointerValue(v_policy_options_policy_statement.Name)

			policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model

			policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model

			policy_options_policy_statement_term_List := make([]Policy_options_Policy_statement_Term_Model, len(v_policy_options_policy_statement.Term))
			for i_policy_options_policy_statement_term, v_policy_options_policy_statement_term := range v_policy_options_policy_statement.Term {
				var policy_options_policy_statement_term_model Policy_options_Policy_statement_Term_Model
				policy_options_policy_statement_term_model.Name = types.StringPointerValue(v_policy_options_policy_statement_term.Name)

				policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model

				policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model

				policy_options_policy_statement_term_from_List := make([]Policy_options_Policy_statement_Term_From_Model, len(v_policy_options_policy_statement_term.From))
				for i_policy_options_policy_statement_term_from, v_policy_options_policy_statement_term_from := range v_policy_options_policy_statement_term.From {
					var policy_options_policy_statement_term_from_model Policy_options_Policy_statement_Term_From_Model
					policy_options_policy_statement_term_from_model.Protocol = types.StringPointerValue(v_policy_options_policy_statement_term_from.Protocol)

					policy_options_policy_statement_term_from_List[i_policy_options_policy_statement_term_from] = policy_options_policy_statement_term_from_model

				}
				policy_options_policy_statement_term_model.From, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_From_Model{}.AttrTypes()}, policy_options_policy_statement_term_from_List)
				policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model

				policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model

				policy_options_policy_statement_term_then_List := make([]Policy_options_Policy_statement_Term_Then_Model, len(v_policy_options_policy_statement_term.Then))
				for i_policy_options_policy_statement_term_then, v_policy_options_policy_statement_term_then := range v_policy_options_policy_statement_term.Then {
					var policy_options_policy_statement_term_then_model Policy_options_Policy_statement_Term_Then_Model

					policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model

					policy_options_policy_statement_term_then_community_List := make([]Policy_options_Policy_statement_Term_Then_Community_Model, len(v_policy_options_policy_statement_term_then.Community))
					for i_policy_options_policy_statement_term_then_community, v_policy_options_policy_statement_term_then_community := range v_policy_options_policy_statement_term_then.Community {
						var policy_options_policy_statement_term_then_community_model Policy_options_Policy_statement_Term_Then_Community_Model
						policy_options_policy_statement_term_then_community_model.Add = types.BoolValue(v_policy_options_policy_statement_term_then_community.Add != nil)

						policy_options_policy_statement_term_then_community_List[i_policy_options_policy_statement_term_then_community] = policy_options_policy_statement_term_then_community_model

						policy_options_policy_statement_term_then_community_model.Community_name = types.StringPointerValue(v_policy_options_policy_statement_term_then_community.Community_name)

						policy_options_policy_statement_term_then_community_List[i_policy_options_policy_statement_term_then_community] = policy_options_policy_statement_term_then_community_model

					}
					policy_options_policy_statement_term_then_model.Community, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Community_Model{}.AttrTypes()}, policy_options_policy_statement_term_then_community_List)
					policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model
					policy_options_policy_statement_term_then_model.Accept = types.BoolValue(v_policy_options_policy_statement_term_then.Accept != nil)

					policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model

					policy_options_policy_statement_term_then_model.Reject = types.BoolValue(v_policy_options_policy_statement_term_then.Reject != nil)

					policy_options_policy_statement_term_then_List[i_policy_options_policy_statement_term_then] = policy_options_policy_statement_term_then_model

				}
				policy_options_policy_statement_term_model.Then, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Then_Model{}.AttrTypes()}, policy_options_policy_statement_term_then_List)
				policy_options_policy_statement_term_List[i_policy_options_policy_statement_term] = policy_options_policy_statement_term_model
			}
			policy_options_policy_statement_model.Term, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Term_Model{}.AttrTypes()}, policy_options_policy_statement_term_List)
			policy_options_policy_statement_List[i_policy_options_policy_statement] = policy_options_policy_statement_model
		}
		policy_options_model.Policy_statement, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Policy_statement_Model{}.AttrTypes()}, policy_options_policy_statement_List)
		policy_options_List[i_policy_options] = policy_options_model
		policy_options_community_List := make([]Policy_options_Community_Model, len(v_policy_options.Community))
		for i_policy_options_community, v_policy_options_community := range v_policy_options.Community {
			var policy_options_community_model Policy_options_Community_Model
			policy_options_community_model.Name = types.StringPointerValue(v_policy_options_community.Name)

			policy_options_community_List[i_policy_options_community] = policy_options_community_model

			policy_options_community_model.Members = types.StringPointerValue(v_policy_options_community.Members)

			policy_options_community_List[i_policy_options_community] = policy_options_community_model

		}
		policy_options_model.Community, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Community_Model{}.AttrTypes()}, policy_options_community_List)
		policy_options_List[i_policy_options] = policy_options_model
	}
	state.Policy_options, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Policy_options_Model{}.AttrTypes()}, policy_options_List)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update implements resource.Resource.
func (r *resource_Apply_Groups) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan Groups_Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}
	var config xml_Configuration
	config.Groups.Name = plan.ResourceName.ValueStringPointer()

	var var_policy_options []Policy_options_Model
	if plan.Policy_options.IsNull() {
		var_policy_options = []Policy_options_Model{}
	} else {
		resp.Diagnostics.Append(plan.Policy_options.ElementsAs(ctx, &var_policy_options, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	config.Groups.Policy_options = make([]xml_Policy_options, len(var_policy_options))

	for i_policy_options, v_policy_options := range var_policy_options {
		var var_policy_options_policy_statement []Policy_options_Policy_statement_Model
		resp.Diagnostics.Append(v_policy_options.Policy_statement.ElementsAs(ctx, &var_policy_options_policy_statement, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Policy_options[i_policy_options].Policy_statement = make([]xml_Policy_options_Policy_statement, len(var_policy_options_policy_statement))
		for i_policy_options_policy_statement, v_policy_options_policy_statement := range var_policy_options_policy_statement {
			config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Name = v_policy_options_policy_statement.Name.ValueStringPointer()
			var var_policy_options_policy_statement_term []Policy_options_Policy_statement_Term_Model
			resp.Diagnostics.Append(v_policy_options_policy_statement.Term.ElementsAs(ctx, &var_policy_options_policy_statement_term, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term = make([]xml_Policy_options_Policy_statement_Term, len(var_policy_options_policy_statement_term))
			for i_policy_options_policy_statement_term, v_policy_options_policy_statement_term := range var_policy_options_policy_statement_term {
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Name = v_policy_options_policy_statement_term.Name.ValueStringPointer()
				var var_policy_options_policy_statement_term_from []Policy_options_Policy_statement_Term_From_Model
				resp.Diagnostics.Append(v_policy_options_policy_statement_term.From.ElementsAs(ctx, &var_policy_options_policy_statement_term_from, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From = make([]xml_Policy_options_Policy_statement_Term_From, len(var_policy_options_policy_statement_term_from))
				for i_policy_options_policy_statement_term_from, v_policy_options_policy_statement_term_from := range var_policy_options_policy_statement_term_from {
					config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].From[i_policy_options_policy_statement_term_from].Protocol = v_policy_options_policy_statement_term_from.Protocol.ValueStringPointer()
				}
				var var_policy_options_policy_statement_term_then []Policy_options_Policy_statement_Term_Then_Model
				resp.Diagnostics.Append(v_policy_options_policy_statement_term.Then.ElementsAs(ctx, &var_policy_options_policy_statement_term_then, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then = make([]xml_Policy_options_Policy_statement_Term_Then, len(var_policy_options_policy_statement_term_then))
				for i_policy_options_policy_statement_term_then, v_policy_options_policy_statement_term_then := range var_policy_options_policy_statement_term_then {
					var var_policy_options_policy_statement_term_then_community []Policy_options_Policy_statement_Term_Then_Community_Model
					resp.Diagnostics.Append(v_policy_options_policy_statement_term_then.Community.ElementsAs(ctx, &var_policy_options_policy_statement_term_then_community, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community = make([]xml_Policy_options_Policy_statement_Term_Then_Community, len(var_policy_options_policy_statement_term_then_community))
					for i_policy_options_policy_statement_term_then_community, v_policy_options_policy_statement_term_then_community := range var_policy_options_policy_statement_term_then_community {
						if v_policy_options_policy_statement_term_then_community.Add.ValueBool() {
							empty := ""
							config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Add = &empty
						}
						config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Community[i_policy_options_policy_statement_term_then_community].Community_name = v_policy_options_policy_statement_term_then_community.Community_name.ValueStringPointer()
					}
					if v_policy_options_policy_statement_term_then.Accept.ValueBool() {
						empty := ""
						config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Accept = &empty
					}
					if v_policy_options_policy_statement_term_then.Reject.ValueBool() {
						empty := ""
						config.Groups.Policy_options[i_policy_options].Policy_statement[i_policy_options_policy_statement].Term[i_policy_options_policy_statement_term].Then[i_policy_options_policy_statement_term_then].Reject = &empty
					}
				}
			}
		}
		var var_policy_options_community []Policy_options_Community_Model
		resp.Diagnostics.Append(v_policy_options.Community.ElementsAs(ctx, &var_policy_options_community, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Policy_options[i_policy_options].Community = make([]xml_Policy_options_Community, len(var_policy_options_community))
		for i_policy_options_community, v_policy_options_community := range var_policy_options_community {
			config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Name = v_policy_options_community.Name.ValueStringPointer()
			config.Groups.Policy_options[i_policy_options].Community[i_policy_options_community].Members = v_policy_options_community.Members.ValueStringPointer()
		}
	}

	err := r.client.SendTransaction(plan.ResourceName.ValueString(), config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	commit_err := r.client.SendCommit()
	if commit_err != nil {
		resp.Diagnostics.AddError("Failed while committing apply-group", commit_err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (r *resource_Apply_Groups) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state Groups_Model
	d := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteConfig(state.ResourceName.ValueString(), false)
	if err != nil {
		if strings.Contains(err.Error(), "ound") {
			return
		}
		resp.Diagnostics.AddError("Failed while deleting configuration", err.Error())
		return
	}
	commit_err := r.client.SendCommit()
	if commit_err != nil {
		resp.Diagnostics.AddError("Failed while committing apply-group", commit_err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
