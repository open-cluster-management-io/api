# Cloudevents Clients

We have implemented the [cloudevents](https://cloudevents.io/)-based clients in this package to assist developers in
easily implementing the [Event Based Manifestwork](https://github.com/open-cluster-management-io/enhancements/tree/main/enhancements/sig-architecture/224-event-based-manifestwork)
proposal.

## Generic Clients

The generic client (`generic.CloudEventsClient`) is used to resync/publish/subscribe resource objects between sources
and agents with cloudevents.

A resource object can be any object that implements the `generic.ResourceObject` interface.

### Building a generic client on a source

Developers can use `generic.NewCloudEventSourceClient` method to build a generic client on the source. To build this
client the developers need to provide

1. A cloudevents source options (`options.CloudEventsSourceOptions`), this options have two parts
    -  `sourceID`, it is a unique identifier for a source, for example, it can generate a source ID by hashing the hub
    cluster URL and appending the controller name. Similarly, a RESTful service can select a unique name or generate a
    unique ID in the associated database for its source identification.
    - `CloudEventsOptions`, it provides cloudevents clients to send/receive cloudevents based on different event
    protocol. We have supported [MQTT protocol (`mqtt.NewSourceOptions`)](./generic/options/mqtt) and [gRPC protocol (`grpc.NewSourceOptions`)](./generic/options/grpc) developers can use it directly.

2. A resource lister (`generic.Lister`), it is used to list the resource objects on the source when resyncing the
resources between sources and agents, for example, a hub controller can list the resources from the resource informers,
and a RESTful service can list its resources from a database.

3. A resource status hash getter method (`generic.StatusHashGetter`), this method will be used to calculate the resource
status hash when resyncing the resource status between sources and agents.

4. Codecs (`generic.Codec`), they are used to encode a resource object into a cloudevent and decode a cloudevent into a
resource object with a given cloudevent data type. We have provided two data types (`io.open-cluster-management.works.v1alpha1.manifests`
that contains a single resource object in the cloudevent payload and `io.open-cluster-management.works.v1alpha1.manifestbundles`
that contains a list of resource objects in the cloudevent payload) for `ManifestWork`, they can be found in the `work/payload`
package.

5. Resource handler methods (`generic.ResourceHandler`), they are used to handle the resources status after the client
received the resources status from agents.

for example, build a generic client on the source using MQTT protocol with the following code:

```golang
// build a client for the source1
client, err := generic.NewCloudEventSourceClient[*CustomerResource](
        ctx,
        mqtt.NewSourceOptions(mqtt.NewMQTTOptions(), "source1"),
        customerResourceLister,
		customerResourceStatusHashGetter,
		customerResourceCodec,
	)

// start a go routine to receive the resources status from agents
go func() {
	if err := client.Subscribe(ctx, customerResourceHandler); err != nil {
		//TODO handle this error when subscribing the cloudevents failed
	}
}()
```

You may refer to the [cloudevents client integration test](../test/integration/cloudevents/source) as an example.

### Building a generic client on a manged cluster

Developers can use `generic.NewCloudEventAgentClient` method to build a generic client on a managed cluster. To build
this client the developers need to provide

1. A cloudevents agent options (`options.CloudEventsAgentOptions`), this options have three parts
    -  `agentID`, it is a unique identifier for an agent, for example, it can consist of a managed cluster name and an
    agent name.
    - `clusterName`, it is the name of a managed cluster on which the agent runs.
    - `CloudEventsOptions`, it provides cloudevents clients to send/receive cloudevents based on different event
    protocol. We have supported [MQTT protocol (`mqtt.NewAgentOptions`)](./generic/options/mqtt) and [gRPC protocol (`grpc.NewAgentOptions`)](./generic/options/grpc) , developers can use it directly.

2. A resource lister (`generic.Lister`), it is used to list the resource objects on a managed cluster when resyncing the
resources between sources and agents, for example, a work agent can list its works from its work informers.

3. A resource status hash getter method (`generic.StatusHashGetter`), this method will be used to calculate the resource
status hash when resyncing the resource status between sources and agents.

4. Codecs (`generic.Codec`), they are used to encode a resource object into a cloudevent and decode a cloudevent into a
resource object with a given cloudevent data type. We have provided two data types (`io.open-cluster-management.works.v1alpha1.manifests`
that contains a single resource object in the cloudevent payload and `io.open-cluster-management.works.v1alpha1.manifestbundles`
that contains a list of resource objects in the cloudevent payload) for `ManifestWork`, they can be found in the `work/payload`
package.

5. Resource handler methods (`generic.ResourceHandler`), they are used to handle the resources after the client received
the resources from sources.

for example, build a generic client on the source using MQTT protocol with the following code:

```golang
// build a client for a work agent on the cluster1
client, err := generic.NewCloudEventAgentClient[*CustomerResource](
        ctx,
        mqtt.NewAgentOptions(mqtt.NewMQTTOptions(), "cluster1", "cluster1-work-agent"),
        &ManifestWorkLister{},
		ManifestWorkStatusHash,
		&ManifestBundleCodec{},
	)

// start a go routine to receive the resources from sources
go func() {
	if err := client.Subscribe(ctx, NewManifestWorkAgentHandler()); err != nil {
		//TODO handle this error when subscribing the cloudevents failed
	}
}()
```

## Work Clients

We have provided a builder to build the `ManifestWork` client (`ManifestWorkInterface`) and informer (`ManifestWorkInformer`)
based on the generic client.

### Building work client for work controllers on the hub cluster

TODO

### Building work client for work agent on the managed cluster

Developers can use the builder to build the `ManifestWork` client and informer with the cluster name.

```golang

clusterName := "cluster1"
// Building the clients based on cloudevents with MQTT
config := mqtt.NewMQTTOptions()

clientHolder, err := work.NewClientHolderBuilder(fmt.Sprintf("%s-work-agent", clusterName), config).
	WithClusterName(clusterName).
    // Supports two event data types for ManifestWork
	WithCodecs(codec.NewManifestBundleCodec(), codec.NewManifestCodec(restMapper)).
	NewClientHolder(ctx)
if err != nil {
	return err
}

manifestWorkClient := clientHolder.ManifestWorks(clusterName)
manifestWorkInformer := clientHolder.ManifestWorkInformer()

// Building controllers with ManifestWork client and informer ...

// Start the ManifestWork informer
go manifestWorkInformer.Informer().Run(ctx.Done())

```