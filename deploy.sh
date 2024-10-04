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

kubectl create deployment $SANITIZED_TAG --image=deecaad/devopsk8s:$SANITIZED_TAG