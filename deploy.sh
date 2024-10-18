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
    
    sed -i "s|\(newName: deecaad/devopsk8s:$SANITIZED_TAG\).*|\1-$VERSION|g" ../k8s/kustomization.yaml
  fi
done
cd ..

if [ -d ./k8s/ ]; then
  kubectl apply -k ./k8s/
fi