package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedClusterAddOn is the Custom Resource object which holds the current state
// of an addon. This object is used by addon operators to convey their state.
// The resource should be created in the ManagedCluster's cluster namespace.
type ManagedClusterAddOn struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// spec holds configuration that could apply to any operator.
	// +kubebuilder:validation:Required
	// +required
	Spec ManagedClusterAddOnSpec `json:"spec"`

	// status holds the information about the state of an operator.  It is consistent with status information across
	// the Kubernetes ecosystem.
	// +optional
	Status ManagedClusterAddOnStatus `json:"status"`
}

// ManagedClusterAddOnSpec is empty for now, but you could imagine holding information like "pause".
type ManagedClusterAddOnSpec struct {
	// UpdateApproved represents that user has approved the update of the particular addon for a specific
	// the managed cluster on the hub.
	// The default value is false, it can only be set to true when the latestVersion and currentVersion
	// is on different version.
	// When the value is set true, the controller/operator that manages the addon will proceed to update the addon.
	// The value will be set to false by mutating webhook when latestVersion matches currentVersion.
	// +required
	UpdateApproved bool `json:"updateApproved"`
}

// ManagedClusterAddOnStatus provides information about the status of the operator.
// +k8s:deepcopy-gen=true
type ManagedClusterAddOnStatus struct {
	// conditions describes the state of the operator's managed and monitored components.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []AddOnStatusCondition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`

	// addonResource is a reference  to an object that contain the configuration of the addon
	// +required
	AddOnResource ObjectReference `json:"addonResource"`

	// LatestVersion indicates the latest available Version for the addon
	// +optional
	LatestVersion Release `json:"latestVersion,omitempty"`

	// CurrentVersion indicates the current Version of the addon
	// During the agent update process this field will be set to same value as latestVersion after the update has been completed
	// +optional
	CurrentVersion Release `json:"currentVersion,omitempty"`

	// UpdateAvailable indicates if there is an version update that is available for the addon
	// it will be set to false if currentVersion.version != latestVersion.version
	// it will be set to true if currentVersion.version == latestVersion.version
	// +optional
	UpdateAvailable boolean `json:"updateAvailable,omitempty"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// group of the referent.
	// +kubebuilder:validation:Required
	// +required
	Group string `json:"group"`
	// resource of the referent.
	// +kubebuilder:validation:Required
	// +required
	Resource string `json:"resource"`
	// name of the referent.
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`
}

// AddOnStatusCondition represents the state of the addon's
// managed and monitored components.
// +k8s:deepcopy-gen=true
type AddOnStatusCondition struct {
	// type specifies the aspect reported by this condition.
	// +kubebuilder:validation:Required
	// +required
	Type AddOnStatusConditionType `json:"type"`

	// status of the condition, one of True, False, Unknown.
	// +kubebuilder:validation:Required
	// +required
	Status metav1.ConditionStatus `json:"status"`

	// lastTransitionTime is the time of the last update to the current status property.
	// +kubebuilder:validation:Required
	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// reason is the CamelCase reason for the condition's current status.
	// +optional
	Reason string `json:"reason,omitempty"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	// +optional
	Message string `json:"message,omitempty"`
}

// Release represents an agent release version and it's related images.
// +k8s:deepcopy-gen=true
type Release struct {
	// version is a semantic versioning identifying the update version. When this
	// field is part of spec, version is optional if image is specified.
	// +required
	Version string `json:"version"`

	// RelatedImages is list of images will be used for the specific version
	// this will be use for backward compatability (i.e hub updated but agent have not)
	// +required
	RelatedImages []RelatedImage `json:"relatedImages"`
}

// RelatedImage represents information for one of the images that the agent uses its associated key.
// +k8s:deepcopy-gen=true
type RelatedImage struct {
	// ImageKey is the key to link the image to specific deployment for the agent
	// this will be use for backward compatability, i.e if the hub updated but agent have not the agent
	// controller can use this information to continue to manage the agent.
	// +required
	ImageKey string `json:"imageKey"`

	// ImagePullSpec represents the desired image of the component for the addon agent.
	// +required
	ImagePullSpec string `json:"image"`
}

// AddOnStatusConditionType is an aspect of agent state.
type AddOnStatusConditionType string

const (
	// Available indicates that the agent is functional and available in the cluster.
	Available AddOnStatusConditionType = "Available"

	// Progressing indicates that the operator is actively rolling out new code,
	// propagating config changes, or otherwise moving from one steady state to
	// another.  Operators should not report progressing when they are reconciling
	// a previously known state.
	Progressing AddOnStatusConditionType = "Progressing"

	// Degraded indicates that the operator's current state does not match its
	// desired state over a period of time resulting in a lower quality of service.
	// The period of time may vary by component, but a Degraded state represents
	// persistent observation of a condition.  As a result, a component should not
	// oscillate in and out of Degraded state.  A service may be Available even
	// if its degraded.  For example, your service may desire 3 running pods, but 1
	// pod is crash-looping.  The service is Available but Degraded because it
	// may have a lower quality of service.  A component may be Progressing but
	// not Degraded because the transition from one state to another does not
	// persist over a long enough period to report Degraded.  A service should not
	// report Degraded during the course of a normal upgrade.  A service may report
	// Degraded in response to a persistent infrastructure failure that requires
	// administrator intervention.  For example, if a control plane host is unhealthy
	// and must be replaced.  An operator should report Degraded if unexpected
	// errors occur over a period, but the expectation is that all unexpected errors
	// are handled as operators mature.
	Degraded AddOnStatusConditionType = "Degraded"
)

// ManagedClusterAddOnList is a list of ManagedClusterAddOn resources.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ManagedClusterAddOnList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ManagedClusterAddOn `json:"items"`
}
