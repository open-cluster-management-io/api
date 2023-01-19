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

	// defaultConfigs represents a list of default add-on configurations.
	// In scenario where all add-ons have the same configuration.
	// User can override the default configuration by defining the configs in the
	// install strategy for specific clusters.
	// +optional
	DefaultConfigs []AddOnConfig `json:"defaultConfigs,omitempty"`

	// InstallStrategy represents the install strategy of the add-on.
	// +optional
	InstallStrategy InstallStrategy `json:"installStrategy,omitempty"`
}

type InstallStrategy struct {
	// Type is the type of the install strategy, it can be:
	// - Manual: no automatic install
	// - Placements: install to clusters selected by placements.
	// +kubebuilder:validation:Enum=Manual;Placements
	// +kubebuilder:default:=Manual
	// +optional
	Type string `json:"type"`

	// Placements is a list of placement references honored when install strategy type is
	// Placements. All clusters selected by these placements will install the addon
	// If one cluster belongs to multiple placements, it will only apply the strategy defined
	// later in the order. That is to say, The latter strategy overrides the previous one.
	// +optional
	Placements []PlacementStrategy `json:"placements,omitempty"`
}

type PlacementStrategy struct {
	// Placement is the reference to a placement
	Placement PlacementRef `json:",inline"`

	// Configs is the configuration of managedClusterAddon during installation.
	// User can override the configuration by updating the managedClusterAddon directly.
	Configs []AddOnConfig `json:"configs,omitempty"`

	// The rollout strategy to apply addon configurations change.
	// The rollout strategy only watches the addon configurations defined in ClusterManagementAddOn.
	// +optional
	RolloutStrategy RolloutStrategy `json:"rolloutStrategy,omitempty"`
}

type PlacementRef struct {
	// Name of the placement
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// Namespace of the placement
	// +kubebuilder:validation:Required
	// +required
	Namespace string `json:"namespace"`
}

// RolloutStrategy represents the rollout strategy of the add-on configuration.
type RolloutStrategy struct {
	// Type is the type of the rollout strategy, it supports UpdateAll, RollingUpdate and RollingUpdateWithCanary:
	// - UpdateAll: when configs change, apply the new configs to all the selected clusters at once.
	//   This is the default strategy.
	// - RollingUpdate: when configs change, apply the new configs to all the selected clusters with
	//   the concurrence rate defined in MaxConcurrentlyUpdating.
	// - RollingUpdateWithCanary: when configs change, wait and check if add-ons on the canary placement
	//   selected clusters have applied the new configs and are healthy, then apply the new configs to
	//   all the selected clusters with the concurrence rate defined in MaxConcurrentlyUpdating.
	//
	//   The field lastKnownGoodConfigSpecHash in the status record the last successfully applied
	//   spec hash of canary placement. If the config spec hash changes after the canary is passed and
	//   before the rollout is done, the current rollout will continue, then roll out to the latest change.
	//
	//   For example, the addon configs have spec hash A. The canary is passed and the lastKnownGoodConfigSpecHash
	//   would be A, and all the selected clusters are rolling out to A.
	//   Then the config spec hash changes to B. At this time, the clusters will continue rolling out to A.
	//   When the rollout is done and canary passed B, the lastKnownGoodConfigSpecHash would be B and
	//   all the clusters will start rolling out to B.
	//
	//   The canary placement does not have to be a subset of the install placement, and it is more like a
	//   reference for finding and checking canary clusters before upgrading all. To trigger the rollout
	//   on the canary clusters, you can define another rollout strategy with the type RollingUpdate, or even
	//   manually upgrade the addons on those clusters.
	//
	// +kubebuilder:validation:Enum=UpdateAll;RollingUpdate;RollingUpdateWithCanary
	// +kubebuilder:default:=UpdateAll
	// +optional
	Type string `json:"type"`

	// Rolling update with placement config params. Present only if the type is RollingUpdate.
	// +optional
	RollingUpdate *RollingUpdate `json:"rollingUpdate,omitempty"`

	// Rolling update with placement config params. Present only if the type is RollingUpdateWithCanary.
	// +optional
	RollingUpdateWithCanary *RollingUpdateWithCanary `json:"rollingUpdateWithCanary,omitempty"`
}

// RollingUpdate represents the behavior to rolling update add-on configurations
// on the selected clusters.
type RollingUpdate struct {
	// The maximum concurrently updating number of addons.
	// Value can be an absolute number (ex: 5) or a percentage of desired addons (ex: 10%).
	// Absolute number is calculated from percentage by rounding up.
	// Defaults to 25%.
	// Example: when this is set to 30%, once the addon configs change, the addon on 30% of the selected clusters
	// will adopt the new configs. When the addons with new configs are healthy, the addon on the remaining clusters
	// will be further updated.
	// +kubebuilder:default:="25%"
	// +optional
	MaxConcurrentlyUpdating *intstr.IntOrString `json:"maxConcurrentlyUpdating,omitempty"`
}

// RollingUpdateWithCanary represents the canary placement and behavior to rolling update add-on configurations
// on the selected clusters.
type RollingUpdateWithCanary struct {
	// Canary placement reference.
	// +kubebuilder:validation:Required
	// +required
	Placement PlacementRef `json:"placement,omitempty"`

	// the behavior to rolling update add-on configurations.
	RollingUpdate `json:",inline"`
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
	// configReferences is a list of current add-on configuration references per placement.
	// +optional
	InstallProgression []InstallProgression `json:"installProgression,omitempty"`
}

type InstallProgression struct {
	// Placement reference.
	// +optional
	Placement PlacementRef `json:"placement,omitempty"`

	// configReferences is a list of current add-on configuration references.
	// +optional
	ConfigReferences []InstallConfigReference `json:"configReferences,omitempty"`

	// conditions describe the state of the managed and monitored components for the operator.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
}

// InstallConfigReference is a reference to the current add-on configuration.
// This resource is used to record the configuration resource for the current add-on.
type InstallConfigReference struct {
	// This field is synced from ClusterManagementAddOn Configurations.
	ConfigGroupResource `json:",inline"`

	// This field is synced from ClusterManagementAddOn Configurations.
	ConfigReferent `json:",inline"`

	// desiredConfigSpecHash record the desired config spec hash.
	DesiredConfigSpecHash string `json:"desiredConfigSpecHash"`

	// lastKnownGoodConfigSpecHash record the last known good config spec hash.
	// For fresh install or rollout with type UpdateAll or RollingUpdate, the
	// lastKnownGoodConfigSpecHash is the same as lastAppliedConfigSpecHash.
	// For rollout with type RollingUpdateWithCanary, the lastKnownGoodConfigSpecHash
	// is the last successfully applied config spec hash of the canary placement.
	LastKnownGoodConfigSpecHash string `json:"lastKnownGoodConfigSpecHash"`

	// lastAppliedConfigSpecHash record the config spec hash when the all the corresponding
	// ManagedClusterAddOn are applied successfully.
	LastAppliedConfigSpecHash string `json:"lastAppliedConfigSpecHash"`
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
