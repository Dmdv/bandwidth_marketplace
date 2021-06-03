#!/bin/sh

PROVIDER_DIR=$(basename "$PWD")
PROVIDER_ID=`echo my directory $PROVIDER_DIR | sed -e 's/.*\(.\)$/\1/'`

echo Starting provider$PROVIDER_ID ...

PROVIDER=$PROVIDER_ID docker-compose -p provider$PROVIDER_ID -f ../prov-docker-compose.yml up