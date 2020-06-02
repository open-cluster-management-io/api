# ManifestWork API

ManifestWork is defined as a workload that hub desires to be deployed on the managed cluster

- ManifestWork must be created in a cluster namespace on hub. The agent on managed cluster then deploys
workload on the managed cluster.
- ManifestWork is declartive, so if the work on hub is updated, the corresponding workload
on managed cluster should also be updated. If work is deleted, the corresponding workload on
the managed cluster is deleted.
- The status of the work contains the conditions of deployed resources on managed cluster.
After agent on managed cluster deploys or patches the resources, it needs to sync the conditions of
deployed resource on managed cluster to the status.WorkloadConditions in corresponding ManifestWork on hub.

Status transition of ManifestWork API.
1. Agent on managed cluster gets `ManagedCluster` and update condition to `WorkloadProgressing`
2. Agent apply each manifest in `spec.Manifests`, and update condition in `statu.WorkloadCondition`.
  - if the resource exists. update the condition of the manifest to `Available`.
  - if the resource is applied successfully, update the condition of the manifest to
  `Applied`.
3. Agent on managed cluster checks `status.WorkloadCondition` and updates `status.Condition`.
  - if all manifests is applied, update condition to `WorkloadApplied`.
  - if all manifests exist, add a condition of `WorkloadAvailable`.

