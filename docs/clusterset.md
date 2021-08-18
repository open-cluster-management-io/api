# ManagedClusterSet

ManagedClusterSet defines a group of ManagedClusters that user's workload can run on.

## Assign a ManagedCluster to a ManagedClusterSet
Put label `cluster.open-cluster-management.io/clusterset=<clusterset_name>` on a ManagedCluster to assign it to a ManagedClusterSet.

## Remove a ManagedCluster from a ManagedClusterSet
To remove a ManagedCluster from a ManagedClusterSet, just remove the label `cluster.open-cluster-management.io/clusterset` from the ManagedCluster.

## Permissions models
User must have a RBAC rule `CREATE` on a virtual subresource of `managedclustersets/join` to  add/remove label `cluster.open-cluster-management.io/clusterset` to/from a ManagedCluster. In order to update this label on a ManagedCluster, user must have the permission on both the old and new ManagedClusterSet.
