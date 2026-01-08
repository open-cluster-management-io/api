// Copyright Contributors to the Open Cluster Management project
package v1beta1

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"open-cluster-management.io/api/addon/v1alpha1"
)

// ConvertTo converts the receiver (v1beta1) into v1alpha1.
func (src *AddOnDeploymentConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha1.AddOnDeploymentConfig)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.AddOnDeploymentConfig, got %T", dstRaw)
	}

	// Copy the ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec - since the types are identical between versions, we can use the generated conversion
	if err := Convert_v1beta1_AddOnDeploymentConfigSpec_To_v1alpha1_AddOnDeploymentConfigSpec(&src.Spec, &dst.Spec, nil); err != nil {
		return fmt.Errorf("failed to convert spec: %w", err)
	}

	return nil
}

// ConvertFrom converts from v1alpha1 into the receiver (v1beta1).
func (dst *AddOnDeploymentConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha1.AddOnDeploymentConfig)
	if !ok {
		return fmt.Errorf("expected *v1alpha1.AddOnDeploymentConfig, got %T", srcRaw)
	}

	// Copy the ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Convert Spec - since the types are identical between versions, we can use the generated conversion
	if err := Convert_v1alpha1_AddOnDeploymentConfigSpec_To_v1beta1_AddOnDeploymentConfigSpec(&src.Spec, &dst.Spec, nil); err != nil {
		return fmt.Errorf("failed to convert spec: %w", err)
	}

	return nil
}
