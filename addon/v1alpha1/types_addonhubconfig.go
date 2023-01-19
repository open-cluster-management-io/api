package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Cluster"

// AddOnHubConfig represents a hub-scoped configuration for an add-on.
type AddOnHubConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec represents a desired configuration for an add-on.
	// +required
	Spec AddOnHubConfigSpec `json:"spec"`

	// status represents the current status of the configuration for an add-on.
	// +optional
	Status AddOnHubConfigStatus `json:"status,omitempty"`
}

type AddOnHubConfigSpec struct {
	// version represents the desired addon version to install.
	// +optional
	DesiredVersion string `json:"desiredVersion,omitempty"`
}

// AddOnHubConfigStatus represents the current status of the configuration for an add-on.
type AddOnHubConfigStatus struct {
	// SupportedVersions lists all the valid addon versions. It's a hint for user to define desired version.
	// +optional
	SupportedVersions []string `json:"supportedVersions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// AddOnHubConfigList is a collection of add-on hub-scoped configuration.
type AddOnHubConfigList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of add-on hub-scoped configuration.
	Items []AddOnHubConfig `json:"items"`
}
