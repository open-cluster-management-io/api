package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope="Cluster"
// +kubebuilder:printcolumn:name="DISPLAY NAME",type=string,JSONPath=`.spec.addOnMeta.displayName`
// +kubebuilder:printcolumn:name="CRD NAME",type=string,JSONPath=`.spec.addOnConfiguration.crdName`

// ClusterManagementAddOn represents the registration of an add-on to the cluster manager.
// This resource allows the user to discover which add-on is available for the cluster manager and
// also provides metadata information about the add-on.
// This resource also provides a linkage to ManagedClusterAddOn, the name of the ClusterManagementAddOn
// resource will be used for the namespace-scoped ManagedClusterAddOn resource.
// ClusterManagementAddOn is a cluster-scoped resource.
type ClusterManagementAddOn struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec represents a desired configuration for the agent on the cluster management add-on.
	// +required
	Spec ClusterManagementAddOnSpec `json:"spec"`

	// status represents the current status of cluster management add-on.
	// +optional
	Status ClusterManagementAddOnStatus `json:"status,omitempty"`
}

// ClusterManagementAddOnSpec provides information for the add-on.
type ClusterManagementAddOnSpec struct {
	// addOnMeta is a reference to the metadata information for the add-on.
	// +optional
	AddOnMeta AddOnMeta `json:"addOnMeta,omitempty"`

	// Deprecated: Use supportedConfigs filed instead
	// addOnConfiguration is a reference to configuration information for the add-on.
	// In scenario where a multiple add-ons share the same add-on CRD, multiple ClusterManagementAddOn
	// resources need to be created and reference the same AddOnConfiguration.
	// +optional
	AddOnConfiguration ConfigCoordinates `json:"addOnConfiguration,omitempty"`

	// Deprecated: move supportedConfigs to ManagedClusterAddOnStatus.
	// supportedConfigs is a list of configuration types supported by add-on.
	// An empty list means the add-on does not require configurations.
	// The default is an empty list
	// +optional
	// +listType=map
	// +listMapKey=group
	// +listMapKey=resource
	SupportedConfigs []ConfigMeta `json:"supportedConfigs,omitempty"`

	// configuration lists the add-on configurations and the rollout strategy when the configurations change.
	// +optional
	Configuration Configuration `json:"configuration,omitempty"`
}

// Configuration represents a list of the add-on configurations and their rollout strategy.
type Configuration struct {
	// configs is a list of add-on configurations.
	// The add-on configurations for each cluster can be overridden by the configs of the ManagedClusterAddon spec.
	// +optional
	Configs []AddOnConfig `json:"configs,omitempty"`

	// The rollout strategy to apply new configs.
	// The rollout strategy only watches the listed configs change.
	// If the rollout strategy is not defined, the default strategy UpdateAll is used.
	// If there are configs change during the rollout process, the rollout will start over. For example, configs list
	// configA and configB. The change in configA triggers the rollout. If the configB is also
	// changed before the rollout is complete, the current rollout stops and the rollout starts over.
	// +optional
	RolloutStrategy RolloutStrategy `json:"rolloutStrategy,omitempty"`
}

// RolloutStrategy represents the rollout strategy of the add-on configuration.
type RolloutStrategy struct {
	// Type is the type of the rollout strategy, it supports UpdateAll and RollingUpdateWithPlacement:
	// - UpdateAll: when configs change, apply the new configs to all the clusters.
	// - RollingUpdateWithPlacement: when configs change, rolling update new configs on the clusters
	//   selected by placements.
	//   If any of the configs are overridden by ManagedClusterAddOn on the specific cluster, the new configs
	//   won't take effect on that cluster.
	//   This rollout strategy is only responsible for applying new configs. When the strategy is modified or
	//   removed, the applied configs won't be deleted from the cluster.
	// +kubebuilder:validation:Enum=RollingUpdateWithPlacement;UpdateAll
	// +kubebuilder:default:=UpdateAll
	// +optional
	Type string `json:"type"`

	// Rolling update with placement config params. Present only if the type is RollingUpdateWithPlacement.
	// +optional
	RollingUpdateWithPlacement *RollingUpdateWithPlacement `json:"rollingUpdateWithPlacement,omitempty"`
}

// RollingUpdateWithPlacement represents the placement and behavior to rolling update add-on configurations
// on the selected clusters.
type RollingUpdateWithPlacement struct {
	// name of the placement
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// namespace of the placement.
	// +kubebuilder:validation:Required
	// +required
	Namespace string `json:"namespace"`

	// The maximum concurrently updating number of addons.
	// Value can be an absolute number (ex: 5) or a percentage of desired addons (ex: 10%).
	// Absolute number is calculated from percentage by rounding up.
	// Defaults to 25%.
	// Example: when this is set to 30%, once the addon configs change, the addon on 30% of the selected clusters
	// will adopt the new configs. When the new configs are ready, the addon on the remaining clusters
	// will be further updated.
	// +optional
	MaxConcurrentlyUpdating *intstr.IntOrString `json:"maxConcurrentlyUpdating,omitempty"`
}

// AddOnMeta represents a collection of metadata information for the add-on.
type AddOnMeta struct {
	// displayName represents the name of add-on that will be displayed.
	// +optional
	DisplayName string `json:"displayName,omitempty"`

	// description represents the detailed description of the add-on.
	// +optional
	Description string `json:"description,omitempty"`
}

// ConfigMeta represents a collection of metadata information for add-on configuration.
type ConfigMeta struct {
	// group and resouce of the add-on configuration.
	ConfigGroupResource `json:",inline"`

	// defaultConfig represents the namespace and name of the default add-on configuration.
	// In scenario where all add-ons have a same configuration.
	// +optional
	DefaultConfig *ConfigReferent `json:"defaultConfig,omitempty"`
}

// ConfigCoordinates represents the information for locating the CRD and CR that configures the add-on.
type ConfigCoordinates struct {
	// crdName is the name of the CRD used to configure instances of the managed add-on.
	// This field should be configured if the add-on have a CRD that controls the configuration of the add-on.
	// +optional
	CRDName string `json:"crdName,omitempty"`

	// crName is the name of the CR used to configure instances of the managed add-on.
	// This field should be configured if add-on CR have a consistent name across the all of the ManagedCluster instaces.
	// +optional
	CRName string `json:"crName,omitempty"`

	// lastObservedGeneration is the observed generation of the custom resource for the configuration of the addon.
	// +optional
	LastObservedGeneration int64 `json:"lastObservedGeneration,omitempty"`
}

// ConfigGroupResource represents the GroupResource of the add-on configuration
type ConfigGroupResource struct {
	// group of the add-on configuration.
	// +optional
	// +kubebuilder:default=""
	Group string `json:"group"`

	// resource of the add-on configuration.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Resource string `json:"resource"`
}

// ConfigReferent represents the namespace and name for an add-on configuration.
type ConfigReferent struct {
	// namespace of the add-on configuration.
	// If this field is not set, the configuration is in the cluster scope.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// name of the add-on configuration.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	Name string `json:"name"`
}

// ClusterManagementAddOnStatus represents the current status of cluster management add-on.
type ClusterManagementAddOnStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ClusterManagementAddOnList is a collection of cluster management add-ons.
type ClusterManagementAddOnList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of cluster management add-ons.
	Items []ClusterManagementAddOn `json:"items"`
}
