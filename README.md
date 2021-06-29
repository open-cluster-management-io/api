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

Various projects under [Open Cluster Management](https://github.com/open-cluster-management-io) leverage this `api` library. For example, see [registration-operator](https://github.com/open-cluster-management-io/registration-operator) for how the `api` is being consumed.

## Architecutre

The following diagram illustrates how the Open Cluster Management API and components can be used to bootstrap a control plane for multicluster management.

![Architecture diagram](https://github.com/open-cluster-management/community/raw/main/assets/ocm-arch.png)

## Community, discussion, contribution, and support

Check the [CONTRIBUTING Doc](CONTRIBUTING.md) for how to contribute to the repo.

### Communication channels

Slack channel: [#open-cluster-mgmt](http://slack.k8s.io/#open-cluster-mgmt)

## License

This code is released under the Apache 2.0 license. See the file LICENSE for more information.
