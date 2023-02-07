//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1beta1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AddOnConfig) DeepCopyInto(out *AddOnConfig) {
	*out = *in
	out.ConfigGroupResource = in.ConfigGroupResource
	out.ConfigReferent = in.ConfigReferent
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddOnConfig.
func (in *AddOnConfig) DeepCopy() *AddOnConfig {
	if in == nil {
		return nil
	}
	out := new(AddOnConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AddOnMeta) DeepCopyInto(out *AddOnMeta) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddOnMeta.
func (in *AddOnMeta) DeepCopy() *AddOnMeta {
	if in == nil {
		return nil
	}
	out := new(AddOnMeta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterManagementAddOn) DeepCopyInto(out *ClusterManagementAddOn) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterManagementAddOn.
func (in *ClusterManagementAddOn) DeepCopy() *ClusterManagementAddOn {
	if in == nil {
		return nil
	}
	out := new(ClusterManagementAddOn)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterManagementAddOn) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterManagementAddOnList) DeepCopyInto(out *ClusterManagementAddOnList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterManagementAddOn, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterManagementAddOnList.
func (in *ClusterManagementAddOnList) DeepCopy() *ClusterManagementAddOnList {
	if in == nil {
		return nil
	}
	out := new(ClusterManagementAddOnList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterManagementAddOnList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterManagementAddOnSpec) DeepCopyInto(out *ClusterManagementAddOnSpec) {
	*out = *in
	out.AddOnMeta = in.AddOnMeta
	if in.DefaultConfigs != nil {
		in, out := &in.DefaultConfigs, &out.DefaultConfigs
		*out = make([]AddOnConfig, len(*in))
		copy(*out, *in)
	}
	in.InstallStrategy.DeepCopyInto(&out.InstallStrategy)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterManagementAddOnSpec.
func (in *ClusterManagementAddOnSpec) DeepCopy() *ClusterManagementAddOnSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterManagementAddOnSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterManagementAddOnStatus) DeepCopyInto(out *ClusterManagementAddOnStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterManagementAddOnStatus.
func (in *ClusterManagementAddOnStatus) DeepCopy() *ClusterManagementAddOnStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterManagementAddOnStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigGroupResource) DeepCopyInto(out *ConfigGroupResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigGroupResource.
func (in *ConfigGroupResource) DeepCopy() *ConfigGroupResource {
	if in == nil {
		return nil
	}
	out := new(ConfigGroupResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigReference) DeepCopyInto(out *ConfigReference) {
	*out = *in
	out.ConfigGroupResource = in.ConfigGroupResource
	out.ConfigReferent = in.ConfigReferent
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigReference.
func (in *ConfigReference) DeepCopy() *ConfigReference {
	if in == nil {
		return nil
	}
	out := new(ConfigReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigReferent) DeepCopyInto(out *ConfigReferent) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigReferent.
func (in *ConfigReferent) DeepCopy() *ConfigReferent {
	if in == nil {
		return nil
	}
	out := new(ConfigReferent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HealthCheck) DeepCopyInto(out *HealthCheck) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HealthCheck.
func (in *HealthCheck) DeepCopy() *HealthCheck {
	if in == nil {
		return nil
	}
	out := new(HealthCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstallStrategy) DeepCopyInto(out *InstallStrategy) {
	*out = *in
	if in.Placements != nil {
		in, out := &in.Placements, &out.Placements
		*out = make([]PlacementStrategy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstallStrategy.
func (in *InstallStrategy) DeepCopy() *InstallStrategy {
	if in == nil {
		return nil
	}
	out := new(InstallStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterAddOn) DeepCopyInto(out *ManagedClusterAddOn) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterAddOn.
func (in *ManagedClusterAddOn) DeepCopy() *ManagedClusterAddOn {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterAddOn)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterAddOn) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterAddOnList) DeepCopyInto(out *ManagedClusterAddOnList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManagedClusterAddOn, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterAddOnList.
func (in *ManagedClusterAddOnList) DeepCopy() *ManagedClusterAddOnList {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterAddOnList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterAddOnList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterAddOnSpec) DeepCopyInto(out *ManagedClusterAddOnSpec) {
	*out = *in
	if in.Configs != nil {
		in, out := &in.Configs, &out.Configs
		*out = make([]AddOnConfig, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterAddOnSpec.
func (in *ManagedClusterAddOnSpec) DeepCopy() *ManagedClusterAddOnSpec {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterAddOnSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterAddOnStatus) DeepCopyInto(out *ManagedClusterAddOnStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RelatedObjects != nil {
		in, out := &in.RelatedObjects, &out.RelatedObjects
		*out = make([]ObjectReference, len(*in))
		copy(*out, *in)
	}
	out.AddOnMeta = in.AddOnMeta
	if in.SupportedConfigs != nil {
		in, out := &in.SupportedConfigs, &out.SupportedConfigs
		*out = make([]ConfigGroupResource, len(*in))
		copy(*out, *in)
	}
	if in.ConfigReferences != nil {
		in, out := &in.ConfigReferences, &out.ConfigReferences
		*out = make([]ConfigReference, len(*in))
		copy(*out, *in)
	}
	if in.Registrations != nil {
		in, out := &in.Registrations, &out.Registrations
		*out = make([]RegistrationConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.HealthCheck = in.HealthCheck
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterAddOnStatus.
func (in *ManagedClusterAddOnStatus) DeepCopy() *ManagedClusterAddOnStatus {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterAddOnStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectReference.
func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementRef) DeepCopyInto(out *PlacementRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementRef.
func (in *PlacementRef) DeepCopy() *PlacementRef {
	if in == nil {
		return nil
	}
	out := new(PlacementRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementStrategy) DeepCopyInto(out *PlacementStrategy) {
	*out = *in
	out.Placement = in.Placement
	if in.Configs != nil {
		in, out := &in.Configs, &out.Configs
		*out = make([]AddOnConfig, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementStrategy.
func (in *PlacementStrategy) DeepCopy() *PlacementStrategy {
	if in == nil {
		return nil
	}
	out := new(PlacementStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegistrationConfig) DeepCopyInto(out *RegistrationConfig) {
	*out = *in
	in.Subject.DeepCopyInto(&out.Subject)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegistrationConfig.
func (in *RegistrationConfig) DeepCopy() *RegistrationConfig {
	if in == nil {
		return nil
	}
	out := new(RegistrationConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Subject) DeepCopyInto(out *Subject) {
	*out = *in
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.OrganizationUnits != nil {
		in, out := &in.OrganizationUnits, &out.OrganizationUnits
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Subject.
func (in *Subject) DeepCopy() *Subject {
	if in == nil {
		return nil
	}
	out := new(Subject)
	in.DeepCopyInto(out)
	return out
}
