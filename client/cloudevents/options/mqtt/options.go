package mqtt

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"strings"

	cloudeventsmqtt "github.com/cloudevents/sdk-go/protocol/mqtt_paho/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/eclipse/paho.golang/paho"
	"github.com/spf13/pflag"
)

const (
	// SpecTopic is a MQTT topic for resource spec.
	SpecTopic = "sources/+/clusters/+/spec"

	// StatusTopic is a MQTT topic for resource status.
	StatusTopic = "sources/+/clusters/+/status"

	// SpecResyncTopic is a MQTT topic for resource spec resync.
	SpecResyncTopic = "sources/clusters/+/specresync"

	// StatusResyncTopic is a MQTT topic for resource status resync.
	StatusResyncTopic = "sources/+/clusters/statusresync"
)

type MQTTOptions struct {
	BrokerHost     string
	Username       string
	Password       string
	CAFile         string
	ClientCertFile string
	ClientKeyFile  string
	KeepAlive      uint16
	PubQoS         int
	SubQoS         int
}

func NewMQTTOptions() *MQTTOptions {
	return &MQTTOptions{
		KeepAlive: 60,
		PubQoS:    1,
		SubQoS:    1,
	}
}

func (o *MQTTOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.BrokerHost, "mqtt-broker-host", o.BrokerHost, "The host of MQTT broker")
	flags.StringVar(&o.Username, "mqtt-username", o.Username, "The username to connect the MQTT broker")
	flags.StringVar(&o.Password, "mqtt-password", o.Password, "The password to connect the MQTT broker")
	flags.StringVar(&o.CAFile, "mqtt-broke-ca", o.CAFile, "A file containing trusted CA certificates MQTT broker")
	flags.StringVar(&o.ClientCertFile, "mqtt-client-certificate", o.ClientCertFile, "The MQTT client certificate file")
	flags.StringVar(&o.ClientKeyFile, "mqtt-client-key", o.ClientKeyFile, "The MQTT client private key file")
	flags.Uint16Var(&o.KeepAlive, "mqtt-keep-alive", o.KeepAlive, "Keep alive in seconds for MQTT clients")
	flags.IntVar(&o.SubQoS, "mqtt-sub-qos", o.SubQoS, "The OoS for subscribe")
	flags.IntVar(&o.PubQoS, "mqtt-pub-qos", o.PubQoS, "The Qos for publish")
}

func (o *MQTTOptions) GetNetConn() (net.Conn, error) {
	if len(o.CAFile) != 0 {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}

		caPEM, err := os.ReadFile(o.CAFile)
		if err != nil {
			return nil, err
		}

		if ok := certPool.AppendCertsFromPEM(caPEM); !ok {
			return nil, fmt.Errorf("invalid CA %s", o.CAFile)
		}

		clientCerts, err := tls.LoadX509KeyPair(o.ClientCertFile, o.ClientKeyFile)
		if err != nil {
			return nil, err
		}

		conn, err := tls.Dial("tcp", o.BrokerHost, &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{clientCerts},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to MQTT broker %s, %v", o.BrokerHost, err)
		}

		return conn, nil
	}

	conn, err := net.Dial("tcp", o.BrokerHost)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker %s, %v", o.BrokerHost, err)
	}

	return conn, nil
}

func (o *MQTTOptions) GetMQTTConnectOption(clientID string) *paho.Connect {
	connect := &paho.Connect{
		ClientID:   clientID,
		KeepAlive:  o.KeepAlive,
		CleanStart: true,
	}

	if len(o.Username) != 0 {
		connect.Username = o.Username
		connect.UsernameFlag = true
	}

	if len(o.Password) != 0 {
		connect.Password = []byte(o.Password)
		connect.PasswordFlag = true
	}

	return connect
}

func (o *MQTTOptions) GetCloudEventsClient(
	ctx context.Context,
	clientID string,
	clientOpt cloudeventsmqtt.Option,
) (cloudevents.Client, error) {
	netConn, err := o.GetNetConn()
	if err != nil {
		return nil, err
	}

	config := &paho.ClientConfig{
		ClientID: clientID,
		Conn:     netConn,
	}

	connectOpt := cloudeventsmqtt.WithConnect(o.GetMQTTConnectOption(clientID))
	protocol, err := cloudeventsmqtt.New(ctx, config, connectOpt, clientOpt)
	if err != nil {
		return nil, err
	}

	return cloudevents.NewClient(protocol)
}

// Replace the nth occurrence of old in str by new.
func replaceNth(str, old, new string, n int) string {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(str[i:], old)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return str[:i] + new + str[i+len(old):]
		}
		i += len(old)
	}
	return str
}
