package grpc

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestBuildGRPCOptionsFromFlags(t *testing.T) {
	file, err := os.CreateTemp("", "grpc-config-test-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	cases := []struct {
		name             string
		config           string
		expectedOptions  *GRPCOptions
		expectedErrorMsg string
	}{
		{
			name:             "empty config",
			config:           "",
			expectedErrorMsg: "url is required",
		},
		{
			name:             "tls config without clientCertFile",
			config:           "{\"url\":\"test\",\"clientCertFile\":\"test\"}",
			expectedErrorMsg: "either both or none of clientCertFile and clientKeyFile must be set",
		},
		{
			name:             "tls config without caFile",
			config:           "{\"url\":\"test\",\"clientCertFile\":\"test\",\"clientKeyFile\":\"test\"}",
			expectedErrorMsg: "setting clientCertFile and clientKeyFile requires caFile",
		},
		{
			name:   "customized options",
			config: "{\"url\":\"test\"}",
			expectedOptions: &GRPCOptions{
				URL: "test",
			},
		},
		{
			name:   "customized options with yaml format",
			config: "url: test",
			expectedOptions: &GRPCOptions{
				URL: "test",
			},
		},
		{
			name:   "customized options with ca",
			config: "{\"url\":\"test\",\"caFile\":\"test\"}",
			expectedOptions: &GRPCOptions{
				URL:    "test",
				CAFile: "test",
			},
		},
		{
			name:   "customized options with ca",
			config: "{\"url\":\"test\",\"caFile\":\"test\",\"clientCertFile\":\"test\",\"clientKeyFile\":\"test\"}",
			expectedOptions: &GRPCOptions{
				URL:            "test",
				CAFile:         "test",
				ClientCertFile: "test",
				ClientKeyFile:  "test",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := os.WriteFile(file.Name(), []byte(c.config), 0644); err != nil {
				t.Fatal(err)
			}

			options, err := BuildGRPCOptionsFromFlags(file.Name())
			if err != nil {
				if err.Error() != c.expectedErrorMsg {
					t.Errorf("unexpected err %v", err)
				}
			}

			if !reflect.DeepEqual(options, c.expectedOptions) {
				t.Errorf("unexpected options %v", options)
			}
		})
	}
}
