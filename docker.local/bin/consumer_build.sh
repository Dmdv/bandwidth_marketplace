#!/bin/sh
set -e

GIT_COMMIT=$(git rev-list -1 HEAD)
echo $GIT_COMMIT

docker build --build-arg GIT_COMMIT=$GIT_COMMIT -f docker.local/cons-dockerfile . -t consumer

for i in $(seq 1 6);
do
  CONSUMER=$i docker-compose -p consumer$i -f docker.local/cons-docker-compose.yml build --force-rm
done

docker.local/bin/sync_clock.sh
