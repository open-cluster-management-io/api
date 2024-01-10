package codec

import (
	"encoding/json"
	"fmt"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cloudeventstypes "github.com/cloudevents/sdk-go/v2/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"

	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/payload"
	workv1 "open-cluster-management.io/api/work/v1"
)

const (
	// CloudEventsDataTypeAnnotationKey is the key of the cloudevents data type annotation.
	CloudEventsDataTypeAnnotationKey = "cloudevents.open-cluster-management.io/datatype"

	// CloudEventsDataTypeAnnotationKey is the key of the cloudevents original source annotation.
	CloudEventsOriginalSourceAnnotationKey = "cloudevents.open-cluster-management.io/originalsource"
)

// ManifestBundleCodec is a codec to encode/decode a ManifestWork/cloudevent with ManifestBundle for a source.
type ManifestBundleCodec struct{}

func NewManifestBundleCodec() *ManifestBundleCodec {
	return &ManifestBundleCodec{}
}

// EventDataType always returns the event data type `io.open-cluster-management.works.v1alpha1.manifestbundles`.
func (c *ManifestBundleCodec) EventDataType() types.CloudEventsDataType {
	return payload.ManifestBundleEventDataType
}

// Encode the spec of a ManifestWork to a cloudevent with ManifestBundle.
func (c *ManifestBundleCodec) Encode(source string, eventType types.CloudEventsType, work *workv1.ManifestWork) (*cloudevents.Event, error) {
	if eventType.CloudEventsDataType != payload.ManifestBundleEventDataType {
		return nil, fmt.Errorf("unsupported cloudevents data type %s", eventType.CloudEventsDataType)
	}

	workMeta := &types.ResourceMeta{
		Group:     workv1.GroupName,
		Resource:  "manifestworks",
		Version:   "v1",
		Name:      work.Name,
		Namespace: work.Namespace,
	}

	workMetaJson, err := json.Marshal(workMeta)
	if err != nil {
		return nil, err
	}

	evt := types.NewEventBuilder(source, eventType).
		WithClusterName(work.Namespace).
		NewEvent()
	evt.SetExtension(types.ExtensionResourceMeta, string(workMetaJson))
	if !work.DeletionTimestamp.IsZero() {
		evt.SetExtension(types.ExtensionDeletionTimestamp, work.DeletionTimestamp.Time)
		return &evt, nil
	}

	manifests := &payload.ManifestBundle{
		Manifests:       work.Spec.Workload.Manifests,
		DeleteOption:    work.Spec.DeleteOption,
		ManifestConfigs: work.Spec.ManifestConfigs,
	}
	if err := evt.SetData(cloudevents.ApplicationJSON, manifests); err != nil {
		return nil, fmt.Errorf("failed to encode manifestwork status to a cloudevent: %v", err)
	}

	return &evt, nil
}

// Decode a cloudevent whose data is ManifestBundle to a ManifestWork.
func (c *ManifestBundleCodec) Decode(evt *cloudevents.Event) (*workv1.ManifestWork, error) {
	eventType, err := types.ParseCloudEventsType(evt.Type())
	if err != nil {
		return nil, fmt.Errorf("failed to parse cloud event type %s, %v", evt.Type(), err)
	}

	if eventType.CloudEventsDataType != payload.ManifestBundleEventDataType {
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

	workMetaExtension, err := cloudeventstypes.ToString(evtExtensions[types.ExtensionResourceMeta])
	if err != nil {
		return nil, fmt.Errorf("failed to get resourcemeta extension: %v", err)
	}

	workMeta := &types.ResourceMeta{}
	if err := json.Unmarshal([]byte(workMetaExtension), workMeta); err != nil {
		return nil, err
	}

	work := &workv1.ManifestWork{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			UID:             kubetypes.UID(resourceID),
			ResourceVersion: resourceVersion,
			Generation:      resourceVersionInt,
			Name:            workMeta.Name,
			Namespace:       workMeta.Namespace,
		},
	}

	manifestStatus := &payload.ManifestBundleStatus{}
	if err := evt.DataAs(manifestStatus); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event data %s, %v", string(evt.Data()), err)
	}

	work.Status = workv1.ManifestWorkStatus{
		Conditions: manifestStatus.Conditions,
		ResourceStatus: workv1.ManifestResourceStatus{
			Manifests: manifestStatus.ResourceStatus,
		},
	}

	return work, nil
}
