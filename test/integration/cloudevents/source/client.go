package source

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cloudeventstypes "github.com/cloudevents/sdk-go/v2/types"
	"k8s.io/klog/v2"

	"open-cluster-management.io/api/cloudevents/generic"
	"open-cluster-management.io/api/cloudevents/generic/options/mqtt"
	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/payload"
	workv1 "open-cluster-management.io/api/work/v1"
)

type resourceCodec struct{}

var _ generic.Codec[*Resource] = &resourceCodec{}

func (c *resourceCodec) EventDataType() types.CloudEventsDataType {
	return payload.ManifestEventDataType
}

func (c *resourceCodec) Encode(source string, eventType types.CloudEventsType, resource *Resource) (*cloudevents.Event, error) {
	if eventType.CloudEventsDataType != payload.ManifestEventDataType {
		return nil, fmt.Errorf("unsupported cloudevents data type %s", eventType.CloudEventsDataType)
	}

	eventBuilder := types.NewEventBuilder(source, eventType).
		WithResourceID(resource.ResourceID).
		WithResourceVersion(resource.ResourceVersion).
		WithClusterName(resource.Namespace)

	if !resource.GetDeletionTimestamp().IsZero() {
		evt := eventBuilder.WithDeletionTimestamp(resource.GetDeletionTimestamp().Time).NewEvent()
		return &evt, nil
	}

	evt := eventBuilder.NewEvent()

	if err := evt.SetData(cloudevents.ApplicationJSON, &payload.Manifest{Manifest: resource.Spec}); err != nil {
		return nil, fmt.Errorf("failed to encode manifests to cloud event: %v", err)
	}

	return &evt, nil
}

func (c *resourceCodec) Decode(evt *cloudevents.Event) (*Resource, error) {
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

	resourceVersion, err := cloudeventstypes.ToString(evtExtensions[types.ExtensionResourceVersion])
	if err != nil {
		return nil, fmt.Errorf("failed to get resourceversion extension: %v", err)
	}

	resourceVersionInt, err := strconv.ParseInt(resourceVersion, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert resourceversion - %v to int64", resourceVersion)
	}

	clusterName, err := cloudeventstypes.ToString(evtExtensions[types.ExtensionClusterName])
	if err != nil {
		return nil, fmt.Errorf("failed to get clustername extension: %v", err)
	}

	manifestStatus := &payload.ManifestStatus{}
	if err := evt.DataAs(manifestStatus); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data %s, %v", string(evt.Data()), err)
	}

	resource := &Resource{
		ResourceID:      resourceID,
		ResourceVersion: resourceVersionInt,
		Namespace:       clusterName,
		Status: ResourceStatus{
			Conditions: manifestStatus.Conditions,
		},
	}

	return resource, nil
}

type resourceLister struct{}

var _ generic.Lister[*Resource] = &resourceLister{}

func (resLister *resourceLister) List(listOpts types.ListOptions) ([]*Resource, error) {
	return GetStore().List(listOpts.ClusterName), nil
}

func StartResourceSourceClient(ctx context.Context, config *mqtt.MQTTOptions) (generic.CloudEventsClient[*Resource], error) {
	client, err := generic.NewCloudEventSourceClient[*Resource](
		ctx,
		mqtt.NewSourceOptions(config, "integration-test"),
		&resourceLister{},
		func(obj *Resource) (string, error) {
			statusBytes, err := json.Marshal(&workv1.ManifestWorkStatus{Conditions: obj.Status.Conditions})
			if err != nil {
				return "", fmt.Errorf("failed to marshal resource status, %v", err)
			}
			return fmt.Sprintf("%x", sha256.Sum256(statusBytes)), nil
		},
		&resourceCodec{},
	)

	if err != nil {
		return nil, err
	}

	go func() {
		if err := client.Subscribe(ctx, func(action types.ResourceAction, resource *Resource) error {
			return GetStore().UpdateStatus(resource)
		}); err != nil {
			klog.Fatalf("failed to subscribe to mqtt broker, %v", err)
		}
	}()

	return client, nil
}
