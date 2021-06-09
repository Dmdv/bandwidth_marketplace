#!/bin/sh

echo Starting magma ...

docker-compose -p magma -f ../magma-docker-compose.yml up
