package utils

import (
	"fmt"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/restmapper"

	workv1 "open-cluster-management.io/api/work/v1"
)

func TestBuildResourceMeta(t *testing.T) {
	restMapper := NewFakeRestMapper()

	cases := []struct {
		name         string
		index        int
		obj          runtime.Object
		expectedErr  error
		expectedGVR  schema.GroupVersionResource
		expectedMeta workv1.ManifestResourceMeta
	}{
		{
			name:        "object nil",
			index:       0,
			obj:         nil,
			expectedErr: nil,
			expectedGVR: schema.GroupVersionResource{},
			expectedMeta: workv1.ManifestResourceMeta{
				Ordinal: int32(0),
			},
		},
		{
			name:        "secret success",
			index:       1,
			obj:         NewUnstructured("v1", "Secret", "ns1", "test"),
			expectedErr: nil,
			expectedGVR: schema.GroupVersionResource{
				Group:    "",
				Version:  "v1",
				Resource: "secrets",
			},
			expectedMeta: workv1.ManifestResourceMeta{
				Ordinal:   int32(1),
				Group:     "",
				Version:   "v1",
				Kind:      "Secret",
				Resource:  "secrets",
				Namespace: "ns1",
				Name:      "test",
			},
		},
		{
			name:        "unknown object type",
			index:       1,
			obj:         NewUnstructured("test/v1", "NewObject", "ns1", "test"),
			expectedErr: fmt.Errorf("the server doesn't have a resource type %q", "NewObject"),
			expectedGVR: schema.GroupVersionResource{},
			expectedMeta: workv1.ManifestResourceMeta{
				Ordinal:   int32(1),
				Group:     "test",
				Version:   "v1",
				Kind:      "NewObject",
				Namespace: "ns1",
				Name:      "test",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			meta, gvr, err := BuildResourceMeta(c.index, c.obj, restMapper)
			if c.expectedErr == nil {
				if err != nil {
					t.Errorf("Case name: %s, expect error nil, but got %v", c.name, err)
				}
			} else if err == nil {
				t.Errorf("Case name: %s, expect error %s, but got nil", c.name, c.expectedErr)
			} else if c.expectedErr.Error() != err.Error() {
				t.Errorf("Case name: %s, expect error %s, but got %s", c.name, c.expectedErr, err)
			}

			if !reflect.DeepEqual(c.expectedGVR, gvr) {
				t.Errorf("Case name: %s, expect gvr %v, but got %v", c.name, c.expectedGVR, gvr)
			}

			if !reflect.DeepEqual(c.expectedMeta, meta) {
				t.Errorf("Case name: %s, expect meta %v, but got %v", c.name, c.expectedMeta, meta)
			}
		})
	}
}

func NewFakeRestMapper() meta.RESTMapper {
	resources := []*restmapper.APIGroupResources{
		{
			Group: metav1.APIGroup{
				Name: "",
				Versions: []metav1.GroupVersionForDiscovery{
					{Version: "v1"},
				},
				PreferredVersion: metav1.GroupVersionForDiscovery{Version: "v1"},
			},
			VersionedResources: map[string][]metav1.APIResource{
				"v1": {
					{Name: "secrets", Namespaced: true, Kind: "Secret"},
					{Name: "pods", Namespaced: true, Kind: "Pod"},
					{Name: "newobjects", Namespaced: true, Kind: "NewObject"},
				},
			},
		},
		{
			Group: metav1.APIGroup{
				Name: "apps",
				Versions: []metav1.GroupVersionForDiscovery{
					{Version: "v1", GroupVersion: "apps/v1"},
				},
				PreferredVersion: metav1.GroupVersionForDiscovery{Version: "v1", GroupVersion: "apps/v1"},
			},
			VersionedResources: map[string][]metav1.APIResource{
				"v1": {
					{Name: "deployments", Group: "apps", Namespaced: true, Kind: "Deployment"},
				},
			},
		},
	}
	return restmapper.NewDiscoveryRESTMapper(resources)
}

func NewUnstructured(apiVersion, kind, namespace, name string, owners ...metav1.OwnerReference) *unstructured.Unstructured {
	u := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": apiVersion,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"namespace": namespace,
				"name":      name,
			},
		},
	}

	u.SetOwnerReferences(owners)

	return u
}
