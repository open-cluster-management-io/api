// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
)

// FakeAddOnPlacementScores implements AddOnPlacementScoreInterface
type FakeAddOnPlacementScores struct {
	Fake *FakeClusterV1alpha1
	ns   string
}

var addonplacementscoresResource = v1alpha1.SchemeGroupVersion.WithResource("addonplacementscores")

var addonplacementscoresKind = v1alpha1.SchemeGroupVersion.WithKind("AddOnPlacementScore")

// Get takes name of the addOnPlacementScore, and returns the corresponding addOnPlacementScore object, and an error if there is any.
func (c *FakeAddOnPlacementScores) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AddOnPlacementScore, err error) {
	emptyResult := &v1alpha1.AddOnPlacementScore{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(addonplacementscoresResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.AddOnPlacementScore), err
}

// List takes label and field selectors, and returns the list of AddOnPlacementScores that match those selectors.
func (c *FakeAddOnPlacementScores) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AddOnPlacementScoreList, err error) {
	emptyResult := &v1alpha1.AddOnPlacementScoreList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(addonplacementscoresResource, addonplacementscoresKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AddOnPlacementScoreList{ListMeta: obj.(*v1alpha1.AddOnPlacementScoreList).ListMeta}
	for _, item := range obj.(*v1alpha1.AddOnPlacementScoreList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested addOnPlacementScores.
func (c *FakeAddOnPlacementScores) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(addonplacementscoresResource, c.ns, opts))

}

// Create takes the representation of a addOnPlacementScore and creates it.  Returns the server's representation of the addOnPlacementScore, and an error, if there is any.
func (c *FakeAddOnPlacementScores) Create(ctx context.Context, addOnPlacementScore *v1alpha1.AddOnPlacementScore, opts v1.CreateOptions) (result *v1alpha1.AddOnPlacementScore, err error) {
	emptyResult := &v1alpha1.AddOnPlacementScore{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(addonplacementscoresResource, c.ns, addOnPlacementScore, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.AddOnPlacementScore), err
}

// Update takes the representation of a addOnPlacementScore and updates it. Returns the server's representation of the addOnPlacementScore, and an error, if there is any.
func (c *FakeAddOnPlacementScores) Update(ctx context.Context, addOnPlacementScore *v1alpha1.AddOnPlacementScore, opts v1.UpdateOptions) (result *v1alpha1.AddOnPlacementScore, err error) {
	emptyResult := &v1alpha1.AddOnPlacementScore{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(addonplacementscoresResource, c.ns, addOnPlacementScore, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.AddOnPlacementScore), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAddOnPlacementScores) UpdateStatus(ctx context.Context, addOnPlacementScore *v1alpha1.AddOnPlacementScore, opts v1.UpdateOptions) (result *v1alpha1.AddOnPlacementScore, err error) {
	emptyResult := &v1alpha1.AddOnPlacementScore{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(addonplacementscoresResource, "status", c.ns, addOnPlacementScore, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.AddOnPlacementScore), err
}

// Delete takes name of the addOnPlacementScore and deletes it. Returns an error if one occurs.
func (c *FakeAddOnPlacementScores) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(addonplacementscoresResource, c.ns, name, opts), &v1alpha1.AddOnPlacementScore{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAddOnPlacementScores) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(addonplacementscoresResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.AddOnPlacementScoreList{})
	return err
}

// Patch applies the patch and returns the patched addOnPlacementScore.
func (c *FakeAddOnPlacementScores) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AddOnPlacementScore, err error) {
	emptyResult := &v1alpha1.AddOnPlacementScore{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(addonplacementscoresResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.AddOnPlacementScore), err
}
