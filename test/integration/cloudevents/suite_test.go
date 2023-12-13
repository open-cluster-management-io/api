package cloudevents

import (
	"context"
	"fmt"
	"testing"

	mochimqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"open-cluster-management.io/api/cloudevents/generic"
	grpcoptions "open-cluster-management.io/api/cloudevents/generic/options/grpc"
	"open-cluster-management.io/api/cloudevents/generic/options/mqtt"
	"open-cluster-management.io/api/test/integration/cloudevents/source"
)

const mqttBrokerHost = "127.0.0.1:1883"
const grpcServerHost = "127.0.0.1"
const grpcServerPort = 8881

var mqttBroker *mochimqtt.Server
var mqttOptions *mqtt.MQTTOptions
var mqttSourceCloudEventsClient generic.CloudEventsClient[*source.Resource]
var cloudEventServer *source.CloudEventServer
var grpcOptions *grpcoptions.GRPCOptions
var grpcSourceCloudEventsClient generic.CloudEventsClient[*source.Resource]
var eventHub *source.EventHub
var store *source.MemoryStore
var consumerStore *source.MemoryStore

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "CloudEvents Client Integration Suite")
}

var _ = ginkgo.BeforeSuite(func(done ginkgo.Done) {
	ginkgo.By("bootstrapping test environment")
	ctx := context.TODO()

	// start a MQTT broker
	mqttBroker = mochimqtt.New(&mochimqtt.Options{})
	// allow all connections.
	err := mqttBroker.AddHook(new(auth.AllowHook), nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	err = mqttBroker.AddListener(listeners.NewTCP("mqtt-test-broker", mqttBrokerHost, nil))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	go func() {
		err := mqttBroker.Serve()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}()

	ginkgo.By("init the event hub")
	eventHub = source.NewEventHub()
	go func() {
		eventHub.Start(ctx)
	}()

	ginkgo.By("init the resource source store")
	store, consumerStore = source.InitStore(eventHub)
	// store.Add(source.NewResource("cluster1", "resource1"))
	// store.Add(source.NewResource("cluster2", "resource1"))

	cloudEventServer = source.NewCloudEventServer(store, eventHub)
	go func() {
		err := cloudEventServer.Start(fmt.Sprintf("%s:%d", grpcServerHost, grpcServerPort))
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}()

	ginkgo.By("start the resource grpc source client")
	grpcOptions = grpcoptions.NewGRPCOptions()
	grpcOptions.Host = grpcServerHost
	grpcOptions.Port = grpcServerPort
	grpcSourceCloudEventsClient, err = source.StartGRPCResourceSourceClient(ctx, grpcoptions.NewSourceOptions(grpcOptions, "integration-grpc-test"))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	ginkgo.By("start the resource mqtt source client")
	mqttOptions = mqtt.NewMQTTOptions()
	mqttOptions.BrokerHost = mqttBrokerHost
	mqttSourceCloudEventsClient, err = source.StartMQTTResourceSourceClient(ctx, mqtt.NewSourceOptions(mqttOptions, "integration-mqtt-test"))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	close(done)
}, 300)

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("tearing down the test environment")

	err := mqttBroker.Close()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
})
