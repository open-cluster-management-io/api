package protobuf

import (
	"net/url"
	"reflect"
	"testing"
	stdtime "time"

	"github.com/stretchr/testify/require"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/types"

	pbv1 "open-cluster-management.io/api/cloudevents/generic/options/grpc/protobuf/v1"
)

func TestProtobufFormatWithoutProtobufCodec(t *testing.T) {
	require := require.New(t)
	const test = "test"
	e := event.New()
	e.SetID(test)
	e.SetTime(stdtime.Date(2021, 1, 1, 1, 1, 1, 1, stdtime.UTC))
	e.SetExtension(test, test)
	e.SetExtension("int", 1)
	e.SetExtension("bool", true)
	e.SetExtension("URI", &url.URL{
		Host: "test-uri",
	})
	e.SetExtension("URIRef", types.URIRef{URL: url.URL{
		Host: "test-uriref",
	}})
	e.SetExtension("bytes", []byte(test))
	e.SetExtension("timestamp", stdtime.Date(2021, 2, 1, 1, 1, 1, 1, stdtime.UTC))
	e.SetSubject(test)
	e.SetSource(test)
	e.SetType(test)
	e.SetDataSchema(test)
	require.NoError(e.SetData(event.ApplicationJSON, "foo"))

	b, err := Protobuf.Marshal(&e)
	require.NoError(err)
	var e2 event.Event
	require.NoError(Protobuf.Unmarshal(b, &e2))
	require.Equal(e, e2)
}

func TestProtobufFormatWithProtobufCodec(t *testing.T) {
	require := require.New(t)
	const test = "test"
	e := event.New()
	e.SetID(test)
	e.SetTime(stdtime.Date(2021, 1, 1, 1, 1, 1, 1, stdtime.UTC))
	e.SetExtension(test, test)
	e.SetExtension("int", 1)
	e.SetExtension("bool", true)
	e.SetExtension("URI", &url.URL{
		Host: "test-uri",
	})
	e.SetExtension("URIRef", types.URIRef{URL: url.URL{
		Host: "test-uriref",
	}})
	e.SetExtension("bytes", []byte(test))
	e.SetExtension("timestamp", stdtime.Date(2021, 2, 1, 1, 1, 1, 1, stdtime.UTC))
	e.SetSubject(test)
	e.SetSource(test)
	e.SetType(test)
	e.SetDataSchema(test)

	// Using the CloudEventAttributeValue because it is convenient and is an
	// independent protobuf message. Any protobuf message would work but this
	// one is already generated and included in the source.
	payload := &pbv1.CloudEventAttributeValue{
		Attr: &pbv1.CloudEventAttributeValue_CeBoolean{
			CeBoolean: true,
		},
	}
	require.NoError(e.SetData(ContentTypeProtobuf, payload))

	b, err := Protobuf.Marshal(&e)
	require.NoError(err)
	var e2 event.Event
	require.NoError(Protobuf.Unmarshal(b, &e2))
	require.Equal(e, e2)

	payload2 := &pbv1.CloudEventAttributeValue{}
	require.NoError(e2.DataAs(payload2))
	require.True(payload2.GetCeBoolean())
}

func TestFromProto(t *testing.T) {
	tests := []struct {
		name    string
		proto   *pbv1.CloudEvent
		want    *event.Event
		wantErr bool
	}{{
		name: "happy binary json",
		proto: &pbv1.CloudEvent{
			SpecVersion: "1.0",
			Id:          "abc-123",
			Source:      "/source",
			Type:        "some.type",
			Attributes: map[string]*pbv1.CloudEventAttributeValue{
				"datacontenttype": {Attr: &pbv1.CloudEventAttributeValue_CeString{CeString: "application/json"}},
				"extra1":          {Attr: &pbv1.CloudEventAttributeValue_CeString{CeString: "extra1 value"}},
				"extra2":          {Attr: &pbv1.CloudEventAttributeValue_CeInteger{CeInteger: 2}},
				"extra3":          {Attr: &pbv1.CloudEventAttributeValue_CeBoolean{CeBoolean: true}},
				"extra4":          {Attr: &pbv1.CloudEventAttributeValue_CeUri{CeUri: "https://example.com"}},
			},
			Data: &pbv1.CloudEvent_BinaryData{
				BinaryData: []byte(`{"unit":"test"}`),
			},
		},
		want: func() *event.Event {
			out := event.New(cloudevents.VersionV1)
			out.SetID("abc-123")
			out.SetSource("/source")
			out.SetType("some.type")
			_ = out.SetData("application/json", map[string]interface{}{"unit": "test"})
			out.SetExtension("extra1", "extra1 value")
			out.SetExtension("extra2", 2)
			out.SetExtension("extra3", true)
			out.SetExtension("extra4", url.URL{Scheme: "https", Host: "example.com"})
			return &out
		}(),
		wantErr: false,
	}, {
		name: "happy text",
		proto: &pbv1.CloudEvent{
			SpecVersion: "1.0",
			Id:          "abc-123",
			Source:      "/source",
			Type:        "some.type",
			Attributes: map[string]*pbv1.CloudEventAttributeValue{
				"datacontenttype": {Attr: &pbv1.CloudEventAttributeValue_CeString{CeString: "text/plain"}},
				"extra1":          {Attr: &pbv1.CloudEventAttributeValue_CeString{CeString: "extra1 value"}},
				"extra2":          {Attr: &pbv1.CloudEventAttributeValue_CeInteger{CeInteger: 2}},
				"extra3":          {Attr: &pbv1.CloudEventAttributeValue_CeBoolean{CeBoolean: true}},
				"extra4":          {Attr: &pbv1.CloudEventAttributeValue_CeUri{CeUri: "https://example.com"}},
			},
			Data: &pbv1.CloudEvent_TextData{
				TextData: `this is some text with a "quote"`,
			},
		},
		want: func() *event.Event {
			out := event.New(cloudevents.VersionV1)
			out.SetID("abc-123")
			out.SetSource("/source")
			out.SetType("some.type")
			_ = out.SetData("text/plain", `this is some text with a "quote"`)
			out.SetExtension("extra1", "extra1 value")
			out.SetExtension("extra2", 2)
			out.SetExtension("extra3", true)
			out.SetExtension("extra4", url.URL{Scheme: "https", Host: "example.com"})
			return &out
		}(),
		wantErr: false,
	}, {
		name: "happy json as text",
		proto: &pbv1.CloudEvent{
			SpecVersion: "1.0",
			Id:          "abc-123",
			Source:      "/source",
			Type:        "some.type",
			Attributes: map[string]*pbv1.CloudEventAttributeValue{
				"datacontenttype": {Attr: &pbv1.CloudEventAttributeValue_CeString{CeString: "application/json"}},
				"extra1":          {Attr: &pbv1.CloudEventAttributeValue_CeString{CeString: "extra1 value"}},
				"extra2":          {Attr: &pbv1.CloudEventAttributeValue_CeInteger{CeInteger: 2}},
				"extra3":          {Attr: &pbv1.CloudEventAttributeValue_CeBoolean{CeBoolean: true}},
				"extra4":          {Attr: &pbv1.CloudEventAttributeValue_CeUri{CeUri: "https://example.com"}},
			},
			Data: &pbv1.CloudEvent_TextData{
				TextData: `{"unit":"test"}`,
			},
		},
		want: func() *event.Event {
			out := event.New(cloudevents.VersionV1)
			out.SetID("abc-123")
			out.SetSource("/source")
			out.SetType("some.type")
			_ = out.SetData("application/json", `{"unit":"test"}`)
			out.SetExtension("extra1", "extra1 value")
			out.SetExtension("extra2", 2)
			out.SetExtension("extra3", true)
			out.SetExtension("extra4", url.URL{Scheme: "https", Host: "example.com"})
			return &out
		}(),
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromProto(tt.proto)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromProto() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromProto() got = %v, want %v", got, tt.want)
			}
		})
	}
}
