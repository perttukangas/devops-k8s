#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: ./clean <relativePath>"
  exit 1
fi

RELATIVE_PATH=$1

cd $RELATIVE_PATH

if [ -d ./manifests ]; then
  kubectl delete -f ./manifests/
fi

if [ -d ./k8s/manifests ]; then
  kubectl delete -f ./k8s/manifests/
fi

if [ -d ./k8s/volumes ]; then
  kubectl delete -f ./k8s/volumes/
fi

if [ -d ./k8s/middleware ]; then
  kubectl delete -f ./k8s/middleware/
fi

if [ -d ./k8s/namespaces ]; then
  kubectl delete -f ./k8s/namespaces/
fi