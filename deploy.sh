#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: ./deploy <relativePath> [version]"
  exit 1
fi

RELATIVE_PATH=$1
VERSION=${2:-1}

cd $RELATIVE_PATH/applications
for DIR in */ ; do
  if [ -d "$DIR" ]; then
    SANITIZED_TAG=$(echo $RELATIVE_PATH$DIR | sed 's|^\./||' | tr '/' '-' | sed 's/-$//')
    docker build --tag "deecaad/devopsk8s:$SANITIZED_TAG-$VERSION" $DIR
    docker push deecaad/devopsk8s:$SANITIZED_TAG-$VERSION
    
    #find ../k8s/ -type f -name "*.yaml" -exec sed -i "s|\(image: deecaad/devopsk8s:$SANITIZED_TAG\).*|\1-$VERSION|g" {} +
  fi
done
cd ..

if [ -d ./k8s/ ]; then
  kubectl apply -k ./k8s/
fi

#if [ -d ./k8s/namespaces ]; then
#  kubectl apply -f ./k8s/namespaces/
#fi
#
#if [ -f ./k8s/secrets/secret.yaml ]; then
#  kubectl apply -f ./k8s/secrets/secret.yaml
#fi
#
#if [ -d ./k8s/crons ]; then
#  kubectl apply -f ./k8s/crons/
#fi
#
#if [ -d ./k8s/volumes ]; then
#  kubectl apply -f ./k8s/volumes/
#fi
#
#if [ -d ./k8s/manifests ]; then
#  kubectl apply -f ./k8s/manifests/
#fi
#
#if [ -d ./manifests ]; then
#  kubectl apply -f ./manifests/
#fi