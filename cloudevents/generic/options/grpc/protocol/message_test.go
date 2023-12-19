package protocol

import (
	"context"
	"testing"

	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/event"

	pbv1 "open-cluster-management.io/api/cloudevents/generic/options/grpc/protobuf/v1"
)

func TestReadStructured(t *testing.T) {
	tests := []struct {
		name    string
		msg     *pbv1.CloudEvent
		wantErr error
	}{
		{
			name:    "nil format",
			msg:     &pbv1.CloudEvent{},
			wantErr: binding.ErrNotStructured,
		},
		{
			name: "json format",
			msg: &pbv1.CloudEvent{
				Attributes: map[string]*pbv1.CloudEventAttributeValue{
					contenttype: {
						Attr: &pbv1.CloudEventAttributeValue_CeString{
							CeString: event.ApplicationCloudEventsJSON,
						},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := NewMessage(tc.msg)
			err := msg.ReadStructured(context.Background(), (*pbEventWriter)(tc.msg))
			if err != tc.wantErr {
				t.Errorf("Error unexpected. got: %v, want: %v", err, tc.wantErr)
			}
		})
	}
}

func TestReadBinary(t *testing.T) {
	msg := &pbv1.CloudEvent{
		SpecVersion: "1.0",
		Id:          "ABC-123",
		Source:      "test-source",
		Type:        "binary.test",
		Attributes:  map[string]*pbv1.CloudEventAttributeValue{},
		Data: &pbv1.CloudEvent_BinaryData{
			BinaryData: []byte("{hello:world}"),
		},
	}

	message := NewMessage(msg)
	err := message.ReadBinary(context.Background(), (*pbEventWriter)(msg))
	if err != nil {
		t.Errorf("Error unexpected. got: %v", err)
	}
}
