#!/bin/bash

set -eo pipefail

BASEDIR=$(realpath $(dirname "$0")/..)

rm -rf "$BASEDIR"/tmp
mkdir -p "$BASEDIR"/tmp/api/operator.kyma-project.io

ln -s "$BASEDIR"/api/v1alpha1 "$BASEDIR"/tmp/api/operator.kyma-project.io/v1alpha1

"$BASEDIR"/bin/client-gen \
  --clientset-name versioned \
  --input-base "" \
  --input github.com/kyma-project/kyma-operator/operator/tmp/api/operator.kyma-project.io/v1alpha1 \
  --go-header-file "$BASEDIR"/hack/boilerplate.go.txt \
  --output-package github.com/kyma-project/kyma-operator/operator/pkg/client/clientset \
  --output-base "$BASEDIR"/tmp/pkg/client

"$BASEDIR"/bin/lister-gen \
  --input-dirs github.com/kyma-project/kyma-operator/operator/tmp/api/operator.kyma-project.io/v1alpha1 \
  --go-header-file "$BASEDIR"/hack/boilerplate.go.txt \
  --output-package github.com/kyma-project/kyma-operator/operator/pkg/client/listers \
  --output-base "$BASEDIR"/tmp/pkg/client

"$BASEDIR"/bin/informer-gen \
  --input-dirs github.com/kyma-project/kyma-operator/operator/tmp/api/operator.kyma-project.io/v1alpha1 \
  --versioned-clientset-package github.com/kyma-project/kyma-operator/operator/pkg/client/clientset/versioned \
  --listers-package github.com/kyma-project/kyma-operator/operator/pkg/client/listers \
  --go-header-file "$BASEDIR"/hack/boilerplate.go.txt \
  --output-package github.com/kyma-project/kyma-operator/operator/pkg/client/informers \
  --output-base "$BASEDIR"/tmp/pkg/client

find "$BASEDIR"/tmp/pkg/client -name "*.go" -exec \
  sed -i "" "s#github\.com/kyma-project/kyma-operator/operator/tmp/api/operator\.kyma-project\.io/v1alpha1#github\.com/kyma-project/kyma-operator/operator/api/v1alpha1#g" \
  {} +

rm -rf "$BASEDIR"/pkg/client && mkdir -p pkg
mv "$BASEDIR"/tmp/pkg/client/github.com/kyma-project/kyma-operator/operator/pkg/client "$BASEDIR"/pkg/client

rm -rf "$BASEDIR"/tmp