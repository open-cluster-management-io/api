<!-- Copyright Contributors to the Open Cluster Management project -->
# Open Cluster Management API

[![Go Reference](https://pkg.go.dev/badge/open-cluster-management.io/api.svg)](https://pkg.go.dev/open-cluster-management.io/api)
[![Go Report Card](https://goreportcard.com/badge/open-cluster-management.io/api)](https://goreportcard.com/report/open-cluster-management.io/api)

The official API definitions and client libraries for [Open Cluster Management (OCM)](https://open-cluster-management.io/), a CNCF sandbox project focused on multicluster and multicloud scenarios for Kubernetes.

## Purpose

This repository serves as the **canonical source** for all Open Cluster Management API definitions and generated client libraries. It contains API types, CRDs, client code, and helper utilities, making it lightweight and easy for external projects to import and integrate with OCM.

### Why This Repository Exists

- **Easy Integration**: External projects can import OCM APIs without pulling in implementation dependencies
- **API Definitions**: Single source of truth for all OCM Custom Resource Definitions (CRDs)  
- **Generated Clients**: Kubernetes-style typed clients for all OCM APIs
- **Lightweight**: No core controllers or operators, just APIs and helper utilities

## API Groups

This repository defines four main API groups that cover the complete OCM ecosystem:

### Cluster Management APIs (`cluster.open-cluster-management.io`)

Manage cluster lifecycle, registration, and placement across your fleet. The `v1` APIs provide core cluster registration and grouping with `ManagedCluster`, `ManagedClusterSet`, and `ManagedClusterSetBinding`. The `v1beta1` APIs offer advanced cluster selection and placement with `Placement` and `PlacementDecision`. The `v1alpha1` APIs include cluster capabilities and scoring through `ClusterClaim` and `AddonPlacementScore`.

### Work Distribution APIs (`work.open-cluster-management.io`)

Deploy and manage Kubernetes resources across multiple clusters. The `v1` APIs provide `ManifestWork` and `AppliedManifestWork` to deploy any Kubernetes resource to managed clusters. The `v1alpha1` APIs include `ManifestWorkReplicaSet` for bulk deployment across cluster sets.

### Addon Management APIs (`addon.open-cluster-management.io`)

Build and manage extensions on top of the OCM foundation. The `v1alpha1` APIs provide addon lifecycle and configuration through `ClusterManagementAddOn`, `ManagedClusterAddOn`, and `AddonDeploymentConfig`.

### Operator APIs (`operator.open-cluster-management.io`)

Install and configure OCM components. The `v1` APIs provide `ClusterManager` and `Klusterlet` for hub and agent installation and configuration.

## Quick Start

### Installation

Add this module to your Go project:

```bash
go get open-cluster-management.io/api@latest
```

### Usage Example

#### Working with ManagedClusters

```go
import (
    clusterv1 "open-cluster-management.io/api/cluster/v1"
    clusterclientset "open-cluster-management.io/api/client/cluster/clientset/versioned"
)

// Create a client
client := clusterclientset.NewForConfigOrDie(config)

// List all managed clusters
clusters, err := client.ClusterV1().ManagedClusters().List(ctx, metav1.ListOptions{})

// Create a new managed cluster
cluster := &clusterv1.ManagedCluster{
    ObjectMeta: metav1.ObjectMeta{
        Name: "my-cluster",
    },
    Spec: clusterv1.ManagedClusterSpec{
        HubAcceptsClient: true,
    },
}
```

## Project Structure

```
api/
├── addon/           # Addon management APIs
├── cluster/         # Cluster management APIs  
├── operator/        # Operator APIs
├── work/            # Work distribution APIs
├── client/          # Generated Kubernetes clients for all APIs
└── utils/           # Utility libraries
```

## Who Uses This Repository

### Official OCM Projects

These projects implement the APIs defined in this repository:

- **[registration](https://github.com/open-cluster-management-io/registration)**: Implements `ManagedCluster` and cluster lifecycle
- **[work](https://github.com/open-cluster-management-io/work)**: Implements `ManifestWork` for resource distribution  
- **[placement](https://github.com/open-cluster-management-io/placement)**: Implements `Placement` for intelligent cluster selection
- **[addon-framework](https://github.com/open-cluster-management-io/addon-framework)**: Implements addon management APIs
- **[registration-operator](https://github.com/open-cluster-management-io/registration-operator)**: Implements operator APIs

### External Projects

This API library is designed for external projects that want to:

- Build OCM-compatible controllers and operators
- Integrate with OCM clusters from existing platforms  
- Create custom tools for multicluster management
- Develop addons and extensions for OCM

## Development

### Prerequisites

- Go 1.23+
- Kubernetes 1.28+

### Building

```bash
# Update generated code
make update

# Verify code generation and formatting  
make verify

# Run tests
make test
```

### Code Generation

This repository uses Kubernetes code generators to create:

- **DeepCopy methods**: For all API types
- **Typed clients**: Kubernetes-style clientsets  
- **Informers and listers**: For efficient caching and watching
- **OpenAPI schemas**: For CRD validation

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## Community

Connect with the OCM community:

- [Website](https://open-cluster-management.io)
- [Slack](https://kubernetes.slack.com/channels/open-cluster-mgmt)
- [Mailing group](https://groups.google.com/g/open-cluster-management)
- [Community meetings](https://calendar.google.com/calendar/u/0/embed?src=openclustermanagement@gmail.com)
- [YouTube channel](https://www.youtube.com/c/OpenClusterManagement)

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

**[Open Cluster Management](https://open-cluster-management.io/)** is a [Cloud Native Computing Foundation](https://cncf.io/) sandbox project.
