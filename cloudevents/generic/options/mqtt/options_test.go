package mqtt

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestBuildMQTTOptionsFromFlags(t *testing.T) {
	file, err := os.CreateTemp("", "mqtt-config-test-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	cases := []struct {
		name             string
		config           string
		expectedOptions  *MQTTOptions
		expectedErrorMsg string
	}{
		{
			name:             "empty config",
			config:           "",
			expectedErrorMsg: "brokerHost is required",
		},
		{
			name:             "tls config without clientCertFile",
			config:           "{\"brokerHost\":\"test\",\"clientCertFile\":\"test\"}",
			expectedErrorMsg: "either both or none of clientCertFile and clientKeyFile must be set",
		},
		{
			name:             "tls config without caFile",
			config:           "{\"brokerHost\":\"test\",\"clientCertFile\":\"test\",\"clientKeyFile\":\"test\"}",
			expectedErrorMsg: "setting clientCertFile and clientKeyFile requires caFile",
		},
		{
			name:   "default options",
			config: "{\"brokerHost\":\"test\"}",
			expectedOptions: &MQTTOptions{
				BrokerHost: "test",
				KeepAlive:  60,
				PubQoS:     1,
				SubQoS:     1,
			},
		},
		{
			name:   "default options with yaml format",
			config: "brokerHost: test",
			expectedOptions: &MQTTOptions{
				BrokerHost: "test",
				KeepAlive:  60,
				PubQoS:     1,
				SubQoS:     1,
			},
		},
		{
			name:   "customized options",
			config: "{\"brokerHost\":\"test\",\"keepAlive\":30,\"pubQoS\":0,\"subQoS\":2}",
			expectedOptions: &MQTTOptions{
				BrokerHost: "test",
				KeepAlive:  30,
				PubQoS:     0,
				SubQoS:     2,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := os.WriteFile(file.Name(), []byte(c.config), 0644); err != nil {
				t.Fatal(err)
			}

			options, err := BuildMQTTOptionsFromFlags(file.Name())
			if err != nil {
				if err.Error() != c.expectedErrorMsg {
					t.Errorf("unexpected err %v", err)
				}
			}

			if !reflect.DeepEqual(options, c.expectedOptions) {
				t.Errorf("unexpected options %v", options)
			}
		})
	}
}
