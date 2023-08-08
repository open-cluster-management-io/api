package cloudevents

import (
	"context"
	"fmt"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	apitypes "k8s.io/apimachinery/pkg/types"

	"open-cluster-management.io/api/client/cloudevents/options/fake"
	"open-cluster-management.io/api/client/cloudevents/payload"
	"open-cluster-management.io/api/client/cloudevents/types"
)

const testSourceName = "mock-source"

func TestSourceResync(t *testing.T) {
	cases := []struct {
		name          string
		resources     []*mockResource
		eventType     types.CloudEventsType
		expectedItems int
	}{
		{
			name:          "no cached resources",
			resources:     []*mockResource{},
			eventType:     types.CloudEventsType{SubResource: types.SubResourceStatus},
			expectedItems: 0,
		},
		{
			name: "has cached resources",
			resources: []*mockResource{
				{UID: apitypes.UID("test1"), ResourceVersion: "2", Status: "test1"},
				{UID: apitypes.UID("test2"), ResourceVersion: "3", Status: "test2"},
			},
			eventType:     types.CloudEventsType{SubResource: types.SubResourceStatus},
			expectedItems: 2,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakeClient := fake.NewCloudEventsFakeClient()
			sourceOptions := fake.NewSourceOptions(fakeClient, testSourceName)
			lister := newMockResourceLister(c.resources...)
			source, err := NewCloudEventSourceClient[*mockResource](context.TODO(), sourceOptions, lister, statusHash, newMockResourceCodec())
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if err := source.Resync(context.TODO(), mockEventDataType); err != nil {
				t.Errorf("unexpected error %v", err)
			}

			evt := fakeClient.GetSentEvents()[0]

			resourceList, err := payload.DecodeStatusResyncRequest(evt)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if len(resourceList.Hashes) != c.expectedItems {
				t.Errorf("expected %d, but got %v", c.expectedItems, resourceList)
			}
		})
	}
}

func TestSourcePublish(t *testing.T) {
	cases := []struct {
		name      string
		resources *mockResource
		eventType types.CloudEventsType
	}{
		{
			name: "publish specs",
			resources: &mockResource{
				UID:             apitypes.UID("1234"),
				ResourceVersion: "2",
				Spec:            "test-spec",
			},
			eventType: types.CloudEventsType{
				CloudEventsDataType: mockEventDataType,
				SubResource:         types.SubResourceSpec,
				Action:              "test_create_request",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakeClient := fake.NewCloudEventsFakeClient()
			sourceOptions := fake.NewSourceOptions(fakeClient, testSourceName)
			lister := newMockResourceLister()
			source, err := NewCloudEventSourceClient[*mockResource](context.TODO(), sourceOptions, lister, statusHash, newMockResourceCodec())
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if err := source.Publish(context.TODO(), c.eventType, c.resources); err != nil {
				t.Errorf("unexpected error %v", err)
			}

			evt := fakeClient.GetSentEvents()[0]
			resourceID, err := evt.Context.GetExtension("resourceid")
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if c.resources.UID != apitypes.UID(fmt.Sprintf("%s", resourceID)) {
				t.Errorf("expected %s, but got %v", c.resources.UID, evt.Context)
			}

			resourceVersion, err := evt.Context.GetExtension("resourceversion")
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if c.resources.ResourceVersion != resourceVersion {
				t.Errorf("expected %s, but got %v", c.resources.ResourceVersion, evt.Context)
			}
		})
	}
}

func TestSpecResyncResponse(t *testing.T) {
	cases := []struct {
		name         string
		requestEvent cloudevents.Event
		resources    []*mockResource
		validate     func([]cloudevents.Event)
	}{
		{
			name: "unsupported event type",
			requestEvent: func() cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetType("unsupported")
				return evt
			}(),
			validate: func(pubEvents []cloudevents.Event) {
				if len(pubEvents) != 0 {
					t.Errorf("unexpected publish events %v", pubEvents)
				}
			},
		},
		{
			name: "unsupported resync event type",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceStatus,
					Action:              types.ResyncRequestAction,
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				return evt
			}(),
			validate: func(pubEvents []cloudevents.Event) {
				if len(pubEvents) != 0 {
					t.Errorf("unexpected publish events %v", pubEvents)
				}
			},
		},
		{
			name: "resync all specs",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              types.ResyncRequestAction,
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				evt.SetExtension("clustername", "cluster1")
				if err := evt.SetData(cloudevents.ApplicationJSON, &payload.ResourceVersionList{}); err != nil {
					t.Fatal(err)
				}
				return evt
			}(),
			resources: []*mockResource{
				{UID: apitypes.UID("test1"), ResourceVersion: "2", Spec: "test1"},
				{UID: apitypes.UID("test2"), ResourceVersion: "3", Spec: "test2"},
			},
			validate: func(pubEvents []cloudevents.Event) {
				if len(pubEvents) != 2 {
					t.Errorf("expected all publish events, but got %v", pubEvents)
				}
			},
		},
		{
			name: "resync specs",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              types.ResyncRequestAction,
				}

				versions := &payload.ResourceVersionList{
					Versions: []payload.ResourceVersion{
						{ResourceID: "test1", ResourceVersion: 1},
						{ResourceID: "test2", ResourceVersion: 2},
					},
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				evt.SetExtension("clustername", "cluster1")
				if err := evt.SetData(cloudevents.ApplicationJSON, versions); err != nil {
					t.Fatal(err)
				}
				return evt
			}(),
			resources: []*mockResource{
				{UID: apitypes.UID("test1"), ResourceVersion: "2", Spec: "test1-updated"},
				{UID: apitypes.UID("test2"), ResourceVersion: "2", Spec: "test2"},
			},
			validate: func(pubEvents []cloudevents.Event) {
				if len(pubEvents) != 1 {
					t.Errorf("expected one publish events, but got %v", pubEvents)
				}
			},
		},
		{
			name: "resync specs - deletion",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              types.ResyncRequestAction,
				}

				versions := &payload.ResourceVersionList{
					Versions: []payload.ResourceVersion{
						{ResourceID: "test1", ResourceVersion: 1},
						{ResourceID: "test2", ResourceVersion: 2},
					},
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				evt.SetExtension("clustername", "cluster1")
				if err := evt.SetData(cloudevents.ApplicationJSON, versions); err != nil {
					t.Fatal(err)
				}
				return evt
			}(),
			resources: []*mockResource{
				{UID: apitypes.UID("test1"), ResourceVersion: "1", Spec: "test1"},
			},
			validate: func(pubEvents []cloudevents.Event) {
				if len(pubEvents) != 1 {
					t.Errorf("expected one publish events, but got %v", pubEvents)
				}

				if _, err := pubEvents[0].Context.GetExtension("deletiontimestamp"); err != nil {
					t.Errorf("expected deletion events, but got %v", pubEvents)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakeClient := fake.NewCloudEventsFakeClient(c.requestEvent)
			sourceOptions := fake.NewSourceOptions(fakeClient, testSourceName)
			lister := newMockResourceLister(c.resources...)
			source, err := NewCloudEventSourceClient[*mockResource](context.TODO(), sourceOptions, lister, statusHash, newMockResourceCodec())
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			if err := source.Subscribe(context.TODO()); err != nil {
				t.Errorf("unexpected error %v", err)
			}

			c.validate(fakeClient.GetSentEvents())
		})
	}
}

func TestReceiveResourceStatus(t *testing.T) {
	cases := []struct {
		name         string
		requestEvent cloudevents.Event
		resources    []*mockResource
		validate     func(event types.ResourceAction, resource *mockResource)
	}{
		{
			name: "unsupported sub resource",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceSpec,
					Action:              "test_create_request",
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				return evt
			}(),
			validate: func(event types.ResourceAction, resource *mockResource) {
				if len(event) != 0 {
					t.Errorf("should not be invoked")
				}
			},
		},
		{
			name: "no registered codec for the resource",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					SubResource: types.SubResourceSpec,
					Action:      "test_create_request",
				}

				evt := cloudevents.NewEvent()
				evt.SetType(eventType.String())
				return evt
			}(),
			validate: func(event types.ResourceAction, resource *mockResource) {
				if len(event) != 0 {
					t.Errorf("should not be invoked")
				}
			},
		},
		{
			name: "update status",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceStatus,
					Action:              "test_update_request",
				}

				evt, _ := newMockResourceCodec().Encode(testAgentName, eventType, &mockResource{UID: apitypes.UID("test1"), ResourceVersion: "1", Status: "update-test1"})
				evt.SetExtension("clustername", "cluster1")
				return *evt
			}(),
			resources: []*mockResource{
				{UID: apitypes.UID("test1"), ResourceVersion: "1", Status: "test1"},
				{UID: apitypes.UID("test2"), ResourceVersion: "1", Status: "test2"},
			},
			validate: func(event types.ResourceAction, resource *mockResource) {
				if event != types.StatusModified {
					t.Errorf("expected modified, but get %s", event)
				}
				if resource.UID != "test1" {
					t.Errorf("unexpected resource %v", resource)
				}
			},
		},
		{
			name: "status no change",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceStatus,
					Action:              "test_update_request",
				}

				evt, _ := newMockResourceCodec().Encode(testAgentName, eventType, &mockResource{UID: apitypes.UID("test1"), ResourceVersion: "1", Status: "test1"})
				evt.SetExtension("clustername", "cluster1")
				return *evt
			}(),
			resources: []*mockResource{
				{UID: apitypes.UID("test1"), ResourceVersion: "1", Status: "test1"},
				{UID: apitypes.UID("test2"), ResourceVersion: "1", Status: "test2"},
			},
			validate: func(event types.ResourceAction, resource *mockResource) {
				if len(event) != 0 {
					t.Errorf("unexpected event %s, %v", event, resource)
				}
			},
		},
		{
			name: "none existing resource",
			requestEvent: func() cloudevents.Event {
				eventType := types.CloudEventsType{
					CloudEventsDataType: mockEventDataType,
					SubResource:         types.SubResourceStatus,
					Action:              "test_update_request",
				}

				evt, _ := newMockResourceCodec().Encode(testAgentName, eventType, &mockResource{UID: apitypes.UID("test1"), ResourceVersion: "1", Status: "test1"})
				evt.SetExtension("clustername", "cluster1")
				return *evt
			}(),
			validate: func(event types.ResourceAction, resource *mockResource) {
				if len(event) != 0 {
					t.Errorf("unexpected event %s, %v", event, resource)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fakeClient := fake.NewCloudEventsFakeClient(c.requestEvent)
			sourceOptions := fake.NewSourceOptions(fakeClient, testSourceName)
			lister := newMockResourceLister(c.resources...)
			source, err := NewCloudEventSourceClient[*mockResource](context.TODO(), sourceOptions, lister, statusHash, newMockResourceCodec())
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}

			var actualEvent types.ResourceAction
			var actualRes *mockResource
			if err := source.Subscribe(context.TODO(), func(event types.ResourceAction, resource *mockResource) error {
				actualEvent = event
				actualRes = resource
				return nil
			}); err != nil {
				t.Errorf("unexpected error %v", err)
			}

			c.validate(actualEvent, actualRes)
		})
	}
}
