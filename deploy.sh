#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: ./deploy <relativePath>"
  exit 1
fi

RELATIVE_PATH=$1

cd $RELATIVE_PATH

SANITIZED_TAG=$(echo $RELATIVE_PATH | sed 's|^\./||' | tr '/' '-' | sed 's/-$//')

docker build . -t $SANITIZED_TAG
docker tag $SANITIZED_TAG:latest deecaad/devopsk8s:$SANITIZED_TAG
docker push deecaad/devopsk8s:$SANITIZED_TAG

# Apply deployment.yaml if it exists
if [ -f ./manifests/deployment.yaml ]; then
  kubectl apply -f ./manifests/deployment.yaml
else
  echo "Warning: ./manifests/deployment.yaml not found, skipping deployment."
fi

# Apply service.yaml if it exists
if [ -f ./manifests/service.yaml ]; then
  kubectl apply -f ./manifests/service.yaml
else
  echo "Warning: ./manifests/service.yaml not found, skipping service."
fi
