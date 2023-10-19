#!/bin/bash

if [ ! -f ./_output/tools/bin/yq ]; then
    mkdir -p ./_output/tools/bin
    curl -s -f -L https://github.com/mikefarah/yq/releases/download/v4.35.2/yq_$(go env GOHOSTOS)_$(go env GOHOSTARCH) -o ./_output/tools/bin/yq
    chmod +x ./_output/tools/bin/yq
fi

FILES="cluster/v1/*.crd.yaml
cluster/v1alpha1/*.crd.yaml
cluster/v1beta1/*.crd.yaml
cluster/v1beta2/*.crd.yaml
work/v1/*.crd.yaml
work/v1alpha1/*crd.yaml
operator/v1/*.crd.yaml
addon/v1alpha1/*.crd.yaml
"

FAILS=false
for f in $FILES
do
    if [[ $(./_output/tools/bin/yq .spec.validation.openAPIV3Schema.properties.metadata.description $f) != "null" ]]; then
        echo "Error: cannot have a metadata description in $f"
        FAILS=true
    fi

    if [[ $(./_output/tools/bin/yq .spec.preserveUnknownFields $f) != "false" ]]; then
        echo "Error: pruning not enabled (spec.preserveUnknownFields != false) in $f"
        FAILS=true
    fi
done

if [ "$FAILS" = true ] ; then
    exit 1
fi
