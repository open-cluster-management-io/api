//go:build tools
// +build tools

// go mod won't pull in code that isn't depended upon, but we have some code we don't depend on from code that must be included
// for our build to work.
package dependencymagnet

import (
	_ "github.com/gogo/protobuf/gogoproto"
	_ "github.com/gogo/protobuf/proto"
	_ "github.com/gogo/protobuf/sortkeys"
	_ "github.com/openshift/build-machinery-go"
	_ "github.com/spf13/pflag"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/api/rbac/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1beta1"
	_ "k8s.io/apimachinery/pkg/apis/testapigroup/v1"
	_ "k8s.io/code-generator"
	_ "k8s.io/code-generator/cmd/go-to-protobuf/protoc-gen-gogo"
)
