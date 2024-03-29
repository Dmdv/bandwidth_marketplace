#!/bin/sh
set -e

GIT_COMMIT=$(git rev-list -1 HEAD)
echo $GIT_COMMIT

docker build --build-arg GIT_COMMIT=$GIT_COMMIT -f docker.local/prov-dockerfile . -t provider

for i in $(seq 1 6);
do
  PROVIDER=$i docker-compose -p provider$i -f docker.local/prov-docker-compose.yml build --force-rm
done

docker.local/bin/sync_clock.sh
