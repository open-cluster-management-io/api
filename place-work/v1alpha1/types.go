/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	work "open-cluster-management.io/api/work/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PlaceManifestWorkSpec defines the desired state of PlaceManifestWork
type PlaceManifestWorkSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ManifestWorkSpec is the ManifestWorkSpec that will be used to generate a per-cluster ManifestWork
	ManifestWorkSpec work.ManifestWorkSpec `json:"manifestWorkSpec"`

	// PacementRef is the name of the Placement resource, from which a PlacementDecision will be found and used
	// to distribute the ManifestWork
	PlacementRef *corev1.LocalObjectReference `json:"placementRef,omitempty"`
}

// PlaceManifestWorkStatus defines the observed state of PlaceManifestWork
type PlaceManifestWorkStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions contains the different condition statuses for distrbution of ManifestWork resources
	// Valid condition types are:
	// 1. AppliedManifestWorks represents ManifestWorks have been distributed as per placement All, Partial, None, Problem
	// 2. PlacementRefValid
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ManifestWorkDelivery tracks each ManifestWork that was created,
	// including that it was applied and its overall status
	PlacedManifestWork []PlacedManifestWork `json:"placedManifestWork,omitempty"`

	// Summary that reflects all relevant ManifestWorks
	PlacedManifestWorkSummary []PlacedManifestWorkSummary `json:"summary"`
}

type PlacedManifestWorkSummary struct {
	Total       int `json:"total"`
	Progressing int `json:"progressing"`
	Available   int `json:"available"`
	Degraded    int `json:"degraded"`
	Applied     int `json:"Applied"`
}

type ManifestWorkStatus string

type PlacedManifestWork struct {

	// Work is an objectReference to the actual ManifestWork resource
	Work corev1.ObjectReference `json:"workRef"`

	// +kubebuilder:default=false
	Applied bool `json:"applied"`

	// ManifestWorkStatus ManifestWorkStatus.Condition can be
	// +kubebuilder:validation:Enum=Progressing;Available;Degraded;Applied
	// +optional
	ManifestWorkStatus ManifestWorkStatus `json:"status"`

	// Timestamp is the time when the ManifestWork was last modified
	Timestamp string `json:"timestamp"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=placementmanifestworks,shortName=pmw;pmws,scope=Namespaced
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Placement",type="string",JSONPath=".status.conditions[?(@.type==\"PlacementVerified\")].reason",description="Reason"
// +kubebuilder:printcolumn:name="Found",type="string",JSONPath=".status.conditions[?(@.type==\"PlacementVerified\")].status",description="Configured"
// +kubebuilder:printcolumn:name="ManifestWorks",type="string",JSONPath=".status.conditions[?(@.type==\"ManifestworkApplied\")].reason",description="Reason"
// +kubebuilder:printcolumn:name="Applied",type="string",JSONPath=".status.conditions[?(@.type==\"ManifestworkApplied\")].status",description="Applied"

// PlaceManifestWork is the Schema for the placementmanifestworks API
type PlaceManifestWork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec reperesents the desired ManifestWork payload and Placement reference to be reconciled
	Spec PlaceManifestWorkSpec `json:"spec,omitempty"`

	// Status represent the current status of Placing ManifestWork resources
	Status PlaceManifestWorkStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PlaceManifestWorkList contains a list of PlaceManifestWork
type PlaceManifestWorkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlaceManifestWork `json:"items"`
}

type ConditionType string

const (
	PlacementDecisionVerifiedAsExpected = "PlacementDecisionVerified"
	PlacementDecisionNotFound           = "PlacementDecisionNotFound"
	PlacementDecisionEmpty              = "PlacementDecisionEmpty"

	AsExpected = "AsExpected"
	Processing = "Processing"
	Partial    = "Partial"

	// PlacementDecisionVerified indicates if Placement is valid
	PlacementDecisionVerified ConditionType = "PlacementVerified"

	// ManifestWorkApplied confirms that a ManifestWork has been created in each cluster defined by PlacementDecision
	ManifestworkApplied ConditionType = "ManifestworkApplied"
)
