#!/bin/sh

PWD=`pwd`
PROVIDER_DIR=$(basename "$PWD")
PROVIDER_ID=`echo my directory $PROVIDER_DIR | sed -e 's/.*\(.\)$/\1/'`

echo Starting provider$PROVIDER_ID ...

PROVIDER=$PROVIDER_ID docker-compose -p provider$PROVIDER_ID -f ../prov-b0docker-compose.yml up