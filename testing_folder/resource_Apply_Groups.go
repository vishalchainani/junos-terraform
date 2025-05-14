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
		XMLName    xml.Name         `xml:"groups"`
		Name       *string          `xml:"name"`
		Interfaces []xml_Interfaces `xml:"interfaces,omitempty"`
	}
}
type xml_Interfaces struct {
	XMLName   xml.Name                   `xml:"interfaces"`
	Interface []xml_Interfaces_Interface `xml:"interface,omitempty"`
}
type xml_Interfaces_Interface struct {
	XMLName      xml.Name                        `xml:"interface"`
	Name         *string                         `xml:"name,omitempty"`
	Description  *string                         `xml:"description,omitempty"`
	Vlan_tagging *string                         `xml:"vlan-tagging,omitempty"`
	Mtu          *int64                          `xml:"mtu,omitempty"`
	Unit         []xml_Interfaces_Interface_Unit `xml:"unit,omitempty"`
}
type xml_Interfaces_Interface_Unit struct {
	XMLName     xml.Name                               `xml:"unit"`
	Name        *string                                `xml:"name,omitempty"`
	Description *string                                `xml:"description,omitempty"`
	Vlan_id     *string                                `xml:"vlan-id,omitempty"`
	Family      []xml_Interfaces_Interface_Unit_Family `xml:"family,omitempty"`
}
type xml_Interfaces_Interface_Unit_Family struct {
	XMLName xml.Name                                     `xml:"family"`
	Inet    []xml_Interfaces_Interface_Unit_Family_Inet  `xml:"inet,omitempty"`
	Inet6   []xml_Interfaces_Interface_Unit_Family_Inet6 `xml:"inet6,omitempty"`
}
type xml_Interfaces_Interface_Unit_Family_Inet struct {
	XMLName xml.Name                                            `xml:"inet"`
	Address []xml_Interfaces_Interface_Unit_Family_Inet_Address `xml:"address,omitempty"`
}
type xml_Interfaces_Interface_Unit_Family_Inet_Address struct {
	XMLName xml.Name `xml:"address"`
	Name    *string  `xml:"name,omitempty"`
}
type xml_Interfaces_Interface_Unit_Family_Inet6 struct {
	XMLName xml.Name                                             `xml:"inet6"`
	Address []xml_Interfaces_Interface_Unit_Family_Inet6_Address `xml:"address,omitempty"`
}
type xml_Interfaces_Interface_Unit_Family_Inet6_Address struct {
	XMLName xml.Name `xml:"address"`
	Name    *string  `xml:"name,omitempty"`
}

// Collecting objects from the .tf file
type Groups_Model struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Interfaces   types.List   `tfsdk:"interfaces"`
}

func (o Groups_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interfaces": types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()}},
	}
}
func (o Groups_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"resource_name": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "xpath is `config.Groups.resource_name`",
		},
		"interfaces": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Model{}.Attributes(),
			},
		},
	}
}

type Interfaces_Model struct {
	Interface types.List `tfsdk:"interface"`
}

func (o Interfaces_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface": types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"interface": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Model{}.Attributes(),
			},
		},
	}
}

type Interfaces_Interface_Model struct {
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Vlan_tagging types.Bool   `tfsdk:"vlan_tagging"`
	Mtu          types.Int64  `tfsdk:"mtu"`
	Unit         types.List   `tfsdk:"unit"`
}

func (o Interfaces_Interface_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"description":  types.StringType,
		"vlan_tagging": types.BoolType,
		"mtu":          types.Int64Type,
		"unit":         types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Name.Interface`",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Description.Interface`",
		},
		"vlan_tagging": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Vlan-tagging.Interface`",
		},
		"mtu": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Mtu.Interface`",
		},
		"unit": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Unit_Model{}.Attributes(),
			},
		},
	}
}

type Interfaces_Interface_Unit_Model struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Vlan_id     types.String `tfsdk:"vlan_id"`
	Family      types.List   `tfsdk:"family"`
}

func (o Interfaces_Interface_Unit_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"vlan_id":     types.StringType,
		"family":      types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Name.Unit`",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Description.Unit`",
		},
		"vlan_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Vlan-id.Unit`",
		},
		"family": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Unit_Family_Model{}.Attributes(),
			},
		},
	}
}

type Interfaces_Interface_Unit_Family_Model struct {
	Inet  types.List `tfsdk:"inet"`
	Inet6 types.List `tfsdk:"inet6"`
}

func (o Interfaces_Interface_Unit_Family_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inet":  types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Model{}.AttrTypes()}},
		"inet6": types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet6_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Family_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"inet": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Unit_Family_Inet_Model{}.Attributes(),
			},
		},
		"inet6": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Unit_Family_Inet6_Model{}.Attributes(),
			},
		},
	}
}

type Interfaces_Interface_Unit_Family_Inet_Model struct {
	Address types.List `tfsdk:"address"`
}

func (o Interfaces_Interface_Unit_Family_Inet_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Family_Inet_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.Attributes(),
			},
		},
	}
}

type Interfaces_Interface_Unit_Family_Inet6_Model struct {
	Address types.List `tfsdk:"address"`
}

func (o Interfaces_Interface_Unit_Family_Inet6_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.ListType{ElemType: types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet6_Address_Model{}.AttrTypes()}},
	}
}
func (o Interfaces_Interface_Unit_Family_Inet6_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Interfaces_Interface_Unit_Family_Inet6_Address_Model{}.Attributes(),
			},
		},
	}
}

type Interfaces_Interface_Unit_Family_Inet_Address_Model struct {
	Name types.String `tfsdk:"name"`
}

func (o Interfaces_Interface_Unit_Family_Inet_Address_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}
func (o Interfaces_Interface_Unit_Family_Inet_Address_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Name.Address`",
		},
	}
}

type Interfaces_Interface_Unit_Family_Inet6_Address_Model struct {
	Name types.String `tfsdk:"name"`
}

func (o Interfaces_Interface_Unit_Family_Inet6_Address_Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}
func (o Interfaces_Interface_Unit_Family_Inet6_Address_Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Name.Address`",
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
			"interfaces": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: Interfaces_Model{}.Attributes(),
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

	var var_interfaces []Interfaces_Model
	if plan.Interfaces.IsNull() {
		var_interfaces = []Interfaces_Model{}
	} else {
		resp.Diagnostics.Append(plan.Interfaces.ElementsAs(ctx, &var_interfaces, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	config.Groups.Interfaces = make([]xml_Interfaces, len(var_interfaces))
	for i_interfaces, v_interfaces := range var_interfaces {
		var var_interfaces_interface []Interfaces_Interface_Model
		resp.Diagnostics.Append(v_interfaces.Interface.ElementsAs(ctx, &var_interfaces_interface, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Interfaces[i_interfaces].Interface = make([]xml_Interfaces_Interface, len(var_interfaces_interface))
		for i_interfaces_interface, v_interfaces_interface := range var_interfaces_interface {
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Name = v_interfaces_interface.Name.ValueStringPointer()
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Description = v_interfaces_interface.Description.ValueStringPointer()
			if v_interfaces_interface.Vlan_tagging.ValueBool() {
				empty := ""
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Vlan_tagging = &empty
			}
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Mtu = v_interfaces_interface.Mtu.ValueInt64Pointer()
			var var_interfaces_interface_unit []Interfaces_Interface_Unit_Model
			resp.Diagnostics.Append(v_interfaces_interface.Unit.ElementsAs(ctx, &var_interfaces_interface_unit, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit = make([]xml_Interfaces_Interface_Unit, len(var_interfaces_interface_unit))
			for i_interfaces_interface_unit, v_interfaces_interface_unit := range var_interfaces_interface_unit {
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Name = v_interfaces_interface_unit.Name.ValueStringPointer()
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Description = v_interfaces_interface_unit.Description.ValueStringPointer()
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Vlan_id = v_interfaces_interface_unit.Vlan_id.ValueStringPointer()
				var var_interfaces_interface_unit_family []Interfaces_Interface_Unit_Family_Model
				resp.Diagnostics.Append(v_interfaces_interface_unit.Family.ElementsAs(ctx, &var_interfaces_interface_unit_family, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family = make([]xml_Interfaces_Interface_Unit_Family, len(var_interfaces_interface_unit_family))
				for i_interfaces_interface_unit_family, v_interfaces_interface_unit_family := range var_interfaces_interface_unit_family {
					var var_interfaces_interface_unit_family_inet []Interfaces_Interface_Unit_Family_Inet_Model
					resp.Diagnostics.Append(v_interfaces_interface_unit_family.Inet.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet = make([]xml_Interfaces_Interface_Unit_Family_Inet, len(var_interfaces_interface_unit_family_inet))
					for i_interfaces_interface_unit_family_inet, v_interfaces_interface_unit_family_inet := range var_interfaces_interface_unit_family_inet {
						var var_interfaces_interface_unit_family_inet_address []Interfaces_Interface_Unit_Family_Inet_Address_Model
						resp.Diagnostics.Append(v_interfaces_interface_unit_family_inet.Address.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet_address, false)...)
						if resp.Diagnostics.HasError() {
							return
						}
						config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address = make([]xml_Interfaces_Interface_Unit_Family_Inet_Address, len(var_interfaces_interface_unit_family_inet_address))
						for i_interfaces_interface_unit_family_inet_address, v_interfaces_interface_unit_family_inet_address := range var_interfaces_interface_unit_family_inet_address {
							config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address[i_interfaces_interface_unit_family_inet_address].Name = v_interfaces_interface_unit_family_inet_address.Name.ValueStringPointer()
						}
					}
					var var_interfaces_interface_unit_family_inet6 []Interfaces_Interface_Unit_Family_Inet6_Model
					resp.Diagnostics.Append(v_interfaces_interface_unit_family.Inet6.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet6, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet6 = make([]xml_Interfaces_Interface_Unit_Family_Inet6, len(var_interfaces_interface_unit_family_inet6))
					for i_interfaces_interface_unit_family_inet6, v_interfaces_interface_unit_family_inet6 := range var_interfaces_interface_unit_family_inet6 {
						var var_interfaces_interface_unit_family_inet6_address []Interfaces_Interface_Unit_Family_Inet6_Address_Model
						resp.Diagnostics.Append(v_interfaces_interface_unit_family_inet6.Address.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet6_address, false)...)
						if resp.Diagnostics.HasError() {
							return
						}
						config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet6[i_interfaces_interface_unit_family_inet6].Address = make([]xml_Interfaces_Interface_Unit_Family_Inet6_Address, len(var_interfaces_interface_unit_family_inet6_address))
						for i_interfaces_interface_unit_family_inet6_address, v_interfaces_interface_unit_family_inet6_address := range var_interfaces_interface_unit_family_inet6_address {
							config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet6[i_interfaces_interface_unit_family_inet6].Address[i_interfaces_interface_unit_family_inet6_address].Name = v_interfaces_interface_unit_family_inet6_address.Name.ValueStringPointer()
						}
					}
				}
			}
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
	state.Interfaces = types.ListNull(types.ObjectType{AttrTypes: Groups_Model{}.AttrTypes()})
	interfaces_List := make([]Interfaces_Model, len(config.Groups.Interfaces))
	for i_interfaces, v_interfaces := range config.Groups.Interfaces {
		var interfaces_model Interfaces_Model
		interfaces_interface_List := make([]Interfaces_Interface_Model, len(v_interfaces.Interface))
		for i_interfaces_interface, v_interfaces_interface := range v_interfaces.Interface {
			var interfaces_interface_model Interfaces_Interface_Model
			interfaces_interface_model.Name = types.StringPointerValue(v_interfaces_interface.Name)
			interfaces_interface_model.Description = types.StringPointerValue(v_interfaces_interface.Description)
			interfaces_interface_model.Vlan_tagging = types.BoolValue(v_interfaces_interface.Vlan_tagging != nil)
			interfaces_interface_model.Mtu = types.Int64PointerValue(v_interfaces_interface.Mtu)

			interfaces_interface_unit_List := make([]Interfaces_Interface_Unit_Model, len(v_interfaces_interface.Unit))
			for i_interfaces_interface_unit, v_interfaces_interface_unit := range v_interfaces_interface.Unit {
				var interfaces_interface_unit_model Interfaces_Interface_Unit_Model
				interfaces_interface_unit_model.Name = types.StringPointerValue(v_interfaces_interface_unit.Name)
				interfaces_interface_unit_model.Description = types.StringPointerValue(v_interfaces_interface_unit.Description)
				interfaces_interface_unit_model.Vlan_id = types.StringPointerValue(v_interfaces_interface_unit.Vlan_id)

				interfaces_interface_unit_family_List := make([]Interfaces_Interface_Unit_Family_Model, len(v_interfaces_interface_unit.Family))
				for i_interfaces_interface_unit_family, v_interfaces_interface_unit_family := range v_interfaces_interface_unit.Family {
					var interfaces_interface_unit_family_model Interfaces_Interface_Unit_Family_Model

					interfaces_interface_unit_family_inet_List := make([]Interfaces_Interface_Unit_Family_Inet_Model, len(v_interfaces_interface_unit_family.Inet))
					for i_interfaces_interface_unit_family_inet, v_interfaces_interface_unit_family_inet := range v_interfaces_interface_unit_family.Inet {
						var interfaces_interface_unit_family_inet_model Interfaces_Interface_Unit_Family_Inet_Model
						interfaces_interface_unit_family_inet_List[i_interfaces_interface_unit_family_inet] = interfaces_interface_unit_family_inet_model

						interfaces_interface_unit_family_inet_address_List := make([]Interfaces_Interface_Unit_Family_Inet_Address_Model, len(v_interfaces_interface_unit_family_inet.Address))
						for i_interfaces_interface_unit_family_inet_address, v_interfaces_interface_unit_family_inet_address := range v_interfaces_interface_unit_family_inet.Address {
							var interfaces_interface_unit_family_inet_address_model Interfaces_Interface_Unit_Family_Inet_Address_Model
							interfaces_interface_unit_family_inet_address_model.Name = types.StringPointerValue(v_interfaces_interface_unit_family_inet_address.Name)
							interfaces_interface_unit_family_inet_address_List[i_interfaces_interface_unit_family_inet_address] = interfaces_interface_unit_family_inet_address_model
						}
						interfaces_interface_unit_family_inet_model.Address, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Address_Model{}.AttrTypes()}, interfaces_interface_unit_family_inet_address_List)
						interfaces_interface_unit_family_inet_List[i_interfaces_interface_unit_family_inet] = interfaces_interface_unit_family_inet_model
					}
					interfaces_interface_unit_family_model.Inet, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet_Model{}.AttrTypes()}, interfaces_interface_unit_family_inet_List)
					interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] = interfaces_interface_unit_family_model

					interfaces_interface_unit_family_inet6_List := make([]Interfaces_Interface_Unit_Family_Inet6_Model, len(v_interfaces_interface_unit_family.Inet6))
					for i_interfaces_interface_unit_family_inet6, v_interfaces_interface_unit_family_inet6 := range v_interfaces_interface_unit_family.Inet6 {
						var interfaces_interface_unit_family_inet6_model Interfaces_Interface_Unit_Family_Inet6_Model
						interfaces_interface_unit_family_inet6_List[i_interfaces_interface_unit_family_inet6] = interfaces_interface_unit_family_inet6_model

						interfaces_interface_unit_family_inet6_address_List := make([]Interfaces_Interface_Unit_Family_Inet6_Address_Model, len(v_interfaces_interface_unit_family_inet6.Address))
						for i_interfaces_interface_unit_family_inet6_address, v_interfaces_interface_unit_family_inet6_address := range v_interfaces_interface_unit_family_inet6.Address {
							var interfaces_interface_unit_family_inet6_address_model Interfaces_Interface_Unit_Family_Inet6_Address_Model
							interfaces_interface_unit_family_inet6_address_model.Name = types.StringPointerValue(v_interfaces_interface_unit_family_inet6_address.Name)
							interfaces_interface_unit_family_inet6_address_List[i_interfaces_interface_unit_family_inet6_address] = interfaces_interface_unit_family_inet6_address_model
						}
						interfaces_interface_unit_family_inet6_model.Address, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet6_Address_Model{}.AttrTypes()}, interfaces_interface_unit_family_inet6_address_List)
						interfaces_interface_unit_family_inet6_List[i_interfaces_interface_unit_family_inet6] = interfaces_interface_unit_family_inet6_model
					}
					interfaces_interface_unit_family_model.Inet6, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Inet6_Model{}.AttrTypes()}, interfaces_interface_unit_family_inet6_List)
					interfaces_interface_unit_family_List[i_interfaces_interface_unit_family] = interfaces_interface_unit_family_model
				}
				interfaces_interface_unit_model.Family, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Family_Model{}.AttrTypes()}, interfaces_interface_unit_family_List)
				interfaces_interface_unit_List[i_interfaces_interface_unit] = interfaces_interface_unit_model
			}
			interfaces_interface_model.Unit, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Unit_Model{}.AttrTypes()}, interfaces_interface_unit_List)
			interfaces_interface_List[i_interfaces_interface] = interfaces_interface_model
		}
		interfaces_model.Interface, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Interface_Model{}.AttrTypes()}, interfaces_interface_List)
		interfaces_List[i_interfaces] = interfaces_model
	}
	state.Interfaces, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: Interfaces_Model{}.AttrTypes()}, interfaces_List)
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

	var var_interfaces []Interfaces_Model
	if plan.Interfaces.IsNull() {
		var_interfaces = []Interfaces_Model{}
	} else {
		resp.Diagnostics.Append(plan.Interfaces.ElementsAs(ctx, &var_interfaces, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	config.Groups.Interfaces = make([]xml_Interfaces, len(var_interfaces))
	for i_interfaces, v_interfaces := range var_interfaces {
		var var_interfaces_interface []Interfaces_Interface_Model
		resp.Diagnostics.Append(v_interfaces.Interface.ElementsAs(ctx, &var_interfaces_interface, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Interfaces[i_interfaces].Interface = make([]xml_Interfaces_Interface, len(var_interfaces_interface))
		for i_interfaces_interface, v_interfaces_interface := range var_interfaces_interface {
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Name = v_interfaces_interface.Name.ValueStringPointer()
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Description = v_interfaces_interface.Description.ValueStringPointer()
			if v_interfaces_interface.Vlan_tagging.ValueBool() {
				empty := ""
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Vlan_tagging = &empty
			}
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Mtu = v_interfaces_interface.Mtu.ValueInt64Pointer()
			var var_interfaces_interface_unit []Interfaces_Interface_Unit_Model
			resp.Diagnostics.Append(v_interfaces_interface.Unit.ElementsAs(ctx, &var_interfaces_interface_unit, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit = make([]xml_Interfaces_Interface_Unit, len(var_interfaces_interface_unit))
			for i_interfaces_interface_unit, v_interfaces_interface_unit := range var_interfaces_interface_unit {
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Name = v_interfaces_interface_unit.Name.ValueStringPointer()
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Description = v_interfaces_interface_unit.Description.ValueStringPointer()
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Vlan_id = v_interfaces_interface_unit.Vlan_id.ValueStringPointer()
				var var_interfaces_interface_unit_family []Interfaces_Interface_Unit_Family_Model
				resp.Diagnostics.Append(v_interfaces_interface_unit.Family.ElementsAs(ctx, &var_interfaces_interface_unit_family, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family = make([]xml_Interfaces_Interface_Unit_Family, len(var_interfaces_interface_unit_family))
				for i_interfaces_interface_unit_family, v_interfaces_interface_unit_family := range var_interfaces_interface_unit_family {
					var var_interfaces_interface_unit_family_inet []Interfaces_Interface_Unit_Family_Inet_Model
					resp.Diagnostics.Append(v_interfaces_interface_unit_family.Inet.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet = make([]xml_Interfaces_Interface_Unit_Family_Inet, len(var_interfaces_interface_unit_family_inet))
					for i_interfaces_interface_unit_family_inet, v_interfaces_interface_unit_family_inet := range var_interfaces_interface_unit_family_inet {
						var var_interfaces_interface_unit_family_inet_address []Interfaces_Interface_Unit_Family_Inet_Address_Model
						resp.Diagnostics.Append(v_interfaces_interface_unit_family_inet.Address.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet_address, false)...)
						if resp.Diagnostics.HasError() {
							return
						}
						config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address = make([]xml_Interfaces_Interface_Unit_Family_Inet_Address, len(var_interfaces_interface_unit_family_inet_address))
						for i_interfaces_interface_unit_family_inet_address, v_interfaces_interface_unit_family_inet_address := range var_interfaces_interface_unit_family_inet_address {
							config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet[i_interfaces_interface_unit_family_inet].Address[i_interfaces_interface_unit_family_inet_address].Name = v_interfaces_interface_unit_family_inet_address.Name.ValueStringPointer()
						}
					}
					var var_interfaces_interface_unit_family_inet6 []Interfaces_Interface_Unit_Family_Inet6_Model
					resp.Diagnostics.Append(v_interfaces_interface_unit_family.Inet6.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet6, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet6 = make([]xml_Interfaces_Interface_Unit_Family_Inet6, len(var_interfaces_interface_unit_family_inet6))
					for i_interfaces_interface_unit_family_inet6, v_interfaces_interface_unit_family_inet6 := range var_interfaces_interface_unit_family_inet6 {
						var var_interfaces_interface_unit_family_inet6_address []Interfaces_Interface_Unit_Family_Inet6_Address_Model
						resp.Diagnostics.Append(v_interfaces_interface_unit_family_inet6.Address.ElementsAs(ctx, &var_interfaces_interface_unit_family_inet6_address, false)...)
						if resp.Diagnostics.HasError() {
							return
						}
						config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet6[i_interfaces_interface_unit_family_inet6].Address = make([]xml_Interfaces_Interface_Unit_Family_Inet6_Address, len(var_interfaces_interface_unit_family_inet6_address))
						for i_interfaces_interface_unit_family_inet6_address, v_interfaces_interface_unit_family_inet6_address := range var_interfaces_interface_unit_family_inet6_address {
							config.Groups.Interfaces[i_interfaces].Interface[i_interfaces_interface].Unit[i_interfaces_interface_unit].Family[i_interfaces_interface_unit_family].Inet6[i_interfaces_interface_unit_family_inet6].Address[i_interfaces_interface_unit_family_inet6_address].Name = v_interfaces_interface_unit_family_inet6_address.Name.ValueStringPointer()
						}
					}
				}
			}
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
}
