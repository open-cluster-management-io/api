#!/bin/bash

CONTROLLER_GEN_VERSION=v0.6.0
CONTROLLER_GEN=_output/bin/controller-gen-$CONTROLLER_GEN_VERSION
controller_gen_dir=$(dirname $CONTROLLER_GEN)
GOHOSTOS=$(go env GOHOSTOS)
GOHOSTARCH=$(go env GOHOSTARCH)

if ! which $CONTROLLER_GEN > /dev/null;
  then echo "Installing controller-gen ...";
  mkdir -p $controller_gen_dir;
  curl -s -f -L https://github.com/openshift/kubernetes-sigs-controller-tools/releases/download/$CONTROLLER_GEN_VERSION/controller-gen-$GOHOSTOS-$GOHOSTARCH -o $CONTROLLER_GEN;
  chmod +x $CONTROLLER_GEN;
fi

$CONTROLLER_GEN schemapatch:manifests="./crdsv1beta1" paths="./operator/v1" 'output:dir="./crdsv1beta1"'
$CONTROLLER_GEN schemapatch:manifests="./crdsv1beta1" paths="./work/v1" 'output:dir="./crdsv1beta1"'
$CONTROLLER_GEN schemapatch:manifests="./crdsv1beta1" paths="./cluster/v1alpha1" 'output:dir="./crdsv1beta1"'

