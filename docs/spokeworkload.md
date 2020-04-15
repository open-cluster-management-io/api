# SpokeWorkload API

SpokeWorkload is defined as a workload that hub desires to be deployed on the spoke cluster

- SpokeWorkload must be created in a cluster namespace on hub. The agent on spoke then deploys
workload on the spoke cluster.
- SpokeWorkload is declartive, so if the work on hub is updated, the corresponding workload
on spoke cluster should also be updated. If work is deleted, the corresponding workload on
the spoke cluster is deleted.
- The status of the work contains the conditions of deployed resources on spoke cluster.
After agent on spoke deploys or patches the resources, it needs to sync the conditions of
deployed resource on spoke to the status.WorkloadConditions in corresponding SpokeWorkload on hub.

Status transition of SpokeWorkload API.
1. Agent on spoke gets `SpokeCluster` and update condition to `WorkloadProgressing`
2. Agent apply each manifest in `spec.Manifests`, and update condition in `statu.WorkloadCondition`.
  - if the resource exists. update the condition of the manifest to `Available`.
  - if the resource is applied successfully, update the condition of the manifest to
  `Applied`.
3. Agent on spoke check `statu.WorkloadCondition` and update `status.Condition`.
  - if all manifests is applied, update condition to `WorkloadApplied`.
  - if all manifests exist, add a condition of `WorkloadAvailable`.
  
