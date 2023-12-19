package protocol

import (
	"context"
	"testing"

	"github.com/cloudevents/sdk-go/v2/binding"
	bindingtest "github.com/cloudevents/sdk-go/v2/binding/test"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/test"
	"github.com/stretchr/testify/require"

	pbv1 "open-cluster-management.io/api/cloudevents/generic/options/grpc/protobuf/v1"
)

func TestEncodeMessage(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name             string
		messageFactory   func(e event.Event) binding.Message
		expectedEncoding binding.Encoding
	}{
		{
			name: "Structured to Structured",
			messageFactory: func(e event.Event) binding.Message {
				return bindingtest.MustCreateMockStructuredMessage(t, e)
			},
			expectedEncoding: binding.EncodingStructured,
		},
		{
			name:             "Binary to Binary",
			messageFactory:   bindingtest.MustCreateMockBinaryMessage,
			expectedEncoding: binding.EncodingBinary,
		},
	}

	test.EachEvent(t, test.Events(), func(t *testing.T, e event.Event) {
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				eventIn := test.ConvertEventExtensionsToString(t, e.Clone())
				// convert the event to binding.Message with specific encoding
				messageIn := tc.messageFactory(eventIn)

				// load the binding.Message into a protobuf event
				pbEvt := &pbv1.CloudEvent{}
				err := WritePBMessage(ctx, messageIn, pbEvt)
				require.NoError(t, err)

				// convert the protobuf event back to binding.Message
				messageOut := NewMessage(pbEvt)
				require.Equal(t, tc.expectedEncoding, messageOut.ReadEncoding())

				// convert the binding.Message back to event.Event
				eventOut, err := binding.ToEvent(ctx, messageOut)
				require.NoError(t, err)

				// check if the event is the same
				test.AssertEventEquals(t, eventIn, *eventOut)
			})
		}
	})
}
