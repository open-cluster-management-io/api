package agent

import (
	"context"

	"open-cluster-management.io/api/cloudevents/generic/options/mqtt"
	"open-cluster-management.io/api/cloudevents/work"
	"open-cluster-management.io/api/cloudevents/work/agent/codec"
)

func StartWorkAgent(ctx context.Context, clusterName string, config *mqtt.MQTTOptions) (*work.ClientHolder, error) {
	clientHolder, err := work.NewClientHolderBuilder(clusterName, config).
		WithClusterName(clusterName).
		WithCodecs(codec.NewManifestCodec(nil)).
		NewClientHolder(ctx)
	if err != nil {
		return nil, err
	}

	go clientHolder.ManifestWorkInformer().Informer().Run(ctx.Done())

	return clientHolder, nil
}
