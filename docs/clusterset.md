# ManagedClusterSet

ManagedClusterSet defines a group of ManagedClusters that user's workload can run on.

## Assign a ManagedCluster to a ManagedClusterSet
Put label `clusters.open-cluster-management.io/clusterset=<clusterset_name>` on a ManagedCluster to assign it to a ManagedClusterSet.

## Remove a ManagedCluster from a ManagedClusterSet
To Remove a ManagedCluster from a ManagedClusterSet, jut remove the label `clusters.open-cluster-management.io/clusterset` from the ManagedCluster.

## Permissions models
User must have a RBAC rule to `CREATE` on a virtual subresource of `managedclustersets/join` to  add/remove label `clusters.open-cluster-management.io/clusterset` to/from a ManagedCluster. In order to update this label on a ManagedCluster, user must have the permission on both the old and new ManagedClusterSet.
