#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: ./deploy <relativePath>"
  exit 1
fi

RELATIVE_PATH=$1

if [ -d $RELATIVE_PATH/applications ]; then
  cd $RELATIVE_PATH/applications
  for DIR in */ ; do
    if [ -d "$DIR" ]; then
      SANITIZED_TAG=$(echo $RELATIVE_PATH$DIR | sed 's|^\./||' | tr '/' '-' | sed 's/-$//')
      docker build $DIR -t $SANITIZED_TAG
      docker tag $SANITIZED_TAG:latest deecaad/devopsk8s:$SANITIZED_TAG
      docker push deecaad/devopsk8s:$SANITIZED_TAG
    fi
  done
  cd ..
else
  cd $RELATIVE_PATH
  SANITIZED_TAG=$(echo $RELATIVE_PATH | sed 's|^\./||' | tr '/' '-' | sed 's/-$//')
  
  docker build . -t $SANITIZED_TAG
  docker tag $SANITIZED_TAG:latest deecaad/devopsk8s:$SANITIZED_TAG
  docker push deecaad/devopsk8s:$SANITIZED_TAG
fi

if [ -d ./k8s/volumes ]; then
  kubectl apply -f ./k8s/volumes/
fi

if [ -d ./k8s/manifests ]; then
  kubectl apply -f ./k8s/manifests/
fi

if [ -d ./manifests ]; then
  kubectl apply -f ./manifests/
fi