package source

import (
	"fmt"

	"github.com/google/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	kubetypes "k8s.io/apimachinery/pkg/types"

	"open-cluster-management.io/api/cloudevents/generic"
)

type ResourceStatus struct {
	Conditions []metav1.Condition
}

type Resource struct {
	ResourceID        string
	ResourceVersion   int64
	Namespace         string
	DeletionTimestamp *metav1.Time
	Spec              unstructured.Unstructured
	Status            ResourceStatus
}

var _ generic.ResourceObject = &Resource{}

func NewResource(namespace, name string) *Resource {
	return &Resource{
		ResourceID:      ResourceID(namespace, name),
		ResourceVersion: 1,
		Namespace:       namespace,
		Spec: unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"namespace": namespace,
					"name":      name,
				},
			},
		},
	}
}

func (r *Resource) GetUID() kubetypes.UID {
	return kubetypes.UID(r.ResourceID)
}

func (r *Resource) GetResourceVersion() string {
	return fmt.Sprintf("%d", r.ResourceVersion)
}

func (r *Resource) GetDeletionTimestamp() *metav1.Time {
	return r.DeletionTimestamp
}

func ResourceID(namespace, name string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("resource-%s-%s", namespace, name))).String()
}
