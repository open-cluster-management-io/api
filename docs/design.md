
# Overview

This repo describes the API implemented by controllers in the `open-cluster-management.io` project.

The API enables a cluster to become a "Hub" that manages a set of other clusters.

When a cluster becomes "managed", an agent is deployed that registers the managed cluster with the hub. Once registration is approved, the managed cluster can receive directives to apply specific Kubernetes manifests to run additional workloads.

# API

## Cluster Manager & Klusterlet

API Group: `operator.open-cluster-management.io/v1`
Kinds:
- `ClusterManager`. Provides lifecycle management for the registration controller on the Hub.
- `Klusterlet`. Provides lifecycle management for the registration and `ManifestWork` controllers on the managed cluster.

## Managed Cluster

API Group: `cluster.open-cluster-management.io/v1`
Kinds:
- `ManagedCluster`. Provides a representation of the managed cluster on the hub. `ManagedCluster` controls the lifecycle of whether the remote cluster has been "accepted" by the Hub for management and can retrieve information from the Hub to direct a set of manifests or actions to apply.

## Cluster

API Group: `work.open-cluster-management.io/v1`
Kinds:
- `ManifestWork`. Provides a representation of manifests that should be applied to a managed cluster.