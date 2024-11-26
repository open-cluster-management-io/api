// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
)

// FakeClusterManagementAddOns implements ClusterManagementAddOnInterface
type FakeClusterManagementAddOns struct {
	Fake *FakeAddonV1alpha1
}

var clustermanagementaddonsResource = v1alpha1.SchemeGroupVersion.WithResource("clustermanagementaddons")

var clustermanagementaddonsKind = v1alpha1.SchemeGroupVersion.WithKind("ClusterManagementAddOn")

// Get takes name of the clusterManagementAddOn, and returns the corresponding clusterManagementAddOn object, and an error if there is any.
func (c *FakeClusterManagementAddOns) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ClusterManagementAddOn, err error) {
	emptyResult := &v1alpha1.ClusterManagementAddOn{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(clustermanagementaddonsResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.ClusterManagementAddOn), err
}

// List takes label and field selectors, and returns the list of ClusterManagementAddOns that match those selectors.
func (c *FakeClusterManagementAddOns) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ClusterManagementAddOnList, err error) {
	emptyResult := &v1alpha1.ClusterManagementAddOnList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(clustermanagementaddonsResource, clustermanagementaddonsKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterManagementAddOnList{ListMeta: obj.(*v1alpha1.ClusterManagementAddOnList).ListMeta}
	for _, item := range obj.(*v1alpha1.ClusterManagementAddOnList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterManagementAddOns.
func (c *FakeClusterManagementAddOns) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(clustermanagementaddonsResource, opts))
}

// Create takes the representation of a clusterManagementAddOn and creates it.  Returns the server's representation of the clusterManagementAddOn, and an error, if there is any.
func (c *FakeClusterManagementAddOns) Create(ctx context.Context, clusterManagementAddOn *v1alpha1.ClusterManagementAddOn, opts v1.CreateOptions) (result *v1alpha1.ClusterManagementAddOn, err error) {
	emptyResult := &v1alpha1.ClusterManagementAddOn{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(clustermanagementaddonsResource, clusterManagementAddOn, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.ClusterManagementAddOn), err
}

// Update takes the representation of a clusterManagementAddOn and updates it. Returns the server's representation of the clusterManagementAddOn, and an error, if there is any.
func (c *FakeClusterManagementAddOns) Update(ctx context.Context, clusterManagementAddOn *v1alpha1.ClusterManagementAddOn, opts v1.UpdateOptions) (result *v1alpha1.ClusterManagementAddOn, err error) {
	emptyResult := &v1alpha1.ClusterManagementAddOn{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(clustermanagementaddonsResource, clusterManagementAddOn, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.ClusterManagementAddOn), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterManagementAddOns) UpdateStatus(ctx context.Context, clusterManagementAddOn *v1alpha1.ClusterManagementAddOn, opts v1.UpdateOptions) (result *v1alpha1.ClusterManagementAddOn, err error) {
	emptyResult := &v1alpha1.ClusterManagementAddOn{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(clustermanagementaddonsResource, "status", clusterManagementAddOn, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.ClusterManagementAddOn), err
}

// Delete takes name of the clusterManagementAddOn and deletes it. Returns an error if one occurs.
func (c *FakeClusterManagementAddOns) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(clustermanagementaddonsResource, name, opts), &v1alpha1.ClusterManagementAddOn{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterManagementAddOns) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(clustermanagementaddonsResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterManagementAddOnList{})
	return err
}

// Patch applies the patch and returns the patched clusterManagementAddOn.
func (c *FakeClusterManagementAddOns) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterManagementAddOn, err error) {
	emptyResult := &v1alpha1.ClusterManagementAddOn{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(clustermanagementaddonsResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.ClusterManagementAddOn), err
}
