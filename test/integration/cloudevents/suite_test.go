package cloudevents

import (
	"context"
	"testing"

	mochimqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"

	"open-cluster-management.io/api/cloudevents/generic"
	"open-cluster-management.io/api/cloudevents/generic/options/mqtt"
	"open-cluster-management.io/api/test/integration/cloudevents/source"
)

const mqttBrokerHost = "127.0.0.1:1883"

var mqttBroker *mochimqtt.Server
var mqttOptions *mqtt.MQTTOptions
var sourceCloudEventsClient generic.CloudEventsClient[*source.Resource]

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "CloudEvents Client Integration Suite")
}

var _ = ginkgo.BeforeSuite(func(done ginkgo.Done) {
	ginkgo.By("bootstrapping test environment")

	// start a MQTT broker
	mqttBroker = mochimqtt.New(&mochimqtt.Options{})
	l := mqttBroker.Log.Level(zerolog.DebugLevel)
	mqttBroker.Log = &l
	// allow all connections.
	err := mqttBroker.AddHook(new(auth.AllowHook), nil)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	err = mqttBroker.AddListener(listeners.NewTCP("mqtt-test-broker", mqttBrokerHost, nil))
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	go func() {
		err := mqttBroker.Serve()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	}()

	mqttOptions = mqtt.NewMQTTOptions()
	mqttOptions.BrokerHost = mqttBrokerHost
	ginkgo.By("start an source")
	sourceCloudEventsClient, err = source.StartResourceSourceClient(context.TODO(), mqttOptions)
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	close(done)
}, 300)

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("tearing down the test environment")

	err := mqttBroker.Close()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())
})
