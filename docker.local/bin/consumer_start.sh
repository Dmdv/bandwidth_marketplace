#!/bin/sh

PWD=`pwd`
CONSUMER_DIR=$(basename "$PWD")
CONSUMER_ID=`echo my directory $CONSUMER_DIR | sed -e 's/.*\(.\)$/\1/'`

echo Starting consumer$CONSUMER_ID ...

CONSUMER=$CONSUMER_ID docker-compose -p consumer$CONSUMER_ID -f ../cons-docker-compose.yml up
