package codec

import (
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"open-cluster-management.io/api/cloudevents/generic/types"
	"open-cluster-management.io/api/cloudevents/work/payload"
	workv1 "open-cluster-management.io/api/work/v1"
)

func TestManifestBundleEventDataType(t *testing.T) {
	codec := NewManifestBundleCodec()

	if codec.EventDataType() != payload.ManifestBundleEventDataType {
		t.Errorf("unexpected event data type %s", codec.EventDataType())
	}
}

func TestManifestBundleEncode(t *testing.T) {
	cases := []struct {
		name        string
		eventType   types.CloudEventsType
		work        *workv1.ManifestWork
		expectedErr bool
	}{
		{
			name: "unsupported cloudevents data type",
			eventType: types.CloudEventsType{
				CloudEventsDataType: types.CloudEventsDataType{
					Group:    "test",
					Version:  "v1",
					Resource: "test",
				},
			},
			expectedErr: true,
		},
		{
			name: "bad resourceversion",
			eventType: types.CloudEventsType{
				CloudEventsDataType: payload.ManifestBundleEventDataType,
				SubResource:         types.SubResourceStatus,
				Action:              "test",
			},
			work: &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion: "abc",
				},
			},
			expectedErr: true,
		},
		{
			name: "no originalsource",
			eventType: types.CloudEventsType{
				CloudEventsDataType: payload.ManifestBundleEventDataType,
				SubResource:         types.SubResourceStatus,
				Action:              "test",
			},
			work: &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					ResourceVersion: "13",
				},
			},
			expectedErr: true,
		},
		{
			name: "encode a manifestwork status",
			eventType: types.CloudEventsType{
				CloudEventsDataType: payload.ManifestBundleEventDataType,
				SubResource:         types.SubResourceStatus,
				Action:              "test",
			},
			work: &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					UID:             "test",
					ResourceVersion: "13",
					Annotations: map[string]string{
						"cloudevents.open-cluster-management.io/originalsource": "source1",
					},
				},
				Status: workv1.ManifestWorkStatus{
					Conditions:     []metav1.Condition{},
					ResourceStatus: workv1.ManifestResourceStatus{},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			codec := NewManifestBundleCodec()

			_, err := codec.Encode("cluster1-work-agent", c.eventType, c.work)
			if c.expectedErr {
				if err == nil {
					t.Errorf("expected an error, but failed")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}

func TestManifestBundleDecode(t *testing.T) {
	cases := []struct {
		name        string
		event       *cloudevents.Event
		expectedErr bool
	}{
		{
			name: "bad cloudevents type",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetType("test")
				return &evt
			}(),
			expectedErr: true,
		},
		{
			name: "unsupported cloudevents data type",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetType("test-group.v1.test.spec.test")
				return &evt
			}(),
			expectedErr: true,
		},
		{
			name: "no resourceid",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetType("io.open-cluster-management.works.v1alpha1.manifestbundles.spec.test")
				return &evt
			}(),
			expectedErr: true,
		},
		{
			name: "no resourceversion",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetType("io.open-cluster-management.works.v1alpha1.manifestbundles.spec.test")
				evt.SetExtension("resourceid", "test")
				return &evt
			}(),
			expectedErr: true,
		},
		{
			name: "no clustername",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetType("io.open-cluster-management.works.v1alpha1.manifestbundles.spec.test")
				evt.SetExtension("resourceid", "test")
				evt.SetExtension("resourceversion", "13")
				return &evt
			}(),
			expectedErr: true,
		},
		{
			name: "bad data",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetSource("source1")
				evt.SetType("io.open-cluster-management.works.v1alpha1.manifestbundles.spec.test")
				evt.SetExtension("resourceid", "test")
				evt.SetExtension("resourceversion", "13")
				return &evt
			}(),
			expectedErr: true,
		},
		{
			name: "has deletion time",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetSource("source1")
				evt.SetType("io.open-cluster-management.works.v1alpha1.manifestbundles.spec.test")
				evt.SetExtension("resourceid", "test")
				evt.SetExtension("resourceversion", "13")
				evt.SetExtension("clustername", "cluster1")
				evt.SetExtension("deletiontimestamp", "1985-04-12T23:20:50.52Z")
				return &evt
			}(),
		},
		{
			name: "decode an invalid cloudevent",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetSource("source1")
				evt.SetType("io.open-cluster-management.works.v1alpha1.manifestbundles.spec.test")
				evt.SetExtension("resourceid", "test")
				evt.SetExtension("resourceversion", "13")
				evt.SetExtension("clustername", "cluster1")
				if err := evt.SetData(cloudevents.ApplicationJSON, &payload.ManifestBundle{}); err != nil {
					t.Fatal(err)
				}
				return &evt
			}(),
			expectedErr: true,
		},
		{
			name: "decode a cloudevent",
			event: func() *cloudevents.Event {
				evt := cloudevents.NewEvent()
				evt.SetSource("source1")
				evt.SetType("io.open-cluster-management.works.v1alpha1.manifestbundles.spec.test")
				evt.SetExtension("resourceid", "test")
				evt.SetExtension("resourceversion", "13")
				evt.SetExtension("clustername", "cluster1")
				if err := evt.SetData(cloudevents.ApplicationJSON, &payload.ManifestBundle{
					Manifests: []workv1.Manifest{
						{
							RawExtension: runtime.RawExtension{
								Raw: toConfigMap(t),
							},
						},
					},
				}); err != nil {
					t.Fatal(err)
				}
				return &evt
			}(),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			codec := NewManifestBundleCodec()

			_, err := codec.Decode(c.event)
			if c.expectedErr {
				if err == nil {
					t.Errorf("expected an error, but failed")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
		})
	}
}
