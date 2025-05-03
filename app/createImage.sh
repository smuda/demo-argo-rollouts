#!/usr/bin/env bash

push image || exit 1
go mod tidy
popd || exit 1

docker build image \
  --no-cache \
  -t docker.io/smuda/demo-rollouts || exit 1

docker push docker.io/smuda/demo-rollouts
