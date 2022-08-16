package v1alpha1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_AddOnMeta = map[string]string{
	"":            "AddOnMeta represents a collection of metadata information for the add-on.",
	"displayName": "displayName represents the name of add-on that will be displayed.",
	"description": "description represents the detailed description of the add-on.",
}

func (AddOnMeta) SwaggerDoc() map[string]string {
	return map_AddOnMeta
}

var map_ClusterManagementAddOn = map[string]string{
	"":       "ClusterManagementAddOn represents the registration of an add-on to the cluster manager. This resource allows the user to discover which add-on is available for the cluster manager and also provides metadata information about the add-on. This resource also provides a linkage to ManagedClusterAddOn, the name of the ClusterManagementAddOn resource will be used for the namespace-scoped ManagedClusterAddOn resource. ClusterManagementAddOn is a cluster-scoped resource.",
	"spec":   "spec represents a desired configuration for the agent on the cluster management add-on.",
	"status": "status represents the current status of cluster management add-on.",
}

func (ClusterManagementAddOn) SwaggerDoc() map[string]string {
	return map_ClusterManagementAddOn
}

var map_ClusterManagementAddOnList = map[string]string{
	"":         "ClusterManagementAddOnList is a collection of cluster management add-ons.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
	"items":    "Items is a list of cluster management add-ons.",
}

func (ClusterManagementAddOnList) SwaggerDoc() map[string]string {
	return map_ClusterManagementAddOnList
}

var map_ClusterManagementAddOnSpec = map[string]string{
	"":                   "ClusterManagementAddOnSpec provides information for the add-on.",
	"addOnMeta":          "addOnMeta is a reference to the metadata information for the add-on.",
	"addOnConfiguration": "addOnConfiguration is a reference to configuration information for the add-on. In scenario where a multiple add-ons share the same add-on CRD, multiple ClusterManagementAddOn resources need to be created and reference the same AddOnConfiguration.",
}

func (ClusterManagementAddOnSpec) SwaggerDoc() map[string]string {
	return map_ClusterManagementAddOnSpec
}

var map_ClusterManagementAddOnStatus = map[string]string{
	"": "ClusterManagementAddOnStatus represents the current status of cluster management add-on.",
}

func (ClusterManagementAddOnStatus) SwaggerDoc() map[string]string {
	return map_ClusterManagementAddOnStatus
}

var map_ConfigCoordinates = map[string]string{
	"":                       "ConfigCoordinates represents the information for locating the CRD and CR that configures the add-on.",
	"crdName":                "crdName is the name of the CRD used to configure instances of the managed add-on. This field should be configured if the add-on have a CRD that controls the configuration of the add-on.",
	"crName":                 "crName is the name of the CR used to configure instances of the managed add-on. This field should be configured if add-on CR have a consistent name across the all of the ManagedCluster instaces.",
	"lastObservedGeneration": "lastObservedGeneration is the observed generation of the custom resource for the configuration of the addon.",
}

func (ConfigCoordinates) SwaggerDoc() map[string]string {
	return map_ConfigCoordinates
}

var map_HealthCheck = map[string]string{
	"mode": "mode indicates which mode will be used to check the healthiness status of the addon.",
}

func (HealthCheck) SwaggerDoc() map[string]string {
	return map_HealthCheck
}

var map_ManagedClusterAddOn = map[string]string{
	"":       "ManagedClusterAddOn is the Custom Resource object which holds the current state of an add-on. This object is used by add-on operators to convey their state. This resource should be created in the ManagedCluster namespace.",
	"spec":   "spec holds configuration that could apply to any operator.",
	"status": "status holds the information about the state of an operator.  It is consistent with status information across the Kubernetes ecosystem.",
}

func (ManagedClusterAddOn) SwaggerDoc() map[string]string {
	return map_ManagedClusterAddOn
}

var map_ManagedClusterAddOnList = map[string]string{
	"": "ManagedClusterAddOnList is a list of ManagedClusterAddOn resources.",
}

func (ManagedClusterAddOnList) SwaggerDoc() map[string]string {
	return map_ManagedClusterAddOnList
}

var map_ManagedClusterAddOnSpec = map[string]string{
	"":                 "ManagedClusterAddOnSpec defines the install configuration of an addon agent on managed cluster.",
	"installNamespace": "installNamespace is the namespace on the managed cluster to install the addon agent. If it is not set, open-cluster-management-agent-addon namespace is used to install the addon agent.",
}

func (ManagedClusterAddOnSpec) SwaggerDoc() map[string]string {
	return map_ManagedClusterAddOnSpec
}

var map_ManagedClusterAddOnStatus = map[string]string{
	"":                   "ManagedClusterAddOnStatus provides information about the status of the operator.",
	"conditions":         "conditions describe the state of the managed and monitored components for the operator.",
	"relatedObjects":     "relatedObjects is a list of objects that are \"interesting\" or related to this operator. Common uses are: 1. the detailed resource driving the operator 2. operator namespaces 3. operand namespaces 4. related ClusterManagementAddon resource",
	"addOnMeta":          "addOnMeta is a reference to the metadata information for the add-on. This should be same as the addOnMeta for the corresponding ClusterManagementAddOn resource.",
	"addOnConfiguration": "addOnConfiguration is a reference to configuration information for the add-on. This resource is use to locate the configuration resource for the add-on.",
	"registrations":      "registrations is the conifigurations for the addon agent to register to hub. It should be set by each addon controller on hub to define how the addon agent on managedcluster is registered. With the registration defined, The addon agent can access to kube apiserver with kube style API or other endpoints on hub cluster with client certificate authentication. A csr will be created per registration configuration. If more than one registrationConfig is defined, a csr will be created for each registration configuration. It is not allowed that multiple registrationConfigs have the same signer name. After the csr is approved on the hub cluster, the klusterlet agent will create a secret in the installNamespace for the registrationConfig. If the signerName is \"kubernetes.io/kube-apiserver-client\", the secret name will be \"{addon name}-hub-kubeconfig\" whose contents includes key/cert and kubeconfig. Otherwise, the secret name will be \"{addon name}-{signer name}-client-cert\" whose contents includes key/cert.",
	"healthCheck":        "healthCheck indicates how to check the healthiness status of the current addon. It should be set by each addon implementation, by default, the lease mode will be used.",
}

func (ManagedClusterAddOnStatus) SwaggerDoc() map[string]string {
	return map_ManagedClusterAddOnStatus
}

var map_ObjectReference = map[string]string{
	"":          "ObjectReference contains enough information to let you inspect or modify the referred object.",
	"group":     "group of the referent.",
	"resource":  "resource of the referent.",
	"namespace": "namespace of the referent.",
	"name":      "name of the referent.",
}

func (ObjectReference) SwaggerDoc() map[string]string {
	return map_ObjectReference
}

var map_RegistrationConfig = map[string]string{
	"":                  "RegistrationConfig defines the configuration of the addon agent to register to hub. The Klusterlet agent will create a csr for the addon agent with the registrationConfig.",
	"signerName":        "signerName is the name of signer that addon agent will use to create csr.",
	"subject":           "subject is the user subject of the addon agent to be registered to the hub. If it is not set, the addon agent will have the default subject \"subject\": {\n\t\"user\": \"system:open-cluster-management:addon:{addonName}:{clusterName}:{agentName}\",\n\t\"groups: [\"system:open-cluster-management:addon\", \"system:open-cluster-management:addon:{addonName}\", \"system:authenticated\"]\n}",
	"certificateStatus": "certificateStatus actively tracks the status of the certificate used by the addon.",
}

func (RegistrationConfig) SwaggerDoc() map[string]string {
	return map_RegistrationConfig
}

var map_RegistrationConfigCertificateStatus = map[string]string{
	"lastRenewedTimestamp": "lastRenewedTimestamp records the last timestamp when we approved/renewed certificates for the addon agents.",
	"expiringTimestamp":    "expiringTimestamp records the next time certificate will expire.",
}

func (RegistrationConfigCertificateStatus) SwaggerDoc() map[string]string {
	return map_RegistrationConfigCertificateStatus
}

var map_Subject = map[string]string{
	"":                 "Subject is the user subject of the addon agent to be registered to the hub.",
	"user":             "user is the user name of the addon agent.",
	"groups":           "groups is the user group of the addon agent.",
	"organizationUnit": "organizationUnit is the ou of the addon agent",
}

func (Subject) SwaggerDoc() map[string]string {
	return map_Subject
}

// AUTO-GENERATED FUNCTIONS END HERE
