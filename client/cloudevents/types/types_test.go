package types

import (
	"fmt"
	"testing"

	"k8s.io/apimachinery/pkg/api/equality"
)

const testManifestsType = "io.open-cluster-management.works.v1alpha1.manifests.spec.create_request"

func TestToString(t *testing.T) {
	dataType := CloudEventsDataType{
		Group:    "io.open-cluster-management.works",
		Version:  "v1alpha1",
		Resource: "manifests",
	}

	eventType := CloudEventsType{
		CloudEventsDataType: dataType,
		SubResource:         "spec",
		Action:              "create_request",
	}

	if eventType.String() != testManifestsType {
		t.Errorf("expected %s, but get %s", testManifestsType, eventType)
	}
}

func TestParseCloudEventType(t *testing.T) {
	cases := []struct {
		name         string
		eventType    string
		expectedType *CloudEventsType
		err          error
	}{
		{
			name:      "manifest creation event",
			eventType: testManifestsType,
			expectedType: &CloudEventsType{
				CloudEventsDataType: CloudEventsDataType{
					Group:    "io.open-cluster-management.works",
					Version:  "v1alpha1",
					Resource: "manifests",
				},
				SubResource: "spec",
				Action:      "create_request",
			},
		},
		{
			name:      "wrong format",
			eventType: "test",
			err:       fmt.Errorf("unsupported cloud event type format"),
		},
		{
			name:      "unsupported subresource",
			eventType: "io.open-cluster-management.works.v1alpha1.manifests.unsupported.create_request",
			err:       fmt.Errorf("unsupported subresource unsupported"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			eventType, err := Parse(c.eventType)
			if err != nil {
				if err.Error() != c.err.Error() {
					t.Errorf("unexpected error %v", err)
				}
			}

			if !equality.Semantic.DeepEqual(eventType, c.expectedType) {
				t.Errorf("unexpected event type %v", eventType)
			}
		})
	}
}
