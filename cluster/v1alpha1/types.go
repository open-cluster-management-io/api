package v1alpha1

import (
	clusterv1 "github.com/open-cluster-management/api/cluster/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status

// ManagedClusterClaim is a user's request for and claim to a managed cluster.
// It is defined as an ownership claim of a managed cluster in a namespace.
type ManagedClusterClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the managed cluster requested by the claim
	Spec ManagedClusterClaimSpec `json:"spec"`

	// Status represents the current status of the claim
	// +optional
	Status ManagedClusterClaimStatus `json:"status,omitempty"`
}

// ManagedClusterClaimSpec describes the common attributes of managed cluster
type ManagedClusterClaimSpec struct {
	// A label query over managed clusters to consider for binding. This selector is
	// ignored when ClusterName is set
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
	// ClusterName is used to match a concrete managed cluster for this
	// claim. When set to non-empty value Selector is not evaluated
	// +optional
	ClusterName string `json:"clusterName,omitempty"`
}

// ManagedClusterClaimStatus represents the status of ManagedCluster claim
type ManagedClusterClaimStatus struct {
	// ClusterName is the binding reference to the ManagedCluster backing this
	// claim
	// +optional
	ClusterName string `json:"clusterName,omitempty"`
	// +optional
	Conditions []ManagedClusterClaimCondition `json:"conditions,omitempty"`
}

// ManagedClusterClaimCondition represents the current condition of ManagedCluster claim
type ManagedClusterClaimCondition struct {
	// Type is the type of the ManagedClusterClaim condition.
	// +required
	Type string `json:"type"`
	// Status is the status of the condition. One of True, False, Unknown.
	// +required
	Status metav1.ConditionStatus `json:"status"`
	// LastTransitionTime is the last time the condition changed from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is a (brief) reason for the condition's last status change.
	// +required
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about the last status change.
	// +required
	Message string `json:"message"`
}

// ManagedClusterClaimConditionType defines the condition of ManagedCluster claim.
type ManagedClusterClaimConditionType string

// These are valid conditions of ManagedClusterClaim
const (
	// ManagedClusterClaim is bound
	ClaimBound ManagedClusterClaimConditionType = "ClaimBound"

	// A mirrored managed cluster is created in the same namespace
	MirroredClusterCreated ManagedClusterClaimConditionType = "MirroredClusterCreated"

	// The mirrored managed cluster is synced with the source managed cluster
	MirroredClusterSynced ManagedClusterClaimConditionType = "MirroredClusterSynced"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedClusterClaimList is a collection of ManagedClusterClaims.
type ManagedClusterClaimList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of ManagedClusterClaims.
	Items []ManagedClusterClaim `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=`.spec.hubAcceptsClient`,name="Hub Accepted",type=boolean
// +kubebuilder:printcolumn:JSONPath=`.spec.managedClusterClientConfigs[*].url`,name="Managed Cluster URLs",type=string
// +kubebuilder:printcolumn:JSONPath=`.status.conditions[?(@.type=="ManagedClusterJoined")].status`,name="Joined",type=string
// +kubebuilder:printcolumn:JSONPath=`.status.conditions[?(@.type=="ManagedClusterConditionAvailable")].status`,name="Available",type=string
// +kubebuilder:printcolumn:JSONPath=`.metadata.creationTimestamp`,name="Age",type=date

// MirroredManagedCluster is a mirror of a managed cluster in namespace scope.
// It is ready only and keeps synced with the managed cluster it mirrored.
type MirroredManagedCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents a desired configuration for the agent on the managed cluster.
	Spec clusterv1.ManagedClusterSpec `json:"spec"`

	// Status represents the current status of joined managed cluster
	// +optional
	Status clusterv1.ManagedClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MirroredManagedClusterList is a collection of mirrored managed clusters.
type MirroredManagedClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of mirrored managed clusters.
	Items []MirroredManagedCluster `json:"items"`
}
