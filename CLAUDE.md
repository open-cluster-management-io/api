# Open Cluster Management API

The canonical source for all Open Cluster Management (OCM) API definitions, client libraries, and CRDs for multicluster and multicloud Kubernetes scenarios.

## Project Overview

**Type**: Go API Library / Kubernetes CRDs  
**Language**: Go  
**Framework**: Kubernetes API Machinery / controller-runtime  
**Primary Purpose**: API definitions and client libraries for OCM ecosystem  
**Target Users**: Platform engineers, multicluster operators, OCM addon developers

## Technical Stack

- **Language**: Go 1.23+ with modern module structure
- **Kubernetes Integration**: API Machinery for CRD definitions
- **Code Generation**: Kubernetes code generators for clients and deepcopy
- **API Versioning**: Multiple API versions (v1, v1beta1, v1alpha1)
- **Client Libraries**: Kubernetes-style typed clientsets
- **Validation**: OpenAPI schema generation for CRD validation

## API Architecture

### Four Core API Groups

#### 1. Cluster Management (`cluster.open-cluster-management.io`)
- **v1**: Core cluster registration and grouping
  - `ManagedCluster`: Represents a cluster joined to the hub
  - `ManagedClusterSet`: Logical grouping of clusters
  - `ManagedClusterSetBinding`: Namespace-level cluster set access
- **v1beta1**: Advanced cluster selection and placement
  - `Placement`: Intelligent cluster selection with constraints
  - `PlacementDecision`: Results of placement decisions
- **v1alpha1**: Extended cluster capabilities
  - `ClusterClaim`: Cluster-specific capability declarations
  - `AddonPlacementScore`: Scoring for placement decisions

#### 2. Work Distribution (`work.open-cluster-management.io`)
- **v1**: Resource deployment across clusters
  - `ManifestWork`: Kubernetes resources to deploy on managed clusters
  - `AppliedManifestWork`: Status of deployed resources
- **v1alpha1**: Bulk deployment capabilities
  - `ManifestWorkReplicaSet`: Deploy to multiple clusters in a set

#### 3. Addon Management (`addon.open-cluster-management.io`)
- **v1alpha1**: Addon lifecycle and configuration
  - `ClusterManagementAddOn`: Hub-side addon definition
  - `ManagedClusterAddOn`: Cluster-specific addon instance
  - `AddonDeploymentConfig`: Addon configuration templates

#### 4. Operator APIs (`operator.open-cluster-management.io`)
- **v1**: OCM component installation
  - `ClusterManager`: Hub cluster OCM installation
  - `Klusterlet`: Managed cluster agent installation

## Development Workflow

### Code Generation
```bash
# Update all generated code (clients, deepcopy, CRDs)
make update

# Verify code generation is up to date
make verify

# Run all tests
make test
```

### Generated Components
- **DeepCopy Methods**: For all API struct types
- **Typed Clients**: Kubernetes-style clientsets for each API group
- **Informers/Listers**: Efficient caching and watching mechanisms
- **OpenAPI Schemas**: CRD validation and documentation

## AI Development Guidelines

### API Design Patterns

#### Resource Structure
```go
// Standard Kubernetes API resource pattern
type ManagedCluster struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
    
    Spec   ManagedClusterSpec   `json:"spec,omitempty"`
    Status ManagedClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster
```

#### Client Usage Patterns
```go
import (
    clusterv1 "open-cluster-management.io/api/cluster/v1"
    clusterclient "open-cluster-management.io/api/client/cluster/clientset/versioned"
)

// Create typed client
config, err := ctrl.GetConfig()
client := clusterclient.NewForConfigOrDie(config)

// List managed clusters
clusters, err := client.ClusterV1().ManagedClusters().List(ctx, metav1.ListOptions{})

// Watch for changes
watchInterface, err := client.ClusterV1().ManagedClusters().Watch(ctx, metav1.ListOptions{})
```

#### Status Management
```go
// Standard status update pattern
managedCluster.Status.Conditions = []metav1.Condition{
    {
        Type:    "ManagedClusterJoined",
        Status:  metav1.ConditionTrue,
        Reason:  "ManagedClusterJoined",
        Message: "Cluster successfully joined the hub",
    },
}
```

### Common Development Tasks

#### Adding New API Types
1. **Define Structure**: Create Go structs with proper kubebuilder tags
2. **Register Types**: Add to scheme registration in `register.go`
3. **Generate Code**: Run `make update` to generate clients and deepcopy
4. **Add Tests**: Create comprehensive unit tests for new types
5. **Update Documentation**: Add examples and usage patterns

#### Working with Placements
```go
// Create placement with cluster selection
placement := &clusterv1beta1.Placement{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-placement",
        Namespace: "default",
    },
    Spec: clusterv1beta1.PlacementSpec{
        ClusterSets: []string{"clusterset1"},
        Predicates: []clusterv1beta1.ClusterPredicate{
            {
                RequiredClusterSelector: clusterv1beta1.ClusterSelector{
                    LabelSelector: metav1.LabelSelector{
                        MatchLabels: map[string]string{
                            "environment": "production",
                        },
                    },
                },
            },
        },
    },
}
```

#### ManifestWork Deployment
```go
// Deploy resources to managed clusters
manifestWork := &workv1.ManifestWork{
    ObjectMeta: metav1.ObjectMeta{
        Name:      "my-workload",
        Namespace: "cluster1", // managed cluster namespace
    },
    Spec: workv1.ManifestWorkSpec{
        Workload: workv1.ManifestsTemplate{
            Manifests: []workv1.Manifest{
                {
                    Object: unstructured.Unstructured{
                        Object: map[string]interface{}{
                            "apiVersion": "apps/v1",
                            "kind":       "Deployment",
                            "metadata": map[string]interface{}{
                                "name":      "my-app",
                                "namespace": "default",
                            },
                            // ... deployment spec
                        },
                    },
                },
            },
        },
    },
}
```

### Integration Patterns

#### Controller Development
```go
// Standard controller setup for OCM APIs
import (
    clusterv1 "open-cluster-management.io/api/cluster/v1"
    "sigs.k8s.io/controller-runtime/pkg/controller"
)

func (r *MyReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&clusterv1.ManagedCluster{}).
        Watches(&workv1.ManifestWork{}, &handler.EnqueueRequestForObject{}).
        Complete(r)
}
```

#### Informer Usage
```go
// Use generated informers for efficient watching
import (
    clusterinformers "open-cluster-management.io/api/client/cluster/informers/externalversions"
)

informerFactory := clusterinformers.NewSharedInformerFactory(client, time.Minute*30)
clusterInformer := informerFactory.Cluster().V1().ManagedClusters()

clusterInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc:    handleClusterAdd,
    UpdateFunc: handleClusterUpdate,
    DeleteFunc: handleClusterDelete,
})
```

### Best Practices

#### API Versioning
- **v1alpha1**: Experimental APIs, may have breaking changes
- **v1beta1**: Stable APIs, backward compatible changes only
- **v1**: Stable APIs, no breaking changes allowed
- Use feature gates for experimental functionality

#### Resource Naming
- Follow Kubernetes naming conventions
- Use descriptive, action-oriented names
- Maintain consistency across API groups
- Consider plural/singular forms carefully

#### Status Conventions
```go
// Use standard Kubernetes condition types
const (
    ConditionReady     = "Ready"
    ConditionAvailable = "Available"
    ConditionDegraded  = "Degraded"
)

// Include helpful error messages
condition := metav1.Condition{
    Type:               ConditionReady,
    Status:             metav1.ConditionFalse,
    Reason:             "InvalidConfiguration",
    Message:            "Cluster configuration validation failed: missing required labels",
    LastTransitionTime: metav1.Now(),
}
```

#### Validation and Defaults
```go
// Use kubebuilder tags for validation
type ManagedClusterSpec struct {
    // +kubebuilder:validation:Required
    // +kubebuilder:default=true
    HubAcceptsClient bool `json:"hubAcceptsClient"`
    
    // +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
    LeaseDurationSeconds *int32 `json:"leaseDurationSeconds,omitempty"`
}
```

## Useful Commands

```bash
# Code Generation
make update                   # Generate all code (clients, deepcopy, etc.)
make verify                   # Verify generated code is up to date
make clean                    # Clean generated code

# Testing
make test                     # Run unit tests
make test-integration        # Run integration tests
make verify-crds             # Verify CRD manifests

# Development
make fmt                      # Format Go code
make vet                      # Run go vet
make lint                     # Run linters
```

## Integration Examples

### External Project Integration
```go
// go.mod
module my-ocm-controller

require (
    open-cluster-management.io/api v0.13.0
    sigs.k8s.io/controller-runtime v0.16.0
)

// main.go
import (
    clusterv1 "open-cluster-management.io/api/cluster/v1"
    clusterscheme "open-cluster-management.io/api/client/cluster/clientset/versioned/scheme"
)

// Add OCM APIs to your scheme
scheme := runtime.NewScheme()
clusterscheme.AddToScheme(scheme)
```

### CRD Installation
```bash
# Install CRDs from this repository
kubectl apply -f https://raw.githubusercontent.com/open-cluster-management-io/api/main/cluster/v1/0000_00_clusters.open-cluster-management.io_managedclusters.crd.yaml
```

## Project Structure

```
├── addon/              # Addon management API definitions
│   └── v1alpha1/      # Addon API types and client code
├── cluster/            # Cluster management API definitions
│   ├── v1/            # Stable cluster APIs
│   ├── v1beta1/       # Beta placement APIs
│   └── v1alpha1/      # Alpha cluster capabilities APIs
├── work/               # Work distribution API definitions
│   ├── v1/            # Stable work APIs
│   └── v1alpha1/      # Alpha work APIs
├── operator/           # Operator installation APIs
│   └── v1/            # Operator APIs
├── client/             # Generated Kubernetes clients
│   ├── addon/         # Addon client packages
│   ├── cluster/       # Cluster client packages
│   ├── work/          # Work client packages
│   └── operator/      # Operator client packages
└── utils/              # Utility libraries and helpers
```

## Contributing Guidelines

- **API Compatibility**: Maintain backward compatibility for stable APIs
- **Code Generation**: Always run `make update` after API changes
- **Testing**: Comprehensive unit tests for all API types
- **Documentation**: Update examples and API documentation
- **Validation**: Add appropriate kubebuilder validation tags
- **Versioning**: Follow Kubernetes API versioning guidelines

This repository serves as the foundation for the entire Open Cluster Management ecosystem, enabling multicluster Kubernetes management at scale.