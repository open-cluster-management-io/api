//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AddOnPlacementScore) DeepCopyInto(out *AddOnPlacementScore) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddOnPlacementScore.
func (in *AddOnPlacementScore) DeepCopy() *AddOnPlacementScore {
	if in == nil {
		return nil
	}
	out := new(AddOnPlacementScore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AddOnPlacementScore) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AddOnPlacementScoreItem) DeepCopyInto(out *AddOnPlacementScoreItem) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddOnPlacementScoreItem.
func (in *AddOnPlacementScoreItem) DeepCopy() *AddOnPlacementScoreItem {
	if in == nil {
		return nil
	}
	out := new(AddOnPlacementScoreItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AddOnPlacementScoreList) DeepCopyInto(out *AddOnPlacementScoreList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AddOnPlacementScore, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddOnPlacementScoreList.
func (in *AddOnPlacementScoreList) DeepCopy() *AddOnPlacementScoreList {
	if in == nil {
		return nil
	}
	out := new(AddOnPlacementScoreList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AddOnPlacementScoreList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AddOnPlacementScoreStatus) DeepCopyInto(out *AddOnPlacementScoreStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Scores != nil {
		in, out := &in.Scores, &out.Scores
		*out = make([]AddOnPlacementScoreItem, len(*in))
		copy(*out, *in)
	}
	if in.ValidUntil != nil {
		in, out := &in.ValidUntil, &out.ValidUntil
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AddOnPlacementScoreStatus.
func (in *AddOnPlacementScoreStatus) DeepCopy() *AddOnPlacementScoreStatus {
	if in == nil {
		return nil
	}
	out := new(AddOnPlacementScoreStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterClaim) DeepCopyInto(out *ClusterClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterClaim.
func (in *ClusterClaim) DeepCopy() *ClusterClaim {
	if in == nil {
		return nil
	}
	out := new(ClusterClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterClaim) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterClaimList) DeepCopyInto(out *ClusterClaimList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterClaimList.
func (in *ClusterClaimList) DeepCopy() *ClusterClaimList {
	if in == nil {
		return nil
	}
	out := new(ClusterClaimList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterClaimList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterClaimSelector) DeepCopyInto(out *ClusterClaimSelector) {
	*out = *in
	if in.MatchExpressions != nil {
		in, out := &in.MatchExpressions, &out.MatchExpressions
		*out = make([]v1.LabelSelectorRequirement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterClaimSelector.
func (in *ClusterClaimSelector) DeepCopy() *ClusterClaimSelector {
	if in == nil {
		return nil
	}
	out := new(ClusterClaimSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterClaimSpec) DeepCopyInto(out *ClusterClaimSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterClaimSpec.
func (in *ClusterClaimSpec) DeepCopy() *ClusterClaimSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterClaimSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterDecision) DeepCopyInto(out *ClusterDecision) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterDecision.
func (in *ClusterDecision) DeepCopy() *ClusterDecision {
	if in == nil {
		return nil
	}
	out := new(ClusterDecision)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterPredicate) DeepCopyInto(out *ClusterPredicate) {
	*out = *in
	in.RequiredClusterSelector.DeepCopyInto(&out.RequiredClusterSelector)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterPredicate.
func (in *ClusterPredicate) DeepCopy() *ClusterPredicate {
	if in == nil {
		return nil
	}
	out := new(ClusterPredicate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSelector) DeepCopyInto(out *ClusterSelector) {
	*out = *in
	in.LabelSelector.DeepCopyInto(&out.LabelSelector)
	in.ClaimSelector.DeepCopyInto(&out.ClaimSelector)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSelector.
func (in *ClusterSelector) DeepCopy() *ClusterSelector {
	if in == nil {
		return nil
	}
	out := new(ClusterSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterSet) DeepCopyInto(out *ManagedClusterSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterSet.
func (in *ManagedClusterSet) DeepCopy() *ManagedClusterSet {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterSetBinding) DeepCopyInto(out *ManagedClusterSetBinding) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterSetBinding.
func (in *ManagedClusterSetBinding) DeepCopy() *ManagedClusterSetBinding {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterSetBinding)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterSetBinding) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterSetBindingList) DeepCopyInto(out *ManagedClusterSetBindingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManagedClusterSetBinding, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterSetBindingList.
func (in *ManagedClusterSetBindingList) DeepCopy() *ManagedClusterSetBindingList {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterSetBindingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterSetBindingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterSetBindingSpec) DeepCopyInto(out *ManagedClusterSetBindingSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterSetBindingSpec.
func (in *ManagedClusterSetBindingSpec) DeepCopy() *ManagedClusterSetBindingSpec {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterSetBindingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterSetList) DeepCopyInto(out *ManagedClusterSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManagedClusterSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterSetList.
func (in *ManagedClusterSetList) DeepCopy() *ManagedClusterSetList {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagedClusterSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterSetSpec) DeepCopyInto(out *ManagedClusterSetSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterSetSpec.
func (in *ManagedClusterSetSpec) DeepCopy() *ManagedClusterSetSpec {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagedClusterSetStatus) DeepCopyInto(out *ManagedClusterSetStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagedClusterSetStatus.
func (in *ManagedClusterSetStatus) DeepCopy() *ManagedClusterSetStatus {
	if in == nil {
		return nil
	}
	out := new(ManagedClusterSetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Placement) DeepCopyInto(out *Placement) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Placement.
func (in *Placement) DeepCopy() *Placement {
	if in == nil {
		return nil
	}
	out := new(Placement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Placement) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementDecision) DeepCopyInto(out *PlacementDecision) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementDecision.
func (in *PlacementDecision) DeepCopy() *PlacementDecision {
	if in == nil {
		return nil
	}
	out := new(PlacementDecision)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlacementDecision) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementDecisionList) DeepCopyInto(out *PlacementDecisionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PlacementDecision, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementDecisionList.
func (in *PlacementDecisionList) DeepCopy() *PlacementDecisionList {
	if in == nil {
		return nil
	}
	out := new(PlacementDecisionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlacementDecisionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementDecisionStatus) DeepCopyInto(out *PlacementDecisionStatus) {
	*out = *in
	if in.Decisions != nil {
		in, out := &in.Decisions, &out.Decisions
		*out = make([]ClusterDecision, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementDecisionStatus.
func (in *PlacementDecisionStatus) DeepCopy() *PlacementDecisionStatus {
	if in == nil {
		return nil
	}
	out := new(PlacementDecisionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementList) DeepCopyInto(out *PlacementList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Placement, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementList.
func (in *PlacementList) DeepCopy() *PlacementList {
	if in == nil {
		return nil
	}
	out := new(PlacementList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PlacementList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementSpec) DeepCopyInto(out *PlacementSpec) {
	*out = *in
	if in.ClusterSets != nil {
		in, out := &in.ClusterSets, &out.ClusterSets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.NumberOfClusters != nil {
		in, out := &in.NumberOfClusters, &out.NumberOfClusters
		*out = new(int32)
		**out = **in
	}
	if in.Predicates != nil {
		in, out := &in.Predicates, &out.Predicates
		*out = make([]ClusterPredicate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.PrioritizerPolicy.DeepCopyInto(&out.PrioritizerPolicy)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementSpec.
func (in *PlacementSpec) DeepCopy() *PlacementSpec {
	if in == nil {
		return nil
	}
	out := new(PlacementSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementStatus) DeepCopyInto(out *PlacementStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementStatus.
func (in *PlacementStatus) DeepCopy() *PlacementStatus {
	if in == nil {
		return nil
	}
	out := new(PlacementStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrioritizerAddOnScore) DeepCopyInto(out *PrioritizerAddOnScore) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrioritizerAddOnScore.
func (in *PrioritizerAddOnScore) DeepCopy() *PrioritizerAddOnScore {
	if in == nil {
		return nil
	}
	out := new(PrioritizerAddOnScore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrioritizerConfig) DeepCopyInto(out *PrioritizerConfig) {
	*out = *in
	out.PrioritizerScoreCoordinate = in.PrioritizerScoreCoordinate
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrioritizerConfig.
func (in *PrioritizerConfig) DeepCopy() *PrioritizerConfig {
	if in == nil {
		return nil
	}
	out := new(PrioritizerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrioritizerPolicy) DeepCopyInto(out *PrioritizerPolicy) {
	*out = *in
	if in.Configurations != nil {
		in, out := &in.Configurations, &out.Configurations
		*out = make([]PrioritizerConfig, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrioritizerPolicy.
func (in *PrioritizerPolicy) DeepCopy() *PrioritizerPolicy {
	if in == nil {
		return nil
	}
	out := new(PrioritizerPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrioritizerScoreCoordinate) DeepCopyInto(out *PrioritizerScoreCoordinate) {
	*out = *in
	out.AddOn = in.AddOn
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrioritizerScoreCoordinate.
func (in *PrioritizerScoreCoordinate) DeepCopy() *PrioritizerScoreCoordinate {
	if in == nil {
		return nil
	}
	out := new(PrioritizerScoreCoordinate)
	in.DeepCopyInto(out)
	return out
}
