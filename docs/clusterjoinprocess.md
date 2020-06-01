# Cluster join processes

This is to describe the various process of cluster join

Actors:
1. cluster-admin on managed cluster
2. cluster-admin on hub cluster
3. hub controller
4. agent on managed cluster

Some rules on cluster join:
- The name of the cluster must be globally unique on hub and conforms to dns label format.

## Cluster join initiated from managed cluster
1. cluster-admin on managed cluster gets a bootstrap kubeconfig to connect to hub,
and deploy the agent on managed cluster.
  - it has the identity to create `ManagedCluster` and create/watch csr.
2. agent on managed cluster creates `ManagedCluster` if it does not exist.
  - The name of `ManagedCluster` is read from Cluster UID in openshift.
  - Otherwise agent generates a UID and saves it in `Configmap` in managed cluster, so restarting agent or redeploying
  agent will not lose the UID.
3. agent on managed cluster creates CSR on hub cluster using bootstrap kubeconfig.
  - The subject in csr is
`{"Organization": ["system:open-cluster-management:clusterName"], "CommonName":"system:open-cluster-management:clusterName:agentName"}`.
  - The name of the csr is the digest of subject and private key, with a common prefix.
  CSR will specify the signer name as the kube-client one.
4. cluster-admin on hub-cluster approve the CSR.
5. hub-controller creates a clusterrolebinding on the hub with the identity of
`system:open-cluster-management:clusterName:agentName`
   - Allows status update of `ManagedCluster`
6. cluster-admin on hub update `spec.hubAcceptsClient` to `true`.
  - Only user on hub who has the RBAC permision to update subresource of `managedclusters/accept`
  can update this field.
7. hub-controller updates condition of `ManagedCluster` to `HubApprovedJoin`.
8. hub-controller creates a namespace as the name of cluster on hub cluster if it does not exist.
  - managed cluster can only join a hub once, and it can join to multiple hubs.
  - The UID of the managed cluster is identical on each of the hub the Klusterlet agent joins.
9. hub-controller creates role/rolebinding on the cluster namespace on the hub
  - Allow the access of agent on managed cluster to the namespace.
10. agent on managed cluster gets certificate in CSR status, uses the certificate to create a new kubeconfig
and saves it as secret.
10. agent on managed cluster connects to hub apiserver using the new kubeconfig.
11. agent on managed cluster updates conditions of `ManagedCluster` as `ManagedClusterJoined`.
12. agent on managed cluster appends updates other fields in status of `ManagedCluster`.

## Certificate renewal
1. agent on managed cluster detects the certificate is going to be expired.
  - it checks if certificate will be expired in 20% of certificate duration.
2. agent on managed cluster generates a new private key and submits a new CSR to hub apiserver.
  - it uses the identity of `system:open-cluster
management:clusterName:agentName` to create the csr
  - the subject in the certificate should be `{"Organization": ["system:open-cluster-management:clusterName"],
  "CommonName":"system:open-cluster-management:clusterName:agentName"}`
3. hub controller auto approves the csr. hub controller checks if the csr can be approved
based on the following steps:
- check if organization field and commonName field is valid.
- check if user name in csr is the same as commonName in certificate to ensure the request
is originated from the same identity.
- check if the corresponding `ManagedCluster` is in the conndition of `HubApprovedJoin`.
4. agent on managed cluster reconstructs the kubeconfig using the new key/certificate
and saves it as a secret on managed cluster.
