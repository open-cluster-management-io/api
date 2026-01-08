// Copyright Contributors to the Open Cluster Management project
package v1beta1

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"open-cluster-management.io/api/addon/v1alpha1"
)

// ConvertTo converts the receiver (v1beta1) into v1alpha1.
func (src *ClusterManagementAddOn) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha1.ClusterManagementAddOn)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.ClusterManagementAddOn, got %T", dstRaw)
	}

	// Copy the ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	if err := Convert_v1beta1_ClusterManagementAddOnSpec_To_v1alpha1_ClusterManagementAddOnSpec(&src.Spec, &dst.Spec, nil); err != nil {
		return fmt.Errorf("failed to convert spec: %w", err)
	}

	// Convert Status
	if err := Convert_v1beta1_ClusterManagementAddOnStatus_To_v1alpha1_ClusterManagementAddOnStatus(&src.Status, &dst.Status, nil); err != nil {
		return fmt.Errorf("failed to convert status: %w", err)
	}

	return nil
}

// ConvertFrom converts from v1alpha1 into the receiver (v1beta1).
func (dst *ClusterManagementAddOn) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha1.ClusterManagementAddOn)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.ClusterManagementAddOn, got %T", srcRaw)
	}

	// Copy the ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec
	if err := Convert_v1alpha1_ClusterManagementAddOnSpec_To_v1beta1_ClusterManagementAddOnSpec(&src.Spec, &dst.Spec, nil); err != nil {
		return fmt.Errorf("failed to convert spec: %w", err)
	}

	// Convert Status
	if err := Convert_v1alpha1_ClusterManagementAddOnStatus_To_v1beta1_ClusterManagementAddOnStatus(&src.Status, &dst.Status, nil); err != nil {
		return fmt.Errorf("failed to convert status: %w", err)
	}

	return nil
}
