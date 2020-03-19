package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SpokeCluster represents the current status of spoke cluster.
// SpokeCluster is cluster scoped resources.
type SpokeCluster struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// Cluster name must conform in the definition of DNS label format
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec represents a desired configuration for the agent on the spoke cluster.
	Spec SpokeClusterSpec `json:"spec"`

	// Status represents the current status of joined spoke cluster
	Status SpokeClusterStatus `json:"status,omitempty"`
}

// SpokeClusterSpec represents a desired configuration for the agents on the spoke cluster.
type SpokeClusterSpec struct{}

// SpokeClusterStatus represents the current status of joined spoke cluster
type SpokeClusterStatus struct {
	// Conditions contains the different condition statuses for this spoke cluster.
	Conditions []StatusCondition `json:"conditions"`
}

const (
	// ClusterOK means that the cluster is "OK".
	SpokeClusterConditionOK string = "ClusterOK"
	// ClusterJoined means the spoke cluster has successfully joined the hub
	SpokeClusterConditionJoined string = "SpokeClusterJoined"
	// ClusterJoinApproved means the request to join the cluster is approved by user or controller
	SpokeClusterClusterConditionJoinApproved string = "HubApprovedJoin"
	// ClusterJoinDenied means the request to join the cluster is denied by user or controller
	SpokeClusterConditionJoinDenied string = "HubDeniedJoin"
)

// StatusCondition contains condition information for a spoke cluster.
type StatusCondition struct {
	// Type is the type of the cluster condition.
	// +required
	Type string `json:"type"`

	// Status is the status of the condition. One of True, False, Unknown.
	// +required
	Status metav1.ConditionStatus `json:"status"`

	// LastTransitionTime is the last time the condition changed from one status to another.
	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Reason is a (brief) reason for the condition's last status change.
	// +required
	Reason string `json:"reason"`

	// Message is a human-readable message indicating details about the last status change.
	// +required
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SpokeClusterList is a collection of spoke cluster
type SpokeClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of spoke cluster.
	Items []SpokeCluster `json:"items"`
}
