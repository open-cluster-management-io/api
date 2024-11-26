// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
	v1alpha1 "open-cluster-management.io/api/addon/v1alpha1"
	scheme "open-cluster-management.io/api/client/addon/clientset/versioned/scheme"
)

// AddOnTemplatesGetter has a method to return a AddOnTemplateInterface.
// A group's client should implement this interface.
type AddOnTemplatesGetter interface {
	AddOnTemplates() AddOnTemplateInterface
}

// AddOnTemplateInterface has methods to work with AddOnTemplate resources.
type AddOnTemplateInterface interface {
	Create(ctx context.Context, addOnTemplate *v1alpha1.AddOnTemplate, opts v1.CreateOptions) (*v1alpha1.AddOnTemplate, error)
	Update(ctx context.Context, addOnTemplate *v1alpha1.AddOnTemplate, opts v1.UpdateOptions) (*v1alpha1.AddOnTemplate, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.AddOnTemplate, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AddOnTemplateList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AddOnTemplate, err error)
	AddOnTemplateExpansion
}

// addOnTemplates implements AddOnTemplateInterface
type addOnTemplates struct {
	*gentype.ClientWithList[*v1alpha1.AddOnTemplate, *v1alpha1.AddOnTemplateList]
}

// newAddOnTemplates returns a AddOnTemplates
func newAddOnTemplates(c *AddonV1alpha1Client) *addOnTemplates {
	return &addOnTemplates{
		gentype.NewClientWithList[*v1alpha1.AddOnTemplate, *v1alpha1.AddOnTemplateList](
			"addontemplates",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1alpha1.AddOnTemplate { return &v1alpha1.AddOnTemplate{} },
			func() *v1alpha1.AddOnTemplateList { return &v1alpha1.AddOnTemplateList{} }),
	}
}
