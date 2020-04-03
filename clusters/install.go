package clusters

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	clustersv1 "github.com/open-cluster-management/api/clusters/v1"
)

const (
	GroupName = "clusters.open-cluster-management.io"
)

var (
	schemeBuilder = runtime.NewSchemeBuilder(clustersv1.Install)
	// Install is a function which adds every version of this group to a scheme
	Install = schemeBuilder.AddToScheme
)

func Resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: GroupName, Resource: resource}
}

func Kind(kind string) schema.GroupKind {
	return schema.GroupKind{Group: GroupName, Kind: kind}
}
