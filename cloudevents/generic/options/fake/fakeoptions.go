package fake

import (
	"context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/protocol"

	"open-cluster-management.io/api/cloudevents/generic/options"
)

type CloudEventsFakeOptions struct {
	client *CloudEventsFakeClient
}

func NewAgentOptions(client *CloudEventsFakeClient, clusterName, agentID string) *options.CloudEventsAgentOptions {
	return &options.CloudEventsAgentOptions{
		CloudEventsOptions: &CloudEventsFakeOptions{client: client},
		AgentID:            agentID,
		ClusterName:        clusterName,
	}
}

func NewSourceOptions(client *CloudEventsFakeClient, sourceID string) *options.CloudEventsSourceOptions {
	return &options.CloudEventsSourceOptions{
		CloudEventsOptions: &CloudEventsFakeOptions{client: client},
		SourceID:           sourceID,
	}
}

func (o *CloudEventsFakeOptions) WithContext(ctx context.Context, evtCtx cloudevents.EventContext) (context.Context, error) {
	return ctx, nil
}

func (o *CloudEventsFakeOptions) Client(ctx context.Context) (cloudevents.Client, error) {
	return o.client, nil
}

type CloudEventsFakeClient struct {
	sentEvents     []cloudevents.Event
	receivedEvents []cloudevents.Event
}

func NewCloudEventsFakeClient(receivedEvents ...cloudevents.Event) *CloudEventsFakeClient {
	return &CloudEventsFakeClient{
		sentEvents:     []cloudevents.Event{},
		receivedEvents: receivedEvents,
	}
}

func (c *CloudEventsFakeClient) Send(ctx context.Context, event cloudevents.Event) protocol.Result {
	c.sentEvents = append(c.sentEvents, event)
	return nil
}

func (c *CloudEventsFakeClient) Request(ctx context.Context, event event.Event) (*cloudevents.Event, protocol.Result) {
	return nil, nil
}

func (c *CloudEventsFakeClient) StartReceiver(ctx context.Context, fn interface{}) error {
	receiver, ok := fn.(func(evt cloudevents.Event))
	if !ok {
		return fmt.Errorf("unsupported receiver %T", fn)
	}

	for _, evt := range c.receivedEvents {
		receiver(evt)
	}
	return nil
}

func (c *CloudEventsFakeClient) GetSentEvents() []cloudevents.Event {
	return c.sentEvents
}
