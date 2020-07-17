package v1

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
var map_ClusterManager = map[string]string{
	"":       "ClusterManager configures the controllers on the hub that govern registration and work distribution for attached Klusterlets. ClusterManager will only be deployed in open-cluster-management-hub namespace.",
	"spec":   "Spec represents a desired deployment configuration of controllers that govern registration and work distribution for attached Klusterlets.",
	"status": "Status represents the current status of controllers that govern the lifecycle of managed clusters.",
}

func (ClusterManager) SwaggerDoc() map[string]string {
	return map_ClusterManager
}

var map_ClusterManagerList = map[string]string{
	"":         "ClusterManagerList is a collection of deployment configurations for registration and work distribution controllers.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
	"items":    "Items is a list of deployment configurations for registration and work distribution controllers.",
}

func (ClusterManagerList) SwaggerDoc() map[string]string {
	return map_ClusterManagerList
}

var map_ClusterManagerSpec = map[string]string{
	"":                          "ClusterManagerSpec represents a desired deployment configuration of controllers that govern registration and work distribution for attached Klusterlets.",
	"registrationImagePullSpec": "RegistrationImagePullSpec represents the desired image of registration controller installed on the hub.",
}

func (ClusterManagerSpec) SwaggerDoc() map[string]string {
	return map_ClusterManagerSpec
}

var map_ClusterManagerStatus = map[string]string{
	"":                   "ClusterManagerStatus represents the current status of the registration and work distribution controllers running on the hub.",
	"observedGeneration": "ObservedGeneration is the last generation change.",
	"conditions":         "Conditions contain the different condition statuses for this ClusterManager. Valid condition types are: Applied: Components in hub are applied. Available: Components in hub are available and ready to serve. Progressing: Components in hub are in a transitioning state. Degraded: Components in hub do not match the desired configuration and only provide degraded service.",
	"generations":        "Generations are used to determine when an item needs to be reconciled or has changed in a way that needs a reaction.",
	"relatedResources":   "RelatedResources are used to track the resources that are related to this ClusterManager.",
}

func (ClusterManagerStatus) SwaggerDoc() map[string]string {
	return map_ClusterManagerStatus
}

var map_GenerationStatus = map[string]string{
	"":               "GenerationStatus keeps track of the generation for a given resource so that decisions about forced updates can be made. The definition matches the GenerationStatus defined in github.com/openshift/api/v1".
	"group":          "The group property is the group of the resource that you're tracking.",
	"version":        "The version property is the version of the resource that you're tracking.",
	"resource":       "The resource property is the resource type of the resource that you're tracking.",
	"namespace":      "The namespace property is the location of the resource that you're tracking.",
	"name":           "The name property is the name of the resource that you're tracking.",
	"lastGeneration": "The lastGeneration property is the last generation of the thing that the controller applies.",
}

func (GenerationStatus) SwaggerDoc() map[string]string {
	return map_GenerationStatus
}

var map_Klusterlet = map[string]string{
	"":       "Klusterlet represents controllers on the managed cluster. When configured, the Klusterlet requires a secret named of bootstrap-hub-kubeconfig in the same namespace to allow API requests to the hub for the registration protocol.",
	"spec":   "Spec represents the desired deployment configuration of Klusterlet agent.",
	"status": "Status represents the current status of Klusterlet agent.",
}

func (Klusterlet) SwaggerDoc() map[string]string {
	return map_Klusterlet
}

var map_KlusterletList = map[string]string{
	"":         "KlusterletList is a collection of Klusterlet agents.",
	"metadata": "Standard list metadata. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
	"items":    "Items is a list of Klusterlet agents.",
}

func (KlusterletList) SwaggerDoc() map[string]string {
	return map_KlusterletList
}

var map_KlusterletSpec = map[string]string{
	"":                          "KlusterletSpec represents the desired deployment configuration of Klusterlet agent.",
	"namespace":                 "Namespace is the namespace to deploy the agent. The namespace must have a prefix of \"open-cluster-management-\", and if it is not set, the namespace of \"open-cluster-management-agent\" is used to deploy agent.",
	"registrationImagePullSpec": "RegistrationImagePullSpec represents the desired image configuration of registration agent.",
	"workImagePullSpec":         "WorkImagePullSpec represents the desired image configuration of work agent.",
	"clusterName":               "ClusterName is the name of the managed cluster to be created on hub. The Klusterlet agent generates a random name if it is not set, or discovers the appropriate cluster name on OpenShift.",
	"externalServerURLs":        "ExternalServerURLs represents the a list of apiserver URLs and ca bundles that is accessible externally. If it is set empty, the managed cluster has no externally accessible URL that the hub cluster can visit.",
}

func (KlusterletSpec) SwaggerDoc() map[string]string {
	return map_KlusterletSpec
}

var map_KlusterletStatus = map[string]string{
	"":                   "KlusterletStatus represents the current status of Klusterlet agent.",
	"observedGeneration": "ObservedGeneration is the last generation change.",
	"conditions":         "Conditions contain the different condition statuses for this Klusterlet. Valid condition types are: Applied: Components have been applied in the managed cluster. Available: Components in the managed cluster are available and ready to serve. Progressing: Components in the managed cluster are in a transitioning state. Degraded: Components in the managed cluster do not match the desired configuration and only provide degraded service.",
	"generations":        "Generations are used to determine when an item needs to be reconciled or has changed in a way that needs a reaction.",
	"relatedResources":   "RelatedResources are used to track the resources that are related to this Klusterlet",
}

func (KlusterletStatus) SwaggerDoc() map[string]string {
	return map_KlusterletStatus
}

var map_RelatedResourceMeta = map[string]string{
	"":          "The RelatedResourceMeta property represents the resource that is managed by an operator.",
	"group":     "The group property is the group of the resource that you're tracking.",
	"version":   "The version property is the version of the resource that you're tracking.",
	"resource":  "The resource property is the resource type of the resource that you're tracking.",
	"namespace": "The namespace is the location of resource that you're tracking.",
	"name":      "The name is the name of the resource that you're tracking.",
}

func (RelatedResourceMeta) SwaggerDoc() map[string]string {
	return map_RelatedResourceMeta
}

var map_ServerURL = map[string]string{
	"":         "ServerURL represents the apiserver URL and ca bundle that is accessible externally.",
	"url":      "URL is the URL of apiserver endpoint of the managed cluster.",
	"caBundle": "CABundle is the ca bundle to connect to apiserver of the managed cluster. System certs are used if it is not set.",
}

func (ServerURL) SwaggerDoc() map[string]string {
	return map_ServerURL
}

var map_StatusCondition = map[string]string{
	"":                   "StatusCondition contains condition information.",
	"type":               "Type is the type of the cluster condition.",
	"status":             "Status is the status of the condition. Valid values are True, False, or Unknown.",
	"lastTransitionTime": "LastTransitionTime is the last time the condition changed from one status to another.",
	"reason":             "Reason is a brief reason for the condition's last status change.",
	"message":            "Message is a human-readable message indicating details about the last status change.",
}

func (StatusCondition) SwaggerDoc() map[string]string {
	return map_StatusCondition
}

// AUTO-GENERATED FUNCTIONS END HERE
