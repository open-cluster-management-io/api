package source

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	cloudeventstypes "github.com/cloudevents/sdk-go/v2/types"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pbv1 "open-cluster-management.io/api/cloudevents/generic/options/grpc/protobuf/v1"
	grpcprotocol "open-cluster-management.io/api/cloudevents/generic/options/grpc/protocol"
	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/payload"
)

type CloudEventServer struct {
	pbv1.UnimplementedCloudEventServiceServer
	store    *MemoryStore
	eventHub *EventHub
}

func NewCloudEventServer(store *MemoryStore, eventHub *EventHub) *CloudEventServer {
	return &CloudEventServer{
		store:    store,
		eventHub: eventHub,
	}
}

func (svr *CloudEventServer) Publish(ctx context.Context, pubReq *pbv1.PublishRequest) (*emptypb.Empty, error) {
	// pbEvt, err := pb.ToProto(evt)
	evt, err := binding.ToEvent(ctx, grpcprotocol.NewMessage(pubReq.Event))
	if err != nil {
		return nil, fmt.Errorf("failed to convert protobuf to cloudevent: %v", err)
	}

	res, err := decode(evt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cloudevent: %v", err)
	}

	store.UpSert(res)
	return &emptypb.Empty{}, nil
}

func (svr *CloudEventServer) Subscribe(subReq *pbv1.SubscriptionRequest, subServer pbv1.CloudEventService_SubscribeServer) error {
	topicSplits := strings.Split(subReq.Topic, "/")
	if len(topicSplits) != 5 {
		return fmt.Errorf("invalid topic %s", subReq.Topic)
	}

	clusterName := topicSplits[3]
	eventClient := NewEventClient(clusterName)
	svr.eventHub.Register(eventClient)
	defer svr.eventHub.Unregister(eventClient)

	for res := range eventClient.Receive() {
		evt, err := encode(res)
		if err != nil {
			return fmt.Errorf("failed to encode resource %s to cloudevent: %v", res.ResourceID, err)
		}

		// pbEvt, err := pb.ToProto(evt)
		pbEvt := &pbv1.CloudEvent{}
		if err = grpcprotocol.WritePBMessage(context.TODO(), binding.ToMessage(evt), pbEvt); err != nil {
			return fmt.Errorf("failed to convert cloudevent to protobuf: %v", err)
		}
		if err := subServer.Send(pbEvt); err != nil {
			return err
		}
	}

	return nil
}

func (svr *CloudEventServer) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return err
	}
	grpcServer := grpc.NewServer()
	pbv1.RegisterCloudEventServiceServer(grpcServer, svr)
	return grpcServer.Serve(lis)
}

func encode(resource *Resource) (*cloudevents.Event, error) {
	source := "test-source"
	eventType := types.CloudEventsType{
		CloudEventsDataType: payload.ManifestEventDataType,
		SubResource:         types.SubResourceStatus,
		Action:              "status_update",
	}

	eventBuilder := types.NewEventBuilder(source, eventType).
		WithResourceID(resource.ResourceID).
		WithResourceVersion(resource.ResourceVersion).
		WithClusterName(resource.Namespace)

	evt := eventBuilder.NewEvent()

	if err := evt.SetData(cloudevents.ApplicationJSON, &payload.ManifestStatus{Conditions: resource.Status.Conditions}); err != nil {
		return nil, fmt.Errorf("failed to encode manifest status to cloud event: %v", err)
	}

	return &evt, nil
}

func decode(evt *cloudevents.Event) (*Resource, error) {
	eventType, err := types.ParseCloudEventsType(evt.Type())
	if err != nil {
		return nil, fmt.Errorf("failed to parse cloud event type %s, %v", evt.Type(), err)
	}

	if eventType.CloudEventsDataType != payload.ManifestEventDataType {
		return nil, fmt.Errorf("unsupported cloudevents data type %s", eventType.CloudEventsDataType)
	}

	evtExtensions := evt.Context.GetExtensions()

	resourceID, err := cloudeventstypes.ToString(evtExtensions[types.ExtensionResourceID])
	if err != nil {
		return nil, fmt.Errorf("failed to get resourceid extension: %v", err)
	}

	resourceVersion, err := cloudeventstypes.ToInteger(evtExtensions[types.ExtensionResourceVersion])
	if err != nil {
		return nil, fmt.Errorf("failed to get resourceversion extension: %v", err)
	}

	clusterName, err := cloudeventstypes.ToString(evtExtensions[types.ExtensionClusterName])
	if err != nil {
		return nil, fmt.Errorf("failed to get clustername extension: %v", err)
	}

	manifest := &payload.Manifest{}
	if err := evt.DataAs(manifest); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data %s, %v", string(evt.Data()), err)
	}

	resource := &Resource{
		ResourceID:      resourceID,
		ResourceVersion: int64(resourceVersion),
		Namespace:       clusterName,
		Spec:            manifest.Manifest,
	}

	return resource, nil
}
