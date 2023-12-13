package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"open-cluster-management.io/api/cloudevents/generic/options/grpc/protocol"
)

const (
	// SpecTopic is a pubsub topic for resource spec.
	SpecTopic = "sources/+/clusters/+/spec"

	// StatusTopic is a pubsub topic for resource status.
	StatusTopic = "sources/+/clusters/+/status"

	// SpecResyncTopic is a pubsub topic for resource spec resync.
	SpecResyncTopic = "sources/clusters/+/specresync"

	// StatusResyncTopic is a pubsub topic for resource status resync.
	StatusResyncTopic = "sources/+/clusters/statusresync"
)

type GRPCOptions struct {
	Host           string
	Port           int
	CAFile         string
	ClientCertFile string
	ClientKeyFile  string
}

func NewGRPCOptions() *GRPCOptions {
	return &GRPCOptions{}
}

func (o *GRPCOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.Host, "grpc-host", o.Host, "The host of grpc server")
	flags.IntVar(&o.Port, "grpc-port", o.Port, "The port of grpc server")
	flags.StringVar(&o.CAFile, "server-ca", o.CAFile, "A file containing trusted CA certificates for server")
	flags.StringVar(&o.ClientCertFile, "client-certificate", o.ClientCertFile, "The grpc client certificate file")
	flags.StringVar(&o.ClientKeyFile, "client-key", o.ClientKeyFile, "The grpc client private key file")
}

func (o *GRPCOptions) GetGRPCClientConn() (*grpc.ClientConn, error) {
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

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{clientCerts},
			RootCAs:      certPool,
			MinVersion:   tls.VersionTLS13,
			MaxVersion:   tls.VersionTLS13,
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", o.Host, o.Port), grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to grpc server %s:%d, %v", o.Host, o.Port, err)
		}

		return conn, nil
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", o.Host, o.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server %s:%d, %v", o.Host, o.Port, err)
	}

	return conn, nil
}

func (o *GRPCOptions) GetCloudEventsClient(ctx context.Context, clientOpts ...protocol.Option) (cloudevents.Client, error) {
	conn, err := o.GetGRPCClientConn()
	if err != nil {
		return nil, err
	}

	opts := []protocol.Option{}
	opts = append(opts, clientOpts...)
	p, err := protocol.NewProtocol(conn, opts...)
	if err != nil {
		return nil, err
	}

	return cloudevents.NewClient(p)
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
