#!/bin/sh
set -e

GIT_COMMIT=$(git rev-list -1 HEAD)
echo $GIT_COMMIT

docker build --build-arg GIT_COMMIT=$GIT_COMMIT -f docker.local/magma-dockerfile . -t magma

docker-compose -p magma -f docker.local/magma-docker-compose.yml build --force-rm

docker.local/bin/sync_clock.sh