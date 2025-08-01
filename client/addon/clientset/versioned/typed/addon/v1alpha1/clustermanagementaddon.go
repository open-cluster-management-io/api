// Copyright Contributors to the Open Cluster Management project
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
	addonv1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	scheme "open-cluster-management.io/api/client/addon/clientset/versioned/scheme"
)

// ClusterManagementAddOnsGetter has a method to return a ClusterManagementAddOnInterface.
// A group's client should implement this interface.
type ClusterManagementAddOnsGetter interface {
	ClusterManagementAddOns() ClusterManagementAddOnInterface
}

// ClusterManagementAddOnInterface has methods to work with ClusterManagementAddOn resources.
type ClusterManagementAddOnInterface interface {
	Create(ctx context.Context, clusterManagementAddOn *addonv1alpha1.ClusterManagementAddOn, opts v1.CreateOptions) (*addonv1alpha1.ClusterManagementAddOn, error)
	Update(ctx context.Context, clusterManagementAddOn *addonv1alpha1.ClusterManagementAddOn, opts v1.UpdateOptions) (*addonv1alpha1.ClusterManagementAddOn, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, clusterManagementAddOn *addonv1alpha1.ClusterManagementAddOn, opts v1.UpdateOptions) (*addonv1alpha1.ClusterManagementAddOn, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*addonv1alpha1.ClusterManagementAddOn, error)
	List(ctx context.Context, opts v1.ListOptions) (*addonv1alpha1.ClusterManagementAddOnList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *addonv1alpha1.ClusterManagementAddOn, err error)
	ClusterManagementAddOnExpansion
}

// clusterManagementAddOns implements ClusterManagementAddOnInterface
type clusterManagementAddOns struct {
	*gentype.ClientWithList[*addonv1alpha1.ClusterManagementAddOn, *addonv1alpha1.ClusterManagementAddOnList]
}

// newClusterManagementAddOns returns a ClusterManagementAddOns
func newClusterManagementAddOns(c *AddonV1alpha1Client) *clusterManagementAddOns {
	return &clusterManagementAddOns{
		gentype.NewClientWithList[*addonv1alpha1.ClusterManagementAddOn, *addonv1alpha1.ClusterManagementAddOnList](
			"clustermanagementaddons",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *addonv1alpha1.ClusterManagementAddOn { return &addonv1alpha1.ClusterManagementAddOn{} },
			func() *addonv1alpha1.ClusterManagementAddOnList { return &addonv1alpha1.ClusterManagementAddOnList{} },
		),
	}
}
