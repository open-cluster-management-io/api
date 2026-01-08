// Copyright Contributors to the Open Cluster Management project
package v1beta1

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"open-cluster-management.io/api/addon/v1alpha1"
)

const (
	// InstallNamespaceAnnotation is the annotation key for storing installNamespace
	// This is used because installNamespace field was removed in v1beta1
	InstallNamespaceAnnotation = "addon.open-cluster-management.io/v1alpha1-install-namespace"
)

// ConvertTo converts the receiver (v1beta1) into v1alpha1.
func (src *ManagedClusterAddOn) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha1.ManagedClusterAddOn)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.ManagedClusterAddOn, got %T", dstRaw)
	}

	// Copy the ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	if err := Convert_v1beta1_ManagedClusterAddOnSpec_To_v1alpha1_ManagedClusterAddOnSpec(&src.Spec, &dst.Spec, nil); err != nil {
		return fmt.Errorf("failed to convert spec: %w", err)
	}

	// Restore installNamespace from annotation
	// This field was removed in v1beta1, so we store it in annotation
	if installNs, ok := src.Annotations[InstallNamespaceAnnotation]; ok {
		dst.Spec.InstallNamespace = installNs
	}

	// Convert Status
	if err := Convert_v1beta1_ManagedClusterAddOnStatus_To_v1alpha1_ManagedClusterAddOnStatus(&src.Status, &dst.Status, nil); err != nil {
		return fmt.Errorf("failed to convert status: %w", err)
	}

	// Manually populate deprecated ConfigReferent field in ConfigReferences
	// The native conversion doesn't handle this deprecated field from v1alpha1.ConfigReference
	// We need to copy from DesiredConfig.ConfigReferent if DesiredConfig exists
	for i := range dst.Status.ConfigReferences {
		if dst.Status.ConfigReferences[i].DesiredConfig != nil {
			dst.Status.ConfigReferences[i].ConfigReferent = dst.Status.ConfigReferences[i].DesiredConfig.ConfigReferent
		}
	}

	return nil
}

// ConvertFrom converts from v1alpha1 into the receiver (v1beta1).
func (dst *ManagedClusterAddOn) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha1.ManagedClusterAddOn)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.ManagedClusterAddOn, got %T", srcRaw)
	}

	// Copy the ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	if err := Convert_v1alpha1_ManagedClusterAddOnSpec_To_v1beta1_ManagedClusterAddOnSpec(&src.Spec, &dst.Spec, nil); err != nil {
		return fmt.Errorf("failed to convert spec: %w", err)
	}

	// Save installNamespace to annotation (removed in v1beta1)
	// This field exists in v1alpha1 but not in v1beta1, so we preserve it in annotations
	if src.Spec.InstallNamespace != "" {
		if dst.Annotations == nil {
			dst.Annotations = make(map[string]string)
		}
		dst.Annotations[InstallNamespaceAnnotation] = src.Spec.InstallNamespace
	}

	// Convert Status
	if err := Convert_v1alpha1_ManagedClusterAddOnStatus_To_v1beta1_ManagedClusterAddOnStatus(&src.Status, &dst.Status, nil); err != nil {
		return fmt.Errorf("failed to convert status: %w", err)
	}

	return nil
}
