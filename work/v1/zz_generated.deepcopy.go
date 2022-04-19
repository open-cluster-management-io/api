//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedManifestResourceMeta) DeepCopyInto(out *AppliedManifestResourceMeta) {
	*out = *in
	out.ResourceIdentifier = in.ResourceIdentifier
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedManifestResourceMeta.
func (in *AppliedManifestResourceMeta) DeepCopy() *AppliedManifestResourceMeta {
	if in == nil {
		return nil
	}
	out := new(AppliedManifestResourceMeta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedManifestWork) DeepCopyInto(out *AppliedManifestWork) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedManifestWork.
func (in *AppliedManifestWork) DeepCopy() *AppliedManifestWork {
	if in == nil {
		return nil
	}
	out := new(AppliedManifestWork)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppliedManifestWork) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedManifestWorkList) DeepCopyInto(out *AppliedManifestWorkList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppliedManifestWork, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedManifestWorkList.
func (in *AppliedManifestWorkList) DeepCopy() *AppliedManifestWorkList {
	if in == nil {
		return nil
	}
	out := new(AppliedManifestWorkList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppliedManifestWorkList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedManifestWorkSpec) DeepCopyInto(out *AppliedManifestWorkSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedManifestWorkSpec.
func (in *AppliedManifestWorkSpec) DeepCopy() *AppliedManifestWorkSpec {
	if in == nil {
		return nil
	}
	out := new(AppliedManifestWorkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppliedManifestWorkStatus) DeepCopyInto(out *AppliedManifestWorkStatus) {
	*out = *in
	if in.AppliedResources != nil {
		in, out := &in.AppliedResources, &out.AppliedResources
		*out = make([]AppliedManifestResourceMeta, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppliedManifestWorkStatus.
func (in *AppliedManifestWorkStatus) DeepCopy() *AppliedManifestWorkStatus {
	if in == nil {
		return nil
	}
	out := new(AppliedManifestWorkStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeleteOption) DeepCopyInto(out *DeleteOption) {
	*out = *in
	if in.SelectivelyOrphan != nil {
		in, out := &in.SelectivelyOrphan, &out.SelectivelyOrphan
		*out = new(SelectivelyOrphan)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeleteOption.
func (in *DeleteOption) DeepCopy() *DeleteOption {
	if in == nil {
		return nil
	}
	out := new(DeleteOption)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FeedbackRule) DeepCopyInto(out *FeedbackRule) {
	*out = *in
	if in.JsonPaths != nil {
		in, out := &in.JsonPaths, &out.JsonPaths
		*out = make([]JsonPath, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FeedbackRule.
func (in *FeedbackRule) DeepCopy() *FeedbackRule {
	if in == nil {
		return nil
	}
	out := new(FeedbackRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FeedbackValue) DeepCopyInto(out *FeedbackValue) {
	*out = *in
	in.Value.DeepCopyInto(&out.Value)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FeedbackValue.
func (in *FeedbackValue) DeepCopy() *FeedbackValue {
	if in == nil {
		return nil
	}
	out := new(FeedbackValue)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FieldValue) DeepCopyInto(out *FieldValue) {
	*out = *in
	if in.Integer != nil {
		in, out := &in.Integer, &out.Integer
		*out = new(int64)
		**out = **in
	}
	if in.String != nil {
		in, out := &in.String, &out.String
		*out = new(string)
		**out = **in
	}
	if in.Boolean != nil {
		in, out := &in.Boolean, &out.Boolean
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FieldValue.
func (in *FieldValue) DeepCopy() *FieldValue {
	if in == nil {
		return nil
	}
	out := new(FieldValue)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JsonPath) DeepCopyInto(out *JsonPath) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JsonPath.
func (in *JsonPath) DeepCopy() *JsonPath {
	if in == nil {
		return nil
	}
	out := new(JsonPath)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Manifest) DeepCopyInto(out *Manifest) {
	*out = *in
	in.RawExtension.DeepCopyInto(&out.RawExtension)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Manifest.
func (in *Manifest) DeepCopy() *Manifest {
	if in == nil {
		return nil
	}
	out := new(Manifest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestCondition) DeepCopyInto(out *ManifestCondition) {
	*out = *in
	out.ResourceMeta = in.ResourceMeta
	in.StatusFeedbacks.DeepCopyInto(&out.StatusFeedbacks)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestCondition.
func (in *ManifestCondition) DeepCopy() *ManifestCondition {
	if in == nil {
		return nil
	}
	out := new(ManifestCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestConfigOption) DeepCopyInto(out *ManifestConfigOption) {
	*out = *in
	out.ResourceIdentifier = in.ResourceIdentifier
	if in.FeedbackRules != nil {
		in, out := &in.FeedbackRules, &out.FeedbackRules
		*out = make([]FeedbackRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.UpdateStrategy.DeepCopyInto(&out.UpdateStrategy)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestConfigOption.
func (in *ManifestConfigOption) DeepCopy() *ManifestConfigOption {
	if in == nil {
		return nil
	}
	out := new(ManifestConfigOption)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestResourceMeta) DeepCopyInto(out *ManifestResourceMeta) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestResourceMeta.
func (in *ManifestResourceMeta) DeepCopy() *ManifestResourceMeta {
	if in == nil {
		return nil
	}
	out := new(ManifestResourceMeta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestResourceStatus) DeepCopyInto(out *ManifestResourceStatus) {
	*out = *in
	if in.Manifests != nil {
		in, out := &in.Manifests, &out.Manifests
		*out = make([]ManifestCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestResourceStatus.
func (in *ManifestResourceStatus) DeepCopy() *ManifestResourceStatus {
	if in == nil {
		return nil
	}
	out := new(ManifestResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestWork) DeepCopyInto(out *ManifestWork) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestWork.
func (in *ManifestWork) DeepCopy() *ManifestWork {
	if in == nil {
		return nil
	}
	out := new(ManifestWork)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManifestWork) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestWorkExecutor) DeepCopyInto(out *ManifestWorkExecutor) {
	*out = *in
	in.Subject.DeepCopyInto(&out.Subject)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestWorkExecutor.
func (in *ManifestWorkExecutor) DeepCopy() *ManifestWorkExecutor {
	if in == nil {
		return nil
	}
	out := new(ManifestWorkExecutor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestWorkExecutorSubject) DeepCopyInto(out *ManifestWorkExecutorSubject) {
	*out = *in
	if in.ServiceAccount != nil {
		in, out := &in.ServiceAccount, &out.ServiceAccount
		*out = new(ManifestWorkSubjectServiceAccount)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestWorkExecutorSubject.
func (in *ManifestWorkExecutorSubject) DeepCopy() *ManifestWorkExecutorSubject {
	if in == nil {
		return nil
	}
	out := new(ManifestWorkExecutorSubject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestWorkList) DeepCopyInto(out *ManifestWorkList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManifestWork, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestWorkList.
func (in *ManifestWorkList) DeepCopy() *ManifestWorkList {
	if in == nil {
		return nil
	}
	out := new(ManifestWorkList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManifestWorkList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestWorkSpec) DeepCopyInto(out *ManifestWorkSpec) {
	*out = *in
	in.Workload.DeepCopyInto(&out.Workload)
	if in.DeleteOption != nil {
		in, out := &in.DeleteOption, &out.DeleteOption
		*out = new(DeleteOption)
		(*in).DeepCopyInto(*out)
	}
	if in.ManifestConfigs != nil {
		in, out := &in.ManifestConfigs, &out.ManifestConfigs
		*out = make([]ManifestConfigOption, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Executor != nil {
		in, out := &in.Executor, &out.Executor
		*out = new(ManifestWorkExecutor)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestWorkSpec.
func (in *ManifestWorkSpec) DeepCopy() *ManifestWorkSpec {
	if in == nil {
		return nil
	}
	out := new(ManifestWorkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestWorkStatus) DeepCopyInto(out *ManifestWorkStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestWorkStatus.
func (in *ManifestWorkStatus) DeepCopy() *ManifestWorkStatus {
	if in == nil {
		return nil
	}
	out := new(ManifestWorkStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestWorkSubjectServiceAccount) DeepCopyInto(out *ManifestWorkSubjectServiceAccount) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestWorkSubjectServiceAccount.
func (in *ManifestWorkSubjectServiceAccount) DeepCopy() *ManifestWorkSubjectServiceAccount {
	if in == nil {
		return nil
	}
	out := new(ManifestWorkSubjectServiceAccount)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManifestsTemplate) DeepCopyInto(out *ManifestsTemplate) {
	*out = *in
	if in.Manifests != nil {
		in, out := &in.Manifests, &out.Manifests
		*out = make([]Manifest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManifestsTemplate.
func (in *ManifestsTemplate) DeepCopy() *ManifestsTemplate {
	if in == nil {
		return nil
	}
	out := new(ManifestsTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OrphaningRule) DeepCopyInto(out *OrphaningRule) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OrphaningRule.
func (in *OrphaningRule) DeepCopy() *OrphaningRule {
	if in == nil {
		return nil
	}
	out := new(OrphaningRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceIdentifier) DeepCopyInto(out *ResourceIdentifier) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceIdentifier.
func (in *ResourceIdentifier) DeepCopy() *ResourceIdentifier {
	if in == nil {
		return nil
	}
	out := new(ResourceIdentifier)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SelectivelyOrphan) DeepCopyInto(out *SelectivelyOrphan) {
	*out = *in
	if in.OrphaningRules != nil {
		in, out := &in.OrphaningRules, &out.OrphaningRules
		*out = make([]OrphaningRule, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SelectivelyOrphan.
func (in *SelectivelyOrphan) DeepCopy() *SelectivelyOrphan {
	if in == nil {
		return nil
	}
	out := new(SelectivelyOrphan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerSideApplyConfig) DeepCopyInto(out *ServerSideApplyConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerSideApplyConfig.
func (in *ServerSideApplyConfig) DeepCopy() *ServerSideApplyConfig {
	if in == nil {
		return nil
	}
	out := new(ServerSideApplyConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatusFeedbackResult) DeepCopyInto(out *StatusFeedbackResult) {
	*out = *in
	if in.Values != nil {
		in, out := &in.Values, &out.Values
		*out = make([]FeedbackValue, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatusFeedbackResult.
func (in *StatusFeedbackResult) DeepCopy() *StatusFeedbackResult {
	if in == nil {
		return nil
	}
	out := new(StatusFeedbackResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UpdateStrategy) DeepCopyInto(out *UpdateStrategy) {
	*out = *in
	if in.ServerSideApply != nil {
		in, out := &in.ServerSideApply, &out.ServerSideApply
		*out = new(ServerSideApplyConfig)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UpdateStrategy.
func (in *UpdateStrategy) DeepCopy() *UpdateStrategy {
	if in == nil {
		return nil
	}
	out := new(UpdateStrategy)
	in.DeepCopyInto(out)
	return out
}
