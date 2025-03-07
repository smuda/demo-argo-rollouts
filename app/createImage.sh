#!/usr/bin/env bash

docker build image \
  --no-cache \
  -t docker.io/smuda/demo-rollouts || exit 1

docker push docker.io/smuda/demo-rollouts
