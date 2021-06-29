# Open Cluster Management API

<a href="https://godoc.org/open-cluster-management.io/api"><img src="https://godoc.org/open-cluster-management.io/api?status.svg"></a> <a href="https://goreportcard.com/report/open-cluster-management.io/api"><img alt="Go Report Card" src="https://goreportcard.com/badge/open-cluster-management.io/api" /></a>

The `api` repository defines relevant concepts and types for problem domains related to managing 0..* Kubernetes clusters.

## Purpose

This library is the canonical location of the Open Cluster Management API definition.

- API for cluster registration independent of cluster CRUD lifecycle.
- API for work distribution across multiple clusters.
- API for dynamic placement of content and behavior across multiple clusters.
- API for building addon/extension based on the foundation components for the purpose of working with multiple clusters.

## Consumers

Various projects under [Open Cluster Management](https://github.com/open-cluster-management-io) leverage this `api` library. 

* [registration](https://github.com/open-cluster-management-io/registration): implements `ManagedCluster`, `ManagedClusterSet`,`ClusterClaim` for cluster registration and lifecycling.
* [work](https://github.com/open-cluster-management-io/work): implements `Manifestwork` for native kubernetes resource distribution to multiple clusters.
* [addon-framework](https://github.com/open-cluster-management-io/addon-framework): implements `ClusterManagementAddOn`, `ManagedClusterAddOn` for addon management.
* [placement](https://github.com/open-cluster-management-io/placement): implements `Placement` for cluster selection with various policies to deploy workloads.
* [registration-operator](https://github.com/open-cluster-management-io/registration-operator): implements `ClusterManager`, `Klusterlet` as an operator to deploy registration, work and placement.   

## Use case

With the Open Cluster Management API and components, you can use the [clusteradm CLI](https://github.com/open-cluster-management-io/clusteradm) to bootstrap a control plane for multicluster management. The following diagram illustrates the deployment architecture for the Open Cluster Management.

![Architecture diagram](https://github.com/open-cluster-management/community/raw/main/assets/ocm-arch.png)

## Community, discussion, contribution, and support

Check the [CONTRIBUTING Doc](CONTRIBUTING.md) for how to contribute to the repo.

### Communication channels

Slack channel: [#open-cluster-mgmt](http://slack.k8s.io/#open-cluster-mgmt)

## License

This code is released under the Apache 2.0 license. See the file LICENSE for more information.
