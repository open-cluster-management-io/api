package payload

import (
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func TestDecodeSpecResyncRequest(t *testing.T) {
	evt := cloudevents.NewEvent()
	if err := evt.SetData(cloudevents.ApplicationJSON, []byte("{\"resourceVersions\":[{\"resourceID\":\"123\",\"resourceVersion\":3}]}")); err != nil {
		t.Errorf("failed to set data %v", err)
	}

	versions, err := DecodeSpecResyncRequest(evt)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(versions.Versions) != 1 {
		t.Errorf("unexpected versions %v", versions)
	}
}

func TestDecodeStatusResyncRequest(t *testing.T) {
	evt := cloudevents.NewEvent()
	if err := evt.SetData(cloudevents.ApplicationJSON, []byte("{\"statusHashes\":[{\"resourceID\":\"123\",\"statusHash\":\"1a2b\"}]}")); err != nil {
		t.Errorf("failed to set data %v", err)
	}

	hashes, err := DecodeStatusResyncRequest(evt)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(hashes.Hashes) != 1 {
		t.Errorf("unexpected versions %v", hashes)
	}
}
