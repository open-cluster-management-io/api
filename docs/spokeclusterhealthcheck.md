# Spoke cluster available check

This document describes a mechanism on how hub cluster check the availability of agent running
on spoke cluster, a new condition `SpokeClusterConditionAvailable` will be introduced to describe
whether a spoke cluster is available

If the kube-apiserver is health and the registration agent is running with the minimum
deployment on a spoke cluster, this spoke cluster will be considered available

The Kubernetes Metrics API may be used to get the status of spoke cluster kube-apiserver or spoke
registration agent, but Metrics API has some limits for our cases

- Reporting metrics will require additional request path from hub to spoke (or spoke to hub), and
if we get the metrics from hub to spoke, the hub may not access the spoke cluster
- To know whether a spoke cluster is available, the hub may need periodically poll all of spoke
clusters, this is an inefficient way

So we propose hub and spoke coordinate a condition to present whether a spoke cluster is available

Actors

- agent on spoke cluster
- hub controller on hub cluster

## Available check

1. after a spoke cluster is accepted by hub cluster admin, the agent on spoke cluster will
create a [`Lease`](https://github.com/kubernetes/api/blob/master/coordination/v1/types.go#L27) with
a label `open-cluster-management.io/cluster-name=<cluster-name>` on its namespace on hub cluster,
then the agent will periodically (by default 60s, the period time should be adjusted, see the scale
problem section)
    - update the `RenewTime` of `Lease` to keep available, and
    - check the spoke cluster kube-apiserver health by `healthz` API, and check the
    registration agent is running with the minimum deployment, if all of them are
    available the agent make sure the status of `SpokeClusterConditionAvailable`
    condition is `true`, otherwise, the agent updates the status to `false`
2. hub controller watches the spoke clusters on hub cluster with its informer, for
each accepted spoke cluster, the controller get its `Lease` from its namespace and check
wether the `Lease` is constantly updated in a duration (by default 5min), if not, the
controller updates the status of this spoke cluster `SpokeClusterConditionAvailable`
condition to `unknown`

> Note: the `SpokeClusterConditionAvailable` condition only represents the availability of
registration agent, other agents should not expand `SpokeCluster` status.

## The scale problem

If there are many spoke clusters, we will face there are a lot of QPS for an idle system due to
spoke lease update request, to solve this problem, we need a way to control the frequency of
lease update, we will add `LeaseDuration` field in `SpokeCluster` API, the hub cluster admin is
able to update it according to system QPS, and spoke agent can adjust the lease update time
accordingly
