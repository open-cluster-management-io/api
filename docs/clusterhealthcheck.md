# Managed cluster availability check

This document describes a mechanism on how hub cluster check the availability of agent running
on managed cluster, a new condition `ManagedClusterConditionAvailable` will be introduced to describe
whether a managed cluster is available

If the kube-apiserver is health and the registration agent is running with the minimum
deployment on a managed cluster, this managed cluster will be considered available

The Kubernetes Metrics API may be used to get the status of managed cluster kube-apiserver or managed
registration agent, but Metrics API has some limits for our cases

- Reporting metrics will require additional request path from hub to managed cluster (or managed cluster to hub), and
if we get the metrics from hub to managed cluster, the hub may not access the managed cluster
- To know whether a managed cluster is available, the hub may need periodically poll all of managed
clusters, this is an inefficient way

So we propose hub and managed coordinate a condition to present whether a managed cluster is available

Actors

- Klusterlet agent on managed cluster
- hub controller on hub cluster

## Available check

1. after a managed cluster is accepted by hub cluster admin, the agent on managed cluster will
create a [`Lease`](https://github.com/kubernetes/api/blob/master/coordination/v1/types.go#L27) with
a label `open-cluster-management.io/cluster-name=<cluster-name>` on its namespace on hub cluster,
then the agent will periodically (by default 60s, the period time should be adjusted, see the scale
problem section)
    - update the `RenewTime` of `Lease` to keep available, and
    - check the managed cluster kube-apiserver health by `healthz` API, and check the
    registration agent is running with the minimum deployment, if all of them are
    available the agent make sure the status of `ManagedClusterConditionAvailable`
    condition is `true`, otherwise, the agent updates the status to `false`
2. hub controller watches the managed clusters on hub cluster with its informer, for
each accepted managed cluster, the controller get its `Lease` from its namespace and check
wether the `Lease` is constantly updated in a duration (by default 5min), if not, the
controller updates the status of this managed cluster `ManagedClusterConditionAvailable`
condition to `unknown`

> Note: the `ManagedClusterConditionAvailable` condition only represents the availability of
registration agent, other agents should not expand `ManagedCluster` status.

## The scale problem

If there are many managed clusters, we will face there are a lot of QPS for an idle system due to
managed lease update request, to solve this problem, we need a way to control the frequency of
lease update, we will add `LeaseDuration` field in `ManagedCluster` API, the hub cluster admin is
able to update it according to system QPS, and managed agent can adjust the lease update time
accordingly
