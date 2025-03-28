package main

import (
	"context"
	"encoding/xml"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Junos XML Hierarchy
type xmlInterfaces struct {
	XMLName xml.Name `xml:"configuration"`
	Groups  struct {
		XMLName   xml.Name       `xml:"groups"`
		Name      string         `xml:"name"`
		Interface []xmlInterface `xml:"interfaces>interface,omitempty"`
	} `xml:"groups"`
}
type xmlInterface struct {
	XMLName      xml.Name           `xml:"interface"`
	Name         *string            `xml:"name,omitempty"`
	Description  *string            `xml:"description,omitempty"`
	Vlan_tagging *string            `xml:"vlan-tagging,omitempty"`
	Mtu          *int64             `xml:"mtu,omitempty"`
	Unit         []xmlInterfaceUnit `xml:"unit,omitempty"`
}
type xmlInterfaceUnit struct {
	XMLName     xml.Name `xml:"unit"`
	Name        *string  `xml:"name,omitempty"`
	Description *string  `xml:"description,omitempty"`
	Vlan_id     *string  `xml:"vlan-id,omitempty"`
	Family      struct {
		XMLName xml.Name `xml:"family"`
		Inet    struct {
			XMLName xml.Name     `xml:"inet"`
			Address []xmlAddress `xml:"address"`
		} `xml:"inet"`
		Inet6 struct {
			XMLName xml.Name     `xml:"inet6"`
			Address []xmlAddress `xml:"address,omitempty"`
		} `xml:"inet6"`
	} `xml:"family"`
}
type xmlAddress struct {
	XMLName xml.Name `xml:"address"`
	Name    *string  `xml:"name,omitempty"`
}

// Collecting objects from the .tf file
type InterfacesModel struct {
	ResourceName types.String `tfsdk:"resource_name"`
	Interface    types.List   `tfsdk:"interface"`
}

func (o InterfacesModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface": types.ListType{ElemType: types.ObjectType{AttrTypes: InterfaceModel{}.AttrTypes()}},
	}
}
func (o InterfacesModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"interface": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: InterfaceModel{}.Attributes(),
			},
		},
	}
}

// InterfacesModel is the model for the resource, which contains the interface configuration
type InterfaceModel struct {
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Vlan_tagging types.Bool   `tfsdk:"vlan_tagging"`
	Mtu          types.Int64  `tfsdk:"mtu"`
	Unit         types.List   `tfsdk:"unit"`
}

func (o InterfaceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"description":  types.StringType,
		"vlan_tagging": types.BoolType,
		"mtu":          types.Int64Type,
		"unit":         types.ListType{ElemType: types.ObjectType{AttrTypes: UnitModel{}.AttrTypes()}},
	}
}
func (o InterfaceModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpth is `config.Groups.Interface.Name",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpth is `config.Groups.Interface.Description",
		},
		"vlan_tagging": schema.BoolAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Vlan_tagging",
		},
		"mtu": schema.Int64Attribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Interface.Mtu",
		},
		"unit": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: UnitModel{}.Attributes(),
			},
		},
	}
}

type UnitModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Vlan_id     types.String `tfsdk:"vlan_id"`
	Family      types.List   `tfsdk:"family"`
}

func (o UnitModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"vlan_id":     types.StringType,
		"family":      types.ListType{ElemType: types.ObjectType{AttrTypes: FamilyModel{}.AttrTypes()}},
	}
}
func (o UnitModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Unit.Name`",
		},
		"description": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Unit.Description`",
		},
		"vlan_id": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Unit.Vlan_id`",
		},
		"family": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: FamilyModel{}.Attributes(),
			},
		},
	}
}

type FamilyModel struct {
	Inet  types.List `tfsdk:"inet"`
	Inet6 types.List `tfsdk:"inet6"`
}

func (o FamilyModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"inet":  types.ListType{ElemType: types.ObjectType{AttrTypes: InetModel{}.AttrTypes()}},
		"inet6": types.ListType{ElemType: types.ObjectType{AttrTypes: Inet6Model{}.AttrTypes()}},
	}
}
func (o FamilyModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"inet": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: InetModel{}.Attributes(),
			},
		},
		"inet6": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: Inet6Model{}.Attributes(),
			},
		},
	}
}

type InetModel struct {
	Address types.List `tfsdk:"address"`
}

func (o InetModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.ListType{ElemType: types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}},
	}
}
func (o InetModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: AddressModel{}.Attributes(),
			},
		},
	}
}

type Inet6Model struct {
	Address types.List `tfsdk:"address"`
}

func (o Inet6Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address": types.ListType{ElemType: types.ObjectType{AttrTypes: AddressModel{}.AttrTypes()}},
	}
}
func (o Inet6Model) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"address": schema.ListNestedAttribute{
			Optional: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: AddressModel{}.Attributes(),
			},
		},
	}
}

type AddressModel struct {
	Name types.String `tfsdk:"name"`
}

func (o AddressModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}
func (o AddressModel) Attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "xpath is `config.Groups.Address.Name`",
		},
	}
}

// Collects the data for the crud work
type resourceInterfaces struct {
	client ProviderConfig
}

func (r *resourceInterfaces) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(ProviderConfig)
}

// Metadata implements resource.Resource.
func (r *resourceInterfaces) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_Interfaces"
}

// Schema implements resource.Resource.
func (r *resourceInterfaces) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"interface": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: InterfaceModel{}.Attributes(),
				},
			},
		},
	}
}

// Create implements resource.Resource.
func (r *resourceInterfaces) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Get the Interfaces Model data and set
	var plan InterfacesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the XML configuration for the interface
	var config xmlInterfaces
	config.Groups.Name = plan.ResourceName.ValueString()

	var var_interface []InterfaceModel
	resp.Diagnostics.Append(plan.Interface.ElementsAs(ctx, &var_interface, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface = make([]xmlInterface, len(var_interface))
	for i, m_interface := range var_interface {
		config.Groups.Interface[i].Name = m_interface.Name.ValueStringPointer()
		config.Groups.Interface[i].Description = m_interface.Description.ValueStringPointer()
		if m_interface.Vlan_tagging.ValueBool() {
			empty := ""
			config.Groups.Interface[i].Vlan_tagging = &empty
		}
		config.Groups.Interface[i].Mtu = m_interface.Mtu.ValueInt64Pointer()

		var var_interface_unit []UnitModel
		resp.Diagnostics.Append(m_interface.Unit.ElementsAs(ctx, &var_interface_unit, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Interface[i].Unit = make([]xmlInterfaceUnit, len(var_interface_unit))
		for y, m_interface_unit := range var_interface_unit {
			config.Groups.Interface[i].Unit[y].Name = m_interface_unit.Name.ValueStringPointer()
			config.Groups.Interface[i].Unit[y].Description = m_interface_unit.Description.ValueStringPointer()
			config.Groups.Interface[i].Unit[y].Vlan_id = m_interface_unit.Vlan_id.ValueStringPointer()

			var var_unit_family []FamilyModel
			resp.Diagnostics.Append(m_interface_unit.Family.ElementsAs(ctx, &var_unit_family, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			if len(var_unit_family) > 0 {
				var var_family_inet []InetModel
				resp.Diagnostics.Append(var_unit_family[0].Inet.ElementsAs(ctx, &var_family_inet, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				if len(var_family_inet) > 0 {
					var var_inet_address []AddressModel
					resp.Diagnostics.Append(var_family_inet[0].Address.ElementsAs(ctx, &var_inet_address, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interface[i].Unit[y].Family.Inet.Address = make([]xmlAddress, len(var_inet_address))
					for z, inet_address := range var_inet_address {
						config.Groups.Interface[i].Unit[y].Family.Inet.Address[z].Name = inet_address.Name.ValueStringPointer()
					}
				}
			}
			var var_family_inet6 []Inet6Model
			resp.Diagnostics.Append(var_unit_family[0].Inet6.ElementsAs(ctx, &var_family_inet6, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			if len(var_family_inet6) > 0 {
				var var_inet6_address []AddressModel
				resp.Diagnostics.Append(var_family_inet6[0].Address.ElementsAs(ctx, &var_inet6_address, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interface[i].Unit[y].Family.Inet6.Address = make([]xmlAddress, len(var_inet6_address))
				for n, inet6_address := range var_inet6_address {
					config.Groups.Interface[i].Unit[y].Family.Inet6.Address[n].Name = inet6_address.Name.ValueStringPointer()
				}
			}
		}
	}
	err := r.client.SendTransaction("", config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read implements resource.Resource.
func (r *resourceInterfaces) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

// Update implements resource.Resource.
func (r *resourceInterfaces) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get the Interfaces Model data and set
	var plan InterfacesModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// Check for errors
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the XML configuration for the interface
	var config xmlInterfaces
	config.Groups.Name = plan.ResourceName.ValueString()

	var var_interface []InterfaceModel
	resp.Diagnostics.Append(plan.Interface.ElementsAs(ctx, &var_interface, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups.Interface = make([]xmlInterface, len(var_interface))
	for i, m_interface := range var_interface {
		config.Groups.Interface[i].Name = m_interface.Name.ValueStringPointer()
		config.Groups.Interface[i].Description = m_interface.Description.ValueStringPointer()
		if m_interface.Vlan_tagging.ValueBool() {
			empty := ""
			config.Groups.Interface[i].Vlan_tagging = &empty
		}
		config.Groups.Interface[i].Mtu = m_interface.Mtu.ValueInt64Pointer()

		var var_interface_unit []UnitModel
		resp.Diagnostics.Append(m_interface.Unit.ElementsAs(ctx, &var_interface_unit, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Groups.Interface[i].Unit = make([]xmlInterfaceUnit, len(var_interface_unit))
		for y, m_interface_unit := range var_interface_unit {
			config.Groups.Interface[i].Unit[y].Name = m_interface_unit.Name.ValueStringPointer()
			config.Groups.Interface[i].Unit[y].Description = m_interface_unit.Description.ValueStringPointer()
			config.Groups.Interface[i].Unit[y].Vlan_id = m_interface_unit.Vlan_id.ValueStringPointer()

			var var_unit_family []FamilyModel
			resp.Diagnostics.Append(m_interface_unit.Family.ElementsAs(ctx, &var_unit_family, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			if len(var_unit_family) > 0 {
				var var_family_inet []InetModel
				resp.Diagnostics.Append(var_unit_family[0].Inet.ElementsAs(ctx, &var_family_inet, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				if len(var_family_inet) > 0 {
					var var_inet_address []AddressModel
					resp.Diagnostics.Append(var_family_inet[0].Address.ElementsAs(ctx, &var_inet_address, false)...)
					if resp.Diagnostics.HasError() {
						return
					}
					config.Groups.Interface[i].Unit[y].Family.Inet.Address = make([]xmlAddress, len(var_inet_address))
					for z, inet_address := range var_inet_address {
						config.Groups.Interface[i].Unit[y].Family.Inet.Address[z].Name = inet_address.Name.ValueStringPointer()
					}
				}
			}
			var var_family_inet6 []Inet6Model
			resp.Diagnostics.Append(var_unit_family[0].Inet6.ElementsAs(ctx, &var_family_inet6, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			if len(var_family_inet6) > 0 {
				var var_inet6_address []AddressModel
				resp.Diagnostics.Append(var_family_inet6[0].Address.ElementsAs(ctx, &var_inet6_address, false)...)
				if resp.Diagnostics.HasError() {
					return
				}
				config.Groups.Interface[i].Unit[y].Family.Inet6.Address = make([]xmlAddress, len(var_inet6_address))
				for n, inet6_address := range var_inet6_address {
					config.Groups.Interface[i].Unit[y].Family.Inet6.Address[n].Name = inet6_address.Name.ValueStringPointer()
				}
			}
		}
	}
	err := r.client.SendTransaction("", config, false)
	if err != nil {
		resp.Diagnostics.AddError("Failed while Sending", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete implements resource.Resource.
func (r *resourceInterfaces) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state InterfacesModel
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
		resp.Diagnostics.AddError("Failed while deleting dile", err.Error())
		return
	}
}
