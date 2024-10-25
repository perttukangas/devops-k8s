#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: ./clean <relativePath>"
  exit 1
fi

RELATIVE_PATH=$1

cd $RELATIVE_PATH

if [ -d ./k8s ]; then
  kubectl delete -k ./k8s/
fi

if [ -d ./k8s/overlays/prod ]; then
  kubectl delete -k ./k8s/overlays/prod
fi

if [ -d ./k8s/overlays/stating ]; then
  kubectl delete -k ./k8s/overlays/staging
fi