#!/usr/bin/env bash
set -eax

APP_HOSTNAME=${APP_HOSTNAME:-test.planet4.dev}

SCALE_OPENRESTY=${SCALE_OPENRESTY:-2}
SCALE_APP=${SCALE_APP:-2}

docker-compose -f ${DOCKER_COMPOSE_FILE:-docker-compose.yml} up -d \
  --scale openresty=$SCALE_OPENRESTY \
  --scale app=$SCALE_APP